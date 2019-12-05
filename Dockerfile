FROM golang:latest AS build-env
RUN export GO111MODULE=on
RUN mkdir -p /app
WORKDIR /app
COPY . ./
RUN env
RUN GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$SOURCE_BRANCH" -o /app/registrar_linux_amd64 ./main.go


VOLUME /app


FROM busybox:1.30-glibc

LABEL maintainer="Kevin Kamps"
LABEL github="https://github.com/kevinkamps/Registrar"
LABEL license="GPL-3.0"

COPY --from=build-env /app/registrar_linux_amd64 /bin/registrar

ENTRYPOINT ["/bin/registrar"]
