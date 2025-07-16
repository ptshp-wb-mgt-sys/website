# Manual API Testing Guide

This guide will help you test the Pet Management API endpoints manually.

## Prerequisites

1. **Set up environment variables**:

   ```bash
   cp env.example .env
   # Edit .env with your Supabase credentials
   ```

2. **Get your Supabase credentials**:
   - Go to your Supabase project dashboard
   - Settings → API
   - Copy the Project URL and anon/service key
   - Settings → JWT Settings → Copy the JWT secret

3. **Get a JWT token**:
   - You'll need a valid JWT token from Supabase auth
   - This can be obtained by signing up/signing in through your frontend

## Running the API

```bash
# Build the API
go build -o pet-mgt-api cmd/api/main.go

# Run the API
./pet-mgt-api
```

The API will start on `http://localhost:3000`

## Testing Endpoints

### 1. Health Check (No Auth Required)

```bash
curl http://localhost:3000/ping
```

### 2. Get Current User Profile

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:3000/api/v1/profile
```

### 3. Create Client Profile

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "+1234567890",
    "address": "123 Main St, City, State",
    "role": "client"
  }' \
  http://localhost:3000/api/v1/users
```

### 4. Create Veterinarian Profile

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Dr. Sarah Smith",
    "email": "dr.sarah@vetclinic.com",
    "phone": "+1987654321",
    "clinic_address": "456 Oak Ave, City, State",
    "role": "veterinarian"
  }' \
  http://localhost:3000/api/v1/users
```

### 5. Get User by ID

```bash
# Replace USER_ID with the actual user ID
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:3000/api/v1/users/USER_ID
```

### 6. Update User

```bash
# Replace USER_ID with the actual user ID
curl -X PUT \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "John Doe Updated",
    "email": "john.updated@example.com",
    "phone": "+1234567890",
    "address": "456 New St, City, State"
  }' \
  http://localhost:3000/api/v1/users/USER_ID
```

### 7. Delete User

```bash
# Replace USER_ID with the actual user ID
curl -X DELETE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:3000/api/v1/users/USER_ID
```

### 8. Create Pet

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Buddy",
    "type": "Dog",
    "breed": "Golden Retriever",
    "date_of_birth": "2020-01-15T00:00:00Z",
    "weight": 25.5,
    "owner_id": "CLIENT_USER_ID"
  }' \
  http://localhost:3000/api/v1/pets
```

### 9. Get Pet by ID

```bash
# Replace PET_ID with the actual ID from the create response
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:3000/api/v1/pets/PET_ID
```

### 10. Update Pet

```bash
# Replace PET_ID with the actual ID
curl -X PUT \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Buddy Updated",
    "type": "Dog",
    "breed": "Golden Retriever",
    "date_of_birth": "2020-01-15T00:00:00Z",
    "weight": 26.0
  }' \
  http://localhost:3000/api/v1/pets/PET_ID
```

### 11. Delete Pet

```bash
# Replace PET_ID with the actual ID
curl -X DELETE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:3000/api/v1/pets/PET_ID
```

### 12. Get Pets by Client

```bash
# Replace CLIENT_ID with the actual client ID
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:3000/api/v1/clients/CLIENT_ID/pets
```

### 13. Create Medical Record

```bash
# Replace PET_ID with the actual pet ID
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "date_of_visit": "2024-01-15T10:00:00Z",
    "reason_for_visit": "Annual checkup",
    "diagnosis": "Healthy",
    "medication_prescribed": ["Vitamin D supplement"],
    "notes": "Pet is in good health"
  }' \
  http://localhost:3000/api/v1/pets/PET_ID/medical-records
```

### 14. Get Medical Records for Pet

```bash
# Replace PET_ID with the actual pet ID
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:3000/api/v1/pets/PET_ID/medical-records
```

### 15. Get Medical Record by ID

```bash
# Replace RECORD_ID with the actual medical record ID
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:3000/api/v1/medical-records/RECORD_ID
```

### 16. Update Medical Record

```bash
# Replace RECORD_ID with the actual medical record ID
curl -X PUT \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "date_of_visit": "2024-01-15T10:00:00Z",
    "reason_for_visit": "Follow-up checkup",
    "diagnosis": "Recovering well",
    "medication_prescribed": ["Vitamin D supplement", "Pain relief"],
    "notes": "Pet is recovering well from previous treatment"
  }' \
  http://localhost:3000/api/v1/medical-records/RECORD_ID
```

### 17. Delete Medical Record

```bash
# Replace RECORD_ID with the actual medical record ID
curl -X DELETE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:3000/api/v1/medical-records/RECORD_ID
```

### 18. List Users (Admin Only)

```bash

curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  "http://localhost:3000/api/v1/users?limit=10&offset=0"
```

## Expected Responses

### Success Response

```json
{
  "success": true,
  "data": {
    // Response data here
  }
}
```

### Error Response

```json
{
  "error": "Error message description"
}
```

## Common HTTP Status Codes

- `200` - Success
- `400` - Bad Request (invalid input)
- `401` - Unauthorized (missing or invalid token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `500` - Internal Server Error

## Testing Tips

1. **Start with the health check** to ensure the API is running
2. **Create a user profile first** before testing other endpoints
3. **Save the IDs** from create responses to use in subsequent requests
4. **Test with different user roles** (client, veterinarian, admin) to verify authorization
5. **Check the response status codes** and error messages for debugging
6. **Note the authorization rules**:
   - Clients can only access their own pets and medical records
   - Veterinarians can access any pet's medical records but cannot create/update pets
   - Admins have full access to all resources

## Using the Test Script

You can also use the automated test script:

1. Edit `test_api.sh` and set your JWT token
2. Run: `./test_api.sh`

This will test all endpoints automatically and show colored output for success/failure.
