FROM golang:alpine as build

ARG Version=0.1.0

COPY src/main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.Version=$Version-$(date +%m%d-%H%M)" -a -installsuffix cgo -o ./heartbeat main.go

FROM busybox

EXPOSE 8080

COPY --from=build /go/heartbeat .

ENTRYPOINT [ "./heartbeat" ]
