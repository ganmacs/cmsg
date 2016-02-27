FROM golang:1.6.0

MAINTAINER Yuta Iwamwa <ganmacs@gmail.com>

RUN go get github.com/motemen/ghq

WORKDIR /go/cmsg
