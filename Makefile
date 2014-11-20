all:

fmt:
	go fmt

run: fmt
	go run cigar.go

build: fmt
	go build -o cigar cigar.go

install: build
	go install

godep:
	godep save -r

.PHONY: fmt run build install godep
