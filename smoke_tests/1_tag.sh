#!/bin/bash

# Smoke test for Tag resource
# This script tests CRUD operations for tags

BASE_URL="http://localhost:8081"
SESSION_NAME="a"

echo "=== Tag Smoke Test ==="
echo "Testing Tag CRUD operations..."
echo

# Test 1: Create first tag
echo "1. Creating first tag..."
TAG1_RESPONSE=$(http --session=$SESSION_NAME POST $BASE_URL/tags name="High Priority created" priority:=2)
TAG1_ID=$(echo "$TAG1_RESPONSE" | grep -o '[0-9]\+' | head -1)
echo "Created tag with ID: $TAG1_ID"
echo

# Test 2: Create second tag
echo "2. Creating second tag..."
TAG2_RESPONSE=$(http --session=$SESSION_NAME POST $BASE_URL/tags name="Medium Priority created" priority:=3)
TAG2_ID=$(echo "$TAG2_RESPONSE" | grep -o '[0-9]\+' | head -1)
echo "Created tag with ID: $TAG2_ID"
echo

# Test 3: Create third tag (for deletion test)
echo "3. Creating third tag..."
TAG3_RESPONSE=$(http --session=$SESSION_NAME POST $BASE_URL/tags name="Low Priority created" priority:=1)
TAG3_ID=$(echo "$TAG3_RESPONSE" | grep -o '[0-9]\+' | head -1)
echo "Created tag with ID: $TAG3_ID"
echo

# Test 4: Update second tag
echo "4. Updating second tag..."
http --session=$SESSION_NAME PUT $BASE_URL/tags/$TAG2_ID name="Medium Priority updated" priority:=5
echo "Updated tag $TAG2_ID"
echo

# Test 5: Delete third tag
echo "5. Deleting third tag..."
http --session=$SESSION_NAME DELETE $BASE_URL/tags/$TAG3_ID
echo "Deleted tag $TAG3_ID"
echo

# Test 6: List all tags
echo "6. Listing all tags..."
http --session=$SESSION_NAME GET $BASE_URL/tags
echo

# Test 7: Get specific tag details
echo "7. Getting tag details for ID $TAG1_ID..."
http --session=$SESSION_NAME GET $BASE_URL/tags/$TAG1_ID
echo

echo "=== Tag Smoke Test Completed ==="
