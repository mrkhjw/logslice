.PHONY: build test lint clean run

BINARY   := logslice
CMD_PATH := ./cmd/logslice
GO       := go

build:
	$(GO) build -o bin/$(BINARY) $(CMD_PATH)

test:
	$(GO) test ./... -v -count=1

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint not found; install from https://golangci-lint.run"; exit 1; }
	golangci-lint run ./...

clean:
	rm -rf bin/

run: build
	./bin/$(BINARY) $(ARGS)

# Usage examples:
#   make run ARGS="-level ERROR myfile.log"
#   make run ARGS="-format json -since 2024-01-01T00:00:00Z myfile.log"
#   cat myfile.log | make run ARGS="-format table"
