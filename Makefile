GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean

VERSION=`git describe --tags`
DATE=`date -u '+%Y-%m-%d_%I:%M:%S%p'`
LDFLAGS=-ldflags "-X main.buildDate=${DATE} -X main.buildVersion=${VERSION}"

all: install

build: grpc
	cd htindex && \
	$(GOCLEAN) && \
	GO111MODULE=on GOOS=linux GOARCH=amd64 $(GOBUILD) ${LDFLAGS}

install:
	cd htindex && \
	GO111MODULE=on $(GOINSTALL) ${LDFLAGS};