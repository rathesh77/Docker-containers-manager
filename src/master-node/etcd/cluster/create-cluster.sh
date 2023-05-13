#!/bin/bash

line=$1

sqlite3 db "insert into cluster(name) VALUES(\"$line\");"