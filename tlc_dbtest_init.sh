#!/usr/bin/env sh

docker cp test_tlc_product_data.sql tealacarte:/
docker exec -it tealacarte psql -U tealacarte -a -f test_tlc_product_data.sql
