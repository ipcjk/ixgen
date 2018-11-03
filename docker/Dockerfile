FROM alpine:latest
RUN apk update && apk upgrade && apk add unzip ca-certificates openssl
WORKDIR /ixgen
RUN wget -O /tmp/ixgen.tar.gz https://github.com/ipcjk/ixgen/releases/download/0.5/release.tar.gz && tar xfz /tmp/ixgen.tar.gz --exclude release/bgpq3.mac --exclude release/ixapiserver.exe --exclude release/ixapiserver.mac --exclude release/ixgen.mac --exclude release/ixgen.exe --strip 1 && rm /tmp/ixgen.tar.gz
RUN mv /ixgen/ixgen.linux /ixgen/ixgen && mv /ixgen/ixapiserver.linux /ixgen/ixapiserver
# RUN mv /ixgen/release/configuration/peering.ini /ixgen/release/configuration/example.ini
# E.g. add your peering.ini from your source
MAINTAINER Joerg Kost <jk@ip-clear.de>
#CMD ["/ixgen/release/ixgen.linux"]
