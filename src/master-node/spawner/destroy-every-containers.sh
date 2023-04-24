#!/bin/bash

docker stop $(docker ps -a --format "{{.ID}}")