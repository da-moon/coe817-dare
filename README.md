# coe817-dare

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io#https://github.com/da-moon/coe817-dare)

#   - $(eval bin_path = $(bin_path)$(PSEP)$(name))
# 	- $(eval command = GO111MODULE=${MOD})
# 	- $(eval command = $(command) CGO_ENABLED=${CGO})
# 	- $(eval command = $(command) GOARCH=${GO_ARCHITECTURE})
# ifneq ($(GOOS), )
# 	- $(eval command = $(command) GOOS=${GOOS})
# endif
# 	- $(info $(name))
	# - $(eval command = $(command) go build -a -installsuffix cgo)
	# - $(eval command = $(command) -o $(bin_path) .$(PSEP)cmd$(PSEP)$(name))
	# - $(info $(command))
	# - $(eval command= $(command) go build -a -installsuffix cgo \
	# 		-o .$(PSEP)bin$(PSEP)$(name) .$(PSEP)cmd$(PSEP)${dir}$(PSEP)$(name) \
	# 	)
    # ifeq ($(DOCKER_ENV),true)
	# - @$(MAKE) --no-print-directory \
	#  -f $(THIS_FILE) shell \
	#  docker_image="${GO_IMAGE}" \
	#  container_name="go_builder_container" \
	#  mount_point="/go/src/${GO_PKG}" \
	#  cmd="$(command)"
    # endif
    # ifeq ($(DOCKER_ENV),false)
	# - @$(MAKE) --no-print-directory \
	#  -f $(THIS_FILE) shell cmd="$(command)"
    # endif