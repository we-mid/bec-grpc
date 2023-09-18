```sh
## Prerequisites
# https://grpc.io/docs/languages/go/quickstart/#prerequisites
# 1.
which protoc
# https://command-not-found.com/protoc
brew install protobuf  # if MacOS
apt-get install protobuf-compiler  # if Ubuntu
# ...
# 2.
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
which protoc-gen-go 
# ~/.bashrc
# golang, grpc
# Update your PATH so that the protoc compiler can find the plugins:
# https://grpc.io/docs/languages/go/quickstart/
export PATH="$PATH:$(go env GOPATH)/bin"

## Develop
go mod tidy
go run server.go &
go run client.go

# Once if *.proto changed
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    */*.proto

## Deploy
go build
./bec-ocr
```
