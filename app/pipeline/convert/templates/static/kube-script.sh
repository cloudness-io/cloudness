: "${CLOUDNESS_DEPLOY_APP_IDENTIFIER:=}"
: "${CLOUDNESS_DEPLOY_APP_NAMESPACE:=}"

# Deployment flags
: "${CLOUDNESS_DEPLOY_FLAG_APP_TYPE:=}"
: "${CLOUDNESS_DEPLOY_FLAG_HAS_VOLUME:=0}"
: "${CLOUDNESS_DEPLOY_FLAG_NEED_REMOUNT:=0}"
: "${CLOUDNESS_DEPLOY_FLAG_HAS_ROUTE:=0}"

# Deployment yaml files
: "${CLOUDNESS_DEPLOY_YAML_COMMON:=}"
: "${CLOUDNESS_DEPLOY_YAML_VOLUME:=}"
: "${CLOUDNESS_DEPLOY_YAML_APP:=}"
: "${CLOUDNESS_DEPLOY_YAML_ROUTE:=}"

# Define color variables
RED='\033[1;31m'    # Bold Red
YELLOW='\033[1;33m' # Bold Yellow
GREEN='\033[1;32m'  # Bold Green
RESET='\033[0m'     # Resets color
CHECK_MARK='\u2714' # Unicode for Heavy Check Mark (✔)

# Helper function for logging
error() { echo -e "${RED}[ERROR]${RESET}  $*"; }
warn() { echo -e "${YELLOW}[WARN]${RESET}  $*"; }
info() { echo -e "$*"; }
success() { echo -e "${GREEN}[SUCCESS]${RESET}  $*"; }
success_info() { echo -e "$* ${GREEN}✔${RESET}"; }

apply_kube_config_from_string() {
	local KUBE_YAML_STRING="$1"
	local ERROR_MESSAGE=""

	# Check if the YAML string is empty
	if [ -z "$KUBE_YAML_STRING" ]; then
		return 0
	fi

	# echo "Attempting to apply Kubernetes configuration..."
	# echo -e "$KUBE_YAML_STRING" | sed 's/\t/  /g'

	# Apply the YAML using kubectl, capturing stderr into ERROR_MESSAGE
	# And redirecting stdout to /dev/null to suppress success messages
	if ! ERROR_MESSAGE=$(echo -e "$KUBE_YAML_STRING" | sed 's/\t/  /g' | kubectl apply -f - 2>&1 >/dev/null); then
		error "Error applying Kubernetes configuration:"
		error "$ERROR_MESSAGE"
		return 1 # Indicate failure
	else
		return 0 # Indicate success
	fi
}

rollout_status() {
	local ROLLOUT_ERROR=""

	if [ "$CLOUDNESS_DEPLOY_FLAG_APP_TYPE" = "Stateless" ]; then
		if ! ROLLOUT_ERROR=$(kubectl rollout status deployment/"$CLOUDNESS_DEPLOY_APP_IDENTIFIER" -n "$CLOUDNESS_DEPLOY_APP_NAMESPACE" --timeout=20s >&1 >/dev/null); then
			error "$ROLLOUT_ERROR"
			error "Error rolling out deployment, reverting..."
			kubectl rollout undo deployment/"$CLOUDNESS_DEPLOY_APP_IDENTIFIER" -n "$CLOUDNESS_DEPLOY_APP_NAMESPACE"
			return 1
		fi
	else
		if ! ROLLOUT_ERROR=$(kubectl rollout status statefulset/"$CLOUDNESS_DEPLOY_APP_IDENTIFIER" -n "$CLOUDNESS_DEPLOY_APP_NAMESPACE" --timeout=2m >&1 >/dev/null); then
			error "$ROLLOUT_ERROR"
			error "Error rolling out deployment, reverting..."
			kubectl rollout undo statefulset/"$CLOUDNESS_DEPLOY_APP_IDENTIFIER" -n "$CLOUDNESS_DEPLOY_APP_NAMESPACE"
			return 1
		fi
	fi
}

cleanup() {
	info "Running clean up..."
	local CLEANUP_ERROR=""

	if [ "$CLOUDNESS_DEPLOY_FLAG_APP_TYPE" = "Stateless" ]; then
		if ! CLEANUP_ERROR=$(kubectl delete statefulset/"$CLOUDNESS_DEPLOY_APP_IDENTIFIER" -n "$CLOUDNESS_DEPLOY_APP_NAMESPACE" --ignore-not-found=true 2>&1 >/dev/null); then
			error "$CLEANUP_ERROR"
			warn "Error cleaning up deployment, Skipping..."
		fi
	else
		if ! CLEANUP_ERROR=$(kubectl delete deployment/"$CLOUDNESS_DEPLOY_APP_IDENTIFIER" -n "$CLOUDNESS_DEPLOY_APP_NAMESPACE" --ignore-not-found=true 2>&1 >/dev/null); then
			error "$CLEANUP_ERROR"
			warn "Error cleaning up deployment, Skipping..."
		fi
	fi
}

wait_for_pvc_resize() {
	PVC_NAME="$1"
	NEW_SIZE="$2"
	NAMESPACE="$CLOUDNESS_DEPLOY_APP_NAMESPACE"

	TIMEOUT_SECONDS="${TIMEOUT_SECONDS:-300}" # 5 minutes
	SLEEP_SECONDS="${SLEEP_SECONDS:-5}"

	local current_time=$(date +%s)
	local deadline=$(expr "$current_time" + "$TIMEOUT_SECONDS")

	check_resize_pending() {
		kubectl get pvc "$PVC_NAME" -n "$NAMESPACE" -o jsonpath='{.status.conditions[?(@.type=="FileSystemResizePending")].status}' 2>/dev/null
	}

	# Function to parse a storage string (e.g., "11Gi") and convert to GiB
	parse_size_to_gib() {
		local size_str=$1
		local value=$(echo "$size_str" | sed 's/Gi//')
		# Use 'bc' for floating-point arithmetic if necessary
		echo "$value"
	}

	check_pvc_status() {
		kubectl get pvc "$PVC_NAME" -n "$NAMESPACE" -o jsonpath='{.status.phase}' 2>/dev/null
	}

	get_pvc_event() {
		kubectl get events -n "$NAMESPACE" --field-selector involvedObject.kind=PersistentVolumeClaim,involvedObject.name="$PVC_NAME" --sort-by=.lastTimestamp -o jsonpath='{.items[-1:].reason}' 2>/dev/null
	}

	while true; do
		PVC_STATUS=$(check_pvc_status)
		# Handle WaitForFirstConsumer specifically
		if [ "$PVC_STATUS" == "Pending" ]; then
			PVC_EVENT=$(get_pvc_event)
			if [ "$PVC_EVENT" == "WaitForFirstConsumer" ]; then
				return 0
			fi
		fi

		CURRENT_SIZE=$(kubectl get pvc $PVC_NAME -n $NAMESPACE -o jsonpath='{.status.capacity.storage}')
		if [ $(parse_size_to_gib "$CURRENT_SIZE") -ge $(parse_size_to_gib "$NEW_SIZE") ]; then
			return 0 # Exit the function with success status
		fi

		# Check for FileSystemResizePending condition
		if [ "$(check_resize_pending)" == "True" ]; then
			info "Volume has been resized. Remounting application to finalize resize."
			return 0
		fi

		# Timeout
		current_time=$(date +%s)
		if [ "$current_time" -ge "$deadline" ]; then
			info "⏱️  Timed out after ${TIMEOUT_SECONDS}s waiting for PVC '$PVC_NAME' to reach $NEW_SIZE."
			return 1
		fi

		echo "Waiting for Volume..."
		sleep "$SLEEP_SECONDS"
	done
}

################################################ MAIN #################################

# Applying common artifacts namespace, service account, role, rolebinding, configmap, secrets...
if apply_kube_config_from_string "$CLOUDNESS_DEPLOY_YAML_COMMON"; then
	success_info "Setting up prequestic artifacts."
else
	error "Error setting up prequestic artifacts. Exiting..."
	return 1
fi

# Apply volume/pvc configuration
if [ "$CLOUDNESS_DEPLOY_FLAG_HAS_VOLUME" -eq 1 ]; then
	if [ "$CLOUDNESS_DEPLOY_FLAG_NEED_REMOUNT" -eq 1 ]; then #if the volume is scaled up and needs remount for changes to apply
		#delete the statefulset before scaling up the volume
		if !(kubectl delete statefulset/"$CLOUDNESS_DEPLOY_APP_IDENTIFIER" -n "$CLOUDNESS_DEPLOY_APP_NAMESPACE" --ignore-not-found=true 2>&1 >/dev/null); then
			error "Error cleaning up deployment, Skipping..."
			return 1
		fi
	fi
	if apply_kube_config_from_string "$CLOUDNESS_DEPLOY_YAML_VOLUME"; then
		# Iterate through each PVC found
		PVC_DATA=$(echo -e "$CLOUDNESS_DEPLOY_YAML_VOLUME" | sed 's/\t/  /g' | yq -r 'select(.kind == "PersistentVolumeClaim") | .metadata.name + " " + .spec.resources.requests.storage' -)
		echo "$PVC_DATA" | while read -r PVC_NAME NEW_SIZE; do
			PVC_NAME=$(echo "$PVC_NAME" | xargs)
			NEW_SIZE=$(echo "$NEW_SIZE" | xargs)
			if [ -z "$PVC_NAME" ] || [ -z "$NEW_SIZE" ]; then #discard ---- and empty lines
				continue
			fi

			if ! wait_for_pvc_resize "$PVC_NAME" "$NEW_SIZE"; then
				error "Failed to resize or confirm PVC '$PVC_NAME'. Check logs above for details."
				return 1
			fi
		done

		success_info "Volume provisioned."
	else
		"Error setting up volume. Exiting..."
		return 1
	fi

fi

if ! apply_kube_config_from_string "$CLOUDNESS_DEPLOY_YAML_APP"; then
	error "Error deploying application. Exiting..."
	return 1
fi

if rollout_status; then
	success_info "Application deployment."
else
	#logs are handles in rollout_status method
	return 1
fi

if [ "$CLOUDNESS_DEPLOY_FLAG_HAS_ROUTE" -eq 1 ]; then
	if apply_kube_config_from_string "$CLOUDNESS_DEPLOY_YAML_ROUTE"; then
		success_info "HTTP routes configured."
	else
		error "Error setting up http routes. Exiting..."
		return 1
	fi
fi

cleanup

success "Deployment completed successfully."
