FROM golang:1.10.0-alpine3.7

RUN apk update && apk add --no-cache git wget xorriso syslinux

RUN mkdir /builder
ADD *.go /go/
RUN go get github.com/imdario/mergo
RUN go get gopkg.in/yaml.v2
RUN go build -o /usr/bin/build
WORKDIR /builder

CMD build