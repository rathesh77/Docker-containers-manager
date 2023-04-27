#!/bin/bash

docker stop $(docker ps --filter status=running -q)>/dev/null

sqlite3 db "delete from machine;delete from node;delete from cluster;"

docker rm $(docker ps --filter status=exited -q) >/dev/null