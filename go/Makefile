VERSION ?= $(shell jq -r .version manifest.json)
LDFLAGS := -s -w
BIN     := server

.PHONY: build package clean test bump-version version

build:
	mkdir -p $(BIN)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BIN)/e-conomic-connector-$(VERSION) ./cmd/e-conomic-connector
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(BIN)/e-conomic-connector-$(VERSION)-darwin ./cmd/e-conomic-connector
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BIN)/e-conomic-connector-$(VERSION).exe ./cmd/e-conomic-connector

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
	jq '.version = "$(V)" | .server.entry_point = "server/e-conomic-connector-$(V)" | .server.mcp_config.command = "$${__dirname}/server/e-conomic-connector-$(V)" | .server.mcp_config.platform_overrides.win32.command = "$${__dirname}/server/e-conomic-connector-$(V).exe" | .server.mcp_config.platform_overrides.darwin.command = "$${__dirname}/server/e-conomic-connector-$(V)-darwin"' manifest.json > tmp.json && mv tmp.json manifest.json
	@echo "Version set to $(V)"

version:
	@echo $(VERSION)
