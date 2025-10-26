#!/bin/bash
# Test: ML Service Endpoints
# Tests ML service anomaly detection and predictions APIs

set -e

NAMESPACE="sentinel-mesh"
FAILED=0

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

echo "Testing ML Service Endpoints..."
echo ""

# Get ML service pod
ML_POD=$(kubectl get pods -n $NAMESPACE -l app=ml-service -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)

if [ -z "$ML_POD" ]; then
    echo -e "${RED}✗ No ML service pods found${NC}"
    exit 1
fi

echo "Using ML pod: $ML_POD"
echo ""

# Test anomaly detection endpoint
echo "Anomaly Detection Endpoint (/api/v1/anomalies):"
ANOMALIES_RESPONSE=$(kubectl exec -n $NAMESPACE $ML_POD -- wget -q -O- --timeout=5 http://localhost:8000/api/v1/anomalies 2>/dev/null || echo "FAILED")

if [[ "$ANOMALIES_RESPONSE" == "FAILED" ]]; then
    echo -e "  ${RED}✗${NC} Failed to connect"
    ((FAILED++))
elif echo "$ANOMALIES_RESPONSE" | grep -q "anomalies"; then
    echo -e "  ${GREEN}✓${NC} Returns anomalies data"

    # Check for anomaly structure
    if echo "$ANOMALIES_RESPONSE" | grep -q "severity"; then
        echo -e "  ${GREEN}✓${NC} Anomalies include severity"
    else
        echo -e "  ${RED}✗${NC} Anomalies missing severity field"
        ((FAILED++))
    fi

    if echo "$ANOMALIES_RESPONSE" | grep -q "resource"; then
        echo -e "  ${GREEN}✓${NC} Anomalies include resource information"
    else
        echo -e "  ${RED}✗${NC} Anomalies missing resource field"
        ((FAILED++))
    fi

    # Count anomalies
    ANOMALY_COUNT=$(echo "$ANOMALIES_RESPONSE" | grep -o "type" | wc -l | tr -d ' ')
    echo "  Found $ANOMALY_COUNT anomalies"
else
    echo -e "  ${RED}✗${NC} Invalid response format"
    ((FAILED++))
fi
echo ""

# Test predictions endpoint
echo "Predictions Endpoint (/api/v1/predictions):"
PREDICTIONS_RESPONSE=$(kubectl exec -n $NAMESPACE $ML_POD -- wget -q -O- --timeout=5 http://localhost:8000/api/v1/predictions 2>/dev/null || echo "FAILED")

if [[ "$PREDICTIONS_RESPONSE" == "FAILED" ]]; then
    echo -e "  ${RED}✗${NC} Failed to connect"
    ((FAILED++))
elif echo "$PREDICTIONS_RESPONSE" | grep -q "predictions"; then
    echo -e "  ${GREEN}✓${NC} Returns predictions data"

    # Check for prediction fields
    if echo "$PREDICTIONS_RESPONSE" | grep -q "cpu_usage"; then
        echo -e "  ${GREEN}✓${NC} Predictions include CPU usage"
    else
        echo -e "  ${RED}✗${NC} Predictions missing CPU usage"
        ((FAILED++))
    fi

    if echo "$PREDICTIONS_RESPONSE" | grep -q "memory_usage"; then
        echo -e "  ${GREEN}✓${NC} Predictions include memory usage"
    else
        echo -e "  ${RED}✗${NC} Predictions missing memory usage"
        ((FAILED++))
    fi

    if echo "$PREDICTIONS_RESPONSE" | grep -q "next_hour"; then
        echo -e "  ${GREEN}✓${NC} Predictions include next_hour forecast"
    else
        echo -e "  ${RED}✗${NC} Predictions missing next_hour forecast"
        ((FAILED++))
    fi
else
    echo -e "  ${RED}✗${NC} Invalid response format"
    ((FAILED++))
fi
echo ""

# Test health endpoint
echo "Health Endpoint:"
HEALTH_RESPONSE=$(kubectl exec -n $NAMESPACE $ML_POD -- wget -q -O- --timeout=5 http://localhost:8000/health 2>/dev/null || echo "FAILED")

if [[ "$HEALTH_RESPONSE" == "FAILED" ]]; then
    echo -e "  ${RED}✗${NC} Failed to connect"
    ((FAILED++))
elif echo "$HEALTH_RESPONSE" | grep -q "healthy"; then
    echo -e "  ${GREEN}✓${NC} ML service is healthy"
else
    echo -e "  ${RED}✗${NC} Health check failed"
    ((FAILED++))
fi
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All ML service tests passed${NC}"
    exit 0
else
    echo -e "${RED}✗ $FAILED ML service test(s) failed${NC}"
    exit 1
fi
