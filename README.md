# TinyBench

> Tiny Kubernetes web server for network testing and monitoring

![License](https://img.shields.io/badge/license-MIT-green.svg)

## Project Description

TinyBench is a very simple http server written in `go` designed to allow you to measure and monitor network performance `into your k8s clusters`. Unlike most network tools, TinyBench runs as a k8s pod which monitors and measures through ingress and service mesh ensuring an end-to-end measurement.

TinyBench will run on any Kubernetes server and is targeted at `edge` scenarios as well as cloud and on-prem. Because we targeted the edge by design, we engineered TinyBench to be, well, `tiny`

We designed TinyBench to run on thousands of edge nodes. During design, here is what `tiny` meant to us

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

TinyBench is also fast. Really fast!

- 64K and less results are measured in micro-seconds
- A 1MB result takes < 5ms
- A 1GB result takes < 4s
  - TinyBench is designed for `small payloads` but, we figured, why not test to 1GB?
- TinyBench is so fast that we turned logging off by default
  - Your `single pane of glass` will tell you if TinyWeb goes down
    - Centralized monitoring of the edge!
  - This also lightens the load on your k8s cluster by not generating extra logs to parse and forward
  - You can turn it on with `--log` and TinyBench will write to stdout (tab delimited)

TinyBench is reliable

- We have already run over ~~100M~~ ~~5000M~~ 1B requests (mixture of 1K, 64K and 1MB) on a mixture of k8s clusters with zero errors!
- While not intended for this load, we are able to provide 10K RPS reliably with < 2ms degredation at 1MB
  - In 10MB of ram!

TinyBench is configurable

- One of the problems, especially on the edge, is network packet fragmentation
- Tiny bench allows you to configure the result to the byte, so you can easily test edge cases
  - You have to do the http overhead math - at least for now

TinyBench supports a `single pane of glass`

- Easily integrate TinyBench with [LodeRunner](https://github.com/retaildevcrews/loderunner) for single pane of glass via logs and metrics
- Any tool that can make http requests will work, so TinyBench likely `drops right in`
  - We use `Fluent Bit` and `Prometheus`

Probes

- For simplicity, we don't provide a dedicated startup or health probe
- We use `/tinybench/1`

Testing Uploads

- You can test upload speed by POSTing to /tinyweb/1
- Payload doesn't matter as we read it and throw it away

### Try TinyBench

From a bash shell

```bash

docker run -d --name tinybench -p 8080:8080 ghcr.io/retaildevcrews/tinybench:beta --log

# check the semver
curl localhost:8080/version

# 0123456789ABCDEF0
curl localhost:8080/tinybench/17

# 1K + 1 (should end with 0)
curl localhost:8080/tinybench/1025

# you probably figured out the last URI segment is the number of bytes

# 1MB
curl localhost:8080/tinybench/1048576

# this will show you how slow stdout is
curl localhost:8080/tinybench/1048576 > /dev/null

# these will return 400
# default config is 1 <= size <= 1MB
curl -i localhost:8080/tinybench/0

curl -i localhost:8080/tinybench/1048577

# check the logs
docker logs tinybench

# other options
docker run -it --rm ghcr.io/retaildevcrews/tinybench:beta -h

```

### Build TinyBench

From a bash shell

```bash

# build with docker
make build

# build and run
make all

```

### CI-CD via GitHub Action

> You will need to edit the repo name and make sure you have the secrets setup

- GitHub Actions (ci-cd pipelines)
  - [Build container image](./.github/workflows/build.yaml)

### Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution. For details, visit <https://cla.opensource.microsoft.com>

When you submit a pull request, a CLA bot will automatically determine whether you need to provide a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/). For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

### Trademarks

This project may contain trademarks or logos for projects, products, or services.

Authorized use of Microsoft trademarks or logos is subject to and must follow [Microsoft's Trademark & Brand Guidelines](https://www.microsoft.com/en-us/legal/intellectualproperty/trademarks/usage/general).

Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship.

Any use of third-party trademarks or logos are subject to those third-party's policies.
