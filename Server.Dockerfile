FROM golang:1.20-rc-bullseye

WORKDIR /beenserve
RUN git clone https://github.com/senorgato/beenserve.git . 
WORKDIR /beenserve/server
COPY .server/test/test.env .
RUN go build
CMD ["./server"]