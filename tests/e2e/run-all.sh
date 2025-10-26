#!/bin/bash
# Sentinel Mesh E2E Test Suite Runner
# Runs all E2E tests and reports results

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results
PASSED=0
FAILED=0
SKIPPED=0

# Test result arrays
declare -a PASSED_TESTS
declare -a FAILED_TESTS
declare -a SKIPPED_TESTS

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║         Sentinel Mesh E2E Test Suite                      ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Helper function to run a test
run_test() {
    local test_name="$1"
    local test_script="$2"

    echo -e "${BLUE}Running: ${test_name}${NC}"

    if [ ! -f "$test_script" ]; then
        echo -e "${YELLOW}⊘ SKIPPED${NC} - Script not found: $test_script"
        ((SKIPPED++))
        SKIPPED_TESTS+=("$test_name")
        return
    fi

    if bash "$test_script"; then
        echo -e "${GREEN}✓ PASSED${NC} - $test_name"
        ((PASSED++))
        PASSED_TESTS+=("$test_name")
    else
        echo -e "${RED}✗ FAILED${NC} - $test_name"
        ((FAILED++))
        FAILED_TESTS+=("$test_name")
    fi
    echo ""
}

# Record start time
START_TIME=$(date +%s)

echo "Test Environment:"
echo "  Kubernetes Cluster: $(kubectl config current-context)"
echo "  Namespace: sentinel-mesh"
echo ""

# Run test suites
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}1. Service Health & Availability Tests${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
run_test "Service Health Checks" "$SCRIPT_DIR/test-service-health.sh"

echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}2. API Integration Tests${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
run_test "API Endpoints" "$SCRIPT_DIR/test-api-endpoints.sh"

echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}3. Frontend Tests${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
run_test "Frontend Availability" "$SCRIPT_DIR/test-frontend.sh"

echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}4. ML Service Tests${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
run_test "ML Service Endpoints" "$SCRIPT_DIR/test-ml-service.sh"

echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}5. Kubernetes Integration Tests${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
run_test "Kubernetes Resources" "$SCRIPT_DIR/test-k8s-resources.sh"

echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}6. Data Flow Tests${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
run_test "Metrics Collection" "$SCRIPT_DIR/test-metrics-flow.sh"

# Calculate duration
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

# Print summary
echo ""
echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                    Test Summary                            ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

TOTAL=$((PASSED + FAILED + SKIPPED))
echo "  Total Tests: $TOTAL"
echo -e "  ${GREEN}✓ Passed: $PASSED${NC}"
echo -e "  ${RED}✗ Failed: $FAILED${NC}"
echo -e "  ${YELLOW}⊘ Skipped: $SKIPPED${NC}"
echo "  Duration: ${DURATION}s"
echo ""

# Print detailed results if there are failures or skipped tests
if [ $FAILED -gt 0 ]; then
    echo -e "${RED}Failed Tests:${NC}"
    for test in "${FAILED_TESTS[@]}"; do
        echo -e "  ${RED}✗${NC} $test"
    done
    echo ""
fi

if [ $SKIPPED -gt 0 ]; then
    echo -e "${YELLOW}Skipped Tests:${NC}"
    for test in "${SKIPPED_TESTS[@]}"; do
        echo -e "  ${YELLOW}⊘${NC} $test"
    done
    echo ""
fi

# Exit with appropriate code
if [ $FAILED -gt 0 ]; then
    echo -e "${RED}❌ Test suite FAILED${NC}"
    exit 1
else
    echo -e "${GREEN}✅ Test suite PASSED${NC}"
    exit 0
fi
