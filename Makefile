BIN_NAME ?= otemplate
BIN_DIR ?= bin
CHART_RELEASE_NAME ?= otemplate
VERSION ?= v0.0.2 # TODO don't hardcode version
IMG_REPO ?= ghcr.io/luukvdm/otemplate
IMG_NAME ?= $(IMG_REPO):$(VERSION)	

SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')
GO_BIN = $(shell go env GOPATH)/bin

kind-install: image kind-load helm-install
kind-reinstall: uninstall kind-install
kind-load:
	kind load docker-image $(IMG_NAME)


install: image helm-install
helm-install:
	helm install $(CHART_RELEASE_NAME) charts/operator-template --set image.repository=$(IMG_REPO),image.tag=$(VERSION)

uninstall:
	helm delete $(CHART_RELEASE_NAME)
	
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

clean: uninstall
	go clean
	rm -rf $(BIN_DIR)
	docker image rm $(IMG_NAME)

fmt:
	gofmt -s -l -w $(SRCS)

generate:
	$(GO_BIN)/controller-gen object paths=./...

install-tools:
	go install sigs.k8s.io/controller-tools/cmd/controller-gen@latest

