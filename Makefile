deps:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.4

lint:
	golangci-lint run --timeout 5m --config .golangci.yml

test:
	go test ./...