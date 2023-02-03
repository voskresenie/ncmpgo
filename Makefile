BINARY_NAME=ncmpgo

all: build test-quiet

build: deps
	go build -o ${BINARY_NAME} src/*.go

deps:
	go mod download

run: build
	./${BINARY_NAME}

test:
	go test -v ./src

test-quiet:
	go test ./src

test-json:
	go test -json ./src

lint:
	gofmt -l ./src

clean:
	go clean
	rm ${BINARY_NAME}
