# this image is built and updated weekly
# https://github.com/cse-labs/codespaces-images

FROM ghcr.io/cse-labs/k3d:latest

# install go
RUN mkdir -p /home/vscode/go/bin && \
    rm -rf /user/local/go && \
    wget -c https://dl.google.com/go/go1.18.linux-amd64.tar.gz -O - | tar -xz -C /usr/local

# add to path
ENV PATH=$PATH:/usr/local/go/bin:/home/vscode/go/bin \
    GOPATH=/home/vscode/go
