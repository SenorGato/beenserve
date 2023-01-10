#!/usr/bin/env bash

nix-shell default.nix

while getopts 'tbr:' OPTION; do
    case "$OPTION" in
        t)
            nix-shell test.nix
            ;;
        b)

            docker cp "$PRODUCT_DATABASE_INIT" tealacarte:/init.sql
            docker cp "$PRODUCT_DATA_INIT" tealacarte:/data.sql
            docker build ./data -t tealacarte:"$TLC_VERSION" -f data/Dockerfile
            ;;
        r)
            docker-compose up -d -f ./data/docker-compose.yml
            cd ./src || exit
            go run main.go
            ;;
        *)
            echo "Invalid flag"
            exit 1
            ;;
    esac
done
