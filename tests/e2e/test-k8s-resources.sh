#!/bin/bash
# Test: Kubernetes Resources
# Tests Kubernetes deployment, services, and resource configuration

set -e

NAMESPACE="sentinel-mesh"
FAILED=0

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "Testing Kubernetes Resources..."
echo ""

# Test namespace exists
echo "Namespace:"
if kubectl get namespace $NAMESPACE &>/dev/null; then
    echo -e "  ${GREEN}✓${NC} Namespace '$NAMESPACE' exists"
else
    echo -e "  ${RED}✗${NC} Namespace '$NAMESPACE' not found"
    ((FAILED++))
fi
echo ""

# Test deployments
echo "Deployments:"
EXPECTED_DEPLOYMENTS=("api" "collector" "ml-service" "web-dashboard")
for deployment in "${EXPECTED_DEPLOYMENTS[@]}"; do
    if kubectl get deployment $deployment -n $NAMESPACE &>/dev/null; then
        # Check replicas
        DESIRED=$(kubectl get deployment $deployment -n $NAMESPACE -o jsonpath='{.spec.replicas}')
        READY=$(kubectl get deployment $deployment -n $NAMESPACE -o jsonpath='{.status.readyReplicas}')

        if [ "$DESIRED" == "$READY" ]; then
            echo -e "  ${GREEN}✓${NC} $deployment - Ready ($READY/$DESIRED replicas)"
        else
            echo -e "  ${RED}✗${NC} $deployment - Not Ready ($READY/$DESIRED replicas)"
            ((FAILED++))
        fi
    else
        echo -e "  ${YELLOW}⊘${NC} $deployment - Not deployed"
    fi
done
echo ""

# Test services
echo "Services:"
EXPECTED_SERVICES=("api" "collector" "ml-service" "web-dashboard")
for service in "${EXPECTED_SERVICES[@]}"; do
    if kubectl get service $service -n $NAMESPACE &>/dev/null; then
        TYPE=$(kubectl get service $service -n $NAMESPACE -o jsonpath='{.spec.type}')
        CLUSTER_IP=$(kubectl get service $service -n $NAMESPACE -o jsonpath='{.spec.clusterIP}')

        echo -e "  ${GREEN}✓${NC} $service - Type: $TYPE, ClusterIP: $CLUSTER_IP"
    else
        echo -e "  ${YELLOW}⊘${NC} $service - Not found"
    fi
done
echo ""

# Test pod resource limits
echo "Resource Limits:"
PODS=$(kubectl get pods -n $NAMESPACE --no-headers -o custom-columns=":metadata.name")
HAS_LIMITS=false

for pod in $PODS; do
    LIMITS=$(kubectl get pod $pod -n $NAMESPACE -o jsonpath='{.spec.containers[0].resources.limits}' 2>/dev/null)

    if [ -n "$LIMITS" ] && [ "$LIMITS" != "{}" ]; then
        HAS_LIMITS=true
        CPU_LIMIT=$(kubectl get pod $pod -n $NAMESPACE -o jsonpath='{.spec.containers[0].resources.limits.cpu}' 2>/dev/null)
        MEM_LIMIT=$(kubectl get pod $pod -n $NAMESPACE -o jsonpath='{.spec.containers[0].resources.limits.memory}' 2>/dev/null)

        if [ -n "$CPU_LIMIT" ] || [ -n "$MEM_LIMIT" ]; then
            echo -e "  ${GREEN}✓${NC} $pod - Limits: CPU=$CPU_LIMIT, Memory=$MEM_LIMIT"
        fi
    fi
done

if [ "$HAS_LIMITS" = false ]; then
    echo -e "  ${YELLOW}⊘${NC} No resource limits configured (recommended for production)"
fi
echo ""

# Test pod restart counts
echo "Pod Stability (Restart Counts):"
while read -r pod restarts; do
    if [ "$restarts" -eq 0 ]; then
        echo -e "  ${GREEN}✓${NC} $pod - No restarts"
    elif [ "$restarts" -lt 3 ]; then
        echo -e "  ${YELLOW}⚠${NC} $pod - $restarts restart(s)"
    else
        echo -e "  ${RED}✗${NC} $pod - $restarts restart(s) (excessive)"
        ((FAILED++))
    fi
done < <(kubectl get pods -n $NAMESPACE --no-headers -o custom-columns=":metadata.name,:status.containerStatuses[*].restartCount")
echo ""

# Test service endpoints
echo "Service Endpoints:"
for service in "${EXPECTED_SERVICES[@]}"; do
    if kubectl get service $service -n $NAMESPACE &>/dev/null; then
        ENDPOINTS=$(kubectl get endpoints $service -n $NAMESPACE -o jsonpath='{.subsets[*].addresses[*].ip}' 2>/dev/null | wc -w | tr -d ' ')

        if [ "$ENDPOINTS" -gt 0 ]; then
            echo -e "  ${GREEN}✓${NC} $service - $ENDPOINTS endpoint(s)"
        else
            echo -e "  ${RED}✗${NC} $service - No endpoints (no ready pods)"
            ((FAILED++))
        fi
    fi
done
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All Kubernetes resource tests passed${NC}"
    exit 0
else
    echo -e "${RED}✗ $FAILED Kubernetes test(s) failed${NC}"
    exit 1
fi
