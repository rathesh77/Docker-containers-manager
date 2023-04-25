#!/bin/bash

docker stop $(docker ps -a --format "{{.ID}}")

sqlite3 db "delete from cluster;delete from machine;"

docker kill $(ps | grep 'docker')