.PHONY: FORCE

######################
##### VARIABLES ######
######################

# Name of current program
PROGRAM_NAME=easyserver

# Binary name
BINARY_CLI = ./build/${PROGRAM_NAME}
ifeq ($(OS),Windows_NT)
	BINARY_CLI := $(BINARY).exe
endif

#  Module name
PKG_PATH=$(shell head -n1 go.mod | cut -d ' ' -f2)

# Git information
COMMIT=$(shell git rev-parse --short HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
TAG=$(shell git describe --tags |cut -d- -f1)

# Go variables
GOPATH=$(shell go env GOPATH)
GOBIN="${GOPATH}/bin"

# Optional colors to beautify output
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

# Make git information availible inside project
LDFLAGS = -ldflags "-X ${PKG_PATH}/internal/version.gitTag=${TAG} -X ${PKG_PATH}/internal/version.gitCommit=${COMMIT} -X ${PKG_PATH}/internal/version.gitBranch=${BRANCH}"

######################
######## HELP ########
######################

.DEFAULT_GOAL := help

help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
.PHONY: help

######################
#### DEPENDENCIES ####
######################

go-deps: ## Download the dependencies for Go.
	go mod download
.PHONY: go-deps

py-deps: ## Download the dependencies for Python.
	pip install -r requirements.txt
.PHONY: py-deps

externals: ## Install external dependencies
	go install golang.org/x/exp/cmd/modgraphviz@latest
	go install github.com/air-verse/air@latest

graph: deps ## Makes dependency graph
	-@mkdir build
	go mod graph | ${GOPATH}/bin/modgraphviz | dot -Tpng -o build/graph.png
	xdg-open build/graph.png
.PHONY: graph

######################
####### BUILD ########
######################

build: $(BINARY_CLI) ## Build program
.PHONY: build

ffi_build:
	go build -o ./build/easyserver.so -buildmode=c-shared ./ffi

run: ## Run program without building
	go run ${LDFLAGS} ./cmd/cli
.PHONY: run

python_run: ffi_build ## Run Python code
	pyuic6 -o ./cmd/python/ui/ui.py ./cmd/python/ui/mainwindow.ui
	pyuic6 -o ./cmd/python/ui/frame.py ./cmd/python/ui/project.ui
	pyuic6 -o ./cmd/python/ui/input_date.py ./cmd/python/ui/date.ui
	pyuic6 -o ./cmd/python/ui/input_time.py ./cmd/python/ui/time.ui
	pyuic6 -o ./cmd/python/ui/input_datetime.py ./cmd/python/ui/datetime.ui
	cd build; python3 ../cmd/python/gui.py

build_race: ## Build program with race detector
	go build -race ${LDFLAGS} -o $(BINARY_CLI) ./cmd/cli
.PHONY: build_race

clean: ## Clean build output
	rm -rf ./build # remove Makefile build 
	rm -rf ./tmp # remove Air build
	go clean
.PHONY: clean

dev: deps ## Start program in dev mode
	$(GOPATH)/bin/air
.PHONY: air

######################
#### CODE QUALITY ####
######################

lint: go-deps ## Lint the source files
	golangci-lint run --timeout 5m -E revive,gosec
.PHONY: lint

test: go-deps ## Run tests
	go test -race -p 1 -timeout 300s -coverprofile=.test_coverage.txt ./... && \
    	go tool cover -func=.test_coverage.txt | tail -n1 | awk '{print "Total test coverage: " $$3}'
	@rm .test_coverage.txt

######################
####### FILES ########
######################

$(BINARY_CLI): FORCE
	go build ${LDFLAGS} -o $0 ./cmd/cli
	$(GOPATH)/bin/gosec

go.mod: FORCE
	go mod tidy
	go mod verify
go.sum: go.mod
