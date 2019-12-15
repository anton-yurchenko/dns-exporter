# global
BINARY := $(notdir $(CURDIR))
GO_BIN_DIR := $(GOPATH)/bin
PKGS := $(go list ./... | grep -v /vendor)

# unit tests
test: lint
	@echo "unit testing..."
	@go test ./... -race -coverprofile codecoverage.out

# lint
.PHONY: lint
lint: $(GO_LINTER)
	@echo "vendoring..."
	@go mod vendor
	@go mod tidy
	@echo "linting..."
	@golangci-lint run ./...

# initialize
.PHONY: init
init:
	@rm -f go.mod
	@rm -f go.sum
	@rm -rf ./vendor
	@go mod init

# linter
GO_LINTER := $(GO_BIN_DIR)/golangci-lint
$(GO_LINTER):
	@echo "installing linter..."
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

# release
OSES := linux
ARCHS := amd64

.PHONY: release
release: test
	@rm -rf ./release
	@mkdir -p release
	@for ARCH in $(ARCHS); do \
		for OS in $(OSES); do \
			if test "$$OS" = "windows"; then \
				GOOS=$$OS GOARCH=$$ARCH go build -o release/$(BINARY)-$$OS-$$ARCH.exe; \
			else \
				GOOS=$$OS GOARCH=$$ARCH go build -o release/$(BINARY)-$$OS-$$ARCH; \
			fi; \
		done; \
	done

.PHONY: codecov
codecov: test
	@go tool cover -html=codecoverage.out