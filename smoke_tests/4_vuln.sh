#!/bin/bash

# Smoke test for Vulnerability resource
# This script tests CRUD operations for vulnerabilities

BASE_URL="http://localhost:8081"
SESSION_NAME="a"

echo "=== Vulnerability Smoke Test ==="
echo "Testing Vulnerability CRUD operations..."
echo

# Test 1: Create first vulnerability
echo "1. Creating first vulnerability..."
VULN1_RESPONSE=$(http --session=$SESSION_NAME POST $BASE_URL/vulns title="SQL Injection created" body="This vulnerability allows attackers to inject malicious SQL queries through user input. The application does not properly sanitize input before constructing database queries." tag_ids:=[1])
VULN1_ID=$(echo "$VULN1_RESPONSE" | grep -o '[0-9]\+' | head -1)
echo "Created vulnerability with ID: $VULN1_ID"
echo

# Test 2: Create second vulnerability
echo "2. Creating second vulnerability..."
VULN2_RESPONSE=$(http --session=$SESSION_NAME POST $BASE_URL/vulns title="XSS Vulnerability created" body="Cross-site scripting vulnerability found in the login form. User input is not properly escaped before being displayed in the response." tag_ids:=[1,2])
VULN2_ID=$(echo "$VULN2_RESPONSE" | grep -o '[0-9]\+' | head -1)
echo "Created vulnerability with ID: $VULN2_ID"
echo

# Test 3: Create third vulnerability (for deletion test)
echo "3. Creating third vulnerability..."
VULN3_RESPONSE=$(http --session=$SESSION_NAME POST $BASE_URL/vulns title="CSRF Vulnerability created" body="Cross-site request forgery vulnerability allows attackers to perform actions on behalf of authenticated users without their knowledge." tag_ids:=[1])
VULN3_ID=$(echo "$VULN3_RESPONSE" | grep -o '[0-9]\+' | head -1)
echo "Created vulnerability with ID: $VULN3_ID"
echo

# Test 4: Create child vulnerability
echo "4. Creating child vulnerability..."
VULN4_RESPONSE=$(http --session=$SESSION_NAME POST $BASE_URL/vulns title="SQL Injection - Authentication Bypass" body="This is a child vulnerability of the main SQL injection issue, specifically related to authentication bypass." parent_id:=$VULN1_ID tag_ids:=[1])
VULN4_ID=$(echo "$VULN4_RESPONSE" | grep -o '[0-9]\+' | head -1)
echo "Created child vulnerability with ID: $VULN4_ID"
echo

# Test 5: Update second vulnerability
echo "5. Updating second vulnerability..."
http --session=$SESSION_NAME PUT $BASE_URL/vulns/$VULN2_ID title="XSS Vulnerability updated" body="Cross-site scripting vulnerability found in the login form. User input is not properly escaped before being displayed in the response. This has been updated with additional details." tag_ids:=[1]
echo "Updated vulnerability $VULN2_ID"
echo

# Test 6: Delete third vulnerability
echo "6. Deleting third vulnerability..."
http --session=$SESSION_NAME DELETE $BASE_URL/vulns/$VULN3_ID
echo "Deleted vulnerability $VULN3_ID"
echo

# Test 7: List all vulnerabilities
echo "7. Listing all vulnerabilities..."
http --session=$SESSION_NAME GET $BASE_URL/vulns
echo

# Test 8: Get specific vulnerability details
echo "8. Getting vulnerability details for ID $VULN1_ID..."
http --session=$SESSION_NAME GET $BASE_URL/vulns/$VULN1_ID
echo

# Test 9: Get vulnerability by slug
echo "9. Getting vulnerability by slug..."
VULN1_SLUG=$(http --session=$SESSION_NAME GET $BASE_URL/vulns/$VULN1_ID | grep -o '"slug":"[^"]*"' | cut -d'"' -f4)
if [ ! -z "$VULN1_SLUG" ]; then
    echo "Getting vulnerability by slug: $VULN1_SLUG"
    http --session=$SESSION_NAME GET $BASE_URL/vulns/slug/$VULN1_SLUG
else
    echo "Could not extract slug from vulnerability details"
fi
echo

echo "=== Vulnerability Smoke Test Completed ==="
