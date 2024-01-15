#!/bin/bash

echo "on-create start" >> ~/status

echo "export GOPATH='$HOME/go'" >> $HOME/.zshrc

# pull docker base images
docker pull golang:latest

echo "on-create complete" >> ~/status
