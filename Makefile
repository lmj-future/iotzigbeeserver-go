# PREFIX?=/usr/local
# MODE?=docker
# MODULES?=agent hub timer remote-mqtt function-manager function-node8 function-python3 function-python2
MODULE=iot-zigbee
SRC_FILES:=$(shell find main.go config constant controllers crc datainit db dgram globalconstant httprequest interactmodule models publicfunction publicstruct rediscache -type f -name '*.go') # TODO use vpath
CFG_FILES:=$(shell find config -type f -name '*.json')
BIN_FILE:=iot-zigbee
COPY_DIR:=./output
PLATFORM_ALL:=linux/amd64 linux/arm/v7
#PLATFORM_ALL:=darwin/amd64 linux/amd64 linux/arm64 linux/386 linux/arm/v7 linux/arm/v6 linux/arm/v5 linux/ppc64le linux/s390x

# GIT_REV:=git-$(shell git rev-parse --short HEAD)
# GIT_TAG:=$(shell git tag --contains HEAD)
VERSION:=v1.0.0
#VERSION:=$(if $(GIT_TAG),$(GIT_TAG),$(GIT_REV))
# CHANGES:=$(if $(shell git status -s),true,false)

GO_OS:=$(shell go env GOOS)
GO_ARCH:=$(shell go env GOARCH)
GO_ARM:=$(shell go env GOARM)
# GO_FLAGS?=-ldflags "-X 'github.com/h3c/oasis/cmd.Revision=$(GIT_REV)' -X 'github.com/h3c/oasis/cmd.Version=$(VERSION)'"
# GO_TEST_FLAGS?=
# GO_TEST_PKGS?=$(shell go list ./... | grep -v oasis-video-infer)

ifndef PLATFORMS
	GO_OS:=$(shell go env GOOS)
	GO_ARCH:=$(shell go env GOARCH)
	GO_ARM:=$(shell go env GOARM)
	PLATFORMS:=$(if $(GO_ARM),$(GO_OS)/$(GO_ARCH)/$(GO_ARM),$(GO_OS)/$(GO_ARCH))
	ifeq ($(GO_OS),darwin)
		PLATFORMS+=linux/amd64
	endif
else ifeq ($(PLATFORMS),all)
	override PLATFORMS:=$(PLATFORM_ALL)
endif

REGISTRY?=
XFLAGS?=--load
XPLATFORMS:=$(shell echo $(filter-out darwin/amd64,$(PLATFORMS)) | sed 's: :,:g')

OUTPUT:=output
OUTPUT_DIRS:=$(PLATFORMS:%=$(OUTPUT)/%)
OUTPUT_BINS:=$(OUTPUT_DIRS:%=%/bin/$(BIN_FILE))
OUTPUT_PKGS:=$(OUTPUT_DIRS:%=%/$(BIN_FILE)-$(VERSION).zip) # TODO: switch to tar
OUTPUT_CFGS:=$(OUTPUT_DIRS:%=%/config/)
OUTPUT_IMAGE:=$(OUTPUT_DIRS:%=%/image/$(MODULE)-$(VERSION).tar)
OUTPUT_FILE:=$(PLATFORMS:%=$(OUTPUT)/%/file/var/image/)

# OUTPUT_MODS:=$(MODULES:%=oasis-%)
# IMAGE_MODS:=$(MODULES:%=image/oasis-%) # a little tricky to add prefix 'image/' in order to distinguish from OUTPUT_MODS
# NATIVE_MODS:=$(MODULES:%=native/oasis-%) # a little tricky to add prefix 'native/' in order to distinguish from OUTPUT_MODS

.PHONY: all
all: $(BIN_FILE)

$(BIN_FILE): $(OUTPUT_BINS) $(OUTPUT_PKGS) $(OUTPUT_CFGS)

$(OUTPUT_BINS): $(SRC_FILES)
	@echo "BUILD $@"
	@mkdir -p $(dir $@)
	@$(shell echo $(@:$(OUTPUT)/%/bin/$(BIN_FILE)=%)  | sed 's:/v:/:g' | awk -F '/' '{print "CGO_ENABLED=0 GOOS="$$1" GOARCH="$$2" GOARM="$$3" go build"}') -o $@ ${GO_FLAGS} .

$(OUTPUT_PKGS):
	@echo "PACKAGE $@"
	@cd $(dir $@) && zip -q -r $(notdir $@) bin

$(OUTPUT_CFGS): $(CFG_FILES)
	@echo "CFGS $@"
	@mkdir -p $(dir $@)
	@cd $(dir $@)
	@install -m 0755 $(CFG_FILES) $(dir $@)

# $(OUTPUT_MODS):
# 	@make -C $@

# .PHONY: image $(IMAGE_MODS)
# image: $(IMAGE_MODS)

# $(IMAGE_MODS):
# 	@make -C $(notdir $@) image
.PHONY: image
image: $(OUTPUT_BINS)
	@echo "BUILDX: $(REGISTRY)$(MODULE):$(VERSION)"
	@-docker buildx create --name $(BIN_FILE)
	@docker buildx use $(BIN_FILE)
	docker buildx build $(XFLAGS) --platform $(XPLATFORMS) -t $(REGISTRY)$(MODULE):$(VERSION) -f Dockerfile $(COPY_DIR)
	@install -d -m 0755 $(OUTPUT_DIRS:%=%/image)
	@docker save -o $(OUTPUT_IMAGE) $(REGISTRY)$(MODULE):$(VERSION)
	@install -d -m 0755 $(OUTPUT_FILE)
	@install -m 0755 $(OUTPUT_IMAGE) $(OUTPUT_FILE)

.PHONY: rebuild
rebuild: clean all

# .PHONY: test
# test:
# 	@cd oasis-function-node8 && npm install && cd -
# 	@cd oasis-function-python2 && pip install -r requirements.txt && cd -
# 	@cd oasis-function-python3 && pip3 install -r requirements.txt && cd -
# 	@go test ${GO_TEST_FLAGS} -coverprofile=coverage.out ${GO_TEST_PKGS}
# 	@go tool cover -func=coverage.out | grep total

# .PHONY: install $(NATIVE_MODS)
# install: all
# 	@install -d -m 0755 ${PREFIX}/bin
# 	@install -m 0755 $(OUTPUT)/$(if $(GO_ARM),$(GO_OS)/$(GO_ARCH)/$(GO_ARM),$(GO_OS)/$(GO_ARCH))/oasis/bin/oasis-iot-edge ${PREFIX}/bin/
# ifeq ($(MODE),native)
# 	@make $(NATIVE_MODS)
# endif
# 	@tar cf - -C example/$(MODE) etc var | tar xvf - -C ${PREFIX}/

# $(NATIVE_MODS):
# 	@install -d -m 0755 ${PREFIX}/var/db/oasis/$(notdir $@)/bin
# 	@install -m 0755 $(OUTPUT)/$(if $(GO_ARM),$(GO_OS)/$(GO_ARCH)/$(GO_ARM),$(GO_OS)/$(GO_ARCH))/$(notdir $@)/bin/* ${PREFIX}/var/db/oasis/$(notdir $@)/bin/
# 	@install -m 0755 $(OUTPUT)/$(if $(GO_ARM),$(GO_OS)/$(GO_ARCH)/$(GO_ARM),$(GO_OS)/$(GO_ARCH))/$(notdir $@)/package.yml ${PREFIX}/var/db/oasis/$(notdir $@)/
# 	cd ${PREFIX}/var/db/oasis/$(notdir $@)/bin && npm install

# .PHONY: uninstall
# uninstall:
# 	@-rm -f ${PREFIX}/bin/oasis-iot-edge
# 	@-rm -rf ${PREFIX}/etc/oasis
# 	@-rm -rf ${PREFIX}/var/db/oasis
# 	@-rm -rf ${PREFIX}/var/log/oasis
# 	@-rm -rf ${PREFIX}/var/run/oasis
# 	@-rmdir ${PREFIX}/bin
# 	@-rmdir ${PREFIX}/etc
# 	@-rmdir ${PREFIX}/var/db
# 	@-rmdir ${PREFIX}/var/log
# 	@-rmdir ${PREFIX}/var/run
# 	@-rmdir ${PREFIX}/var
# 	@-rmdir ${PREFIX}

# .PHONY: generate
# generate:
# 	go generate ./...

.PHONY: clean
clean:
	@-rm -rf $(OUTPUT)

# .PHONY: fmt
# fmt:
# 	go fmt  ./...