# Smoke Tests for RequesterBackend API

This directory contains bash scripts for smoke testing the RequesterBackend API using httpie.

## Prerequisites

1. **httpie** - Install with: `pip install httpie` or `brew install httpie`
2. **API Server** - The RequesterBackend API must be running on `http://localhost:8081`
3. **Database** - The API must be connected to a database with proper schema

## Usage

### Run All Tests
```bash
./run_all_tests.sh
```

### Run Individual Tests
```bash
./1_tag.sh      # Test Tag CRUD operations
./2_program.sh  # Test Program CRUD operations
./3_endpoint.sh # Test Endpoint CRUD operations
./4_vuln.sh     # Test Vulnerability CRUD operations
./5_note.sh     # Test Note CRUD operations
./6_request.sh  # Test Request read operations
```

## Test Structure

Each test script follows this pattern:
1. **Create** - Creates 2-3 resources with meaningful data
2. **Update** - Updates one of the created resources
3. **Delete** - Deletes one of the created resources
4. **List** - Lists all resources
5. **Get Detail** - Gets details of a specific resource

## Test Data

The tests use descriptive names and data to make it easy to verify that operations worked:
- Resources created have "created" in their names/descriptions
- Resources updated have "updated" in their names/descriptions
- Resources deleted are marked as "to be deleted"

## Authentication

All tests use session-based authentication with session name "a". The master script automatically starts a session before running tests.

## Notes

- Tests are designed to run in order (1-6) as later tests may depend on data created by earlier tests
- Request tests focus on read operations since requests are typically created via import operations
- Import operations (HAR/Burp XML) are not tested as they require file uploads
- The `start-session` endpoint is not tested as it's used for authentication setup
