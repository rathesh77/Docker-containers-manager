#!/bin/bash


#sudo sh ./watcher/start-watcher.sh &

sudo apt-get -y install nginx

cd kubelet

go run .


