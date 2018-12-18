# consignment-service/Dockerfile

# BUILD BINARY FIRST WITH OFFICIAL GOLANG IMAGE
FROM golang:alpine as builder

# Set working directory and copy code into working directory
WORKDIR /go/src/github.com/edwintcloud/shippy/consignment-service
COPY . .

# Install needed system packages
RUN apk update && apk add git protobuf

# Install needed go packages
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u github.com/micro/protoc-gen-micro
RUN go get -u github.com/kardianos/govendor

# Fetch dependencies
RUN govendor sync

# Generate go proto files using protoc
RUN protoc --proto_path=$GOPATH/src/github.com/edwintcloud/shippy/consignment-service/proto/consignment:. \
    --micro_out=proto/consignment/. --go_out=proto/consignment/. consignment.proto

# Build binary
RUN go build

# SETUP OUR ALPINE ENV FOR RUNNING BUILT EXECUTABLE
FROM alpine:latest

# Security related package that's nice to have
# RUN apk --no-cache add ca-certificates

# Create app directory and set it as working directory
RUN mkdir /app
WORKDIR /app

# Pull binary from container we labeled as builder
COPY --from=builder /go/src/github.com/edwintcloud/shippy/consignment-service/consignment-service .

# Pull proto files from container we labeled as builder
COPY --from=builder /go/src/github.com/edwintcloud/shippy/consignment-service/proto/consignment/consignment.pb.go .
COPY --from=builder /go/src/github.com/edwintcloud/shippy/consignment-service/proto/consignment/consignment.micro.go .

# Run the binary
CMD ["./consignment-service"]