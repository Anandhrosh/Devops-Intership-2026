
#!/usr/bin/env bash
set -euo pipefail

NAMESPACE="default"
LOAD_POD="sprint6-http-load"

cleanup() {
  echo "Stopping HTTP load generator..."
  kubectl delete pod "$LOAD_POD" \
    -n "$NAMESPACE" --ignore-not-found
}

trap cleanup EXIT INT TERM

echo "Creating HTTP load generator..."

kubectl run "$LOAD_POD" \
  --namespace="$NAMESPACE" \
  --image=busybox:1.36 \
  --restart=Never \
  --command -- sh -c '
    echo "Sending repeated requests to sprint6-service"
    while true; do
      wget -q -O /dev/null http://sprint6-service/
    done
  '

echo "Waiting for load pod..."
kubectl wait \
  --for=condition=Ready \
  "pod/$LOAD_POD" \
  -n "$NAMESPACE" \
  --timeout=120s

echo "HTTP load is running."
echo "Press Ctrl+C to stop the load generator."

while true; do
  sleep 5
done
