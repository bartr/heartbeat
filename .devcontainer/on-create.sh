#!/bin/bash

echo "on-create start" >> ~/status

# this runs when container is initially created

go get -v golang.org/x/tools/gopls

# pull docker base images
docker pull golang:alpine
docker pull busybox:latest

echo "on-create complete" >> ~/status
