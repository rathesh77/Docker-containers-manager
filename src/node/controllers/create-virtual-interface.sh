#!/bin/bash

default_interface=$(ip route | grep '^default' | awk '/default/ {print $5}')
default_gateway=$(ip route | grep '^default' | awk '/default/ {print $3}')

ipaddr=$1
mask=$2
cidr=$3
gateway=$4


#sudo /usr/sbin/ifconfig $default_interface $ipaddr $mask

sudo ip addr add $ipaddr/$cidr dev $default_interface



echo "done"