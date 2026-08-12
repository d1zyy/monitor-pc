VERSION := $(shell git describe --tags --always)
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -Iseconds)

.PHONY: build test info run clean

info:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"

build: test
	go build \
		-ldflags="-X github.com/d1zyy/monitor-pc/internal/buildinfo.version=$(VERSION) -X github.com/d1zyy/monitor-pc/internal/buildinfo.commit=$(COMMIT) -X github.com/d1zyy/monitor-pc/internal/buildinfo.buildTime=$(BUILD_TIME)" \
		-o monitor \
		./cmd/monitor

test:
	go test -v ./...

run: build
	./monitor

clean:
	rm -f monitor