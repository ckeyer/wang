PWD := $(shell pwd)
APP := wang
PKG := github.com/ckeyer/$(APP)
CMS_PKG := ${PKG}/vendor/github.com/ckeyer/commons
GO := CGO_ENABLED=0 go
HASH := $(shell which sha1sum || which shasum)

OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
VERSION := $(shell cat VERSION)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_AT := $(shell date "+%Y-%m-%dT%H:%M:%SZ%z")
PACKAGE_NAME := $(APP)$(VERSION).$(OS)-$(ARCH)

LD_FLAGS := -X $(CMS_PKG)/version.version=$(VERSION) \
 -X $(CMS_PKG)/version.gitCommit=$(GIT_COMMIT) \
 -X $(CMS_PKG)/version.buildAt=$(BUILD_AT) -w

BUILD_IMAGE := ckeyer/dev:go
IMAGE_NAME := ckeyer/$(APP):$(VERSION)

build:
	$(GO) build -v -ldflags="$(LD_FLAGS)" -o bundles/${APP} .

test:
	${GO} test -ldflags="$(LD_FLAGS)" $$(go list ./... |grep -v "vendor")

patch:
	wang patch

tag:
	git tag -f v$(VERSION)
	-git push origin v$(VERSION)
	-hub release create -m "v$(VERSION)" v$(VERSION)

release: clean tag build
	mkdir -p bundles/$(PACKAGE_NAME)
	mv bundles/$(APP) bundles/$(PACKAGE_NAME)
	cd bundles ;\
	 echo $(VERSION) > $(PACKAGE_NAME)/release.txt ;\
	 $(HASH) $(PACKAGE_NAME)/$(APP) > $(PACKAGE_NAME)/sha1.txt ;\
	 tar zcvf $(PACKAGE_NAME).tar.gz $(PACKAGE_NAME);
	hub release edit -m "v$(VERSION)" -a bundles/$(PACKAGE_NAME).tar.gz v${VERSION}

clean:
	rm -rf bundles/*
