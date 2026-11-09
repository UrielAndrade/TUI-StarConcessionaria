.PHONY: all build build-linux build-windows test clean run

BINARY_NAME=star-tui

all: test build

build: build-linux build-windows

build-linux:
	@echo "Construindo binário Linux (AMD64)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY_NAME) ./cmd/star-tui

build-windows:
	@echo "Construindo executável Windows (AMD64)..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY_NAME).exe ./cmd/star-tui

test:
	@echo "Executando testes automatizados com validação da API..."
	go test -v ./...

run: build-linux
	./$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe
