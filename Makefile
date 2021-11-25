up:
	docker-compose up -d --build

down:
	docker-compose down

test:
	go test -cover -v -parallel 8 ./...

integration-test: up
	go test -v -tags integration -count 1 -p 1 ./tests/...
	make down

lint:
	golangci-lint run --verbose

dep:
	go mod tidy
