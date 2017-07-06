FROM alpine:latest
RUN apk update
RUN apk upgrade
RUN apk add unzip
RUN apk add ca-certificates 
RUN apk add openssl
#COPY ixgen.zip /tmp
RUN wget -O /tmp/ixgen.zip https://github.com/ipcjk/ixgen/releases/download/0.3b/release.zip
RUN unzip /tmp/ixgen.zip -d /ixgen
RUN rm /tmp/ixgen.zip
WORKDIR /ixgen/release
RUN /ixgen/release/ixgen.linux -buildcache
RUN mv /ixgen/release/configuration/peering.ini /ixgen/release/configuration/example.ini
# E.g. add your peering.ini from your source
MAINTAINER Joerg Kost <jk@ip-clear.de>
#CMD ["/ixgen/release/ixgen.linux"]
