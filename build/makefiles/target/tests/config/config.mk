.SILENT: encrypted-size-test decrypted-size-test
.PHONY: encrypted-size-test decrypted-size-test
encrypted-size-test: 
	- $(call print_running_target)
ifeq ($(DOCKER_ENV),false)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd=" GO111MODULE=${MOD}\
       CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go test -timeout -30s github.com/da-moon/coe817-dare/internal/config -run TestEncryptedSize -v \
    "
	- exit 
endif
ifeq ($(DOCKER_ENV),true)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go get -v -d ./... & \
	go test -timeout -30s github.com/da-moon/coe817-dare/internal/config -run TestEncryptedSize -v "
	- exit 
endif
	- $(call print_completed_target)

decrypted-size-test: 
	- $(call print_running_target)
ifeq ($(DOCKER_ENV),false)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd=" GO111MODULE=${MOD}\
       CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go test -timeout -30s github.com/da-moon/coe817-dare/internal/config -run TestDecryptedSize -v \
    "
	- exit 
endif
ifeq ($(DOCKER_ENV),true)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go get -v -d ./... & \
	go test -timeout -30s github.com/da-moon/coe817-dare/internal/config -run TestDecryptedSize -v "
	- exit 
endif
	- $(call print_completed_target)
