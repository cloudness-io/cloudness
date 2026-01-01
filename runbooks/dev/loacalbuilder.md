## Quick Start (All-in-One)

```bash
# Run from plugins directory
cd plugins/helper
export IMAGE_TAG="1.0.0" && \
  docker build -f Dockerfile -t cloudnessio/helper:${IMAGE_TAG} . && \
  sudo docker save cloudnessio/helper:${IMAGE_TAG} | \
  sudo ctr -a /run/k3s/containerd/containerd.sock -n=k8s.io images import -
```
