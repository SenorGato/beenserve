FROM golang:1.20-rc-bullseye

WORKDIR /beenserve

WORKDIR /beenserve/server/
COPY ./test.env .
RUN go build
#CD ["serve"]
