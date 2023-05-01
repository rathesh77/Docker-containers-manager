#!/bin/bash


sudo sh ./watcher/start-watcher.sh &

cd kubelet

go run .


