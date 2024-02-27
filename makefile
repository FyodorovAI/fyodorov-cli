# Compiler and linker options
GO := go
GOFLAGS := -ldflags="-s -w"

# Target executable names
TARGET_UNIX := fyodorov
TARGET_WINDOWS := fyodorov.exe

# Build targets
build-mac:
	$(GO) build $(GOFLAGS) -o $(TARGET_UNIX) ./cmd/cli

build-unix-amd64:
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(TARGET_UNIX)-amd64 ./cmd/cli

build-windows:
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(TARGET_WINDOWS) ./cmd/cli

build: build-mac build-windows build-unix-amd64

# Clean target
clean:
	rm -f $(TARGET_UNIX) $(TARGET_WINDOWS)
