# Browser Images
[![Build Status](https://github.com/aerokube/images/workflows/build/badge.svg)](https://github.com/aerokube/images/actions?query=workflow%3Abuild)
[![Release](https://img.shields.io/github/release/aerokube/images.svg)](https://github.com/aerokube/images/releases/latest)

This repository contains [Docker](http://docker.com/) build files to be used for [Selenoid](http://github.com/aerokube/selenoid) and [Moon](http://github.com/aerokube/moon) projects. You can find prebuilt images [here](https://hub.docker.com/u/selenoid/).

## Download Statistics

### Firefox: [![Firefox Docker Pulls](https://img.shields.io/docker/pulls/selenoid/firefox.svg)](https://hub.docker.com/r/selenoid/firefox)

### Chrome: [![Chrome Docker Pulls](https://img.shields.io/docker/pulls/selenoid/chrome.svg)](https://hub.docker.com/r/selenoid/chrome)

### Opera: [![Opera Docker Pulls](https://img.shields.io/docker/pulls/selenoid/opera.svg)](https://hub.docker.com/r/selenoid/opera)

### Android: [![Android Docker Pulls](https://img.shields.io/docker/pulls/selenoid/android.svg)](https://hub.docker.com/r/selenoid/android)

## How images are built

![layers](layers.png)

Each image consists of 3 or 4 layers:
1) **Base layer** - contains stuff needed in every image: Xvfb, fonts, cursor blinking fix, timezone definition and so on. This layer is always built manually.
2) **Optional Java layer** - contains latest Java Runtime Environment. Only needed for old Firefox versions incompatible with Geckodriver. This layer is always built manually.
3) **Browser layer** - contains browser binary. We create two versions: with APT cache and without it. The latter is then used to add driver layer.
4) **Driver layer** - contains either respective web driver binary or corresponding Selenium server version.

## How to build images yourself

Building procedure is automated with Go binary with all Docker build files embedded to . To build it from source:

```
$ go get github.com/markbates/pkger/cmd/pkger
$ go generate github.com/aerokube/images
$ go build
```

To show help:

```
$ ./images --help
$ ./images firefox --help
```

* Before building images you can optionally clone tests repository:
```
$ git clone https://github.com/aerokube/selenoid-container-tests.git
```
These tests require Java and Maven 3 to be installed. Tests directory should be cloned to this repository parent directory:
```
images/ # <== this repo
selenoid-container-tests/ # <== optional tests repo
```
* To build a Firefox image use the following command:
```
$ ./images firefox -b 79.0+build1-0ubuntu0.18.04.1 -d 0.26.0 --selenoid-version 1.10.0 -t selenoid/firefox:79.0
```
Here `79.0+build1-0ubuntu0.18.04.1` is `firefox` package version for Ubuntu 18.04, `1.10.10` is [Selenoid](https://github.com/aerokube/selenoid/releases) version to use inside image (just use latest release version here), `selenoid/firefox:79.0` is Docker tag to be applied, `0.26.0` is [Geckodriver](http://github.com/mozilla/geckodriver/releases) version to use.

If you wish to automatically use the latest Selenoid and Geckodriver versions - just omit respective flags or use **latest** as value:
```
$ ./images firefox -b 79.0+build1-0ubuntu0.18.04.1 -t selenoid/firefox:79.0
```

If you wish to pack a local Debian package instead of APT - just replace package version with full path to **deb** file:
```
$ ./images firefox -b /path/to/firefox_79.0+build1-0ubuntu0.18.04.1_amd64.deb -t selenoid/firefox:79.0
``` 
It is important to use package files with full version specified name because automation scripts determine browser version by parsing package file name!

To run the tests after building the image add `--test` flag:

```
$ ./images firefox -b 79.0+build1-0ubuntu0.18.04.1 -t selenoid/firefox:79.0 --test
```

To push image after building add `--push` flag:

```
$ ./images firefox -b 79.0+build1-0ubuntu0.18.04.1 -t selenoid/firefox:79.0 --push
```

* To build a Chrome image use the following command:
```
$ ./images chrome -b 78.0.3904.97-1 -d 78.0.3904.70 -t selenoid/chrome:78.0
```
Here `78.0.3904.97-1` is `google-chrome-stable` package version for Ubuntu 18.04, `78.0.3904.70` is [Chromedriver](https://chromedriver.storage.googleapis.com/index.html) version, `selenoid/chrome:78.0` is Docker tag to be applied.  

If you wish to automatically use the latest [compatible](https://chromedriver.chromium.org/downloads/version-selection) Chromedriver version - just omit respective flag or use **latest** as value:
```
$ ./images chrome -b 78.0.3904.97-1 -d latest -t selenoid/chrome:78.0
```
* To build an Opera image use the following command:
```
$ ./images opera -b 64.0.3417.92 -d 77.0.3865.120 -t selenoid/opera:64.0
```
Here `64.0.3417.92` is `opera-stable` package version for Ubuntu 18.04, `77.0.3865.120` is [Operadriver](https://github.com/operasoftware/operachromiumdriver/releases) version, `64.0` is Docker tag to be applied.  

* To build a Yandex image use the following command:
```
$ ./images yandex -b 20.4.3.268-1 -d 20.4.3.321 -t selenoid/yandex-browser:20.4
```
Here `20.4.3.268-1` is `yandex-browser-beta` package version for Ubuntu 18.04, `20.4.3.321` is [Yandexdriver](https://github.com/yandex/YandexDriver/releases) Linux asset version, `20.4` is Docker tag to be applied.

* To build an Android image use the following command:
```
$ cd selenium
$ ./automate_android.sh
```
This command is interactive - just answer the questions and it will build an image for you. In order to bundle custom APK to image - put it to `selenium/android` directory before running the script.

## How to build images for non-default channels

Apart from the default stable release channel, the following ones are also supported:

| Browser | Channel | Package |
| :--- | :--- | :--- |
| firefox | beta | firefox [(PPA)](http://launchpad.net/~mozillateam/+archive/firefox-next/+packages) |
| firefox | dev | firefox-trunk [(PPA)](http://launchpad.net/~ubuntu-mozilla-daily/+archive/ppa/+packages) |
| firefox | esr | firefox-esr [(PPA)](http://launchpad.net/~mozillateam/+archive/ppa/+packages) |
| chrome | beta | google-chrome-beta |
| chrome | dev | google-chrome-unstable |
| opera | beta | opera-beta | |
| opera | dev | opera-developer | |

* To build an image for one of the channels above use the `--channel` flag as follows:
```
$ ./images firefox -b 72.0~a1~hg20191114r501767-0ubuntu0.18.04.1~umd1 --channel dev -t selenoid/firefox:72.0a1
```

## Image information
Moved to: http://aerokube.com/selenoid/latest/#_browser_image_information
