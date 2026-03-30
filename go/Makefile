VERSION ?= $(shell jq -r .version manifest.json)
LDFLAGS := -s -w
BIN     := server

.PHONY: build package clean test bump-version version

build:
	mkdir -p $(BIN)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BIN)/e-conomic-connector ./cmd/e-conomic-connector
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(BIN)/e-conomic-connector-darwin ./cmd/e-conomic-connector
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BIN)/e-conomic-connector.exe ./cmd/e-conomic-connector

package: build
	rm -f e-conomic-connector.mcpb
	zip -j e-conomic-connector.mcpb manifest.json
	zip -r e-conomic-connector.mcpb $(BIN)/

test:
	go test ./...

clean:
	rm -rf $(BIN) e-conomic-connector.mcpb

bump-version:
	@test -n "$(V)" || (echo "Usage: make bump-version V=x.y.z" && exit 1)
	jq '.version = "$(V)"' manifest.json > tmp.json && mv tmp.json manifest.json
	@echo "Version set to $(V)"

version:
	@echo $(VERSION)
