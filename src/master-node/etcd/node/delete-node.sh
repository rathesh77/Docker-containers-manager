#!/bin/bash

sqlite3 db "delete from node where network=\"$1\""