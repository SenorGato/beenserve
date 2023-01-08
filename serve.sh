#!/bin/sh

export DATABASE_URL="postgres://tealacarte:smoke@localhost:5432/tealacarte"
docker build . -t tealacarte -f data/Dockerfile
docker-compose up -d -f ./data/docker-compose.yml
cd ./src || exit
go build
./src
