#!/bin/bash

# Test script for the authentication API
# Usage: ./scripts/test-api.sh [API_URL]

API_URL="${1:-http://localhost:8080}"
EMAIL="test-$(date +%s)@example.com"
PASSWORD="SecurePassword123!"

echo "Testing API at: $API_URL"
echo "=============================="

# Test 1: Health Check
echo -e "\n1. Health Check"
curl -s "${API_URL}/api/health" | jq '.' || echo "Health check failed"

# Test 2: Register User
echo -e "\n2. Register User: $EMAIL"
REGISTER_RESPONSE=$(curl -s -X POST "${API_URL}/api/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")
echo "$REGISTER_RESPONSE" | jq '.'

# Test 3: Login
echo -e "\n3. Login with registered user"
LOGIN_RESPONSE=$(curl -s -X POST "${API_URL}/api/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" \
  -c cookies.txt)
echo "$LOGIN_RESPONSE" | jq '.'

# Test 4: Login with wrong password
echo -e "\n4. Login with wrong password (should fail)"
curl -s -X POST "${API_URL}/api/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"WrongPassword\"}" | jq '.'

# Test 5: Logout
echo -e "\n5. Logout"
LOGOUT_RESPONSE=$(curl -s -X POST "${API_URL}/api/logout" \
  -b cookies.txt)
echo "$LOGOUT_RESPONSE" | jq '.'

# Test 6: Try logout without session (should fail)
echo -e "\n6. Logout without session (should fail)"
curl -s -X POST "${API_URL}/api/logout" | jq '.'

# Cleanup
rm -f cookies.txt

echo -e "\n=============================="
echo "Tests completed!"
