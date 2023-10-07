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

## rabbitmq (if needed)
docker compose up -d

## Develop
go mod tidy
./codegen.sh
go run server.go &
go run client.go

# Once if *.proto changed
./codegen.sh

## Test
go test ./...

## Deploy
go build
./bec-grpc

# Deploy via PM2
./pm2.sh
```

### Troubleshootings

- `go test` caching unexpectedly

  ```sh
  go clean -testcache
  go test ./...
  ```

- `go run/build/test` throws `ld: library not found for -lleptonica`

  ```plain
  ζ go test ./...
  ?       bec-grpc        [no test files]
  ?       bec-grpc/mq     [no test files]
  ok      bec-grpc/helloworld     (cached)
  # bec-grpc/ocr.test
  /usr/local/Cellar/go/1.20.2/libexec/pkg/tool/darwin_amd64/link: running clang++ failed: exit status 1
  ld: library not found for -lleptonica
  clang: error: linker command failed with exit code 1 (use -v to see invocation)

  FAIL    bec-grpc/ocr [build failed]
  FAIL

  # if MacOS
  # https://github.com/otiai10/gosseract#installation
  # https://github.com/tesseract-ocr/tessdoc/blob/main/Installation.md
  brew install tesseract
  ```

- `go install xxx` throws: `can only use path@version syntax with 'go get'`

  Should upgrade GO from `<=1.13.x`<br>
  How to Install latest version of GO on Ubuntu 20.04 LTS (Focal Fossa)<br>
  <https://www.cyberithub.com/how-to-install-latest-version-of-go-on-ubuntu-20-04/>

- `go build` throws: `leptonica/allheaders.h: No such file or directory`

  ```plain
  go build
  # github.com/otiai10/gosseract/v2
  tessbridge.cpp:5:10: fatal error: leptonica/allheaders.h: No such file or directory
      5 | #include <leptonica/allheaders.h>
        |          ^~~~~~~~~~~~~~~~~~~~~~~~`

  # if Ubuntu
  # https://github.com/otiai10/gosseract#installation
  # https://github.com/tesseract-ocr/tessdoc/blob/main/Installation.md
  sudo apt install tesseract-ocr
  sudo apt install libtesseract-dev
  ```
