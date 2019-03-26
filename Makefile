
src !=	find . -type f -name '*.go'
version !=	git describe --always || echo "0.0.0"
platforms =	darwin linux windows freebsd
ldflags =	-ldflags "-w -s -X main.version=$(version)"
binary ?=	server

# This is to allow Go to work in Alpine docker image as well as bare metal
GOPATH ?=	$(HOME)/go

# Targets

build: $(binary)

.PHONY: all
all: clean test build-all ## Clean build for current architecture

# Default build for host architecture
$(binary): $(src) ## Build for current architecture
	go build -ldflags "-w -s -X main.version=$(version)" -o $@ ./cmd/api

build-all: $(platforms) ## Build for all architectures

.PHONY: $(platforms)
$(platforms): $(src)
	GOOS=$@ go build $(ldflags) -o $(binary)-$@ ./cmd

lambda: $(src)
	go build -ldflags "-w -s -X main.version=$(version)" -o $@ ./cmd/lambda

handler.zip: lambda
	zip $@ $<

.PHONY: docker
docker:
	docker build .

.PHONY: test
test: lint ## Run tests and create coverage report
	go test -short -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html

.PHONY: lint
lint: $(GOPATH)/bin/golint ## Run the code linter
	@for file in $$(find . -name 'vendor' -prune -o -type f -name '*.go'); do \
		golint $$file; done

$(GOPATH)/bin/golint:
	go get -u golang.org/x/lint/golint

.PHONY: clean
clean: ## Clean up temp files and binaries
	@rm -rf coverage*
	@rm -rf lambda handler.zip
	@rm -f $(binary) $(binary)-*

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) |sort \
		|awk 'BEGIN{FS=":.*?## "};{printf "\033[36m%-30s\033[0m %s\n",$$1,$$2}'
