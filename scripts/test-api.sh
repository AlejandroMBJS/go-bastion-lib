#!/bin/bash

set -euo pipefail

# --- Configuration ---
BASE_URL="http://localhost:9876"
SERVER_CMD="go run ./examples/basic/main.go"
SERVER_PID=""
ACCESS_TOKEN=""
TEST_COUNT=0
PASS_COUNT=0
FAIL_COUNT=0

# --- Helper Functions ---

log_info() {
  echo "INFO: $*"
}

log_error() {
  echo "ERROR: $*" >&2
}

start_server() {
  log_info "Starting server: ${SERVER_CMD}"
  # Start server in background, redirecting output to a log file
  ${SERVER_CMD} > server.log 2>&1 &
  SERVER_PID=$!
  log_info "Server started with PID: ${SERVER_PID}"

  log_info "Waiting for server to be ready..."
  for i in $(seq 1 10);
  do
    if curl -s "${BASE_URL}/api/health" > /dev/null;
    then
      log_info "Server is ready."
      return 0
    fi
    sleep 1
  done
  log_error "Server did not become ready in time."
  return 1
}

stop_server() {
  if [ -n "${SERVER_PID}" ]; then
    log_info "Stopping server with PID: ${SERVER_PID}"
    kill "${SERVER_PID}"
    wait "${SERVER_PID}" || true # Wait for the process to terminate, ignore error if already dead
    log_info "Server stopped."
  fi
}

run_test() {
  local name="$1"
  local command="$2"
  local expected_status="$3"
  local assertion_jq="${4:-}" # Optional jq assertion

  TEST_COUNT=$((TEST_COUNT + 1))
  log_info "Running test [${TEST_COUNT}]: ${name}"
  log_info "Command: ${command}"

  local http_code
  local response_body
  
  # Execute curl command, capture HTTP code and body
  # Use a temporary file for headers to extract status code
  response_body=$(eval "${command}" -s -o /dev/null -w '%{http_code}' -D /dev/stderr)
  http_code=$(tail -n 1 /dev/stderr | tr -d '\n') # Extract last line (HTTP code) from stderr
  
  # Check status code
  if [ "${http_code}" -eq "${expected_status}" ]; then
    log_info "Status Code: ${http_code} (Expected: ${expected_status}) - PASS"
    
    # If jq assertion is provided, run it
    if [ -n "${assertion_jq}" ]; then
      if echo "${response_body}" | jq -e "${assertion_jq}" > /dev/null;
      then
        log_info "JQ Assertion: '${assertion_jq}' - PASS"
        PASS_COUNT=$((PASS_COUNT + 1))
      else
        log_error "JQ Assertion: '${assertion_jq}' - FAIL"
        log_error "Response Body: ${response_body}"
        FAIL_COUNT=$((FAIL_COUNT + 1))
      fi
    else
      PASS_COUNT=$((PASS_COUNT + 1))
    fi
  else
    log_error "Status Code: ${http_code} (Expected: ${expected_status}) - FAIL"
    log_error "Response Body: ${response_body}"
    FAIL_COUNT=$((FAIL_COUNT + 1))
  fi
  echo "" # Newline for readability
}

# --- Test Cases ---

test_health_endpoint() {
  log_info "--- Testing /api/health ---"
  run_test "Health Check (Happy Path)" "curl -X GET \"${BASE_URL}/api/health\"" 200 ".status == \"ok\""
  # Rate limit error path is hard to automate reliably without waiting or complex setup
}

test_login_endpoint() {
  log_info "--- Testing /api/auth/login ---"
  local login_response
  login_response=$(curl -s -X POST "${BASE_URL}/api/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"testuser","password":"testpass"}')
  
  # Extract token for subsequent tests
  ACCESS_TOKEN=$(echo "${login_response}" | jq -r '.access_token')
  
  run_test "Login (Happy Path)" \
    "curl -X POST \"${BASE_URL}/api/auth/login\" -H \"Content-Type: application/json\" -d '{\"username\":\"testuser\",\"password\":\"testpass\"}'" \
    200 \
    ".access_token != null and .token_type == \"bearer\""

  run_test "Login (Error Path: Invalid JSON Body)" \
    "curl -X POST \"${BASE_URL}/api/auth/login\" -H \"Content-Type: application/json\" -d '{\"username\":\"testuser\",\"password\":}'" \
    400 \
    ".error.code == \"invalid_request\""

  run_test "Login (Error Path: Empty JSON Body)" \
    "curl -X POST \"${BASE_URL}/api/auth/login\" -H \"Content-Type: application/json\" -d ''" \
    400 \
    ".error.code == \"invalid_request\""
}

test_users_list_endpoint() {
  log_info "--- Testing /api/users (GET) ---"
  if [ -z "${ACCESS_TOKEN}" ]; then
    log_error "Skipping /api/users (GET) tests: ACCESS_TOKEN not available."
    return
  fi

  run_test "List Users (Happy Path)" \
    "curl -X GET \"${BASE_URL}/api/users\" -H \"Authorization: Bearer ${ACCESS_TOKEN}\"" \
    200 \
    ".data.users | length > 0"

  run_test "List Users (Error Path: Missing JWT)" \
    "curl -X GET \"${BASE_URL}/api/users\"" \
    401 \
    ".error == \"unauthorized\""

  run_test "List Users (Error Path: Invalid JWT)" \
    "curl -X GET \"${BASE_URL}/api/users\" -H \"Authorization: Bearer invalid.jwt.token\"" \
    401 \
    ".error == \"unauthorized\""
}

test_users_create_endpoint() {
  log_info "--- Testing /api/users (POST) ---"
  if [ -z "${ACCESS_TOKEN}" ]; then
    log_error "Skipping /api/users (POST) tests: ACCESS_TOKEN not available."
    return
  fi

  run_test "Create User (Happy Path)" \
    "curl -X POST \"${BASE_URL}/api/users\" -H \"Authorization: Bearer ${ACCESS_TOKEN}\" -H \"Content-Type: application/json\" -d '{\"username\":\"newuser_$(date +%s)\",\"email\":\"newuser_$(date +%s)@example.com\",\"password\":\"securepassword123\"}'" \
    201 \
    ".data.id != null"

  run_test "Create User (Error Path: Missing JWT)" \
    "curl -X POST \"${BASE_URL}/api/users\" -H \"Content-Type: application/json\" -d '{\"username\":\"newuser_no_jwt\",\"email\":\"newuser_no_jwt@example.com\",\"password\":\"securepassword123\"}'" \
    401 \
    ".error == \"unauthorized\""

  run_test "Create User (Error Path: Invalid JSON Body)" \
    "curl -X POST \"${BASE_URL}/api/users\" -H \"Authorization: Bearer ${ACCESS_TOKEN}\" -H \"Content-Type: application/json\" -d '{\"username\":\"newuser_bad_json\",\"email\":\"newuser_bad_json@example.com\",\"password\":}'" \
    400 \
    ".error.code == \"invalid_request\""

  run_test "Create User (Error Path: Missing Required Fields)" \
    "curl -X POST \"${BASE_URL}/api/users\" -H \"Authorization: Bearer ${ACCESS_TOKEN}\" -H \"Content-Type: application/json\" -d '{\"username\":\"\",\"email\":\"newuser_missing_field@example.com\",\"password\":\"securepassword123\"}'" \
    400 \
    ".error.code == \"validation_error\""
}

# --- Main Execution ---

# Ensure server is stopped on exit
trap stop_server EXIT

# Start the server
if ! start_server; then
  log_error "Failed to start server. Exiting."
  exit 1
fi

# Run all test functions
test_health_endpoint
test_login_endpoint
test_users_list_endpoint
test_users_create_endpoint

# --- Summary ---
echo "--- Test Summary ---"
echo "Total Tests: ${TEST_COUNT}"
echo "Passed: ${PASS_COUNT}"
echo "Failed: ${FAIL_COUNT}"

if [ "${FAIL_COUNT}" -eq 0 ]; then
  log_info "All API tests passed!"
  exit 0
else
  log_error "Some API tests failed."
  exit 1
fi
