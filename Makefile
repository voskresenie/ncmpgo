BINARY_NAME=ncmpgo

all: build test-quiet

build: deps
	go build -o ${BINARY_NAME} src/*.go

deps:
	go mod download

run:
	go build -o ${BINARY_NAME} src/*.go
	./${BINARY_NAME}

test:
	go test -v ./src

test-quiet:
	go test ./src

clean:
	go clean
	rm ${BINARY_NAME}
