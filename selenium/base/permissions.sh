#!/bin/bash
if ! id -un &> /dev/null; then
    export HOME=/home/selenium
    export USER_ID=$(id -u)
    export GROUP_ID=$(id -g)
    cat /etc/passwd | awk -v uid=${USER_ID} -v gid=${GROUP_ID} 'BEGIN { FS = OFS = ":" } ; $1=="selenium"{$3=uid; $4=gid}1' > ${HOME}/passwd
    export LD_PRELOAD=/usr/lib/libnss_wrapper.so
    export NSS_WRAPPER_PASSWD=${HOME}/passwd
    export NSS_WRAPPER_GROUP=/etc/group
fi
