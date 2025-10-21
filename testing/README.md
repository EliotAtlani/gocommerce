# API Testing Suite

Automated testing suite for the GoCommerce API Gateway.

## Quick Start

### Run All Tests

```bash
# From project root
./testing/scripts/test-api.sh
```

### With Verbose Output

```bash
VERBOSE=1 ./testing/scripts/test-api.sh
```

### Test Against Different Environment

```bash
API_URL=http://staging.example.com ./testing/scripts/test-api.sh
```

## Prerequisites

### Required Tools

1. **curl** - HTTP client (usually pre-installed)
2. **jq** - JSON processor

Install jq:
```bash
# macOS
brew install jq

# Ubuntu/Debian
sudo apt-get install jq

# Alpine (Docker)
apk add jq
```

### Running Services

Ensure the API Gateway and backend services are running:

```bash
# Start all services
docker-compose up -d

# Verify services are running
docker-compose ps

# Check API Gateway is responding
curl http://localhost:8080/health
```

## Test Script Features

### Automated Validation

The test script (`test-api.sh`) automatically:

✅ **Checks dependencies** (curl, jq)
✅ **Tests health endpoint**
✅ **Registers a new user** with unique email
✅ **Logs in** and obtains JWT token
✅ **Tests authentication failures** (invalid credentials)
✅ **Validates protected endpoints** require auth
✅ **Tests invalid/malformed tokens** are rejected
✅ **Verifies authorization** (users can't access others' data)
✅ **Provides colored output** (pass/fail indicators)
✅ **Generates summary report**

### Test Coverage

| Category | Tests | Description |
|----------|-------|-------------|
| **Health** | 2 tests | Service availability |
| **Auth** | 3 tests | Register, login, invalid credentials |
| **Security** | 3 tests | Missing auth, invalid token, malformed header |
| **User API** | 2 tests | Get user, forbidden access |
| **Total** | 10 tests | Comprehensive API validation |

## Test Output Example

```
╔════════════════════════════════════════╗
║   API Gateway Integration Tests        ║
║   Target: http://localhost:8080
╚════════════════════════════════════════╝

=========================================
  Checking Dependencies
=========================================
✓ PASS: curl is installed
✓ PASS: jq is installed

=========================================
  Health Check
=========================================
✓ PASS: Health endpoint returns 200
✓ PASS: Health endpoint returns 'OK'

=========================================
  Authentication - Register
=========================================
✓ PASS: Register returns 201 Created
✓ PASS: Register returns user_id: abc-123-def

=========================================
  Authentication - Login
=========================================
✓ PASS: Login returns 200 OK
✓ PASS: Login returns JWT token

=========================================
  TEST SUMMARY
=========================================
Total Tests:  10
Passed:       10
Failed:       0
=========================================
✓ ALL TESTS PASSED
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `API_URL` | API Gateway base URL | `http://localhost:8080` |
| `VERBOSE` | Show detailed error messages | `0` (off) |

## Usage Examples

### Basic Usage

```bash
./testing/scripts/test-api.sh
```

### Debug Failed Tests

```bash
VERBOSE=1 ./testing/scripts/test-api.sh
```

### Test Production

```bash
API_URL=https://api.production.com ./testing/scripts/test-api.sh
```

### Run in CI/CD

```bash
#!/bin/bash
# .github/workflows/integration-tests.sh

set -e

# Start services
docker-compose up -d

# Wait for services to be healthy
sleep 10

# Run tests
./testing/scripts/test-api.sh

# Cleanup
docker-compose down
```

## How It Works

### 1. Dependency Check
Verifies `curl` and `jq` are installed before running tests.

### 2. Test Execution
Each test function:
- Makes HTTP request with `curl`
- Extracts HTTP status code
- Parses JSON response with `jq`
- Validates expected behavior
- Reports PASS/FAIL with colored output

### 3. State Management
Tests build on each other:
- Register creates unique user → stores `TEST_USER_ID` and `TEST_USER_EMAIL`
- Login uses stored credentials → stores `TEST_JWT_TOKEN`
- Subsequent tests use stored token for authentication

### 4. Result Tracking
- Counts total tests, passes, and failures
- Exits with code 0 (success) or 1 (failure) for CI/CD integration

## Test Scenarios

### Happy Path
1. ✅ Health check succeeds
2. ✅ Register new user
3. ✅ Login with credentials
4. ✅ Access protected resource with valid token

### Error Cases
1. ✅ Login with invalid credentials → 401
2. ✅ Access protected route without token → 401
3. ✅ Access protected route with invalid token → 401
4. ✅ Access protected route with malformed header → 401
5. ✅ Access other user's data → 403

## Troubleshooting

### Services Not Running

```bash
# Check if API Gateway is up
curl http://localhost:8080/health

# If not, start services
docker-compose up -d

# Check logs
docker-compose logs api-gateway
```

### jq Not Found

```bash
# Install jq
brew install jq  # macOS
# or
sudo apt-get install jq  # Ubuntu
```

### All Tests Failing

```bash
# Verify API URL is correct
echo $API_URL

# Test health endpoint manually
curl http://localhost:8080/health

# Check Docker services
docker-compose ps
```

### Timeouts

If tests timeout, increase curl timeout:
```bash
# Edit test-api.sh and add --max-time option
curl -s --max-time 30 ...
```

## Extending the Tests

### Add New Test

```bash
# Add to test-api.sh

test_new_feature() {
    print_test_header "New Feature Test"

    response=$(curl -s -w "\n%{http_code}" -X GET "$API_URL/api/v1/new-endpoint")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    if [ "$http_code" = "200" ]; then
        pass "New endpoint returns 200"
    else
        fail "New endpoint failed" "Expected 200, got $http_code"
    fi
}

# Add to main() function
main() {
    # ... existing tests ...
    test_new_feature  # ← Add here
    print_summary
}
```

### Test User Update

```bash
test_user_update() {
    print_test_header "User Endpoints - Update User"

    if [ -z "$TEST_JWT_TOKEN" ] || [ -z "$TEST_USER_ID" ]; then
        fail "Skipping update test - no auth token"
        return
    fi

    response=$(curl -s -w "\n%{http_code}" -X PUT \
        "$API_URL/api/v1/users/$TEST_USER_ID" \
        -H "Authorization: Bearer $TEST_JWT_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"name":"Updated Name","phone":"+1234567890"}')

    http_code=$(echo "$response" | tail -n1)

    if [ "$http_code" = "200" ]; then
        pass "Update user returns 200"
    else
        fail "Update user failed" "Got $http_code"
    fi
}
```

## Best Practices

### 1. **Unique Test Data**
Always use unique identifiers (timestamps, UUIDs) for test data:
```bash
TEST_EMAIL="test-$(date +%s)@example.com"
```

### 2. **Cleanup**
Consider adding cleanup after tests:
```bash
cleanup() {
    if [ -n "$TEST_USER_ID" ]; then
        curl -s -X DELETE "$API_URL/api/v1/users/$TEST_USER_ID" \
            -H "Authorization: Bearer $TEST_JWT_TOKEN"
    fi
}
trap cleanup EXIT
```

### 3. **Parallel Testing**
Avoid parallel execution of this script - tests share state (JWT token, user ID).

### 4. **CI/CD Integration**
Script exits with proper codes:
- `0` = All tests passed
- `1` = One or more tests failed

## Future Enhancements

- [ ] Add user update/delete tests
- [ ] Add address management tests
- [ ] Test pagination and filtering
- [ ] Performance/load testing
- [ ] Test concurrent requests
- [ ] Add cleanup after tests
- [ ] Generate HTML test report
- [ ] Integration with test frameworks (BATS)

## Related Documentation

- [API Gateway README](../api-gateway/README.md) - API endpoint documentation
- [Project Description](../docs/PROJECT_DESCRIPTION.md) - Overall architecture

---

**Pro Tip:** Run this script in your CI/CD pipeline to catch API regressions automatically!
