#!/bin/sh

ip=$(/usr/sbin/ifconfig | grep -A 10  -e "enp0s3" | grep  -e '^ *inet .*$' | sed -E 's/( *(netmask|broadcast) [0-9\.]+)+//' | tr -d '[:lower:]|[:blank:]' | head -n 1)

nmap -T5 -sn -oN scan.txt $ip/24

default_interface=$(ip route | grep '^default' | awk '/default/ {print $5}')

/usr/sbin/ifconfig | grep -A 10  -e "$default_interface" | grep  -e '^ *inet .*$' | sed -E 's/( *(netmask|inet) [0-9\.]+)+//' | tr -d '[:lower:]|[:blank:]'

cat scan.txt | grep 'for' | tr -d 'Nmap scan report for ' | while IFS= read -r addr ; do
    echo "$addr"

    response=$(curl -sSf -m 5 "http://$addr:3001/healthcheck" 2>&1 | grep '401')
    if [ "$response" != "" ] 
    then
        echo "$addr: kubelet is running"
        sudo bash ./spawner/spawn-node.sh $addr
    else
        echo "$addr: no heartbeat"
    fi
done

rm scan.txt