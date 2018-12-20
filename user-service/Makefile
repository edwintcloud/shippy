build:
	docker build -t user-service .
	rm -f proto/consignment/*.go
	docker run --entrypoint cat user-service user.micro.go > proto/user/user.micro.go
	docker run --entrypoint cat user-service user.pb.go > proto/user/user.pb.go
run:
	docker run -p 50053:50051 \
		-e MICRO_SERVER_ADDRESS=:50051 \
		-e MICRO_REGISTRY=mdns user-service