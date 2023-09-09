SOURCEDIR = cmd/lambda
BUILDDIR = bin

all: clean build

mod:
	go mod tidy

build: mod bin/bootstrap

$(BUILDDIR)/bootstrap: $(SOURCEDIR)/proxy/main.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o $@ $<

clean:
	@rm -rvf bin