FROM golang:1.20-rc-bullseye

#RUN apt-get -y update && apt-get -y upgrade
#RUN apt-get -y install git && rm -rf /var/cache/apk/*
WORKDIR /beenserve
RUN git clone https://github.com/senorgato/beenserve.git . 
COPY ./test.env ./server
#RUN go build
#RUN go build -v -o /usr/local/bin/app ./...
#
#CMD ["app"]
