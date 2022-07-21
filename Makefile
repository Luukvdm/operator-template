BIN_NAME ?= operator-template
BIN_DIR ?= bin

SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')
	
build:
	GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(BIN_NAME) src/main.go

manifests:
	controller-gen rbac:roleName=my-role crd paths="./..." output:dir=./charts/operator-template/templates

clean:
	go clean
	rm -rf $(BIN_DIR)

fmt:
	gofmt -s -l -w $(SRCS)

generate:
	controller-gen object paths=./...

