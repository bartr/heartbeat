#!/bin/bash

echo "post-start start" >> ~/status

# this runs in background each time the container starts

# pull docker base image
docker pull golang:latest

echo "post-start complete" >> ~/status
