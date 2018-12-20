# shippy
An exercise building a shipping platform using micro-services in go. Based on [this](https://ewanvalentine.io/microservices-in-golang-part-1/) blog series.

## Run Stack With Docker Compose
1. Build stack: `docker-compose build`
2. Run stack: `docker-compose up`

## user-service
Microservice for managing users of the platform.
1. Have docker installed and configured for your machine
2. Build and run: `make build && make run`

## vessel-service
Microservice for for generating a vessel id based on specifications.
1. Have docker installed and configured for your machine
2. Build and run: `make build && make run`

## consignment-service
Microservice for matching a consignment of containers to a vessel.
1. Have docker installed and configured for your machine
2. Build and run: `make build && make run`

## consignment-cli
Client for testing vessel and consignment microservices via command-line interface.
1. Have docker installed and configured for your machine
2. Build and run: `make build && make run`

## user-cli
Client for testing user microservice via command-line interface.
1. Have docker installed and configured for your machine
2. Build and run: `make build && make run`
