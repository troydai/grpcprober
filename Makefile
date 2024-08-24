.PHONY: bin tools gen image push integration

# Override with setting these two and run make with option -e
ARCH=$(shell uname -m | tr '[:upper:]' '[:lower:]')
OS=$(shell uname -s | tr '[:upper:]' '[:lower:]')

OUTPUT_DIR=bin
OUTPUT_NAME=server
MAIN_FILE=cmd/main.go
GO_FILES=$(shell find . -name '*.go' -type f -not -path "./vendor/*")
PROTO_FILES=$(shell find . -name '*.proto' -type f -not -path "./vendor/*")

bin: gen $(GO_FILES)
	GOOS=$(OS) GOARCH=$(ARCH) go build -v -o $(OUTPUT_DIR)/$(OUTPUT_NAME) $(MAIN_FILE)

run: bin
	$(OUTPUT_DIR)/$(OUTPUT_NAME)

tools:
	@ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	@ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

gen: $(PROTO_FILES)
	@ mkdir -p gen
	@ protoc --go_out=gen --go_opt=paths=source_relative \
    --go-grpc_out=gen --go-grpc_opt=paths=source_relative \
	$(PROTO_FILES)
