#!/bin/bash

# Smoke test for Endpoint resource
# This script tests CRUD operations for endpoints

BASE_URL="http://localhost:8081"
SESSION_NAME="a"

echo "=== Endpoint Smoke Test ==="
echo "Testing Endpoint CRUD operations..."
echo

# Test 1: Create first endpoint
echo "1. Creating first endpoint..."
ENDPOINT1_RESPONSE=$(http --session=$SESSION_NAME POST $BASE_URL/endpoints domain="api.example.com" program_id:=1 method="GET" uri="/api/v1/users" endpoint_type="api" description="Get users endpoint created" tag_ids:=[1])
ENDPOINT1_ID=$(echo "$ENDPOINT1_RESPONSE" | grep -o '[0-9]\+' | head -1)
echo "Created endpoint with ID: $ENDPOINT1_ID"
echo

# Test 2: Create second endpoint
echo "2. Creating second endpoint..."
ENDPOINT2_RESPONSE=$(http --session=$SESSION_NAME POST $BASE_URL/endpoints domain="www.example.com" program_id:=1 method="POST" uri="/login" endpoint_type="web" description="Login endpoint created" tag_ids:=[1,2])
ENDPOINT2_ID=$(echo "$ENDPOINT2_RESPONSE" | grep -o '[0-9]\+' | head -1)
echo "Created endpoint with ID: $ENDPOINT2_ID"
echo

# Test 3: Create third endpoint (for deletion test)
echo "3. Creating third endpoint..."
ENDPOINT3_RESPONSE=$(http --session=$SESSION_NAME POST $BASE_URL/endpoints domain="admin.example.com" program_id:=1 method="PUT" uri="/admin/settings" endpoint_type="web" description="Admin settings endpoint to be deleted")
ENDPOINT3_ID=$(echo "$ENDPOINT3_RESPONSE" | grep -o '[0-9]\+' | head -1)
echo "Created endpoint with ID: $ENDPOINT3_ID"
echo

# Test 4: Update second endpoint
echo "4. Updating second endpoint..."
http --session=$SESSION_NAME PUT $BASE_URL/endpoints/$ENDPOINT2_ID domain="www.updated-example.com" program_id:=1 method="POST" uri="/updated-login" endpoint_type="web" description="Login endpoint updated" tag_ids:=[1]
echo "Updated endpoint $ENDPOINT2_ID"
echo

# Test 5: Delete third endpoint
echo "5. Deleting third endpoint..."
http --session=$SESSION_NAME DELETE $BASE_URL/endpoints/$ENDPOINT3_ID
echo "Deleted endpoint $ENDPOINT3_ID"
echo

# Test 6: List all endpoints
echo "6. Listing all endpoints..."
http --session=$SESSION_NAME GET $BASE_URL/endpoints
echo

# Test 7: Get specific endpoint details
echo "7. Getting endpoint details for ID $ENDPOINT1_ID..."
http --session=$SESSION_NAME GET $BASE_URL/endpoints/$ENDPOINT1_ID
echo

echo "=== Endpoint Smoke Test Completed ==="
