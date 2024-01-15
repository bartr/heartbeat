# build the app
FROM golang:latest as build

ARG Version=0.3.0

COPY src/main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.Version=$Version-$(date +%m%d-%H%M)" -a -installsuffix cgo -o ./heartbeat main.go

## create the image
FROM scratch

EXPOSE 8080
ENTRYPOINT [ "./heartbeat" ]

COPY --from=build /go/heartbeat .
