.PHONY: help all docker run test load-test linux windows

help :
	@echo "Usage:"
	@echo "   make all         - build, run and test Heartbeat"
	@echo ""
	@echo "   make docker      - build a Heartbeat docker image"
	@echo "   make linux       - build local Linux executable (/bin/heartbeat)"
	@echo "   make windows     - build local Windows executable (/bin/heartbeat.exe)"
	@echo ""
	@echo "   make run         - run Heartbeat from docker image"
	@echo "   make test        - run a test against Heartbeat docker image"
	@echo "   make load-test   - run a 60 second load test against Heartbeat docker images (100 req/sec)"

all : docker run test

linux :
	# Building Linux binary
	@CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.Version=0.2.0" -o bin/heartbeat -a src/main.go

windows :
	# Building Windows binary
	@CGO_ENABLED=0 GOOS=windows go build -ldflags="-X main.Version=0.2.0" -o bin/heartbeat.exe -a src/main.go

docker :
	docker build . -t heartbeat

run :
	docker run -d --rm --name heartbeat -p 8080:8080 heartbeat
	docker logs heartbeat

test :
	@cd webv && webv --server http://localhost:8080 --files heartbeat.json --verbose --verbose-errors

load-test :
	@cd webv && webv --server http://localhost:8080 --port 8088 --files load.json -r --sleep 10 --duration 30 --verbose
