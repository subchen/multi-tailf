BASEPATH := $(shell pwd)
BASENAME := mtailf
VERSION  := 1.0.0

LDFLAGS := -s -w
LDFLAGS := $(LDFLAGS) -X 'main.BuildVersion=$(shell git rev-list HEAD --count)'
LDFLAGS := $(LDFLAGS) -X 'main.BuildGitCommit=$(shell git describe --abbrev=0 --always)'
LDFLAGS := $(LDFLAGS) -X 'main.BuildDate=$(shell date -u -R)'

#GOPATH := $(BASEPATH)/../../../../

clean:
	rm -rf ./bin $(BASENAME)

fmt:
	go fmt ./...

test: clean
	go test -v ./...

run: fmt
	go build -o $(BASENAME)
	./$(BASENAME)

build-linux: clean
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$(BASENAME)-linux-$(VERSION)

build-darwin: clean
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$(BASENAME)-darwin-$(VERSION)

build-windows: clean
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$(BASENAME)-windows-$(VERSION).exe

build: build-linux build-darwin build-windows

