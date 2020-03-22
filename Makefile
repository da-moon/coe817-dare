include build/makefiles/pkg/base/base.mk
include build/makefiles/pkg/string/string.mk
include build/makefiles/pkg/color/color.mk
include build/makefiles/pkg/functions/functions.mk
include build/makefiles/target/buildenv/buildenv.mk
include build/makefiles/target/go/go.mk
include build/makefiles/target/tests/header/header.mk
include build/makefiles/target/tests/config/config.mk
include build/makefiles/target/tests/dare/dare.mk
THIS_FILE := $(firstword $(MAKEFILE_LIST))
SELF_DIR := $(dir $(THIS_FILE))
.PHONY: test build clean
.SILENT: test build clean
build: 
	- $(call print_running_target)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) go-build
	- $(call print_completed_target)

test:
	- $(CLEAR)
	- $(call print_running_target)
	# - @$(MAKE) --no-print-directory -f $(THIS_FILE) encrypted-size-test
	# - @$(MAKE) --no-print-directory -f $(THIS_FILE) decrypted-size-test
	# - @$(MAKE) --no-print-directory -f $(THIS_FILE) header-length-test
	# - @$(MAKE) --no-print-directory -f $(THIS_FILE) header-nonce-test
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) basic-encryption-test
	- $(call print_completed_target)

clean:
	- $(call print_running_target)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) go-clean
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) test-clean
	- $(call print_completed_target)
