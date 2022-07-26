BIN_NAME ?= otemplate
BIN_DIR ?= bin
# IMG_NAME ?= otemplate
VERSION ?= v0.0.1 # TODO don't hardcode version
IMG_NAME ?= ghcr.io/luukvdm/otemplate:$(VERSION)

SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')
GO_BIN = $(shell go env GOPATH)/bin

install: image
	# kind load docker-image $(IMG_NAME)
	helm install otemplate charts/operator-template
	
build: generate manifests
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(BIN_NAME) src/main.go

image:
	# rm $(BIN_DIR)/$(BIN_NAME)
	docker build -t $(IMG_NAME) .

image.local:
	# This goal builds an image with the binary from the bin folder
	# Compiling the go code outside of the container gives a huge speed boost
	docker build -t $(IMG_NAME) --file=hack/Dockerfile.local .

manifests:
	$(GO_BIN)/controller-gen rbac:roleName=my-role crd paths="./..." output:dir=./charts/operator-template/templates

clean:
	helm delete otemplate
	go clean
	rm -rf $(BIN_DIR)
	docker image rm $(IMG_NAME)

fmt:
	gofmt -s -l -w $(SRCS)

generate:
	$(GO_BIN)/controller-gen object paths=./...

install-tools:
	go install sigs.k8s.io/controller-tools/cmd/controller-gen@latest

