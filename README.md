# Selenoid Docker Images
This repository contains [Docker](http://docker.com/) build files to be used for [Selenoid](http://github.com/aerokube/selenoid) project. You can find prebuilt images [here](https://hub.docker.com/u/selenoid/).

## How image are built

![layers](layers.png)

Each image consists of 3 or 4 layers:
1) **Base layer** - contains stuff needed in every image: Xvfb, fonts, cursor blinking fix, timezone definition and so on. This layer is always built manually.
2) **Optional Java layer** - contains latest Java Runtime Environment. Only needed for old Firefox versions incompatible with Geckodriver. This layer is always built manually.
3) **Browser layer** - contains browser binary. We create two versions: with APT cache and without it. The latter is then used to add driver layer.
4) **Driver layer** - contains either respective web driver binary or corresponding Selenium server version.

Building procedure is automated with shell scripts ```selenium/build-dev.sh``` and ```selenium/build.sh``` that generate Dockerfile and then create browser and driver layers respectively. Before push each image is tested with these [tests](https://github.com/aerokube/selenoid-container-tests).

## Image information
Moved to: http://aerokube.com/selenoid/latest/#_browser_image_information
