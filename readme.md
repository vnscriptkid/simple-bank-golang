## module
- go mod init github.com/vnscriptkid/simple-bank-golang
- go mod tidy

## cli
- migrate cli tool: https://github.com/golang-migrate/migrate
- gen code from sql: https://github.com/kyleconroy/sqlc

## libs
- driver to connect to pg: github.com/lib/pq
- assertion lib: https://github.com/stretchr/testify

## proto
- brew install protobuf
- protoc --version
- https://grpc.io/docs/languages/go/quickstart/
    - go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
    - go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
    - export PATH="$PATH:$(go env GOPATH)/bin"

## grpc client
- evans --host localhost --port 8002 -r repl
- show service
- call CreateUser