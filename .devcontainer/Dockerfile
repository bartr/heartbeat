# this image is built and updated weekly
# https://github.com/cse-labs/codespaces-images

FROM ghcr.io/cse-labs/dind:latest

# install go
RUN wget -c https://dl.google.com/go/go1.16.linux-amd64.tar.gz -O - | tar -xz -C /usr/local

# add to path
ENV PATH=$PATH:/usr/local/go/bin
