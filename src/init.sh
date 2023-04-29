#!/bin/bash

sudo sh ./master-node/spawner/spawn-cluster.sh
sudo sh ./master-node/spawner/spawn-node.sh
#sudo sh ./master-node/watcher/start-watcher.sh &


while true; do nc -l -p 2000 ; done


