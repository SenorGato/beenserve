FROM golang:1.20-rc-bullseye as build-env

WORKDIR /
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . ./beenserve-admin
WORKDIR /beenserve-admin/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/server

FROM scratch
COPY --from=build-env /go/bin/server /go/bin/server
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /go/bin/
COPY ./client/ ./client/
ENTRYPOINT ["/go/bin/server"]
