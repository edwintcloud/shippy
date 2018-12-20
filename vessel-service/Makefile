build:
	docker build -t vessel-service .
	rm -f proto/vessel/*.go
	docker run --entrypoint cat vessel-service vessel.micro.go > proto/vessel/vessel.micro.go
	docker run --entrypoint cat vessel-service vessel.pb.go > proto/vessel/vessel.pb.go
run:
	docker run -p 50052:50051 \
		-e MICRO_SERVER_ADDRESS=:50051 \
		-e MICRO_REGISTRY=mdns \
		-e DB_HOST=datastore:27017 vessel-service