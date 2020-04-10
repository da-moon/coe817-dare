# coe817-dare

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io#https://github.com/da-moon/coe817-dare)

```bash
rm -rf /tmp/go-build* && rm /tmp/plain && rm /tmp/encrypted && make go-build-dare && make run && make demo-encrypt
```

	# - $(eval src_md5=$(shell md5sum $(ENCRYPT_PATH)))
	# - $(eval dst_md5=$(shell md5sum $(DECRYPT_PATH)))
	# - $(eval plain_md5=$(shell md5sum $(PLAIN_PATH)))
	# - $(eval src_sha=$(shell sha256sum $(ENCRYPT_PATH)))
	# - $(eval dst_sha=$(shell sha256sum $(DECRYPT_PATH)))
	# - $(eval plain_sha=$(shell sha256sum $(PLAIN_PATH)))
	# - $(call print_completed_target,src md5 : $(src_md5))
	# - $(call print_completed_target,src sha256 : $(src_sha))
	# - $(call print_completed_target,dst md5 : $(dst_md5))
	# - $(call print_completed_target,dst sha256 : $(dst_sha))
	# - $(call print_completed_target,original plaintext md5 : $(plain_md5))
	# - $(call print_completed_target,original plaintext sha256 : $(plain_sha))