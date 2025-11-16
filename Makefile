.PHONY: all build clean install

# Binary names
BINARIES = al algo alinit alnote allink

# Build directory
BUILD_DIR = build

# Go build flags
GOFLAGS = -ldflags="-s -w"

all: build

build:
	@echo "Building al CLI..."
	@mkdir -p $(BUILD_DIR)
	@go build $(GOFLAGS) -o $(BUILD_DIR)/al .
	@echo "✓ Built al"
	@cd $(BUILD_DIR) && for binary in $(BINARIES); do \
		if [ "$$binary" != "al" ]; then \
			ln -sf al $$binary; \
			echo "✓ Created symlink $$binary"; \
		fi \
	done
	@echo "All binaries ready in $(BUILD_DIR)/"

clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
	@echo "✓ Clean complete"

install: build
	@echo "Installing al CLI..."
	@cd $(BUILD_DIR) && sudo ./al install
	@echo "✓ Installation complete"

test:
	@go test -v ./...

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies ready"

help:
	@echo "Available targets:"
	@echo "  make build    - Build all binaries"
	@echo "  make clean    - Remove build directory"
	@echo "  make install  - Build and install to system"
	@echo "  make test     - Run tests"
	@echo "  make deps     - Download and tidy dependencies"
	@echo "  make help     - Show this help message"
