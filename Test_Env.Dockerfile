FROM golang:1.20-rc-bullseye as build-env

WORKDIR /
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . ./beenserve
WORKDIR /beenserve/server

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/server

FROM scratch
COPY --from=build-env /go/bin/server /go/bin/server
WORKDIR /go/bin/
COPY test.env .
ENTRYPOINT ["/go/bin/server"]
