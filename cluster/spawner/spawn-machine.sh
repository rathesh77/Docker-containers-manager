#!/bin/bash


docker run --expose "3000:3000" --name test-container alpine:3.14

filename = "$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')"

docker stats $filename

echo "machine $filename running on port 3000\n"