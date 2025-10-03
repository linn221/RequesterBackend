#!/bin/bash

# Master script to run all smoke tests
# This script runs all smoke tests in the correct order

BASE_URL="http://localhost:8081"
SESSION_NAME="a"

echo "=========================================="
echo "    RequesterBackend Smoke Tests"
echo "=========================================="
echo

# Check if the API is running
echo "Checking if API is running..."
if ! curl -s "$BASE_URL/start_session?secret=super_secret_key" > /dev/null; then
    echo "ERROR: API is not running at $BASE_URL"
    echo "Please start the API server first"
    exit 1
fi
echo "API is running ✓"
echo

# Start a session
echo "Starting session..."
http --session=$SESSION_NAME GET "$BASE_URL/start_session?secret=super_secret_key"
echo

# Run tests in order
echo "Running smoke tests..."
echo

echo "=========================================="
echo "Running Tag Tests..."
echo "=========================================="
./1_tag.sh
echo

echo "=========================================="
echo "Running Program Tests..."
echo "=========================================="
./2_program.sh
echo

echo "=========================================="
echo "Running Endpoint Tests..."
echo "=========================================="
./3_endpoint.sh
echo

echo "=========================================="
echo "Running Vulnerability Tests..."
echo "=========================================="
./4_vuln.sh
echo

echo "=========================================="
echo "Running Note Tests..."
echo "=========================================="
./5_note.sh
echo

echo "=========================================="
echo "Running Request Tests..."
echo "=========================================="
./6_request.sh
echo

echo "=========================================="
echo "    All Smoke Tests Completed!"
echo "=========================================="
