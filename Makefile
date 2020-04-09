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
.PHONY: test build clean run demo-encrypt demo-decrypt kill dd
.SILENT: test build clean run demo-encrypt demo-decrypt kill dd
PORT:=8080
RPC_ENDPOINT:=rpc
# demo file size in megabytes
FILE_SIZE=50
PLAIN_PATH:= /tmp/plain
ENCRYPT_PATH:= /tmp/encrypted
DECRYPT_PATH:= /tmp/decrypted
build: 
	- $(call print_running_target)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) go-build
	- $(call print_completed_target)
dd: 
	- $(call print_running_target)
	- dd if=/dev/urandom of=${PLAIN_PATH} bs=1048576 count=${FILE_SIZE}
	- $(call print_completed_target)

run: kill
	- $(call print_running_target)
	- bin$(PSEP)dare daemon --log-level=info --api-addr=127.0.0.1:${PORT} > $(PWD)/server.log 2>&1 &
	- $(call print_completed_target)
kill :
	- $(call print_running_target)
	- $(RM) $(PWD)/server.log
	- for pid in $(shell ps  | grep "dare" | awk '{print $$1}'); do kill -9 "$$pid"; done
	- $(call print_completed_target)
sample-file:
demo-encrypt: sample-file
	- $(call print_running_target)
	- $(eval request=$(shell jq -n \
  --arg source "${PLAIN_PATH}" \
  --arg destination "${ENCRYPT_PATH}" \
  --arg id "1" \
  --arg method "Service.Encrypt" \
 ' { "jsonrpc": "2.0" , "method":$$method ,"params":{ "source": $$source, "destination":$$destination },"id": $$id } '))
	- $(call print_completed_target,request ---> $(request))
	- jq -n \
  --arg source "${PLAIN_PATH}" \
  --arg destination "${ENCRYPT_PATH}" \
  --arg id "1" \
  --arg method "Service.Encrypt" \
 '{"jsonrpc": "2.0", "method":$$method,"params":{"source": $$source, "destination":$$destination},"id": $$id}' | curl \
    -X POST  \
 	--silent \
    --header "Authorization: 12445" \
	--header "Content-type: application/json" \
    --data @- \
    http://127.0.0.1:${PORT}/${RPC_ENDPOINT}  | jq -r 
	- $(call print_completed_target)

demo-decrypt: 
	- $(call print_running_target)
	- $(eval request=$(shell jq -n \
  --arg source "${ENCRYPT_PATH}" \
  --arg destination "${DECRYPT_PATH}" \
  --arg id "2" \
  --arg method "Service.Decrypt" \
 ' { "jsonrpc": "2.0" , "method":$$method ,"params":{ "source": $$source, "destination":$$destination },"id": $$id } '))
	- $(call print_completed_target,request ---> $(request))

	- jq -n \
  --arg source "${ENCRYPT_PATH}" \
  --arg destination "${DECRYPT_PATH}" \
  --arg id "2" \
  --arg method "Service.Decrypt" \
 '{"jsonrpc": "2.0", "method":$$method,"params":{"source": $$source, "destination":$$destination},"id": $$id}' | curl \
    -X POST  \
 	--silent \
    --header "Authorization: 12445" \
	--header "Content-type: application/json" \
    --data @- \
    http://127.0.0.1:${PORT}/${RPC_ENDPOINT}  | jq -r
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
