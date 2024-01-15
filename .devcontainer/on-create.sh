#!/bin/bash

echo "on-create start" >> ~/status

# Change shell to zsh for vscode
sudo chsh --shell /bin/zsh vscode

#echo "export GOPATH='$HOME/go'" >> $HOME/.zshrc

# pull docker base images
docker pull golang:latest

# configure git
git config --global core.whitespace blank-at-eol,blank-at-eof,space-before-tab
git config --global pull.rebase false
git config --global init.defaultbranch main
git config --global fetch.prune true
git config --global core.pager more
git config --global diff.colorMoved zebra
git config --global devcontainers-theme.show-dirty 1
git config --global core.editor "nano -w"

echo "generating completions"
gh completion -s zsh > ~/.oh-my-zsh/completions/_gh
kubectl completion zsh > "$HOME/.oh-my-zsh/completions/_kubectl"
k3d completion zsh > "$HOME/.oh-my-zsh/completions/_k3d"

sudo apt-get update

# only run apt upgrade on pre-build
if [ "$CODESPACE_NAME" = "null" ]
then
    echo "$(date +'%Y-%m-%d %H:%M:%S')    upgrading" >> "$HOME/status"
    sudo apt-get upgrade -y
fi

echo "on-create complete" >> ~/status
