# Input Variables
DEPLOYMENT_NAME="{{ .Identifier }}"
NAMESPACE="{{ .Namespace }}"
MANIFEST_PATH="{{ .ManifestPath }}"
HAS_HTTPROUTE="{{ .HasHTTPRoute }}"

# set -x  #For debugging

# Helper function for logging
log() { echo -e "\033[1;32m[INFO]\033[0m $*"; }

# 1- Starting deployment
log "Starting deployment of $DEPLOYMENT_NAME in namespace $NAMESPACE."

# 2- Provisioning application
log "Applying Kubernetes manifests from $MANIFEST_PATH..."
kubectl apply -n "$NAMESPACE" -f "$MANIFEST_PATH"
if [[ $? -ne 0 ]]; then
  log "ERROR" "kubectl apply failed. Exiting..."
  exit 1
fi

# 4- Initiate rolling update
log "Initiating rolling update..."
kubectl rollout restart statefulset/$DEPLOYMENT_NAME -n "$NAMESPACE"

# 5- Check rollout status
log "Checking rollout status..."
if ! kubectl rollout status statefulset/$DEPLOYMENT_NAME -n "$NAMESPACE" --timeout=2m; then # Add a timeout for the rollout itself
  log "ERROR" "StatefulSet rollout did not succeed within rollout time. Rolling back..."
  kubectl rollout undo statefulset/$DEPLOYMENT_NAME -n "$NAMESPACE"
  exit 1
fi

# 6- Verify health check
log "Verifying health check..."
sleep 5
if ! kubectl get pods -n "$NAMESPACE" | grep "$DEPLOYMENT_NAME" | grep -q "Running"; then
  log "ERROR" "Health check failed! Rolling back..."
  kubectl rollout undo statefulset/$DEPLOYMENT_NAME -n "$NAMESPACE"
  exit 1
fi

# 7- Cleaning up Deployments if any
kubectl delete deployment/$DEPLOYMENT_NAME -n "$NAMESPACE" --ignore-not-found=true 

# 8- Remove HTTPRoute if not present
if [ "$HAS_HTTPROUTE" = "false" ]; then
kubectl delete httproute -n "$NAMESPACE" -l app.kubernetes.io/instance="$DEPLOYMENT_NAME"
fi

# 9- Deployment success
log "SUCCESS" "Deployment of $DEPLOYMENT_NAME completed successfully."
exit 0
