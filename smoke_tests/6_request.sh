#!/bin/bash

# Smoke test for Request resource
# This script tests read operations for requests (requests are typically created via imports)

BASE_URL="http://localhost:8081"
SESSION_NAME="a"

echo "=== Request Smoke Test ==="
echo "Testing Request read operations..."
echo

# Test 1: List all requests
echo "1. Listing all requests..."
http --session=$SESSION_NAME GET $BASE_URL/requests
echo

# Test 2: List requests with program filter
echo "2. Listing requests filtered by program ID 1..."
http --session=$SESSION_NAME GET $BASE_URL/requests program_id=1
echo

# Test 3: List requests with endpoint filter
echo "3. Listing requests filtered by endpoint ID 1..."
http --session=$SESSION_NAME GET $BASE_URL/requests endpoint_id=1
echo

# Test 4: List requests with job filter
echo "4. Listing requests filtered by job ID 1..."
http --session=$SESSION_NAME GET $BASE_URL/requests job_id=1
echo

# Test 5: Search requests
echo "5. Searching requests for 'GET'..."
http --session=$SESSION_NAME GET $BASE_URL/requests search=GET
echo

# Test 6: List requests with ordering
echo "6. Listing requests ordered by method..."
http --session=$SESSION_NAME GET $BASE_URL/requests order_by1=method asc1:=true
echo

# Test 7: List requests with multiple ordering
echo "7. Listing requests with multiple ordering (method, then size)..."
http --session=$SESSION_NAME GET $BASE_URL/requests order_by1=method asc1:=true order_by2=size asc2:=false
echo

# Test 8: Get specific request details (if any exist)
echo "8. Attempting to get request details for ID 1..."
http --session=$SESSION_NAME GET $BASE_URL/requests/1
echo

# Test 9: Get request details for ID 2 (if any exist)
echo "9. Attempting to get request details for ID 2..."
http --session=$SESSION_NAME GET $BASE_URL/requests/2
echo

echo "=== Request Smoke Test Completed ==="
echo "Note: Requests are typically created via import operations (HAR/Burp XML)"
echo "This test focuses on read operations and filtering capabilities"
