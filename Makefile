TAILSCALE_EXPORTER := $(shell git describe --tags)
LDFLAGS += -X "main.BuildTimestamp=$(shell date -u "+%Y-%m-%d %H:%M:%S")"
LDFLAGS += -X "main.tailscaleExporterVersion=$(tailscale_exporter)"
LDFLAGS += -X "main.goVersion=$(shell go version | sed -r 's/go version go(.*)\ .*/\1/')"

GO := GO111MODULE=on CGO_ENABLED=0 go

.PHONY: build
build:
	$(GO) build -ldflags '$(LDFLAGS)'

.PHONY: install
install:
	@echo "Installing tailscale-exporter..."
	@$(GO) install -ldflags '$(LDFLAGS)'

.PHONY: release
release:
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags '$(LDFLAGS)' -o bin/darwin/tailscale-exporter
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags '$(LDFLAGS)' -o bin/linux/tailscale-exporter
