export GOPATH = $(shell pwd)

all: eavesdrop

eavesdrop: eavesdrop.go
	go get -d -x
	go build eavesdrop.go

clean:
	go clean
	rm -rf src

