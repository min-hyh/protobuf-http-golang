#!/bin/bash

# Test script to demonstrate error handling middleware
BASE_URL="http://localhost:8080"

echo "Testing Error Handling Middleware"
echo "================================="
echo ""

# Test 1: Valid request
echo "1. Testing valid request:"
curl -s -X GET "$BASE_URL/v1/get-param-in-body/test-id?content=test-content" | jq .
echo ""

# Test 2: Missing required field (id)
echo "2. Testing missing required field (id):"
curl -s -X GET "$BASE_URL/v1/get-param-in-body/?content=test-content" | jq .
echo ""

# Test 3: Missing required field (content)
echo "3. Testing missing required field (content):"
curl -s -X GET "$BASE_URL/v1/get-param-in-body/test-id" | jq .
echo ""

# Test 4: Resource not found
echo "4. Testing resource not found:"
curl -s -X GET "$BASE_URL/v1/get-param-in-body/not-found?content=test-content" | jq .
echo ""

# Test 5: Permission denied
echo "5. Testing permission denied:"
curl -s -X GET "$BASE_URL/v1/get-param-in-body/unauthorized?content=test-content" | jq .
echo ""

# Test 6: Internal server error
echo "6. Testing internal server error:"
curl -s -X GET "$BASE_URL/v1/get-param-in-body/error?content=test-content" | jq .
echo ""

# Test 7: Missing required header
echo "7. Testing missing required header:"
curl -s -X GET "$BASE_URL/v1/get-param-in-header?content=test-content" | jq .
echo ""

# Test 8: Invalid authentication token
echo "8. Testing invalid authentication token:"
curl -s -H "X-Custom-Header-Id: invalid-token" -X GET "$BASE_URL/v1/get-param-in-header?content=test-content" | jq .
echo ""

# Test 9: Valid header request
echo "9. Testing valid header request:"
curl -s -H "X-Custom-Header-Id: valid-token" -X GET "$BASE_URL/v1/get-param-in-header?content=test-content" | jq .
echo ""

# Test 10: Duplicate resource (POST)
echo "10. Testing duplicate resource (POST):"
curl -s -X POST "$BASE_URL/v1/post/unstructured-data" \
  -H "Content-Type: application/json" \
  -d '{"id": "duplicate", "data": {"@type": "type.googleapis.com/google.protobuf.StringValue", "value": "test"}}' | jq .
echo ""

# Test 11: Rate limit exceeded
echo "11. Testing rate limit exceeded:"
curl -s -X POST "$BASE_URL/v1/post/unstructured-data" \
  -H "Content-Type: application/json" \
  -d '{"id": "rate-limit", "data": {"@type": "type.googleapis.com/google.protobuf.StringValue", "value": "test"}}' | jq .
echo ""

# Test 12: Valid POST request
echo "12. Testing valid POST request:"
curl -s -X POST "$BASE_URL/v1/post/unstructured-data" \
  -H "Content-Type: application/json" \
  -d '{"id": "valid-id", "data": {"@type": "type.googleapis.com/google.protobuf.StringValue", "value": "test"}}' | jq .
echo ""

# Test 13: Request with X-Request-ID header
echo "13. Testing request with X-Request-ID header:"
curl -s -H "X-Request-ID: test-request-123" -X GET "$BASE_URL/v1/get-param-in-body/test-id?content=test-content" | jq .
echo ""

echo "Error handling tests completed!" 