export GO111MODULE=on
APP_NAME := pir
CMD_DIR := ./cmd
OUT_DIR := out
WASM_OUTPUT := $(OUT_DIR)/main.wasm
NATIVE_OUTPUT := $(OUT_DIR)/$(APP_NAME)
SHELL := /bin/bash

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

check-quality:
	make fmt
	make vet

vet:
	go vet ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

test:
	make tidy
	go test -v -timeout 10m ./... -coverprofile=coverage.out -json > report.json

coverage:
	make test
	go tool cover -html=coverage.out

build: $(OUT_DIR)
	go build -o $(NATIVE_OUTPUT) $(CMD_DIR)


wasm: $(OUT_DIR)
	GOOS=js GOARCH=wasm go build -o $(WASM_OUTPUT) $(CMD_DIR)

run: build
	chmod +x ./$(NATIVE_OUTPUT)
	./$(NATIVE_OUTPUT) -r

clean:
	go clean
	rm -rf out/
	rm -f coverage*.out

.PHONY: all test build wasm
all: 
	make check-quality
	make test
	make build

.PHONY: help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

