#!/bin/bash

docker stop $(docker ps --filter status=running -q)

sqlite3 db "delete from cluster;delete from machine;"

docker rm $(docker ps --filter status=exited -q)