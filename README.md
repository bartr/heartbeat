# Heartbeat

> Tiny Kubernetes web server for network testing and monitoring

![License](https://img.shields.io/badge/license-MIT-green.svg)

> New: Linux and Windows binaries are in the ./bin directory and the release
>
> `make local` will build the Linux and Windows binaries locally

## Project Description

Heartbeat is a simple http server written in `go` designed to allow you to measure and monitor network performance `into your k8s clusters`. Unlike most network tools, Heartbeat runs as a k8s pod which monitors and measures all the way through ingress and service mesh ensuring an end-to-end measurement.

Heartbeat will run on any Kubernetes cluster and is targeted at `edge` scenarios as well as cloud and on-prem. Because we targeted the edge by design, we engineered Heartbeat to be, well, `tiny`

We designed Heartbeat to run on thousands of edge nodes. During design, here is what `tiny` meant to us

- The code has to run 7x24x365
- The code has to be simple to deploy
- The container has to be tiny
- The runtime memory footprint has to be tiny
- The network overhead has to be tiny
- The code has to be tiny
- The response has to be configurable to the byte
- The results need to be on a single pane of glass

Here are our current results

- Docker image size - 7.46 MB
- Memory usage - 10 MB
- Less than 200 lines of code
- Deploy as a pod, via Helm or via GitOps
- Deploy it and forget it - your `single pane of glass` will alert you to any problems

Heartbeat is also fast. Really fast!

- 64K and less results are measured in micro-seconds
- A 1MB result takes < 5ms
- A 1GB result takes < 4s
  - Heartbeat is designed for `small payloads` but, we figured, why not test to 1GB?
- Heartbeat is so fast that we turned logging off by default
  - Your `single pane of glass` will tell you if Heartbeat goes down
    - Centralized monitoring of the edge!
  - This also lightens the load on your k8s cluster by not generating extra logs to parse and forward
  - You can turn it on with `--log` and Heartbeat will write to stdout (tab delimited)

Heartbeat is reliable

- We have already run over ~~100M~~ ~~5000M~~ 1B requests (mixture of 1K, 64K and 1MB) on a mixture of k8s clusters with zero errors!
- While not intended for this load, we are able to provide 10K RPS reliably with < 2ms degredation at 1MB
  - In 10MB of ram!

Heartbeat is configurable

- One of the problems, especially on the edge, is network packet fragmentation
- Heartbeat allows you to configure the result to the byte, so you can easily test edge cases
  - You have to do the http overhead math - at least for now

Heartbeat supports a `single pane of glass`

- Easily integrate Heartbeat with [WebValidate](https://github.com/microsoft/webvalidate) for single pane of glass via logs and metrics
  - This is how we run Heartbeat every day
- We use `Fluent Bit` and `Prometheus`
  - Any tool that can make http requests will work, so Heartbeat likely `drops right in`

Probes

- For simplicity, we don't provide a dedicated ready or health probe
- We use `/heartbeat/1`

Testing Uploads

- You can test upload speed by POSTing to /heartbeat/1
  - Payload doesn't matter as we read it and throw it away

### Try Heartbeat

From a bash shell

```bash

# build the docker image
make build

# run the local docker image
make run

# check the semver
curl localhost:8080/version

# 0123456789ABCDEF0
curl localhost:8080/heartbeat/17

# 1K + 1 (should end with 0)
curl localhost:8080/heartbeat/1025

# you probably figured out the last URI segment is the number of bytes

# 1MB
curl localhost:8080/heartbeat/1048576

# this will show you how slow stdout is
curl localhost:8080/heartbeat/1048576 > /dev/null

# these will return 400
# default config is 1 <= size <= 1MB
curl -i localhost:8080/heartbeat/0

curl -i localhost:8080/heartbeat/1048577

# check the logs
docker logs heartbeat

# other options
docker run -it --rm heartbeat -h

# test heartbeat with WebV
make test

# generate 100 req/sec for 30 seconds with WebV
make load-test

```

### CI-CD via GitHub Action

> You will need to edit the repo name and make sure you have the secrets setup

- GitHub Actions (ci-cd pipelines)
  - [Build container image](./.github/workflows/build.yaml)

## Contributing

This project welcomes contributions and suggestions and has adopted the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/1/code_of_conduct.html).

For more information see the [Code of Conduct FAQ](https://www.contributor-covenant.org/faq).

## Trademarks

This project may contain trademarks or logos for projects, products, or services. Any use of third-party trademarks or logos are subject to those third-party's policies.
