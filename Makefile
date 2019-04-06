# init project path
HOMEDIR := $(shell pwd)
OUTDIR  := $(HOMEDIR)/output
# init GO & GOD path
export GOPATH  := $(HOMEDIR)/../../../../../
export GOROOT  := $(HOMEDIR)/../../../baidu/go-env/go1-10-3-linux-amd64/
export GODPATH := $(HOMEDIR)/../../../baidu/god-env/god-v0-6-0-linux-amd64/
export PATH    := $(GODPATH)/bin:$(GOPATH)/bin:$(GOROOT)/bin:$(PATH)
# init command params
GO      := go
GOD     := god
GOBUILD := $(GO) build
GOTEST  := $(GO) test
GOPKGS  := $$($(GO) list ./...| grep -vE "vendor")
# make, make all
all: prepare compile package
# make prepare, download dependencies
prepare: prepare-dep
prepare-dep:
	comake2 -UB
	$(GOD) restore -v
# make compile, go build
compile: build
build:
	$(GOBUILD) -o $(HOMEDIR)/minispider
# make test, test your code
test: test-case
test-case:
	$(GOTEST) -v -cover $(GOPKGS)
# make package
package: package-bin
package-bin:
	mkdir -p $(OUTDIR)
	mv minispider  $(OUTDIR)/
# make clean
clean:
	rm -rf $(OUTDIR)
	rm -rf $(HOMEDIR)/minispider
	rm -rf $(GOPATH)/pkg/darwin_amd64
# avoid filename conflict and speed up build 
.PHONY: all prepare compile test package clean build
