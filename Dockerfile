# Build terraprism without a local Go toolchain.
#
#   docker build -o dist .                 # -> dist/terraprism (linux, host arch)
#   docker build -o dist --build-arg VERSION=0.12.0-local .
#   docker build -o dist --build-arg GOOS=darwin --build-arg GOARCH=arm64 .
#
# The container is always Linux, so the default output is a Linux binary even on
# a Mac. Set GOOS/GOARCH to cross-compile - Go does this without extra toolchains.
#
# The final stage is scratch, so `-o` writes the binary to the host instead of
# producing an image. There is no "run terraprism in a container" target on
# purpose: it needs your terraform binary, credentials, and TTY.

FROM golang:1.24-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
ARG TARGETOS
ARG TARGETARCH
ARG GOOS=${TARGETOS}
ARG GOARCH=${TARGETARCH}

RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} \
  go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" \
  -o /out/terraprism ./cmd/terraprism

FROM scratch AS export
COPY --from=build /out/terraprism /terraprism
