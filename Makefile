SOURCEDIR = cmd/lambda
BUILDDIR = bin

VERSION = $(shell date --utc '+%Y%m%d%H%M%S')_$(shell git rev-parse --short HEAD)

SOURCES = $(shell find $(SOURCEDIR) -name '*.go' -type f)
DIRS = $(shell find $(SOURCEDIR) -maxdepth 1 -mindepth 1 -type d)
OBJECTS = $(patsubst $(SOURCEDIR)/%, $(BUILDDIR)/%, $(DIRS)) 

all: clean build

mod:
	go mod tidy

build: mod $(OBJECTS)

$(BUILDDIR)/%: $(SOURCEDIR)/%/main.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o $@ $<

clean:
	@rm -rvf bin