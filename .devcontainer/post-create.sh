#!/bin/bash

echo "post-create start" >> ~/status

# this runs in background after UI is available

# (optional) upgrade packages
sudo apt update
#sudo apt upgrade -y
#sudo apt autoremove -y
#sudo apt clean -y

echo "post-create complete" >> ~/status
