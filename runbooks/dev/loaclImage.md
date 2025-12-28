# Dev Runbook (Local image with k3s cluster)

### Set IMAGE_TAG

`export IMAGE_TAG="1.0.0"`

### Build image

Build a new docker image with below command

`docker build . -t cloudnessio/cloudness-dev:"${IMAGE_TAG}"`

### Publish it k3s cluster

Below command saves the local image and uploads

`docker save cloudnessio/cloudness-dev:${IMAGE_TAG} | ctr -a /run/k3s/containerd/containerd.sock -n=k8s.io images import -`


#### All in one

`export IMAGE_TAG="1.0.40" && docker build . -t cloudnessio/cloudness-dev:"${IMAGE_TAG}" && sudo docker save cloudnessio/cloudness-dev:${IMAGE_TAG} | sudo ctr -a /run/k3s/containerd/containerd.sock -n=k8s.io images import -`