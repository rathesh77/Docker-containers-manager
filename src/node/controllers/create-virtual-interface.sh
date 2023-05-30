#!/bin/bash

default_interface=$(ip route | grep '^default' | awk '/default/ {print $5}')
default_gateway=$(ip route | grep '^default' | awk '/default/ {print $3}')

ipaddr=$1
mask=$2
cidr=$3
gateway=$4
service_port=$5
service=$6

#sudo /usr/sbin/ifconfig $default_interface:1 down

#sudo /usr/sbin/ifconfig $default_interface:1 $ipaddr

mkdir -p /etc/nginx/servers/

touch /etc/nginx/servers/$service.conf

server_ips=""

for pod in $7; do
    ip=$(docker network inspect -f '{{range .IPAM.Config}}{{.IPAddress}}{{end}}' "$pod")
    server_ips+="server $ip:$service_port;"
done

echo "
#include /etc/nginx/upstreams/node/$service/default.conf;
server {
    #listen 80;
    #listen $ipaddr:$service_port;
    listen $service_port;
    listen [::]:$service_port;

    #server_name $service-service.org
    include /etc/nginx/locations/node/$service/default.conf;

}" > /etc/nginx/servers/$service.conf

mkdir -p /etc/nginx/locations/node/$service

touch /etc/nginx/locations/node/$service/default.conf

echo "location /$service/ {
    #include proxy_params;

    proxy_pass http://$service/;

    #proxy_set_header Host \$http_host;
    proxy_set_header X-Real-IP \$remote_addr;
    proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto \$scheme;

}" > /etc/nginx/locations/node/$service/default.conf

mkdir -p /etc/nginx/upstreams/node/$service
touch /etc/nginx/upstreams/node/$service/default.conf

echo "upstream $service {
   $server_ips
}" > /etc/nginx/upstreams/node/$service/default.conf

sudo systemctl restart nginx

echo "done"

