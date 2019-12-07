FROM busybox:1.30-glibc

LABEL maintainer="Kevin Kamps"
LABEL github="https://github.com/kevinkamps/Registrar"
LABEL license="GPL-3.0"

COPY bin/registrar_linux_amd64 /bin/registrar

ENTRYPOINT ["/bin/registrar"]
