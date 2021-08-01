test:
	go test -race ./...

cover:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

fmt:
	goimports -l -w .
	gofmt -l -w .