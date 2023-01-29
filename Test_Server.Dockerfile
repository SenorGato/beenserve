FROM golang:1.20-rc-bullseye

WORKDIR /beenserve
COPY . .
WORKDIR /beenserve/server/
RUN go build
#CD ["serve"]
