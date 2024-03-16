# Compiler and linker options
GO := go
GOFLAGS := -ldflags="-s -w"

# Target executable names
TARGET_UNIX := toolbar
TARGET_MAC := toolbar-arm
TARGET_WINDOWS := toolbar.exe

# Build targets
build-mac:
	$(GO) build $(GOFLAGS) -o $(TARGET_MAC) .

build-windows:
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(TARGET_WINDOWS) .

build-unix-amd64:
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(TARGET_UNIX) .

build: build-mac build-windows

# Clean target
clean:
	rm -f $(TARGET_MAC) $(TARGET_WINDOWS)
