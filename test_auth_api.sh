#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Base URL
BASE_URL="http://localhost:8080/api"

# Variables to store data between tests
TOKEN=""
USER_EMAIL=""
USER_ID=""
RESPONSE_FILE="/tmp/api_response_$$"
HTTP_CODE_FILE="/tmp/api_code_$$"

# Cleanup temp files on exit
trap "rm -f $RESPONSE_FILE $HTTP_CODE_FILE" EXIT

# Function to print section headers
print_header() {
    echo -e "\n${BLUE}================================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}================================================${NC}\n"
}

# Function to print test results
print_result() {
    if [ "$1" = "pass" ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
    fi
}

# Function to make API calls and display results
api_call() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    local auth_header=$5

    echo -e "${YELLOW}Testing: $description${NC}"
    echo "Request: $method $endpoint"

    if [ -n "$data" ]; then
        echo "Data: $data"
    fi

    echo -e "${BLUE}CURL:${NC} curl -X \"$method\" -H \"Content-Type: application/json\"${auth_header:+ -H \"Authorization: Bearer $auth_header\"} -d '$data' \"$BASE_URL$endpoint\""

    if [ -n "$auth_header" ]; then
        http_code=$(curl -s -w "%{http_code}" -o "$RESPONSE_FILE" -X "$method" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $auth_header" \
            -d "$data" \
            "$BASE_URL$endpoint")
    else
        http_code=$(curl -s -w "%{http_code}" -o "$RESPONSE_FILE" -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$BASE_URL$endpoint")
    fi

    body=$(cat "$RESPONSE_FILE")

    echo "Response Code: $http_code"
    echo "Response Body: $body" | jq '.' 2>/dev/null || echo "$body"
    echo ""

    # Return the http code
    echo "$http_code" > "$HTTP_CODE_FILE"
}

# Test 1: Register a new user
test_register() {
    print_header "TEST 1: Register New User"

    # Generate unique email with timestamp
    timestamp=$(date +%s)
    USER_EMAIL="testuser${timestamp}@example.com"

    api_call "POST" "/auth/register" \
        "{\"name\":\"Test User\",\"email\":\"$USER_EMAIL\",\"password\":\"password123\"}" \
        "Register new user with valid data"

    http_code=$(cat "$HTTP_CODE_FILE")

    if [ "$http_code" = "201" ]; then
        # Extract token from response
        TOKEN=$(cat "$RESPONSE_FILE" | jq -r '.token' 2>/dev/null)
        USER_ID=$(cat "$RESPONSE_FILE" | jq -r '.user.id' 2>/dev/null)
        print_result "pass" "User registered successfully"
        echo "Token: $TOKEN"
        echo "User ID: $USER_ID"
    else
        print_result "fail" "Failed to register user (expected 201, got $http_code)"
    fi
}

# Test 2: Register with duplicate email
test_register_duplicate() {
    print_header "TEST 2: Register with Duplicate Email"

    api_call "POST" "/auth/register" \
        "{\"name\":\"Test User 2\",\"email\":\"$USER_EMAIL\",\"password\":\"password123\"}" \
        "Attempt to register with duplicate email"

    http_code=$(cat "$HTTP_CODE_FILE")

    if [ "$http_code" = "500" ] || [ "$http_code" = "400" ] || [ "$http_code" = "409" ]; then
        print_result "pass" "Correctly rejected duplicate email"
    else
        print_result "fail" "Should reject duplicate email (got $http_code)"
    fi
}

# Test 3: Register with invalid email
test_register_invalid_email() {
    print_header "TEST 3: Register with Invalid Email"

    api_call "POST" "/auth/register" \
        "{\"name\":\"Test User\",\"email\":\"invalid-email\",\"password\":\"password123\"}" \
        "Register with invalid email format"

    http_code=$(cat "$HTTP_CODE_FILE")

    if [ "$http_code" = "400" ]; then
        print_result "pass" "Correctly rejected invalid email"
    else
        print_result "fail" "Should reject invalid email (expected 400, got $http_code)"
    fi
}

# Test 4: Register with short password
test_register_short_password() {
    print_header "TEST 4: Register with Short Password"

    api_call "POST" "/auth/register" \
        "{\"name\":\"Test User\",\"email\":\"newuser@example.com\",\"password\":\"12345\"}" \
        "Register with password less than 6 characters"

    http_code=$(cat "$HTTP_CODE_FILE")

    if [ "$http_code" = "400" ]; then
        print_result "pass" "Correctly rejected short password"
    else
        print_result "fail" "Should reject short password (expected 400, got $http_code)"
    fi
}

# Test 5: Register with missing fields
test_register_missing_fields() {
    print_header "TEST 5: Register with Missing Fields"

    api_call "POST" "/auth/register" \
        "{\"email\":\"test@example.com\"}" \
        "Register with missing name and password"

    http_code=$(cat "$HTTP_CODE_FILE")

    if [ "$http_code" = "400" ]; then
        print_result "pass" "Correctly rejected missing fields"
    else
        print_result "fail" "Should reject missing fields (expected 400, got $http_code)"
    fi
}

# Test 6: Login with valid credentials
test_login_valid() {
    print_header "TEST 6: Login with Valid Credentials"

    api_call "POST" "/auth/login" \
        "{\"email\":\"$USER_EMAIL\",\"password\":\"password123\"}" \
        "Login with correct credentials"

    http_code=$(cat "$HTTP_CODE_FILE")

    if [ "$http_code" = "200" ]; then
        # Update token with login token
        TOKEN=$(cat "$RESPONSE_FILE" | jq -r '.token' 2>/dev/null)
        print_result "pass" "Login successful"
        echo "New Token: $TOKEN"
    else
        print_result "fail" "Login failed (expected 200, got $http_code)"
    fi
}

# Test 7: Login with wrong password
test_login_wrong_password() {
    print_header "TEST 7: Login with Wrong Password"

    api_call "POST" "/auth/login" \
        "{\"email\":\"$USER_EMAIL\",\"password\":\"wrongpassword\"}" \
        "Login with incorrect password"

    http_code=$(cat "$HTTP_CODE_FILE")

    if [ "$http_code" = "401" ]; then
        print_result "pass" "Correctly rejected wrong password"
    else
        print_result "fail" "Should reject wrong password (expected 401, got $http_code)"
    fi
}

# Test 8: Login with non-existent email
test_login_nonexistent() {
    print_header "TEST 8: Login with Non-existent Email"

    api_call "POST" "/auth/login" \
        "{\"email\":\"nonexistent@example.com\",\"password\":\"password123\"}" \
        "Login with non-existent email"

    http_code=$(cat "$HTTP_CODE_FILE")

    if [ "$http_code" = "401" ]; then
        print_result "pass" "Correctly rejected non-existent email"
    else
        print_result "fail" "Should reject non-existent email (expected 401, got $http_code)"
    fi
}

# Test 9: Login with invalid email format
test_login_invalid_email() {
    print_header "TEST 9: Login with Invalid Email Format"

    api_call "POST" "/auth/login" \
        "{\"email\":\"invalid-email\",\"password\":\"password123\"}" \
        "Login with invalid email format"

    http_code=$(cat "$HTTP_CODE_FILE")

    if [ "$http_code" = "400" ]; then
        print_result "pass" "Correctly rejected invalid email format"
    else
        print_result "fail" "Should reject invalid email (expected 400, got $http_code)"
    fi
}

# Test 10: Get profile with valid token
test_get_profile_valid() {
    print_header "TEST 10: Get Profile with Valid Token"

    api_call "GET" "/auth/profile" "" "Get user profile with valid token" "$TOKEN"

    http_code=$(cat "$HTTP_CODE_FILE")

    if [ "$http_code" = "200" ]; then
        print_result "pass" "Successfully retrieved profile"
    else
        print_result "fail" "Failed to get profile (expected 200, got $http_code)"
    fi
}

# Test 11: Get profile without token
test_get_profile_no_token() {
    print_header "TEST 11: Get Profile without Token"

    api_call "GET" "/auth/profile" "" "Get user profile without token"

    http_code=$(cat "$HTTP_CODE_FILE")

    if [ "$http_code" = "401" ]; then
        print_result "pass" "Correctly rejected request without token"
    else
        print_result "fail" "Should reject request without token (expected 401, got $http_code)"
    fi
}

# Test 12: Get profile with invalid token
test_get_profile_invalid_token() {
    print_header "TEST 12: Get Profile with Invalid Token"

    api_call "GET" "/auth/profile" "" "Get user profile with invalid token" "invalid.token.here"

    http_code=$(cat "$HTTP_CODE_FILE")

    if [ "$http_code" = "401" ]; then
        print_result "pass" "Correctly rejected invalid token"
    else
        print_result "fail" "Should reject invalid token (expected 401, got $http_code)"
    fi
}

# Main execution
main() {
    echo -e "${GREEN}======================================${NC}"
    echo -e "${GREEN}   Auth API Testing Script${NC}"
    echo -e "${GREEN}======================================${NC}"
    echo "Base URL: $BASE_URL"
    echo "Started at: $(date)"

    # Check if jq is installed
    if ! command -v jq &> /dev/null; then
        echo -e "${YELLOW}Warning: jq is not installed. JSON output will not be formatted.${NC}"
    fi

    # Run all tests
    test_register
    test_register_duplicate
    test_register_invalid_email
    test_register_short_password
    test_register_missing_fields
    test_login_valid
    test_login_wrong_password
    test_login_nonexistent
    test_login_invalid_email
    test_get_profile_valid
    test_get_profile_no_token
    test_get_profile_invalid_token

    # Summary
    print_header "TEST SUMMARY"
    echo "Completed at: $(date)"
    echo -e "\n${GREEN}All tests completed!${NC}"
    echo -e "Review the results above for any failures.\n"
}

# Run main function
main
