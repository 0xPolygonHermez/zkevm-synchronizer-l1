ARCH := $(shell arch)

ifeq ($(ARCH),x86_64)
	ARCH = amd64
else
	ifeq ($(ARCH),aarch64)
		ARCH = arm64
	endif
endif
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/dist
GOENVVARS := GOBIN=$(GOBIN) CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH)
GOBINARY := zkevm-synchronizer-l1
GOCMD := $(GOBASE)/cmd


# Variables
VENV           = .venv
VENV_PYTHON    = $(VENV)/bin/python
SYSTEM_PYTHON  = $(or $(shell which python3), $(shell which python))
PYTHON         = $(or $(wildcard $(VENV_PYTHON)), "install_first_venv")
GENERATE_SCHEMA_DOC = $(VENV)/bin/generate-schema-doc
GENERATE_DOC_PATH   = "docs/config-file/"
GENERATE_DOC_TEMPLATES_PATH = "docs/config-file/templates/"

# Check dependencies
# Check for Go
.PHONY: check-go
check-go:
	@which go > /dev/null || (echo "Error: Go is not installed" && exit 1)

# Check for Docker
.PHONY: check-docker
check-docker:
	@which docker > /dev/null || (echo "Error: docker is not installed" && exit 1)

# Check for Docker-compose
.PHONY: check-docker-compose
check-docker-compose:
	@which docker-compose > /dev/null || (echo "Error: docker-compose is not installed" && exit 1)

# Check for Protoc
.PHONY: check-protoc
check-protoc:
	@which protoc > /dev/null || (echo "Error: Protoc is not installed" && exit 1)

# Check for Python
.PHONY: check-python
check-python:
	@which python3 > /dev/null || which python > /dev/null || (echo "Error: Python is not installed" && exit 1)

# Check for Curl
.PHONY: check-curl
check-curl:
	@which curl > /dev/null || (echo "Error: curl is not installed" && exit 1)

# Targets that require the checks
build: check-go
lint: check-go
build-docker: check-docker
build-docker-nc: check-docker
run-rpc: check-docker check-docker-compose
stop: check-docker check-docker-compose
install-linter: check-go check-curl
install-config-doc-gen: check-python
config-doc-node: check-go check-python
config-doc-custom_network: check-go check-python
update-external-dependencies: check-go
generate-code-from-proto: check-protoc

.PHONY: build
build: ## Builds the binary locally into ./dist
	$(GOENVVARS) go build  -o $(GOBIN)/$(GOBINARY) $(GOCMD)

.PHONY: build-docker
build-docker: ## Builds a docker image with the node binary
	docker build -t zkevm-synchronizer-l1 -f ./Dockerfile .

.PHONY: build-docker-nc
build-docker-nc: ## Builds a docker image with the node binary - but without build cache
	docker build --no-cache=true -t zkevm-synchronizer-l1 -f ./Dockerfile .

.PHONY: run-rpc
run-rpc: ## Runs all the services needed to run a local zkEVM RPC node
	docker-compose up -d zkevm-state-db zkevm-pool-db
	sleep 2
	docker-compose up -d zkevm-prover
	sleep 5
	docker-compose up -d zkevm-sync
	sleep 2
	docker-compose up -d zkevm-rpc

.PHONY: stop
stop: ## Stops all services
	docker-compose down

.PHONY: install-linter
install-linter: ## Installs the linter
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2

.PHONY: lint
lint:  ## Runs the linter
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/golangci-lint run --timeout=5m





.PHONY: update-external-dependencies
update-external-dependencies: ## Updates external dependencies like images, test vectors or proto files
	go run ./scripts/cmd/... updatedeps

.PHONY: install-git-hooks
install-git-hooks: ## Moves hook files to the .git/hooks directory
	cp .github/hooks/* .git/hooks

.PHONY: generate-code-from-proto
generate-code-from-proto: ## Generates code from proto files
	cd proto/src/proto/hashdb/v1 && protoc --proto_path=. --proto_path=../../../../include --go_out=../../../../../merkletree/hashdb --go-grpc_out=../../../../../merkletree/hashdb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative hashdb.proto
	cd proto/src/proto/executor/v1 && protoc --proto_path=. --go_out=../../../../../state/runtime/executor --go-grpc_out=../../../../../state/runtime/executor --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative executor.proto
	cd proto/src/proto/aggregator/v1 && protoc --proto_path=. --proto_path=../../../../include --go_out=../../../../../aggregator/prover --go-grpc_out=../../../../../aggregator/prover --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative aggregator.proto

.PHONY: generate-mocks
generate-mocks:
	(cd test; make generate-mocks)

.PHONY: clean-mocks
clean-mocks:
	(cd test; make clean-mocks)

.PHONY: unittest
unittest:  ## Runs the unittest
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test  -short -race -failfast -covermode=atomic -coverprofile=./coverage_unittest.out  -coverpkg ./... -timeout 70s ./... 

.PHONY: unittest-report
unittest-report:  ## Runs the unittest and generate json report
	rm report.json || true
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test  -short -race -failfast -covermode=atomic -coverprofile=./coverage_unittest.out  -coverpkg ./... -timeout 70s ./... -json > report_unittest.json

.PHONY: test-db
test-db:  ## Runs the tests-db
	(cd test; make run-dbs)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test  -race -failfast -covermode=atomic  -coverprofile=./coverage_db.out -timeout 180s ./state/... 
	(cd test; make stop)

.PHONY: test-db-report
test-db-report: ## Runs the tests-db generate json report
	(cd test; make run-dbs)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test  -race -failfast -covermode=atomic  -coverprofile=./coverage_db.out -timeout 180s ./state/... -json | tee report_db.json
	(cd test; make stop)


## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.
.DEFAULT_GOAL := help

.PHONY: help
help: ## Prints this help
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sort \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
