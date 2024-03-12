# Compiler and linker options
GO := go
GOFLAGS := -ldflags="-s -w"

# Target executable names
TARGET_UNIX := toolbar
TARGET_WINDOWS := toolbar.exe

# Build targets
build-mac:
	$(GO) build $(GOFLAGS) -o $(TARGET_UNIX) .

build-windows:
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(TARGET_WINDOWS) .

build: build-mac build-windows

# Clean target
clean:
	rm -f $(TARGET_UNIX) $(TARGET_WINDOWS)
