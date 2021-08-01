test:
	go test -v ./pkg/...

cover:
	go test -v -coverprofile=coverage.txt -covermode=atomic ./pkg/...

fmt:
	goimports -l -w .
	gofmt -l -w .

dep:
	go get -u ./...