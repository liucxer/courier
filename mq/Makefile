test: download
	go test -v -race ./...

cover: download
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

gen:
	protoc --go_out=paths=source_relative:. *.proto

download:
	go mod download -x
