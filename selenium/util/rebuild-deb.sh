#!/bin/bash
set -e

usage()
{
cat << EOF
usage: $0 options

OPTIONS:
   -h      Show this message
   -f      Filename
   -r      Remove version postfix
   -a      Add version postfix
EOF
}

filename=''
remove_postfix=''
add_postfix=''

while getopts 'f:r:a:' OPTION
do
     case $OPTION in
         h)
             usage
             exit 1
             ;;
         f)
             filename=$OPTARG
             ;;
         r)
             remove_postfix=$OPTARG
             ;;
         a)
             add_postfix=$OPTARG
             ;;
         ?)
             usage
             exit
             ;;
     esac
done

if [[ -z $filename ]]; then
     usage
     exit 1
fi
set -x

cp "$filename" .
spackage=$(echo "$filename" | awk -F '/' '{print $NF}')
package=$(echo $spackage | awk -F '_' '{print $1}')
version=$(echo $spackage | awk -F '_' '{print $2}')
rm -Rf "$package" control changelog

# Set variables
export DH_ALWAYS_EXCLUDE="CVS:.svn:.git"

# Extract package content
dpkg-deb -x $spackage $package/

# Extract package control
dpkg-deb -e $spackage $package/DEBIAN

# Extract package changelog
changelog=$(find $package/usr/share/doc -name 'changelog*.gz' 2>/dev/null| head -n 1)
if [ -n $changelog -a -f $changelog ]; then
    rm -f $changelog
fi
old_version=$version
if [ -n $remove_postfix ]; then
    version=$(echo $version | sed -e "s|+$remove_postfix||g")
    old_version=$version"+"$remove_postfix
fi 
new_version=$version
if [ -n $add_postfix ]; then
    new_version=$new_version"+"$add_postfix
fi
dch --create --force-distribution --distribution unstable --package "$package" --newversion "$new_version" -c changelog 'Package rebuilt'

# Fix control
section=$(awk '/Section:/ {print $NF}' $package/DEBIAN/control 2>/dev/null || echo 'misc')
priority=$(awk '/Priority:/ {print $NF}' $package/DEBIAN/control 2>/dev/null || echo 'extra')
echo "$(echo $spackage | sed "s/$old_version/$new_version/") $section $priority" > files
echo "Source: $package" > control
grep -iv "Source:" $package/DEBIAN/control >> control
cp control $package/DEBIAN/control

# Build package
sed -i "s/$old_version/$new_version/g" $package/DEBIAN/control
dpkg-deb -b $package/ ./

# Remove source package
rm -f "./$spackage"
