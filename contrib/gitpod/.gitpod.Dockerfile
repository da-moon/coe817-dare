FROM gitpod/workspace-full-vnc
USER root
ARG SHELLCHECK_VERSION=stable
ARG SHELLCHECK_FORMAT=gcc
# installing base deps
RUN curl -fsSL \
    https://raw.githubusercontent.com/da-moon/core-utils/master/bin/fast-apt | sudo bash -s -- \
    --init
# installing shellcheck
RUN aria2c "https://storage.googleapis.com/shellcheck/shellcheck-${SHELLCHECK_VERSION}.linux.x86_64.tar.xz"
RUN tar -xvf shellcheck-"${SHELLCHECK_VERSION}".linux.x86_64.tar.xz
RUN cp shellcheck-"${SHELLCHECK_VERSION}"/shellcheck /usr/bin/
RUN shellcheck --version
# adding shellcheck running
RUN wget -q -O \
    /usr/bin/run-sc \
    https://raw.githubusercontent.com/da-moon/core-utils/master/bin/run-sc
RUN curl -fsSL \
    https://raw.githubusercontent.com/da-moon/core-utils/master/bin/get-hashi | sudo bash -s --
RUN wget -q -O /usr/bin/gitt https://raw.githubusercontent.com/da-moon/core-utils/master/bin/gitt
RUN chmod +x "/usr/bin/gitt"
RUN gitt --init
RUN echo 'export PATH="/workspace/bifrost/bin:$PATH"' >>~/.bashrc
CMD ["bash"]
