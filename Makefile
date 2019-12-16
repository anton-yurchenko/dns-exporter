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

.PHONY: release
release: test
	@rm -rf ./release
	@mkdir -p release
	@GOOS=linux GOARCH=amd64 go build -o ./release/app

.PHONY: codecov
codecov: test
	@go tool cover -html=codecoverage.out