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
.PHONY: test build clean run demo-encrypt encrypt-first demo-decrypt decrypt-first kill linux-dd dd temp-clean
.SILENT: test build clean run demo-encrypt encrypt-first demo-decrypt decrypt-first kill linux-dd dd temp-clean
PORT:=8082
RPC_ENDPOINT:=rpc
# demo file size in megabytes
FILE_SIZE=50
PLAIN_PATH:= /tmp/plain
ENCRYPT_PATH:= /tmp/encrypted
DECRYPT_PATH:= /tmp/decrypted
NONCE=f8a9c278c70c0fd7083eb4050895ceeb8c3748b9fd16f920
KEY=4d6a9b0b51c29adeb6c1b4f25606796b6dd94d6bcf1244d5f05871a99c1750e1
build: 
	- $(call print_running_target)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) go-build
	- $(call print_completed_target)
linux-dd: 
	- $(call print_running_target)
	- dd if=/dev/urandom of=${PLAIN_PATH} bs=1048576 count=${FILE_SIZE}
	- $(call print_completed_target)
dd : 
	- $(call print_running_target)
	- bin$(PSEP)dare dd --size=${FILE_SIZE}MB --path=${PLAIN_PATH}
	- $(call print_completed_target)
run: kill
	- $(call print_running_target)
	- bin$(PSEP)dare daemon --api-addr=127.0.0.1:${PORT} > $(PWD)/server.log 2>&1 &
	- $(call print_completed_target)
encrypt-first: 
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

demo-encrypt: encrypt-first
	- $(call print_running_target)
	- $(eval demo_enc_src_md5=$(shell md5sum $(PLAIN_PATH)))
	- $(eval demo_enc_dst_md5=$(shell md5sum $(ENCRYPT_PATH)))
	- $(eval demo_enc_src_sha=$(shell sha256sum $(PLAIN_PATH)))
	- $(eval demo_enc_dst_sha=$(shell sha256sum $(ENCRYPT_PATH)))
	- $(call print_completed_target,plaintext md5 : $(demo_enc_src_md5))
	- $(call print_completed_target,encrypted md5 : $(demo_enc_dst_md5))
	- $(call print_completed_target,plaintext sha256 : $(demo_enc_src_sha))
	- $(call print_completed_target,encrypted sha256 : $(demo_enc_dst_sha))
	- $(call print_completed_target)


decrypt-first:
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

demo-decrypt: decrypt-first
	- $(call print_running_target)
	- $(eval demo_dec_src_md5=$(shell md5sum $(ENCRYPT_PATH)))
	- $(eval demo_dec_dst_md5=$(shell md5sum $(DECRYPT_PATH)))
	- $(eval demo_dec_plain_md5=$(shell md5sum $(PLAIN_PATH)))
	- $(eval demo_dec_src_sha=$(shell sha256sum $(ENCRYPT_PATH)))
	- $(eval demo_dec_dst_sha=$(shell sha256sum $(DECRYPT_PATH)))
	- $(eval demo_dec_plain_sha=$(shell sha256sum $(PLAIN_PATH)))
	- $(call print_completed_target,encrypted md5 : $(demo_dec_src_md5))
	- $(call print_completed_target,decrypted md5 : $(demo_dec_dst_md5))
	- $(call print_completed_target,original plaintext md5 : $(demo_dec_plain_md5))
	- $(call print_completed_target,encrypted sha256 : $(demo_dec_src_sha))
	- $(call print_completed_target,decrypted sha256 : $(demo_dec_dst_sha))
	- $(call print_completed_target,original plaintext sha256 : $(demo_dec_plain_sha))
clean:
	- $(call print_running_target)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) go-clean
	- $(call print_completed_target)
kill : temp-clean
	- $(call print_running_target)
	- $(RM) $(PWD)/server.log
	- for pid in $(shell ps  | grep "dare" | awk '{print $$1}'); do kill -9 "$$pid"; done
	- $(call print_completed_target)
temp-clean:
	- $(call print_running_target)
	- $(RM) /tmp/go-build*
	- $(call print_completed_target)
