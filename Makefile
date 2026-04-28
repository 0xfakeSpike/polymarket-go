GO ?= go

.PHONY: fmt vet test build examples

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...

build:
	$(GO) build ./...

examples:
	$(GO) run ./examples/public-search
