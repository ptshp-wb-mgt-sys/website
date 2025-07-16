#!/bin/bash

# Pet Management API Test Script
# Make sure to set your JWT token and API base URL

# Configuration
API_BASE_URL="http://localhost:3000/api/v1"
JWT_TOKEN="your_jwt_token_here"  # Replace with actual JWT token

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to make API calls
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo -e "${YELLOW}Testing: $description${NC}"
    echo "Endpoint: $method $API_BASE_URL$endpoint"
    
    if [ -n "$data" ]; then
        echo "Data: $data"
        response=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
            -X "$method" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $JWT_TOKEN" \
            -d "$data" \
            "$API_BASE_URL$endpoint")
    else
        response=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
            -X "$method" \
            -H "Authorization: Bearer $JWT_TOKEN" \
            "$API_BASE_URL$endpoint")
    fi
    
    # Extract status code
    http_status=$(echo "$response" | grep "HTTP_STATUS:" | cut -d: -f2)
    response_body=$(echo "$response" | sed '/HTTP_STATUS:/d')
    
    echo "Status: $http_status"
    echo "Response: $response_body"
    
    if [ "$http_status" -ge 200 ] && [ "$http_status" -lt 300 ]; then
        echo -e "${GREEN}✓ Success${NC}"
    else
        echo -e "${RED}✗ Failed${NC}"
    fi
    echo "----------------------------------------"
}

# Check if JWT token is set
if [ "$JWT_TOKEN" = "your_jwt_token_here" ]; then
    echo -e "${RED}Error: Please set your JWT token in the script${NC}"
    echo "Replace 'your_jwt_token_here' with your actual JWT token"
    exit 1
fi

echo -e "${GREEN}Starting API Tests...${NC}"
echo "API Base URL: $API_BASE_URL"
echo "========================================"

# Test 1: Health Check (Public endpoint)
echo -e "${YELLOW}Test 1: Health Check${NC}"
curl -s -w "\nHTTP_STATUS:%{http_code}" http://localhost:3000/ping
echo -e "\n----------------------------------------"

# Test 2: Get Current User Profile
make_request "GET" "/profile" "" "Get Current User Profile"

# Test 3: Create Client Profile
CLIENT_DATA='{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "+1234567890",
    "address": "123 Main St, City, State",
    "role": "client"
}'
make_request "POST" "/users" "$CLIENT_DATA" "Create Client Profile"

# Test 4: Create Veterinarian Profile
VET_DATA='{
    "name": "Dr. Sarah Smith",
    "email": "dr.sarah@vetclinic.com",
    "phone": "+1987654321",
    "clinic_address": "456 Oak Ave, City, State",
    "role": "veterinarian"
}'
make_request "POST" "/users" "$VET_DATA" "Create Veterinarian Profile"

# Test 5: Get User by ID (you'll need to replace USER_ID with actual ID from previous response)
# make_request "GET" "/users/USER_ID" "" "Get User by ID"

# Test 6: Update User (you'll need to replace USER_ID with actual ID)
USER_UPDATE_DATA='{
    "name": "John Doe Updated",
    "email": "john.updated@example.com",
    "phone": "+1234567890",
    "address": "456 New St, City, State"
}'
# make_request "PUT" "/users/USER_ID" "$USER_UPDATE_DATA" "Update User"

# Test 7: Create Pet
PET_DATA='{
    "name": "Buddy",
    "type": "Dog",
    "breed": "Golden Retriever",
    "date_of_birth": "2020-01-15T00:00:00Z",
    "weight": 25.5,
    "owner_id": "CLIENT_USER_ID"
}'
make_request "POST" "/pets" "$PET_DATA" "Create Pet"

# Test 8: Get Pet by ID (you'll need to replace PET_ID with actual ID from previous response)
# make_request "GET" "/pets/PET_ID" "" "Get Pet by ID"

# Test 9: Update Pet (you'll need to replace PET_ID with actual ID)
PET_UPDATE_DATA='{
    "name": "Buddy Updated",
    "type": "Dog",
    "breed": "Golden Retriever",
    "date_of_birth": "2020-01-15T00:00:00Z",
    "weight": 26.0
}'
# make_request "PUT" "/pets/PET_ID" "$PET_UPDATE_DATA" "Update Pet"

# Test 10: Delete Pet (you'll need to replace PET_ID with actual ID)
# make_request "DELETE" "/pets/PET_ID" "" "Delete Pet"

# Test 11: Get Pets by Client (you'll need to replace CLIENT_ID with actual client ID)
# make_request "GET" "/clients/CLIENT_ID/pets" "" "Get Pets by Client"

# Test 12: Create Medical Record (you'll need to replace PET_ID with actual pet ID)
MEDICAL_RECORD_DATA='{
    "date_of_visit": "2024-01-15T10:00:00Z",
    "reason_for_visit": "Annual checkup",
    "diagnosis": "Healthy",
    "medication_prescribed": ["Vitamin D supplement"],
    "notes": "Pet is in good health"
}'
# make_request "POST" "/pets/PET_ID/medical-records" "$MEDICAL_RECORD_DATA" "Create Medical Record"

# Test 13: Get Medical Records for Pet (you'll need to replace PET_ID with actual pet ID)
# make_request "GET" "/pets/PET_ID/medical-records" "" "Get Medical Records for Pet"

# Test 14: Get Medical Record by ID (you'll need to replace RECORD_ID with actual medical record ID)
# make_request "GET" "/medical-records/RECORD_ID" "" "Get Medical Record by ID"

# Test 15: Update Medical Record (you'll need to replace RECORD_ID with actual medical record ID)
MEDICAL_RECORD_UPDATE_DATA='{
    "date_of_visit": "2024-01-15T10:00:00Z",
    "reason_for_visit": "Follow-up checkup",
    "diagnosis": "Recovering well",
    "medication_prescribed": ["Vitamin D supplement", "Pain relief"],
    "notes": "Pet is recovering well from previous treatment"
}'
# make_request "PUT" "/medical-records/RECORD_ID" "$MEDICAL_RECORD_UPDATE_DATA" "Update Medical Record"

# Test 16: Delete Medical Record (you'll need to replace RECORD_ID with actual medical record ID)
# make_request "DELETE" "/medical-records/RECORD_ID" "" "Delete Medical Record"

# Test 17: List Users (Admin only)
make_request "GET" "/users?limit=10&offset=0" "" "List Users"

echo -e "${GREEN}API Tests Completed!${NC}"
echo ""
echo -e "${YELLOW}Note: Some tests are commented out because they require IDs from previous responses.${NC}"
echo "To test those endpoints:"
echo "1. Run the create endpoints first"
echo "2. Extract the IDs from the responses"
echo "3. Uncomment and update the relevant test lines with actual IDs"
echo ""
echo -e "${YELLOW}Authorization Rules:${NC}"
echo "- Clients can only access their own pets and medical records"
echo "- Veterinarians can access any pet's medical records but cannot create/update pets"
echo "- Admins have full access to all resources" 
