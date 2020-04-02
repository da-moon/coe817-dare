FROM gitpod/workspace-full-vnc
USER root
ARG SHELLCHECK_VERSION=stable
ARG SHELLCHECK_FORMAT=gcc
RUN curl -fsSL \
    https://raw.githubusercontent.com/da-moon/core-utils/master/bin/fast-apt | sudo bash -s -- \
    --init || true;
# installing and setting up protobuf compiler
RUN apt-get install -yq protobuf-compiler
RUN go get -v github.com/golang/protobuf/protoc-gen-go
RUN go get -v github.com/gogo/protobuf/proto
RUN go get -v github.com/gogo/protobuf/jsonpb
RUN go get -v github.com/gogo/protobuf/protoc-gen-gogo
RUN go get -v github.com/gogo/protobuf/gogoproto
RUN go get -v github.com/gogo/protobuf/protoc-gen-gofast
RUN go get -v github.com/gogo/protobuf/protoc-gen-gogofast
RUN go get -v github.com/gogo/protobuf/protoc-gen-gogofaster
RUN go get -v github.com/gogo/protobuf/protoc-gen-gogoslick
# installing shellcheck
RUN aria2c "https://storage.googleapis.com/shellcheck/shellcheck-${SHELLCHECK_VERSION}.linux.x86_64.tar.xz"
RUN tar -xvf shellcheck-"${SHELLCHECK_VERSION}".linux.x86_64.tar.xz
RUN cp shellcheck-"${SHELLCHECK_VERSION}"/shellcheck /usr/bin/
RUN shellcheck --version
# adding shellcheck running
RUN wget -q -O \
    /usr/bin/run-sc \
    https://raw.githubusercontent.com/da-moon/core-utils/master/bin/run-sc && \
    chmod +x "/usr/bin/run-sc"
RUN curl -fsSL \
    https://raw.githubusercontent.com/da-moon/core-utils/master/bin/get-hashi | sudo bash -s -- || true;
RUN wget -q -O \
    /usr/bin/gitt \
    https://raw.githubusercontent.com/da-moon/core-utils/master/bin/gitt && \
    chmod +x "/usr/bin/gitt"
RUN gitt --init || true;
RUN wget -O ~/vsls-reqs https://aka.ms/vsls-linux-prereq-script && chmod +x ~/vsls-reqs && ~/vsls-reqs
RUN echo 'export PATH="$PATH:/workspace/go/src/github.com/da-moon/coe817-dare/bin"' >>~/.bashrc
RUN echo 'export GO111MODULE=on' >>~/.bashrc
CMD ["bash"]
