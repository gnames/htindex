GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
FLAGS_SHARED = $(FLAG_MODULE) CGO_ENABLED=0 GOARCH=amd64
FLAGS_LINUX = $(FLAGS_SHARED) GOOS=linux
FLAGS_MAC = $(FLAGS_SHARED) GOOS=darwin

VERSION=`git describe --tags`
VER = $(shell git describe --tags --abbrev=0)
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

release:
	cd htindex; \
	$(GOCLEAN); \
	$(FLAGS_LINUX) $(GOBUILD); \
	tar zcf /tmp/htindex-$(VER)-linux.tar.gz htindex; \
	$(GOCLEAN); \
	$(FLAGS_MAC) $(GOBUILD); \
	tar zcf /tmp/htindex-$(VER)-mac.tar.gz htindex; \
	$(GOCLEAN);
