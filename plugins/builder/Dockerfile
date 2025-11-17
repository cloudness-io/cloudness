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

FROM bitnami/kubectl:${KUBECTL_VERSION} AS kubectl
FROM ${BASE_IMAGE} AS base


USER root
RUN apk update && apk add --no-cache \
    bash curl git git-lfs openssh-client tar tini musl ca-certificates openssh-client yq \
    && rm -rf /var/cache/apk/*

RUN curl -sSL https://nixpacks.com/install.sh | bash

COPY --from=kubectl /opt/bitnami/kubectl/bin/kubectl /usr/local/bin/kubectl
RUN chmod +x /usr/local/bin/kubectl

ENV PATH="/usr/local/bin:${PATH}"
USER user 
ENTRYPOINT ["/sbin/tini", "--"]