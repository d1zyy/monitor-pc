VERSION := $(shell git describe --tags --always)
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -Iseconds)
OUTPUT ?= monitor

.PHONY: build check info run clean

info:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"

build: 
	go build \
		-ldflags="-X github.com/d1zyy/monitor-pc/internal/buildinfo.version=$(VERSION) -X github.com/d1zyy/monitor-pc/internal/buildinfo.commit=$(COMMIT) -X github.com/d1zyy/monitor-pc/internal/buildinfo.buildTime=$(BUILD_TIME)" \
		-o $(OUTPUT) \
		./cmd/monitor

check:
	go vet ./...
	go fmt ./...
	go test ./...
	go test -race ./...

run: build
	./$(OUTPUT)

clean:
	rm -f $(OUTPUT)