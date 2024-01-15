#!/bin/bash

echo "post-create start" >> ~/status

# this runs in background after UI is available

echo "update oh-my-zsh"
git -C "$HOME/.oh-my-zsh" pull

echo "post-create complete" >> ~/status
