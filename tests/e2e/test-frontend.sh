#!/bin/bash
# Test: Frontend Availability
# Tests frontend service and basic functionality

set -e

NAMESPACE="sentinel-mesh"
FAILED=0

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

echo "Testing Frontend Availability..."
echo ""

# Get frontend pod
FRONTEND_POD=$(kubectl get pods -n $NAMESPACE -l app=web-dashboard -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)

if [ -z "$FRONTEND_POD" ]; then
    echo -e "${RED}✗ No frontend pods found${NC}"
    exit 1
fi

echo "Using frontend pod: $FRONTEND_POD"
echo ""

# Test frontend is accessible via NodePort
echo "Frontend Accessibility:"
NODE_PORT=$(kubectl get svc web-dashboard -n $NAMESPACE -o jsonpath='{.spec.ports[0].nodePort}')
if [ -z "$NODE_PORT" ]; then
    echo -e "  ${RED}✗${NC} NodePort not found"
    ((FAILED++))
else
    FRONTEND_RESPONSE=$(curl -s --max-time 5 http://localhost:$NODE_PORT/ 2>/dev/null || echo "FAILED")

    if [[ "$FRONTEND_RESPONSE" == "FAILED" ]]; then
        echo -e "  ${RED}✗${NC} Failed to connect to frontend on port $NODE_PORT"
        ((FAILED++))
    elif echo "$FRONTEND_RESPONSE" | grep -qi "<!DOCTYPE html>"; then
        echo -e "  ${GREEN}✓${NC} Frontend returns HTML on port $NODE_PORT"
    else
        echo -e "  ${RED}✗${NC} Invalid HTML response"
        ((FAILED++))
    fi
fi
echo ""

# Check for Vue.js application
echo "Vue.js Application:"
if echo "$FRONTEND_RESPONSE" | grep -q "id=\"app\""; then
    echo -e "  ${GREEN}✓${NC} Vue.js app container found"
else
    echo -e "  ${RED}✗${NC} Vue.js app container not found"
    ((FAILED++))
fi
echo ""

# Check for static assets
echo "Static Assets:"
if echo "$FRONTEND_RESPONSE" | grep -qE "(\.js|\.css)"; then
    echo -e "  ${GREEN}✓${NC} JavaScript/CSS assets referenced"
else
    echo -e "  ${RED}✗${NC} No static assets found"
    ((FAILED++))
fi
echo ""

# Test NodePort accessibility
echo "NodePort Service:"
NODE_PORT=$(kubectl get svc web-dashboard -n $NAMESPACE -o jsonpath='{.spec.ports[0].nodePort}')
if [ -n "$NODE_PORT" ]; then
    echo -e "  ${GREEN}✓${NC} NodePort configured: $NODE_PORT"
    echo "  Access URL: http://localhost:$NODE_PORT"
else
    echo -e "  ${RED}✗${NC} NodePort not configured"
    ((FAILED++))
fi
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All frontend tests passed${NC}"
    exit 0
else
    echo -e "${RED}✗ $FAILED frontend test(s) failed${NC}"
    exit 1
fi
