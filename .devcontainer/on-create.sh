#!/bin/bash

echo "on-create start" >> ~/status

echo "export GOPATH='$HOME/go'" >> $HOME/.zshrc

# go install golang.org/x/tools/gopls@latest

# pull docker base images
docker pull golang:alpine
docker pull busybox:latest

echo "on-create complete" >> ~/status
