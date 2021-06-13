.PHONY: help all build run load-test

help :
	@echo "Usage:"
	@echo "   make all         - build, run and test TinyBench"
	@echo "   make build       - build a TinyBench docker image"
	@echo "   make run         - run TinyBench from docker image"
	@echo "   make test        - run a test against TinyBench"

all : build run test

build :
	docker pull golang:alpine
	docker pull busybox
	docker build . -t tinybench

run :
	docker run -d --rm --name tinybench -p 8080:8080 tinybench
	docker logs tinybench

test :
	webv -s http://localhost:8080 -f tinybench.json
