FROM gitpod/workspace-full-vnc
USER root
ARG SHELLCHECK_VERSION=stable
ARG SHELLCHECK_FORMAT=gcc
RUN apt-get update && apt-get install -y \
                                libssl1.0.0 \
                                libkrb5-3 \
                                zlib1g \
                                libicu[0-9][0-9] \
                                gnome-keyring \
                                libsecret-1-0 \
                                desktop-file-utils \
                                x11-utils
# installing base deps
RUN curl -fsSL \
    https://raw.githubusercontent.com/da-moon/core-utils/master/bin/fast-apt | sudo bash -s -- \
    --init || true;
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
RUN echo 'export PATH="/workspace/coe817-dare/bin:$PATH"' >>~/.bashrc
RUN echo 'export GO111MODULE=on' >>~/.bashrc
CMD ["bash"]
