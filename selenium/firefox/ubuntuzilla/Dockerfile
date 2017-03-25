FROM selenoid/base:1.0

ARG VERSION
ARG CLEANUP=false

RUN \
        ( \
            wget -O firefox-mozilla-build.deb https://master.dl.sourceforge.net/project/ubuntuzilla/mozilla/apt/pool/main/f/firefox-mozilla-build/firefox-mozilla-build_$VERSION-0ubuntu1_amd64.deb || \
            ( \ 
                wget -O firefox-mozilla-build.deb https://master.dl.sourceforge.net/project/ubuntuzilla/mozilla/apt/pool/main/f/firefox-mozilla-build/firefox-mozilla-build_$VERSION-0ubuntu1_i386.deb && \
                dpkg --add-architecture i386 && \
                apt-get update && \
                apt-get -y install libc6:i386 libncurses5:i386 libstdc++6:i386 libgtk2.0-0:i386 libasound2:i386 libdbus-glib-1-2:i386 libxt6:i386 \                
            ) \
        ) && \
        dpkg -i firefox-mozilla-build.deb && \
        firefox --version && \ # Just to test it launches successfully
        ($CLEANUP && apt-get clean && rm -Rf /tmp/*)