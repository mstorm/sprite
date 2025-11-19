.PHONY: help install-deps build test clean run

help:
	@echo "Available targets:"
	@echo "  make install-deps  - Install resvg dependency"
	@echo "  make build         - Build the sprite binary"
	@echo "  make test          - Run tests"
	@echo "  make run           - Run with test data"
	@echo "  make clean         - Clean build artifacts"

install-deps:
	@echo "Checking for resvg..."
	@which resvg > /dev/null 2>&1 || { \
		echo "resvg not found. Installing..."; \
		if which cargo > /dev/null 2>&1; then \
			cargo install resvg; \
		else \
			echo "Error: cargo (Rust) is required to install resvg"; \
			echo "Please install Rust from https://rustup.rs/"; \
			echo "Or use pre-built binaries from https://github.com/RazrFalcon/resvg/releases"; \
			exit 1; \
		fi \
	}
	@echo "resvg is installed at: $$(which resvg)"

build:
	@echo "Building sprite..."
	go build -o sprite .

test:
	@echo "Running tests..."
	go test ./...

clean:
	@echo "Cleaning..."
	rm -f sprite
	rm -f sprites*.png sprites*.json sprites.svg

run: build
	./sprite testdata/dot.svg testdata/restaurant.svg testdata/mountain.svg testdata/airport.svg

all: install-deps build test
