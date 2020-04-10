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
PORT:=8082
RPC_ENDPOINT:=rpc
# demo file size in megabytes
FILE_SIZE=1
PLAIN_PATH:= /tmp/plain
ENCRYPT_PATH:= /tmp/encrypted
DECRYPT_PATH:= /tmp/decrypted
NONCE=6846aba2350ad80a050c2824117acda9bca9c1afeebc160a
KEY=1c0390e0b14b61885fe4cb38ad935eb67f22be8f96c3e7c8c431f412b9cdf328
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
	- bin$(PSEP)dare daemon --log-level=debug --api-addr=127.0.0.1:${PORT} > $(PWD)/server.log 2>&1 &
	- $(call print_completed_target)
demo-encrypt: dd  
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
	- $(eval src_md5=$(shell md5sum $(PLAIN_PATH)))
	- $(eval dst_md5=$(shell md5sum $(ENCRYPT_PATH)))
	- $(eval src_sha=$(shell sha256sum $(PLAIN_PATH)))
	- $(eval dst_sha=$(shell sha256sum $(ENCRYPT_PATH)))
	- $(call print_completed_target,src md5 : $(src_md5))
	- $(call print_completed_target,src sha256 : $(src_sha))
	- $(call print_completed_target,dst md5 : $(dst_md5))
	- $(call print_completed_target,dst sha256 : $(dst_sha))
	- $(call print_completed_target)

demo-decrypt: 
	- $(call print_running_target)
	- $(eval request=$(shell jq -n \
  --arg source "${ENCRYPT_PATH}" \
  --arg destination "${DECRYPT_PATH}" \
  --arg nonce "${NONCE}" \
  --arg key "${KEY}" \
  --arg id "2" \
  --arg method "Service.Decrypt" \
 ' { "jsonrpc": "2.0" , "method":$$method ,"params":{ "source": $$source, "destination":$$destination, "nonce":$$nonce, "key":$$key },"id": $$id } '))
	- $(call print_completed_target,request ---> $(request))
	- jq -n \
		--arg source "${ENCRYPT_PATH}" \
		--arg destination "${DECRYPT_PATH}" \
		--arg nonce "${NONCE}" \
		--arg key "${KEY}" \
		--arg id "2" \
		--arg method "Service.Decrypt" \
	'{"jsonrpc": "2.0", "method":$$method,"params":{"source": $$source, "destination":$$destination, "nonce":$$nonce,"key":$$key},"id": $$id}' | curl \
		-X POST  \
		--silent \
		--header "Authorization: 12445" \
		--header "Content-type: application/json" \
		--data @- \
		http://127.0.0.1:${PORT}/${RPC_ENDPOINT}  | jq -r

	- $(call print_completed_target)

clean:
	- $(call print_running_target)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) go-clean
	- $(call print_completed_target)
kill :
	- $(call print_running_target)
	- $(RM) $(PWD)/server.log
	- for pid in $(shell ps  | grep "dare" | awk '{print $$1}'); do kill -9 "$$pid"; done
	- $(call print_completed_target)
