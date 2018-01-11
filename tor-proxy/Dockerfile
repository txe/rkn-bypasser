FROM alpine:3.7
MAINTAINER Vadim Chernov <dimuls@yandex.ru>

ENV TOR_VERSION 0.3.2.9

RUN apk update
RUN apk add make automake autoconf gcc libtool curl libevent-dev libssl1.0 musl musl-dev libgcc openssl openssl-dev openssh
RUN curl https://dist.torproject.org/tor-${TOR_VERSION}.tar.gz | tar xz

WORKDIR tor-${TOR_VERSION}

RUN ./configure
RUN make
RUN make install

WORKDIR /

RUN rm -rf tor-${TOR_VERSION}
RUN apk del make automake autoconf gcc

RUN echo "Log notice stdout" >> /etc/torrc
RUN echo "SocksPort 0.0.0.0:9150" >> /etc/torrc

RUN addgroup -S tor
RUN adduser -S -g tor tor

USER tor

EXPOSE 9150

ENTRYPOINT /usr/local/bin/tor -f /etc/torrc