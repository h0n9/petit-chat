.PHONY: build

all: build

GO := $(shell command -v go 2> /dev/null)
GOSRCS=$(shell find . -name \*.go)
BUILDENV=
BUILDTAGS=

ifeq ($(GO),)
  $(error could not find go. Is it in PATH? $(GO))
endif

ifneq ($(TARGET),)
  BUILDENV += GOOS=$(TARGET)
endif

PROFCMD=go test -cpuprofile=cpu.prof -bench .

build:
	$(BUILDENV) go build ./cmd/petit-chat

install:
	$(BUILDENV) go install ./cmd/petit-chat

test:
	go test ./... --cover

clean:
	rm -f petit-chat *.test
