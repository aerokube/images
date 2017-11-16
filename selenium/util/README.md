# Rebuilding Debian Packages

1. Download package from APT-repository:
```
$ apt-get download google-chrome-stable=48.0.2564.109-1+suffix0
```
2. Rebuild package to new owner using script:
```
$ rebuild-deb.sh -f google-chrome-stable=48.0.2564.109-1+suffix0_amd64.deb -r 'suffix0' -a 'aerokube0'
```