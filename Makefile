deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.6

lint:
	golangci-lint run --timeout 5m --config .golangci.yml

test:
	go test ./...