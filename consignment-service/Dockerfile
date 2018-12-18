# consignment-service/Dockerfile

# BUILD BINARY FIRST WITH OFFICIAL GOLANG IMAGE
FROM golang:alpine as builder

# Set working directory and copy code into working directory
WORKDIR /go/src/shippy/consignment-service
COPY . .

# Install govendor and git then pull in dependencies using govendor
RUN apk add git
RUN go get -u github.com/kardianos/govendor
RUN govendor sync

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
COPY --from=builder /go/src/shippy/consignment-service/consignment-service .

# Run the binary
CMD ["./consignment-service"]