IMAGE   ?= keppel.eu-de-1.cloud.sap/ccloud/baremetal_temper
VERSION = $(shell git rev-parse --verify HEAD | head -c 8)

GOOS    ?= $(shell go env | grep GOOS | cut -d'"' -f2)
BINARIES := temper scheduler

LDFLAGS := -X github.com/sapcc/baremetal_temper/pkg/baremetal_temper.VERSION=$(VERSION)
GOFLAGS := -ldflags "$(LDFLAGS)"

SRCDIRS  := cmd pkg internal
PACKAGES := $(shell find $(SRCDIRS) -type d)
GOFILES  := $(addsuffix /*.go,$(PACKAGES))
GOFILES  := $(wildcard $(GOFILES))


all: $(BINARIES:%=bin/$(GOOS)/%)

bin/%: $(GOFILES) Makefile
	GOOS=$(*D) GOARCH=amd64 go build $(GOFLAGS) -v -i -o $(@D)/$(@F) ./cmd/$(basename $(@F))

build:
	docker build -t $(IMAGE):$(VERSION) .

push: build
	docker push $(IMAGE):$(VERSION)

clean:
	rm -rf bin/*

vendor:
	GO111MODULE=on go get -u ./... && go mod tidy && go mod vendor