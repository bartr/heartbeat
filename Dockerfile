FROM golang:alpine as build

COPY src/main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.Version=0.0.1-$(date +%m%d-%H%M)" -a -installsuffix cgo -o ./heartbeat main.go

FROM busybox

EXPOSE 8080

COPY --from=build /go/heartbeat .

ENTRYPOINT [ "./heartbeat" ]
