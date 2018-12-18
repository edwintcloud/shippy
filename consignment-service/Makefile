build:
	protoc -I. --go_out=plugins=grpc:$(GOPATH)/src/shippy/consignment-service \
	  proto/consignment/consignment.proto
	GOOS=linux GOARCH=amd64 go build
	docker build -t consignment-service .
run:
	docker run -p 50051:50051 consignment-service