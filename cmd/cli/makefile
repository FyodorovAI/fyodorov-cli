# Compiler and linker options
GO := go
GOFLAGS := -ldflags="-s -w"

# Target executable names
TARGET_UNIX := fyodorov
TARGET_WINDOWS := fyodorov.exe

VERSION := 0.4.7

# Build targets
build-mac:
	$(GO) build $(GOFLAGS) -o $(TARGET_UNIX) .

build-unix-amd64:
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(TARGET_UNIX)-amd64 .

build-windows:
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(TARGET_WINDOWS) .

build: build-mac build-unix-amd64 # build-windows

bump-patch:
	@NEW_VERSION=$(shell echo $(VERSION) | awk -F. '{print $$1"."$$2"."$$3+1}') && \
	echo "Bumping version from $(VERSION) to $${NEW_VERSION}" && \
	sed -i '' "s/$(VERSION)/$${NEW_VERSION}/g" update_command.go Makefile && \
	echo "Version bumped to $${NEW_VERSION}"

# Clean target
clean:
	rm -f $(TARGET_UNIX) $(TARGET_WINDOWS)
