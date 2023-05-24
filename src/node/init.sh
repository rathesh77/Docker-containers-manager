#!/bin/bash


#sudo sh ./watcher/start-watcher.sh &

sudo apt-get -y install nginx

rm /etc/nginx/sites-enabled/default

mkdir -p /etc/nginx/locations
mkdir -p /etc/nginx/servers
mkdir -p /etc/nginx/upstreams

touch /etc/nginx/sites-available/custom.config.conf
ln -s /etc/nginx/sites-available/custom.config.conf /etc/nginx/sites-enabled/custom.config.conf

echo "error_log /var/log/nginx/error.log debug;
include /etc/nginx/upstreams/*/*.conf;
include /etc/nginx/servers/*/*.conf;" > /etc/nginx/sites-available/custom.config.conf

systemctl restart nginx

cd kubelet

go run .


