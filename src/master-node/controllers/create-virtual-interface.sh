#!/bin/bash

default_interface=$(ip route | grep '^default' | awk '/default/ {print $5}')
default_gateway=$(ip route | grep '^default' | awk '/default/ {print $3}')

ipaddr=$(ip route get 8.8.8.8 | head -1 | cut -d' ' -f7)
mask=$2
cidr=$3
gateway=$4
service_port=$5
service=$6


server_ips=""

for pod in $7; do
    server_ips+="server $pod;"
done


mkdir -p /etc/nginx/locations/$service

touch /etc/nginx/locations/$service/default.conf

echo "
    include /etc/nginx/upstreams/$service/*.conf;

server {
    listen 80;

    server_name www.$service.fr $service.fr;


location / {
    #include proxy_params;

    #rewrite ^/$service/$ / break;
    proxy_pass http://$service/;

    proxy_set_header Host \$server_name;
    proxy_set_header X-Real-IP \$remote_addr;
    proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto \$scheme;

}
}" > /etc/nginx/locations/$service/default.conf

mkdir -p /etc/nginx/upstreams/$service
touch /etc/nginx/upstreams/$service/default.conf

echo "upstream $service {
   $server_ips
}" > /etc/nginx/upstreams/$service/default.conf

sudo systemctl restart nginx

echo "done"

