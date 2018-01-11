FROM golang:1.9.2
MAINTAINER Vadim Chernov <dimuls@yandex.ru>

WORKDIR /go/src/github.com/someanon/rkn-bypasser
COPY . .

RUN go-wrapper download
RUN go-wrapper install

WORKDIR /

RUN rm -rf /go/src

ENV ADDR 0.0.0.0:8000

RUN groupadd -r rkn-bypasser
RUN useradd -r -g rkn-bypasser rkn-bypasser

USER rkn-bypasser

EXPOSE 8000

ENTRYPOINT /go/bin/rkn-bypasser