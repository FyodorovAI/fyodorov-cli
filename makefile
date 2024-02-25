# Compiler and flags
CC = gcc
CFLAGS = -Wall -Wextra

# Source files
SRCS = main.c utils.c

# Object files
OBJS = $(SRCS:.c=.o)

# Target executable
TARGET = cli_tool

# Default target
all: $(TARGET)

# Compile source files into object files
%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

# Link object files into the target executable
$(TARGET): $(OBJS)
	$(CC) $(CFLAGS) $^ -o $@

# Clean up object files and the target executable
clean:
	rm -f $(OBJS) $(TARGET)
# Compiler and linker options
GO := go
GOFLAGS := -ldflags="-s -w"

# Build target
build:
	$(GO) build $(GOFLAGS) -o mycli

build-windows:
	$(GO) build $(GOFLAGS) -o mycli.exe

# Clean target
clean:
	rm -f mycli
	
build:
	go build -o fyodorov ./cmd/cli
