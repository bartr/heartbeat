#!/bin/bash

echo "on-create start" >> ~/status

# this runs when container is initially created

# add your commands here
go get -v golang.org/x/tools/gopls

echo "on-create complete" >> ~/status
