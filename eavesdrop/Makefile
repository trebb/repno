export GOPATH = $(shell pwd)

all: eavesdrop

eavesdrop: eavesdrop.go
	go get -d -x
	go build $^

clean:
	go clean
	rm -rf src

