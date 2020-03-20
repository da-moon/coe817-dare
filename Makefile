include build/makefiles/functions/functions.mk
include build/makefiles/git/git.mk

ifeq ($(OS),Windows_NT)
    GO_PATH = $(subst \,/,${GOPATH})
else
    GO_PATH = ${GOPATH}
endif
THIS_FILE := $(lastword $(MAKEFILE_LIST))
SELF_DIR := $(dir $(THIS_FILE))

GO_TARGET = $(notdir $(patsubst %/,%,$(dir $(wildcard ./cmd/*/.))))
CGO=0
GO_ARCHITECTURE=amd64
GO_IMAGE=golang:alpine
MOD=on
GO_PKG=github.com/da-moon/coe817-dare
.PHONY: build full-build build-mac-os build-linux build-windows go-clean go-print
.SILENT: build full-build build-mac-os build-linux build-windows go-clean go-print
go-print:
	- $(info some random stuff)
	- $(info $(GO_PKG))

build: go-clean
    ifeq ($(DOCKER_ENV),true)
ifeq (${MOD},on)
ifeq ($(wildcard ./go.mod),)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod init"
endif
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod tidy"
endif
	for target in $(GO_TARGET); do \
            $(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd=" \
            GO111MODULE=${MOD} \
            CGO_ENABLED=${CGO} \
            GOARCH=${GO_ARCHITECTURE} \
            GOOS=${GOOS} \
            go build -a -installsuffix cgo -ldflags \
            '-X ${GO_PKG}/version.Version=${VERSION} \
			-X ${GO_PKG}/version.Revision=${REVISION} \
			-X ${GO_PKG}/version.Branch=${BRANCH} \
			-X ${GO_PKG}/version.BuildUser=${BUILDUSER} \
			-X ${GO_PKG}/version.BuildDate=${BUILDTIME}' \
			-o .$(PSEP)bin$(PSEP)$$target .$(PSEP)cmd$(PSEP)$$target"; \
	done

    endif
    ifeq ($(DOCKER_ENV),false)
ifeq (${MOD},on)
ifeq ($(wildcard ./go.mod),)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod init"
endif
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod tidy"
endif
	for target in $(GO_TARGET); do \
            $(RM) .$(PSEP)bin$(PSEP)$$target$(PSEP)${GOOS}$(PSEP)${VERSION}; \
			$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="CGO_ENABLED=${CGO} \
            GO111MODULE=$(MOD) \
            GOARCH=${GO_ARCHITECTURE} \
            GOOS=${GOOS} \
            go build -a -installsuffix cgo -ldflags \
            '-X ${GO_PKG}/version.Version=${VERSION} \
			-X ${GO_PKG}/version.Revision=${REVISION} \
			-X ${GO_PKG}/version.Branch=${BRANCH} \
			-X ${GO_PKG}/version.BuildUser=${BUILDUSER} \
			-X ${GO_PKG}/version.BuildDate=${BUILDTIME}' \
			-o .$(PSEP)bin$(PSEP)$$target .$(PSEP)cmd$(PSEP)$$target"; \
	done
    endif

full-build:
	- $(CLEAR)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) build-linux
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) build-windows
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) build-mac-os

build-linux:  
	- $(eval GOOS := linux)
    ifeq ($(DOCKER_ENV),true)
ifeq (${MOD},on)
ifeq ($(wildcard ./go.mod),)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod init"
endif
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod tidy"
endif
	for target in $(GO_TARGET); do \
            $(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="rm -rf /go/src/github.com/da-moon/go-packages/bin/$$target/${GOOS}/${VERSION} && \
            GO111MODULE=${MOD} \
            CGO_ENABLED=${CGO} \
            GOARCH=${GO_ARCHITECTURE} \
            GOOS=${GOOS} \
            go build -a -installsuffix cgo -ldflags \
            '-X ${GO_PKG}/version.Version=${VERSION} \
			-X ${GO_PKG}/version.Revision=${REVISION} \
			-X ${GO_PKG}/version.Branch=${BRANCH} \
			-X ${GO_PKG}/version.BuildUser=${BUILDUSER} \
			-X ${GO_PKG}/version.BuildDate=${BUILDTIME}' \
			-o .$(PSEP)bin$(PSEP)$$target$(PSEP)${GOOS}$(PSEP)${VERSION}$(PSEP)$$target .$(PSEP)cmd$(PSEP)$$target"; \
	done

    endif
    ifeq ($(DOCKER_ENV),false)
ifeq (${MOD},on)
ifeq ($(wildcard ./go.mod),)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod init"
endif
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod tidy"
endif
	for target in $(GO_TARGET); do \
            $(RM) .$(PSEP)bin$(PSEP)$$target$(PSEP)${GOOS}$(PSEP)${VERSION}; \
			$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="CGO_ENABLED=${CGO} \
            GO111MODULE=$(MOD) \
            GOARCH=${GO_ARCHITECTURE} \
            GOOS=${GOOS} \
            go build -a -installsuffix cgo -ldflags \
            '-X ${GO_PKG}/version.Version=${VERSION} \
			-X ${GO_PKG}/version.Revision=${REVISION} \
			-X ${GO_PKG}/version.Branch=${BRANCH} \
			-X ${GO_PKG}/version.BuildUser=${BUILDUSER} \
			-X ${GO_PKG}/version.BuildDate=${BUILDTIME}' \
			-o .$(PSEP)bin$(PSEP)$$target$(PSEP)${GOOS}$(PSEP)${VERSION}$(PSEP)$$target .$(PSEP)cmd$(PSEP)$$target"; \
	done
    endif

build-windows:
	- $(eval GOOS := windows)
    ifeq ($(DOCKER_ENV),true)
ifeq (${MOD},on)
ifeq ($(wildcard ./go.mod),)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod init"
endif
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod tidy"
endif
	for target in $(GO_TARGET); do \
            $(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="rm -rf /go/src/github.com/da-moon/go-packages/bin/$$target/${GOOS}/${VERSION} && \
            GO111MODULE=${MOD} \
            CGO_ENABLED=${CGO} \
            GOARCH=${GO_ARCHITECTURE} \
            GOOS=${GOOS} \
            go build -a -installsuffix cgo -ldflags \
            '-X ${GO_PKG}/version.Version=${VERSION} \
			-X ${GO_PKG}/version.Revision=${REVISION} \
			-X ${GO_PKG}/version.Branch=${BRANCH} \
			-X ${GO_PKG}/version.BuildUser=${BUILDUSER} \
			-X ${GO_PKG}/version.BuildDate=${BUILDTIME}' \
			-o .$(PSEP)bin$(PSEP)$$target$(PSEP)${GOOS}$(PSEP)${VERSION}$(PSEP)$$target.exe .$(PSEP)cmd$(PSEP)$$target"; \
	done

    endif
    ifeq ($(DOCKER_ENV),false)
ifeq (${MOD},on)
ifeq ($(wildcard ./go.mod),)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod init"
endif
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod tidy"
endif
	for target in $(GO_TARGET); do \
            $(RM) .$(PSEP)bin$(PSEP)$$target$(PSEP)${GOOS}$(PSEP)${VERSION}; \
			$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="CGO_ENABLED=${CGO} \
            GO111MODULE=$(MOD) \
            GOARCH=${GO_ARCHITECTURE} \
            GOOS=${GOOS} \
            go build -a -installsuffix cgo -ldflags \
            '-X ${GO_PKG}/version.Version=${VERSION} \
			-X ${GO_PKG}/version.Revision=${REVISION} \
			-X ${GO_PKG}/version.Branch=${BRANCH} \
			-X ${GO_PKG}/version.BuildUser=${BUILDUSER} \
			-X ${GO_PKG}/version.BuildDate=${BUILDTIME}' \
			-o .$(PSEP)bin$(PSEP)$$target$(PSEP)${GOOS}$(PSEP)${VERSION}$(PSEP)$$target .$(PSEP)cmd$(PSEP)$$target"; \
	done
    endif

build-mac-os:
	- $(eval GOOS := darwin)
    ifeq ($(DOCKER_ENV),true)
ifeq (${MOD},on)
ifeq ($(wildcard ./go.mod),)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod init"
endif
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod tidy"
endif
	for target in $(GO_TARGET); do \
            $(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="rm -rf /go/src/github.com/da-moon/go-packages/bin/$$target/${GOOS}/${VERSION} && \
            GO111MODULE=${MOD} \
            CGO_ENABLED=${CGO} \
            GOARCH=${GO_ARCHITECTURE} \
            GOOS=${GOOS} \
            go build -a -installsuffix cgo -ldflags \
            '-X ${GO_PKG}/version.Version=${VERSION} \
			-X ${GO_PKG}/version.Revision=${REVISION} \
			-X ${GO_PKG}/version.Branch=${BRANCH} \
			-X ${GO_PKG}/version.BuildUser=${BUILDUSER} \
			-X ${GO_PKG}/version.BuildDate=${BUILDTIME}' \
			-o .$(PSEP)bin$(PSEP)$$target$(PSEP)${GOOS}$(PSEP)${VERSION}$(PSEP)$$target .$(PSEP)cmd$(PSEP)$$target"; \
	done

    endif
    ifeq ($(DOCKER_ENV),false)
ifeq (${MOD},on)
ifeq ($(wildcard ./go.mod),)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod init"
endif
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod tidy"
endif
	for target in $(GO_TARGET); do \
            $(RM) .$(PSEP)bin$(PSEP)$$target$(PSEP)${GOOS}$(PSEP)${VERSION}; \
			$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="CGO_ENABLED=${CGO} \
            GO111MODULE=$(MOD) \
            GOARCH=${GO_ARCHITECTURE} \
            GOOS=${GOOS} \
            go build -a -installsuffix cgo -ldflags \
            '-X ${GO_PKG}/version.Version=${VERSION} \
			-X ${GO_PKG}/version.Revision=${REVISION} \
			-X ${GO_PKG}/version.Branch=${BRANCH} \
			-X ${GO_PKG}/version.BuildUser=${BUILDUSER} \
			-X ${GO_PKG}/version.BuildDate=${BUILDTIME}' \
			-o .$(PSEP)bin$(PSEP)$$target$(PSEP)${GOOS}$(PSEP)${VERSION}$(PSEP)$$target .$(PSEP)cmd$(PSEP)$$target"; \
	done
    endif

go-clean:
	- $(CLEAR)
    ifeq ($(DOCKER_ENV),true)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="rm -rf /go/src/${GO_PKG}/bin/"
    else
	- $(RM) ./bin/
    endif
