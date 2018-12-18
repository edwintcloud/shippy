build:
	docker build -t consignment-service .
	rm -f proto/consignment/*.go
	docker run --entrypoint cat consignment-service consignment.micro.go > proto/consignment/consignment.micro.go
	docker run --entrypoint cat consignment-service consignment.pb.go > proto/consignment/consignment.pb.go
run:
	docker run -p 50051:50051 \
		-e MICRO_SERVER_ADDRESS=:50051 \
		-e MICRO_REGISTRY=mdns consignment-service