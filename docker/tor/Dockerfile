FROM alpine:3.8 AS builder

ARG TOR_VERSION=0.3.2.11

RUN apk update && \
apk add make automake autoconf gcc libtool curl libevent-dev \
libssl1.0 musl musl-dev libgcc openssl openssl-dev openssh && \
curl https://dist.torproject.org/tor-${TOR_VERSION}.tar.gz | tar xz

WORKDIR tor-${TOR_VERSION}

RUN ./configure && make && make install

RUN echo "Log notice stdout" >> /etc/torrc && \
echo "SocksPort 0.0.0.0:9150" >> /etc/torrc && \
echo 'DataDirectory /var/lib/tor' >> /etc/torrc



FROM alpine:3.8

RUN apk update && \
    apk add libevent libssl1.0 musl openssl openssh

WORKDIR /

EXPOSE 9150

COPY --from=builder /usr/local/bin/tor /
COPY --from=builder /etc/torrc /torrc

ENTRYPOINT ["/tor", "-f", "/torrc"]
