up:
	docker-compose up -d --build

down:
	docker-compose down

test:
	LOG_LEVEL=panic go test -cover -v -parallel 8 ./...

lint:
	golangci-lint run --verbose

dep:
	go mod tidy
