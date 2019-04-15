FROM golang:1.11 as builder
WORKDIR /usr/app
COPY . .
RUN GO111MODULE=on go get -d -v ./...
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -v ./...

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /usr/app .
CMD /root/members-test