FROM busybox:1.30-glibc

COPY bin/registrar_linux_amd64 /bin/registrar

ENTRYPOINT ["/bin/registrar"]
