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
GO_IMAGE=golang:buster
MOD=on
GO_PKG=github.com/da-moon/coe817-dare
.PHONY: go-build full-build build-mac-os build-linux build-windows go-clean go-dependancy go-print
.SILENT: go-build full-build build-mac-os build-linux build-windows go-clean go-dependancy go-print
go-print:
	- $(info some random stuff)
	- $(info $(GO_PKG))

go-dependancy:
	- $(call print_running_target)
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
ifeq (${MOD},off)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go get -v -d ./..."
endif
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
ifeq (${MOD},off)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go get -v -d ./..."
endif
    endif
	- $(call print_completed_target)
go-build: go-clean go-dependancy
	- $(call print_running_target)
    ifeq ($(DOCKER_ENV),true)

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
	- $(call print_completed_target)

full-build: go-clean go-dependancy
	- $(CLEAR)
	- $(call print_running_target)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) build-linux
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) build-windows
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) build-mac-os
	- $(call print_completed_target)

build-linux:  
	- $(call print_running_target)
	- $(eval GOOS := linux)
    ifeq ($(DOCKER_ENV),true)
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
	- $(call print_completed_target)

build-windows:
	- $(call print_running_target)
	- $(eval GOOS := windows)
    ifeq ($(DOCKER_ENV),true)
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
	- $(call print_completed_target)

build-mac-os:
	- $(call print_running_target)
	- $(eval GOOS := darwin)
    ifeq ($(DOCKER_ENV),true)
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
	- $(call print_completed_target)

go-clean:
	- $(CLEAR)
	- $(call print_running_target)
    ifeq ($(DOCKER_ENV),true)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="rm -rf /go/src/${GO_PKG}/bin/"
    else
	- $(RM) ./bin/
    endif
	- $(call print_completed_target)
