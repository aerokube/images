# Selenoid Docker Containers
This repository contains [Docker](http://docker.com/) build files to be used for [Selenoid](http://github.com/aandryashin/selenoid) project. You can find prebuilt containers [here](https://hub.docker.com/u/selenoid/dashboard/).

## How containers are built

![layers](layers.png)

Each container consists of 3 or 4 layers:
1) **Base layer** - contains stuff needed in every container: Xvfb, fonts, cursor blinking fix, timezone definition and so on. This layer is always built manually.
2) **Optional Java layer** - contains latest Java Runtime Environment. Only needed for old Firefox versions incompatible with Geckodriver. This layer is always built manually.
3) **Browser layer** - contains browser binary. We create two versions: with APT cache and without it. The latter is then used to add driver layer.
4) **Driver layer** - container either respective web driver binary or corresponding Selenium server version.

Building procedure is automated with shell scripts ```selenium/build-dev.sh``` and ```selenium/build.sh``` that generate Dockerfile and then create browser and driver layers respectively. Before push each container is tested with these [tests](https://github.com/aerokube/selenoid-container-tests).

## Container information
### Firefox

| Container | Selenium version | Firefox version |
| --------------- | ---------------- | ---------------- |
| selenoid/firefox:3.6 | 2.20.0 | 3.6.16 i386 (dialogs may not work) |
| selenoid/firefox:4.0 | 2.20.0 | 4.0.1 i386 |
| selenoid/firefox:5.0 | 2.20.0 | 5.0.1 i386 |
| selenoid/firefox:6.0 | 2.20.0 | 6.0.2 i386 |
| selenoid/firefox:7.0 | 2.20.0 | 7.0.1 i386 |
| selenoid/firefox:8.0 | 2.20.0 | 8.0.1 i386 |
| selenoid/firefox:9.0 | 2.20.0 | 9.0.1 |
| selenoid/firefox:10.0 | 2.32.0 | 10.0.2 |
| selenoid/firefox:11.0 | 2.32.0 | 11.0 |
| selenoid/firefox:12.0 | 2.32.0 | 12.0 |
| selenoid/firefox:13.0 | 2.32.0 | 13.0 |
| selenoid/firefox:14.0 | 2.32.0 | 14.0.1 |
| selenoid/firefox:15.0 | 2.32.0 | 15.0.1 |
| selenoid/firefox:16.0 | 2.32.0 | 16.0.2 |
| selenoid/firefox:17.0 | 2.32.0 | 17.0.1 |
| selenoid/firefox:18.0 | 2.32.0 | 18.0.2 |
| selenoid/firefox:19.0 | 2.32.0 | 19.0.2 |
| selenoid/firefox:20.0 | 2.32.0 | 20.0 |
| selenoid/firefox:21.0 | 2.32.0 | 21.0 |
| selenoid/firefox:22.0 | 2.32.0 | 22.0 |
| selenoid/firefox:23.0 | 2.35.0 | 23.0.1 |
| selenoid/firefox:24.0 | 2.39.0 | 24.0 |
| selenoid/firefox:25.0 | 2.39.0 | 25.0.1 |
| selenoid/firefox:26.0 | 2.39.0 | 26.0 |
| selenoid/firefox:27.0 | 2.40.0 | 27.0.1 |
| selenoid/firefox:28.0 | 2.41.0 | 28.0 |
| selenoid/firefox:29.0 | 2.43.1 | 29.0.1 |
| selenoid/firefox:30.0 | 2.43.1 | 30.0 |
| selenoid/firefox:31.0 | 2.44.0 | 31.0 |
| selenoid/firefox:32.0 | 2.44.0 | 32.0.3 |
| selenoid/firefox:33.0 | 2.44.0 | 33.0.3 |
| selenoid/firefox:34.0 | 2.45.0 | 34.0.5 |
| selenoid/firefox:35.0 | 2.45.0 | 35.0.1 |
| selenoid/firefox:36.0 | 2.45.0 | 36.0.1 |
| selenoid/firefox:37.0 | 2.45.0 |  |
| selenoid/firefox:38.0 | 2.45.0 |  |
| selenoid/firefox:39.0 | 2.45.0 |  |
| selenoid/firefox:40.0 | 2.45.0 |  |
| selenoid/firefox:41.0 | 2.45.0 |  |
| selenoid/firefox:42.0 | 2.47.1 |  |
| selenoid/firefox:43.0 | 2.47.1 |  |
| selenoid/firefox:44.0 | 2.53.1 |  |
| selenoid/firefox:45.0 | 2.53.1 |  |
| selenoid/firefox:46.0 | 2.53.1 |  |
| selenoid/firefox:47.0 | 2.53.1 |  |
| selenoid/firefox:48.0 | 3.3.1 + GD 0.15.0 |  |
| selenoid/firefox:49.0 | 3.3.1 + GD 0.15.0 |  |
| selenoid/firefox:50.0 | 3.3.1 + GD 0.15.0 |  |
| selenoid/firefox:51.0 | 3.3.1 + GD 0.15.0 |  |
| selenoid/firefox:52.0 | 3.3.1 + GD 0.15.0 |  |

### Chrome

| Chrome version | Chromedriver version |
| -------------- | -------------------- |
| 29 | 2.6 |
| 30 | 2.8 |
| 31 | 2.9 |
| 32 | 2.9 |
| 33 | 2.10 |
| 34 | 2.10 |
| 35 | 2.10 |
| 36 | 2.12 |
| 37 | 2.12 |
| 38 | 2.13 |
| 39 | 2.14 |
| 40 | 2.15 |
| 41 | 2.15 |
| 42 | 2.16 |
| 43 | 2.20 |
| 44 | 2.20 |
| 45 | 2.20 |
| 46 | 2.21 |
| 47 | 2.21 |
| 48 | 2.21 |
| 49 | 2.22 |
| 50 | 2.22 |
| 51 | 2.23 |
| 52 | 2.24 |
| 53 | 2.26 |
| 54 | 2.27 |
| 55 | 2.28 |
| 56 | 2.28 |
| 57 | 2.28 |

### Opera

| Opera Presto version | Selenium version |
| --------------------- | ---------------- |
| 12.16 | 2.35.0 |

| Opera Blink version | Operadriver version |
| ------------------- | ------------------- |
| 15 | 0.2.2 |
| 16 | 0.2.2 |
| 17 | 0.2.2 |
| 18 | 0.2.2 |
| 19 | 0.2.2 |
| 20 | 0.2.2 |
| 21 | 0.2.2 |
| 22 | 0.2.2 |
| 23 | 0.2.2 |
| 24 | 0.2.2 |
| 25 | 0.2.2 |
| 26 | 0.2.2 |
| 27 | 0.2.2 |
| 28 | 0.2.2 |
| 29 | 0.2.2 |
| 30 | 0.2.2 |
| 32 | 0.2.2 |
| 33 | 0.2.2 |
| 34 | 0.2.2 |
| 35 | 0.2.2 |
| 36 | 0.2.2 |
| 37 | 0.2.2 |
| 38 | 0.2.2 |
| 39 | 0.2.2 |
| 40 | 0.2.2 |
| 41 | 0.2.2 |
| 42 | 0.2.2 |
| 43 | 0.2.2 |
| 44 | 0.2.2 |
