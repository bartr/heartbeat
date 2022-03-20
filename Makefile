.PHONY: help all build run test load-test local

help :
	@echo "Usage:"
	@echo "   make all         - build, run and test Heartbeat"
	@echo "   make build       - build a Heartbeat docker image"
	@echo "   make run         - run Heartbeat from docker image"
	@echo "   make test        - run a test against Heartbeat"
	@echo "   make load-test   - run a 60 second load test against Heartbeat (100 req/sec)"
	@echo "   make local       - build local executables"

all : build run test

local :
	# Building Linux binary
	@CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.Version=0.0.1" -o bin/heartbeat -a src/main.go

	# Building Windows binary
	@CGO_ENABLED=0 GOOS=windows go build -ldflags="-X main.Version=0.0.1" -o bin/heartbeat.exe -a src/main.go

build :
	docker build . -t heartbeat

run :
	docker run -d --rm --name heartbeat -p 8080:8080 tinybheartbeatench
	docker logs heartbeat

test :
	@cd webv && webv --server http://localhost:8080 --files heartbeat.json --verbose

load-test :
	@cd webv && webv --server http://localhost:8080 --files load.json -r --sleep 100 --duration 30 --verbose
