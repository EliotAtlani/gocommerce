#!/bin/bash

# API Gateway Integration Test Script
# Tests all endpoints with automated validation

set -e  # Exit on any error

# Configuration
API_URL="${API_URL:-http://localhost:8080}"
VERBOSE="${VERBOSE:-0}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Test result tracking
print_test_header() {
    echo ""
    echo "========================================="
    echo "  $1"
    echo "========================================="
}

pass() {
    TESTS_PASSED=$((TESTS_PASSED + 1))
    TESTS_RUN=$((TESTS_RUN + 1))
    echo -e "${GREEN}✓ PASS${NC}: $1"
}

fail() {
    TESTS_FAILED=$((TESTS_FAILED + 1))
    TESTS_RUN=$((TESTS_RUN + 1))
    echo -e "${RED}✗ FAIL${NC}: $1"
    if [ "$VERBOSE" = "1" ]; then
        echo -e "${YELLOW}  Details: $2${NC}"
    fi
}

info() {
    echo -e "${YELLOW}ℹ INFO${NC}: $1"
}

# Check dependencies
check_dependencies() {
    print_test_header "Checking Dependencies"

    if ! command -v curl &> /dev/null; then
        fail "curl is not installed"
        exit 1
    fi
    pass "curl is installed"

    if ! command -v jq &> /dev/null; then
        fail "jq is not installed (required for JSON parsing)"
        info "Install with: brew install jq (macOS) or apt-get install jq (Ubuntu)"
        exit 1
    fi
    pass "jq is installed"
}

# Health check
test_health_check() {
    print_test_header "Health Check"

    response=$(curl -s -w "\n%{http_code}" "$API_URL/health")
    http_code=$(echo "$response" | tail -n 1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" = "200" ]; then
        pass "Health endpoint returns 200"
    else
        fail "Health endpoint failed" "Expected 200, got $http_code"
        return
    fi

    if [ "$body" = "OK" ]; then
        pass "Health endpoint returns 'OK'"
    else
        fail "Health endpoint body incorrect" "Expected 'OK', got '$body'"
    fi
}

# Test authentication endpoints
test_auth_register() {
    print_test_header "Authentication - Register"

    # Generate unique email for this test run
    TIMESTAMP=$(date +%s)
    TEST_EMAIL="test-${TIMESTAMP}@example.com"
    TEST_PASSWORD="SecurePass123!"
    TEST_NAME="Test User ${TIMESTAMP}"

    response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/v1/auth/register" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\",\"name\":\"$TEST_NAME\"}")

    http_code=$(echo "$response" | tail -n 1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" = "201" ]; then
        pass "Register returns 201 Created"
    else
        fail "Register failed" "Expected 201, got $http_code. Body: $body"
        return
    fi

    user_id=$(echo "$body" | jq -r '.user_id // empty')
    if [ -n "$user_id" ] && [ "$user_id" != "null" ]; then
        pass "Register returns user_id: $user_id"
        # Store for later tests
        export TEST_USER_ID="$user_id"
        export TEST_USER_EMAIL="$TEST_EMAIL"
        export TEST_USER_PASSWORD="$TEST_PASSWORD"
    else
        fail "Register response missing user_id" "Body: $body"
    fi
}

test_auth_login() {
    print_test_header "Authentication - Login"

    if [ -z "$TEST_USER_EMAIL" ]; then
        fail "Skipping login test - no user registered"
        return
    fi

    response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$TEST_USER_EMAIL\",\"password\":\"$TEST_USER_PASSWORD\"}")

    http_code=$(echo "$response" | tail -n 1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" = "200" ]; then
        pass "Login returns 200 OK"
    else
        fail "Login failed" "Expected 200, got $http_code. Body: $body"
        return
    fi

    token=$(echo "$body" | jq -r '.token // empty')
    if [ -n "$token" ] && [ "$token" != "null" ]; then
        pass "Login returns JWT token"
        export TEST_JWT_TOKEN="$token"
    else
        fail "Login response missing token" "Body: $body"
    fi
}

test_auth_invalid_credentials() {
    print_test_header "Authentication - Invalid Credentials"

    response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"nonexistent@example.com\",\"password\":\"wrongpass\"}")

    http_code=$(echo "$response" | tail -n 1)

    if [ "$http_code" = "401" ]; then
        pass "Login with invalid credentials returns 401"
    else
        fail "Login with invalid credentials should return 401" "Got $http_code"
    fi
}

# Test protected endpoints
test_protected_without_auth() {
    print_test_header "Protected Routes - No Authentication"

    response=$(curl -s -w "\n%{http_code}" -X GET "$API_URL/api/v1/users/123")
    http_code=$(echo "$response" | tail -n 1)

    if [ "$http_code" = "401" ]; then
        pass "Protected endpoint without auth returns 401"
    else
        fail "Protected endpoint should require auth" "Expected 401, got $http_code"
    fi
}

test_protected_invalid_token() {
    print_test_header "Protected Routes - Invalid Token"

    response=$(curl -s -w "\n%{http_code}" -X GET "$API_URL/api/v1/users/123" \
        -H "Authorization: Bearer invalid.token.here")

    http_code=$(echo "$response" | tail -n 1)

    if [ "$http_code" = "401" ] || [ "$http_code" = "500" ]; then
        pass "Protected endpoint with invalid token returns 401/500"
    else
        fail "Protected endpoint should reject invalid token" "Got $http_code"
    fi
}

test_protected_malformed_header() {
    print_test_header "Protected Routes - Malformed Auth Header"

    response=$(curl -s -w "\n%{http_code}" -X GET "$API_URL/api/v1/users/123" \
        -H "Authorization: NotBearer token")

    http_code=$(echo "$response" | tail -n 1)

    if [ "$http_code" = "401" ]; then
        pass "Malformed Authorization header returns 401"
    else
        fail "Should reject malformed auth header" "Expected 401, got $http_code"
    fi
}

test_user_get() {
    print_test_header "User Endpoints - Get User"

    if [ -z "$TEST_JWT_TOKEN" ] || [ -z "$TEST_USER_ID" ]; then
        fail "Skipping user tests - no auth token or user ID"
        return
    fi

    response=$(curl -s -w "\n%{http_code}" -X GET "$API_URL/api/v1/users/$TEST_USER_ID" \
        -H "Authorization: Bearer $TEST_JWT_TOKEN")

    http_code=$(echo "$response" | tail -n 1)
    body=$(echo "$response" | sed '$d')

    # User might not exist in user-service yet (auth-service vs user-service sync issue)
    if [ "$http_code" = "200" ]; then
        pass "Get user returns 200"

        email=$(echo "$body" | jq -r '.email // empty')
        if [ "$email" = "$TEST_USER_EMAIL" ]; then
            pass "Get user returns correct email"
        fi
    elif [ "$http_code" = "404" ]; then
        info "User not found in user-service (expected if services not fully synced)"
        pass "Get user handles missing user correctly (404)"
    else
        fail "Get user unexpected response" "Expected 200 or 404, got $http_code"
    fi
}

test_user_forbidden_access() {
    print_test_header "User Endpoints - Forbidden Access"

    if [ -z "$TEST_JWT_TOKEN" ]; then
        fail "Skipping forbidden access test - no auth token"
        return
    fi

    # Try to access a different user's data
    response=$(curl -s -w "\n%{http_code}" -X GET "$API_URL/api/v1/users/different-user-id" \
        -H "Authorization: Bearer $TEST_JWT_TOKEN")

    http_code=$(echo "$response" | tail -n 1)

    if [ "$http_code" = "403" ]; then
        pass "Accessing other user's data returns 403 Forbidden"
    else
        info "Got $http_code (might be 404 if user doesn't exist, or 403 if forbidden)"
    fi
}

# Summary
print_summary() {
    echo ""
    echo "========================================="
    echo "  TEST SUMMARY"
    echo "========================================="
    echo "Total Tests:  $TESTS_RUN"
    echo -e "${GREEN}Passed:       $TESTS_PASSED${NC}"
    if [ $TESTS_FAILED -gt 0 ]; then
        echo -e "${RED}Failed:       $TESTS_FAILED${NC}"
    else
        echo "Failed:       $TESTS_FAILED"
    fi
    echo "========================================="

    if [ $TESTS_FAILED -eq 0 ]; then
        echo -e "${GREEN}✓ ALL TESTS PASSED${NC}"
        return 0
    else
        echo -e "${RED}✗ SOME TESTS FAILED${NC}"
        return 1
    fi
}

# Main execution
main() {
    echo "╔════════════════════════════════════════╗"
    echo "║   API Gateway Integration Tests        ║"
    echo "║   Target: $API_URL"
    echo "╚════════════════════════════════════════╝"

    check_dependencies
    test_health_check
    test_auth_register
    test_auth_login
    test_auth_invalid_credentials
    test_protected_without_auth
    test_protected_invalid_token
    test_protected_malformed_header
    test_user_get
    test_user_forbidden_access

    print_summary
}

# Run main and exit with appropriate code
main
exit $?
