#!/bin/bash

id=$1

echo "container_restart:$1";

docker restart $id