FROM golang:1.20-rc-bullseye as build-env

WORKDIR /
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . ./beenserve
WORKDIR /beenserve/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/server

COPY .env .
COPY ./client/ ../client/
ENTRYPOINT ["server"]
