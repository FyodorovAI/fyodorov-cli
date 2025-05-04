# Compiler and linker options
GO := go
GOFLAGS := -ldflags="-s -w"

# Target executable names
TARGET_UNIX := fyodorov
TARGET_WINDOWS := fyodorov.exe

# Build targets
build-mac:
	$(GO) build $(GOFLAGS) -o $(TARGET_UNIX) .

build-unix-amd64:
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(TARGET_UNIX)-amd64 .

build-windows:
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(TARGET_WINDOWS) .

build: build-mac build-unix-amd64 # build-windows

# Clean target
clean:
	rm -f $(TARGET_UNIX) $(TARGET_WINDOWS)
