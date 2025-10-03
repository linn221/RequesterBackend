#!/bin/bash

# Smoke test for Note resource
# This script tests CRUD operations for notes

BASE_URL="http://localhost:8081"
SESSION_NAME="a"

echo "=== Note Smoke Test ==="
echo "Testing Note CRUD operations..."
echo

# Test 1: Create first note (for program)
echo "1. Creating first note for program..."
NOTE1_RESPONSE=$(http --session=$SESSION_NAME POST $BASE_URL/notes reference_type="programs" reference_id:=1 value="This is a note created for the program" tag_ids:=[1])
NOTE1_ID=$(echo "$NOTE1_RESPONSE" | grep -o '[0-9]\+' | head -1)
echo "Created note with ID: $NOTE1_ID"
echo

# Test 2: Create second note (for endpoint)
echo "2. Creating second note for endpoint..."
NOTE2_RESPONSE=$(http --session=$SESSION_NAME POST $BASE_URL/notes reference_type="endpoints" reference_id:=1 value="This is a note created for the endpoint" tag_ids:=[1,2])
NOTE2_ID=$(echo "$NOTE2_RESPONSE" | grep -o '[0-9]\+' | head -1)
echo "Created note with ID: $NOTE2_ID"
echo

# Test 3: Create third note (for vulnerability)
echo "3. Creating third note for vulnerability..."
NOTE3_RESPONSE=$(http --session=$SESSION_NAME POST $BASE_URL/notes reference_type="vulns" reference_id:=1 value="This is a note created for the vulnerability" tag_ids:=[1])
NOTE3_ID=$(echo "$NOTE3_RESPONSE" | grep -o '[0-9]\+' | head -1)
echo "Created note with ID: $NOTE3_ID"
echo

# Test 4: Create fourth note (for deletion test)
echo "4. Creating fourth note for deletion..."
NOTE4_RESPONSE=$(http --session=$SESSION_NAME POST $BASE_URL/notes reference_type="programs" reference_id:=1 value="This note will be deleted" tag_ids:=[1])
NOTE4_ID=$(echo "$NOTE4_RESPONSE" | grep -o '[0-9]\+' | head -1)
echo "Created note with ID: $NOTE4_ID"
echo

# Test 5: Update second note
echo "5. Updating second note..."
http --session=$SESSION_NAME PATCH $BASE_URL/notes/$NOTE2_ID value="This is a note updated for the endpoint"
echo "Updated note $NOTE2_ID"
echo

# Test 6: Delete fourth note
echo "6. Deleting fourth note..."
http --session=$SESSION_NAME DELETE $BASE_URL/notes/$NOTE4_ID
echo "Deleted note $NOTE4_ID"
echo

# Test 7: List all notes
echo "7. Listing all notes..."
http --session=$SESSION_NAME GET $BASE_URL/notes
echo

# Test 8: List notes by reference type
echo "8. Listing notes by reference type (programs)..."
http --session=$SESSION_NAME GET $BASE_URL/notes reference_type=programs
echo

# Test 9: List notes by reference ID
echo "9. Listing notes by reference ID (program ID 1)..."
http --session=$SESSION_NAME GET $BASE_URL/notes reference_id=1
echo

# Test 10: Search notes
echo "10. Searching notes for 'created'..."
http --session=$SESSION_NAME GET $BASE_URL/notes search=created
echo

# Test 11: Get specific note details
echo "11. Getting note details for ID $NOTE1_ID..."
http --session=$SESSION_NAME GET $BASE_URL/notes/$NOTE1_ID
echo

echo "=== Note Smoke Test Completed ==="
