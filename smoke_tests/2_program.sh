#!/bin/bash

# Smoke test for Program resource
# This script tests CRUD operations for programs

BASE_URL="http://localhost:8081"
SESSION_NAME="a"

echo "=== Program Smoke Test ==="
echo "Testing Program CRUD operations..."
echo

# Test 1: Create first program
echo "1. Creating first program..."
PROGRAM1_RESPONSE=$(http --session=$SESSION_NAME POST $BASE_URL/programs name="Test Program created" url="https://example.com" scope="internal" domains="example.com" note="Initial program note" tag_ids:=[1])
PROGRAM1_ID=$(echo "$PROGRAM1_RESPONSE" | grep -o '[0-9]\+' | head -1)
echo "Created program with ID: $PROGRAM1_ID"
echo

# Test 2: Create second program
echo "2. Creating second program..."
PROGRAM2_RESPONSE=$(http --session=$SESSION_NAME POST $BASE_URL/programs name="Another Program created" url="https://test.com" scope="external" domains="test.com" note="Another program note" tag_ids:=[1,2])
PROGRAM2_ID=$(echo "$PROGRAM2_RESPONSE" | grep -o '[0-9]\+' | head -1)
echo "Created program with ID: $PROGRAM2_ID"
echo

# Test 3: Create third program (for deletion test)
echo "3. Creating third program..."
PROGRAM3_RESPONSE=$(http --session=$SESSION_NAME POST $BASE_URL/programs name="Delete Program created" url="https://delete.com" scope="internal" domains="delete.com" note="Program to be deleted")
PROGRAM3_ID=$(echo "$PROGRAM3_RESPONSE" | grep -o '[0-9]\+' | head -1)
echo "Created program with ID: $PROGRAM3_ID"
echo

# Test 4: Update second program
echo "4. Updating second program..."
http --session=$SESSION_NAME PUT $BASE_URL/programs/$PROGRAM2_ID name="Another Program updated" url="https://updated-test.com" scope="internal" domains="updated-test.com" note="Updated program note" tag_ids:=[1]
echo "Updated program $PROGRAM2_ID"
echo

# Test 5: Delete third program
echo "5. Deleting third program..."
http --session=$SESSION_NAME DELETE $BASE_URL/programs/$PROGRAM3_ID
echo "Deleted program $PROGRAM3_ID"
echo

# Test 6: List all programs
echo "6. Listing all programs..."
http --session=$SESSION_NAME GET $BASE_URL/programs
echo

# Test 7: Get specific program details
echo "7. Getting program details for ID $PROGRAM1_ID..."
http --session=$SESSION_NAME GET $BASE_URL/programs/$PROGRAM1_ID
echo

echo "=== Program Smoke Test Completed ==="
