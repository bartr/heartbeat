.PHONY: help all build run test load-test

help :
	@echo "Usage:"
	@echo "   make all         - build, run and test TinyBench"
	@echo "   make build       - build a TinyBench docker image"
	@echo "   make run         - run TinyBench from docker image"
	@echo "   make test        - run a test against TinyBench"
	@echo "   make load-test   - run a 60 second load test against TinyBench (100 req/sec)"

all : build run test

build :
	docker pull golang:alpine
	docker pull busybox
	docker build . -t tinybench

run :
	docker run -d --rm --name tinybench -p 8080:8080 tinybench
	docker logs tinybench

test :
	webv --server http://localhost:8080 --files tinybench.json

load-test :
	webv --server http://localhost:8080 --files tinybench.json -r --sleep 10 --duration 60
