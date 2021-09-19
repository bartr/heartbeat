#!/bin/bash

echo "post-start start" >> ~/status

# this runs in background each time the container starts

# pull docker base images
docker pull golang:alpine
docker pull busybox:latest

echo "post-start complete" >> ~/status
