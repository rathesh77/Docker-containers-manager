#!/bin/bash

read -r line

sqlite3 db "select from machine where id = $line;"

#redis-cli get "$line"