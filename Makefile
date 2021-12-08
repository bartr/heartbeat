.PHONY: help all build run test load-test local

help :
	@echo "Usage:"
	@echo "   make all         - build, run and test TinyBench"
	@echo "   make build       - build a TinyBench docker image"
	@echo "   make run         - run TinyBench from docker image"
	@echo "   make test        - run a test against TinyBench"
	@echo "   make load-test   - run a 60 second load test against TinyBench (100 req/sec)"
	@echo "   make local       - build local executables"

all : build run test

local :
	# Building Linux binary
	@CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.Version=0.0.1" -o bin/tinybench -a src/main.go

	# Building Windows binary
	@CGO_ENABLED=0 GOOS=windows go build -ldflags="-X main.Version=0.0.1" -o bin/tinybench.exe -a src/main.go

build :
	docker pull golang:alpine
	docker pull busybox
	docker build . -t tinybench

run :
	docker run -d --rm --name tinybench -p 8080:8080 tinybench
	docker logs tinybench

test :
	@cd webv && webv --server http://localhost:8080 --files tinybench.json --verbose

load-test :
	@cd webv && webv --server http://localhost:8080 --files load.json -r --sleep 100 --duration 30 --verbose
