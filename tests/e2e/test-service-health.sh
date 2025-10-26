#!/bin/bash
# Test: Service Health & Availability
# Tests all deployed services for health, readiness, and metrics endpoints

set -e

NAMESPACE="sentinel-mesh"
FAILED=0

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "Testing Service Health & Availability..."
echo ""

# Helper function to test service endpoint
test_endpoint() {
    local service_name="$1"
    local endpoint="$2"
    local expected_code="${3:-200}"

    # Get service port
    local pod=$(kubectl get pods -n $NAMESPACE -l app=$service_name -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)

    if [ -z "$pod" ]; then
        echo -e "  ${RED}✗${NC} No pods found for service: $service_name"
        return 1
    fi

    # Test endpoint
    local response=$(kubectl exec -n $NAMESPACE $pod -- wget -q -O- --timeout=5 http://localhost:8080$endpoint 2>/dev/null || echo "FAILED")

    if [[ "$response" == "FAILED" ]] || [ -z "$response" ]; then
        echo -e "  ${RED}✗${NC} $service_name $endpoint - Failed to connect"
        return 1
    else
        echo -e "  ${GREEN}✓${NC} $service_name $endpoint - OK"
        return 0
    fi
}

# Test API service
echo "API Service:"
test_endpoint "api" "/health" || ((FAILED++))
if ! test_endpoint "api" "/ready"; then
    echo -e "  ${YELLOW}⚠${NC} /ready endpoint not implemented (optional)"
fi
test_endpoint "api" "/metrics" || ((FAILED++))
echo ""

# Test Collector service
echo "Collector Service:"
test_endpoint "collector" "/health" || ((FAILED++))
if ! test_endpoint "collector" "/ready"; then
    echo -e "  ${YELLOW}⚠${NC} /ready endpoint not implemented (optional)"
fi
test_endpoint "collector" "/metrics" || ((FAILED++))
echo ""

# Test ML service (port 8000)
echo "ML Service:"
ml_pod=$(kubectl get pods -n $NAMESPACE -l app=ml-service -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
if [ -n "$ml_pod" ]; then
    if kubectl exec -n $NAMESPACE $ml_pod -- python3 -c "import urllib.request; urllib.request.urlopen('http://localhost:8000/health')" &>/dev/null; then
        echo -e "  ${GREEN}✓${NC} ml-service /health - OK"
    else
        echo -e "  ${RED}✗${NC} ml-service /health - Failed"
        ((FAILED++))
    fi

    if kubectl exec -n $NAMESPACE $ml_pod -- python3 -c "import urllib.request; urllib.request.urlopen('http://localhost:8000/ready')" &>/dev/null; then
        echo -e "  ${GREEN}✓${NC} ml-service /ready - OK"
    else
        echo -e "  ${YELLOW}⚠${NC} /ready endpoint not implemented (optional)"
    fi

    if kubectl exec -n $NAMESPACE $ml_pod -- python3 -c "import urllib.request; urllib.request.urlopen('http://localhost:8000/metrics')" &>/dev/null; then
        echo -e "  ${GREEN}✓${NC} ml-service /metrics - OK"
    else
        echo -e "  ${RED}✗${NC} ml-service /metrics - Failed"
        ((FAILED++))
    fi
else
    echo -e "  ${YELLOW}⊘${NC} ML service not deployed"
fi
echo ""

# Check pod status
echo "Pod Status:"
all_running=true
while read -r pod status; do
    if [ "$status" == "Running" ]; then
        echo -e "  ${GREEN}✓${NC} $pod - Running"
    else
        echo -e "  ${RED}✗${NC} $pod - $status"
        all_running=false
        ((FAILED++))
    fi
done < <(kubectl get pods -n $NAMESPACE --no-headers -o custom-columns=":metadata.name,:status.phase")
echo ""

# Check if all pods are ready
echo "Pod Readiness:"
while read -r pod; do
    ready_count=$(kubectl get pod $pod -n $NAMESPACE -o jsonpath='{.status.containerStatuses[*].ready}' | tr ' ' '\n' | grep -c "true" || echo "0")
    total_count=$(kubectl get pod $pod -n $NAMESPACE -o jsonpath='{.spec.containers[*].name}' | wc -w | tr -d ' ')

    if [ "$ready_count" == "$total_count" ]; then
        echo -e "  ${GREEN}✓${NC} $pod - Ready ($ready_count/$total_count)"
    else
        echo -e "  ${RED}✗${NC} $pod - Not Ready ($ready_count/$total_count)"
        ((FAILED++))
    fi
done < <(kubectl get pods -n $NAMESPACE --no-headers -o custom-columns=":metadata.name")
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All service health checks passed${NC}"
    exit 0
else
    echo -e "${RED}✗ $FAILED health check(s) failed${NC}"
    exit 1
fi
