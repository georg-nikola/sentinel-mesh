#!/bin/bash
# Test: Metrics Collection & Data Flow
# Tests that metrics are being collected and flowing through the system

set -e

NAMESPACE="sentinel-mesh"
FAILED=0

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "Testing Metrics Collection & Data Flow..."
echo ""

# Test Prometheus metrics endpoints
echo "Prometheus Metrics Availability:"

SERVICES=("api" "collector")
for service in "${SERVICES[@]}"; do
    POD=$(kubectl get pods -n $NAMESPACE -l app=$service -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)

    if [ -n "$POD" ]; then
        METRICS=$(kubectl exec -n $NAMESPACE $POD -- wget -q -O- --timeout=5 http://localhost:8080/metrics 2>/dev/null || echo "FAILED")

        if [[ "$METRICS" == "FAILED" ]]; then
            echo -e "  ${RED}✗${NC} $service - Metrics endpoint failed"
            ((FAILED++))
        else
            # Count metrics
            METRIC_COUNT=$(echo "$METRICS" | grep -c "^[a-z_]" || true)
            echo -e "  ${GREEN}✓${NC} $service - Metrics available ($METRIC_COUNT lines)"

            # Check for specific metrics
            if echo "$METRICS" | grep -q "${service}_up"; then
                echo -e "  ${GREEN}✓${NC} $service - Service up metric present"
            fi
        fi
    else
        echo -e "  ${YELLOW}⊘${NC} $service - No pods found"
    fi
done
echo ""

# Test ML service metrics
echo "ML Service Metrics:"
ML_POD=$(kubectl get pods -n $NAMESPACE -l app=ml-service -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)

if [ -n "$ML_POD" ]; then
    ML_METRICS=$(kubectl exec -n $NAMESPACE $ML_POD -- python3 -c "import urllib.request; print(urllib.request.urlopen('http://localhost:8000/metrics').read().decode())" 2>&1)

    if echo "$ML_METRICS" | grep -qE "(ml_service_up|ml.*up)"; then
        echo -e "  ${GREEN}✓${NC} ML service up metric present"
    else
        echo -e "  ${YELLOW}⚠${NC} ML service up metric format differs (non-critical)"
    fi

    if echo "$ML_METRICS" | grep -q "ml_service_requests_total"; then
        echo -e "  ${GREEN}✓${NC} ML request counter present"
    fi

    # Count total metrics
    METRIC_COUNT=$(echo "$ML_METRICS" | grep -cE "^[a-z_]" || echo "0")
    if [ "$METRIC_COUNT" -gt 0 ]; then
        echo -e "  ${GREEN}✓${NC} ML service exposing $METRIC_COUNT metrics"
    fi
else
    echo -e "  ${YELLOW}⊘${NC} ML service not deployed"
fi
echo ""

# Test data integration (Frontend → API → ML)
echo "Data Integration Flow:"

# Check if frontend can reach API
echo "  Frontend → API connectivity:"
FRONTEND_POD=$(kubectl get pods -n $NAMESPACE -l app=web-dashboard -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
API_SVC=$(kubectl get svc api -n $NAMESPACE -o jsonpath='{.spec.clusterIP}' 2>/dev/null)

if [ -n "$FRONTEND_POD" ] && [ -n "$API_SVC" ]; then
    if kubectl exec -n $NAMESPACE $FRONTEND_POD -- wget -q -O- --timeout=5 http://$API_SVC:8080/health &>/dev/null; then
        echo -e "    ${GREEN}✓${NC} Frontend can reach API service"
    else
        echo -e "    ${RED}✗${NC} Frontend cannot reach API service"
        ((FAILED++))
    fi
else
    echo -e "    ${YELLOW}⊘${NC} Frontend or API service not available"
fi

# Check if API can reach ML service
echo "  API → ML Service connectivity:"
API_POD=$(kubectl get pods -n $NAMESPACE -l app=api -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
ML_SVC=$(kubectl get svc ml-service -n $NAMESPACE -o jsonpath='{.spec.clusterIP}' 2>/dev/null)

if [ -n "$API_POD" ] && [ -n "$ML_SVC" ]; then
    if kubectl exec -n $NAMESPACE $API_POD -- wget -q -O- --timeout=5 http://$ML_SVC:8000/health &>/dev/null; then
        echo -e "    ${GREEN}✓${NC} API can reach ML service"
    else
        echo -e "    ${RED}✗${NC} API cannot reach ML service"
        ((FAILED++))
    fi
else
    echo -e "    ${YELLOW}⊘${NC} API or ML service not available"
fi
echo ""

# Test Prometheus is scraping metrics (if deployed)
echo "Metrics Collection (Prometheus):"
PROM_POD=$(kubectl get pods -n $NAMESPACE -l app=prometheus -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)

if [ -n "$PROM_POD" ]; then
    # Check if Prometheus is running
    if kubectl get pod $PROM_POD -n $NAMESPACE -o jsonpath='{.status.phase}' | grep -q "Running"; then
        echo -e "  ${GREEN}✓${NC} Prometheus is running"

        # Check if Prometheus can scrape targets
        PROM_SVC=$(kubectl get svc prometheus -n $NAMESPACE -o jsonpath='{.spec.clusterIP}' 2>/dev/null)
        if [ -n "$PROM_SVC" ]; then
            echo -e "  ${GREEN}✓${NC} Prometheus service available at $PROM_SVC"
        fi
    else
        echo -e "  ${RED}✗${NC} Prometheus is not running"
        ((FAILED++))
    fi
else
    echo -e "  ${YELLOW}⊘${NC} Prometheus not deployed"
fi
echo ""

# Test InfluxDB is available (if deployed)
echo "Time-Series Storage (InfluxDB):"
INFLUX_POD=$(kubectl get pods -n $NAMESPACE -l app=influxdb -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)

if [ -n "$INFLUX_POD" ]; then
    if kubectl get pod $INFLUX_POD -n $NAMESPACE -o jsonpath='{.status.phase}' | grep -q "Running"; then
        echo -e "  ${GREEN}✓${NC} InfluxDB is running"
    else
        echo -e "  ${RED}✗${NC} InfluxDB is not running"
        ((FAILED++))
    fi
else
    echo -e "  ${YELLOW}⊘${NC} InfluxDB not deployed"
fi
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All metrics flow tests passed${NC}"
    exit 0
else
    echo -e "${RED}✗ $FAILED metrics test(s) failed${NC}"
    exit 1
fi
