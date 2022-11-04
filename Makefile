LOCAL_BIN=$(CURDIR)/bin

# lint
GOLANGCI_BIN=$(LOCAL_BIN)/golangci-lint
$(GOLANGCI_BIN):
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.47.2


.PHONY: lint
lint: $(GOLANGCI_BIN)
	$(GOLANGCI_BIN) run ./...
