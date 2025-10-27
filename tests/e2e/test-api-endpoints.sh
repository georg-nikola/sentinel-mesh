#!/bin/bash
# Test: API Endpoints
# Tests API service endpoints and responses

set -e

NAMESPACE="sentinel-mesh"
FAILED=0

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

echo "Testing API Endpoints..."
echo ""

# Get API pod
API_POD=$(kubectl get pods -n $NAMESPACE -l app=api -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)

if [ -z "$API_POD" ]; then
    echo -e "${RED}✗ No API pods found${NC}"
    exit 1
fi

echo "Using API pod: $API_POD"
echo ""

# Test health endpoint returns JSON
echo "Health Endpoint:"
HEALTH_RESPONSE=$(kubectl exec -n $NAMESPACE $API_POD -- wget -q -O- http://localhost:8080/health 2>/dev/null || echo "FAILED")

if [[ "$HEALTH_RESPONSE" == "FAILED" ]]; then
    echo -e "  ${RED}✗${NC} Failed to connect"
    ((FAILED++))
elif echo "$HEALTH_RESPONSE" | grep -q "healthy"; then
    echo -e "  ${GREEN}✓${NC} Returns healthy status"
    echo "  Response: $HEALTH_RESPONSE"
else
    echo -e "  ${RED}✗${NC} Invalid response"
    ((FAILED++))
fi
echo ""

# Test ready endpoint (optional - not implemented in current version)
echo "Readiness Endpoint:"
READY_RESPONSE=$(kubectl exec -n $NAMESPACE $API_POD -- wget -q -O- http://localhost:8080/ready 2>/dev/null || echo "FAILED")

if [[ "$READY_RESPONSE" == "FAILED" ]]; then
    echo -e "  ${YELLOW}⚠${NC} Not implemented (optional endpoint)"
elif echo "$READY_RESPONSE" | grep -q "ready"; then
    echo -e "  ${GREEN}✓${NC} Returns ready status"
    echo "  Response: $READY_RESPONSE"
else
    echo -e "  ${YELLOW}⚠${NC} Unexpected response format"
fi
echo ""

# Test metrics endpoint
echo "Metrics Endpoint:"
METRICS_RESPONSE=$(kubectl exec -n $NAMESPACE $API_POD -- wget -q -O- http://localhost:8080/metrics 2>/dev/null || echo "FAILED")

if [[ "$METRICS_RESPONSE" == "FAILED" ]]; then
    echo -e "  ${RED}✗${NC} Failed to connect"
    ((FAILED++))
elif echo "$METRICS_RESPONSE" | grep -qE "^# (TYPE|HELP)"; then
    echo -e "  ${GREEN}✓${NC} Returns Prometheus metrics"
    METRIC_COUNT=$(echo "$METRICS_RESPONSE" | grep -cE "^# TYPE" || true)
    echo "  Metrics found: $METRIC_COUNT metric types"
else
    echo -e "  ${RED}✗${NC} Invalid Prometheus format"
    ((FAILED++))
fi
echo ""

# Test version information
echo "Version Information:"
if echo "$HEALTH_RESPONSE" | grep -q "version"; then
    VERSION=$(echo "$HEALTH_RESPONSE" | grep -o '"version":"[^"]*"' | cut -d'"' -f4)
    echo -e "  ${GREEN}✓${NC} Version present: $VERSION"
else
    echo -e "  ${RED}✗${NC} Version not found in health response"
    ((FAILED++))
fi
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All API endpoint tests passed${NC}"
    exit 0
else
    echo -e "${RED}✗ $FAILED API test(s) failed${NC}"
    exit 1
fi
