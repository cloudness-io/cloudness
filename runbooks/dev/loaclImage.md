# Dev Runbook: Local Image with k3s Cluster

This guide explains how to build and deploy a local Docker image to a k3s cluster for development.

## Prerequisites

- Docker installed and running
- k3s cluster running locally
- `sudo` access (required for k3s containerd socket)

## Quick Start (All-in-One)

```bash
export IMAGE_TAG="1.0.0" && \
  docker build . -t cloudnessio/cloudness-dev:${IMAGE_TAG} && \
  sudo docker save cloudnessio/cloudness-dev:${IMAGE_TAG} | \
  sudo ctr -a /run/k3s/containerd/containerd.sock -n=k8s.io images import -
```

## Step-by-Step

### 1. Set the Image Tag

```bash
export IMAGE_TAG="1.0.0"
```

### 2. Build the Docker Image

```bash
docker build . -t cloudnessio/cloudness-dev:${IMAGE_TAG}
```

### 3. Import Image to k3s

k3s uses containerd instead of Docker, so we need to import the image directly:

```bash
sudo docker save cloudnessio/cloudness-dev:${IMAGE_TAG} | \
  sudo ctr -a /run/k3s/containerd/containerd.sock -n=k8s.io images import -
```

### 4. Verify the Image

```bash
sudo ctr -a /run/k3s/containerd/containerd.sock -n=k8s.io images list | grep cloudness
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Permission denied on socket | Ensure you're using `sudo` |
| Image not found in k3s | Verify the namespace is `k8s.io` |
| Build fails | Run from the repository root where `Dockerfile` exists |

## Notes

- The `k8s.io` namespace is where k3s stores images used by Kubernetes
- Images imported this way are only available on the local node
- For multi-node clusters, repeat the import on each node or use a registry
