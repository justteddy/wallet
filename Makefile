GO_CI_LINT=golangci-lint

test:
	LOG_LEVEL=panic go test -cover -v -parallel 8 ./...

lint:
	$(GO_CI_LINT) run --verbose

dep:
	go mod tidy
