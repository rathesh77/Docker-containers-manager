#!/bin/bash

line=$1

sqlite3 db "select from machine where id = $line;"

#redis-cli get "$line"