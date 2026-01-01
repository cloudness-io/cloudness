# Versions
# https://hub.docker.com/_/alpine
ARG BASE_IMAGE=moby/buildkit:rootless
# https://github.com/railwayapp/nixpacks/releases
ARG NIXPACKS_VERSION=1.33.0
# https://dl.k8s.io/release/stable.txt
ARG KUBECTL_VERSION=latest

ARG TARGETPLATFORM
ARG NIXPACKS_VERSION
ARG KUBECTL_VERSION

# Build the deployer binary (has its own go.mod)
FROM golang:1.25-alpine AS deployer-builder
WORKDIR /build
COPY ./deployer/go.mod ./deployer/go.sum ./
RUN go mod download
COPY ./deployer/*.go ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o ./deployer .

FROM bitnami/kubectl:${KUBECTL_VERSION} AS kubectl
FROM ${BASE_IMAGE} AS base


USER root
RUN apk update && apk add --no-cache \
    bash curl git git-lfs openssh-client tar tini musl ca-certificates openssh-client yq \
    && rm -rf /var/cache/apk/*

RUN curl -sSL https://nixpacks.com/install.sh | bash

COPY --from=kubectl /opt/bitnami/kubectl/bin/kubectl /usr/local/bin/kubectl
RUN chmod +x /usr/local/bin/kubectl

# Add deployer binary
COPY --from=deployer-builder /build/deployer /usr/local/bin/cloudness-deploy
RUN chmod +x /usr/local/bin/cloudness-deploy

# Add shared utilities and scripts
COPY scripts/cloudness-utils.sh /usr/local/lib/cloudness-utils.sh
COPY scripts/init-script.sh /usr/local/lib/init-script.sh
COPY scripts/build-script.sh /usr/local/lib/build-script.sh
COPY scripts/deploy-script.sh /usr/local/lib/deploy-script.sh
RUN chmod +x /usr/local/lib/cloudness-utils.sh /usr/local/lib/init-script.sh /usr/local/lib/build-script.sh /usr/local/lib/deploy-script.sh

# Auto-source utilities for all sh scripts
ENV ENV="/usr/local/lib/cloudness-utils.sh"
ENV PATH="/usr/local/bin:${PATH}"
USER user 
ENTRYPOINT ["/sbin/tini", "--"]