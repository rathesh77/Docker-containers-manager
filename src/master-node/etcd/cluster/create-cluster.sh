#!/bin/bash

line=$1

echo "cluster:$line"
sqlite3 db "insert into cluster(name) VALUES(\"$line\");"