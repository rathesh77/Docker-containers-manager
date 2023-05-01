#!/bin/bash

nodes=$(sqlite3 db "select * from node");

echo $nodes | while IFS= read -r node ; do 

    IFS='|'

    fields=($node)

    id=${fields[0]}
    name=${fields[1]}
    network=${fields[3]}
    mask=${fields[4]}
    IFS=' '

    if ping -q -c 1 $network &>/dev/null
    then

        echo "node at $network is reachable"
        machines=$(sqlite3 db "select * from machine where node_id=$id" | tr -d "\n[:blank:]")
        if [ "$machines" != "" ] 
        then
            echo $machines | while IFS= read -r machine ; do 
                IFS='|'
                fields=($machine)
                id=${fields[0]}
                curl -v --location "http://${network}:3001/healthcheck" \
                -H 'Content-Type: application/json' \
                -d  "{\"id\": \"$id\"}"
            done
        fi
    fi

done