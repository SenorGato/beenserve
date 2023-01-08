#!/bin/sh
#docker run --name postgres-db -e POSTGRES_PASSWORD=docker -p 5432:5432 -d postgres 
docker cp teaalacarte.sql b942dfb83a63:/
docker exec -it b942dfb83a63 psql -U postgres -W postgres -f teaalacarte.sql -q
