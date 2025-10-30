BINARY := keysej
PKG := github.com/repsejnworb/keysej

.PHONY: build run tidy fmt vet

build:
	GOFLAGS="-trimpath" go build -ldflags "-s -w -X $(PKG)/internal/version.Version=$$(git describe --tags --always --dirty 2>/dev/null || echo 0.0.0) -X $(PKG)/internal/version.Commit=$$(git rev-parse --short HEAD 2>/dev/null || echo dev)" -o $(BINARY)

run: build
	./$(BINARY) --help

tidy:
	go mod tidy

fmt:
	gofmt -w .

vet:
	go vet ./...