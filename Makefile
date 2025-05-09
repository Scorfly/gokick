deps:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6

lint:
	golangci-lint run --timeout 5m --config .golangci.yml

test:
	go test ./...