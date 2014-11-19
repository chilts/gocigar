all:

fmt:
	go fmt

run: fmt
	go run cigar.go

build: fmt
	go build

install: build
	go install

.PHONY: fmt run build install
