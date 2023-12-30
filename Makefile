BINARY_NAME=go-notes

MAIN_PATH=cmd/main.go

GO=go

LDFLAGS=-ldflags "-s -w"

all: build

build:
	$(GO) build $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)

clean:
	@rm -f $(BINARY_NAME)

.PHONY: all build clean
