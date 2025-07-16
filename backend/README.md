# Pet Management System Backend API

A Go-based backend API for a comprehensive pet-shop and veterinary clinic management system. This API provides robust endpoints for managing users, pets, and medical records with proper role-based access control.

## Features

- **User Management**: Support for Clients (pet owners), Veterinarians, and Admins
- **Pet Management**: Complete CRUD operations for pet records
- **Medical Records**: Veterinary visit tracking and medical history
- **Role-Based Access Control**: Secure authorization based on user roles
- **JWT Authentication**: Integration with Supabase authentication
- **RESTful API**: Clean, consistent API design

## Architecture

- **Language**: Go 1.24.4
- **Framework**: Chi router
- **Database**: Supabase (PostgreSQL)
- **Authentication**: JWT tokens from Supabase
- **Middleware**: CORS, rate limiting, authentication, database injection

## Project Structure

```plaintext
backend/
├── cmd/api/main.go              # Application entry point
├── internal/
│   ├── config/config.go         # Configuration management
│   ├── handlers/                # HTTP request handlers
│   │   ├── response.go          # Response utilities
│   │   ├── user.go              # User profile handler
│   │   ├── users.go             # User management handlers
│   │   ├── pets.go              # Pet management handlers
│   │   └── medical_records.go   # Medical record handlers
│   ├── middleware/              # HTTP middleware
│   │   ├── auth.go              # JWT authentication
│   │   ├── cors.go              # CORS handling
│   │   └── db.go                # Database injection
│   ├── routes/routes.go         # Route definitions
│   └── store/                   # Data access layer
│       ├── db.go                # Supabase client
│       └── models.go            # Data models
├── database/schema/tables.sql   # Database schema
└── README.md                    # This file
```

## User Roles & Permissions

### Admin

- Full access to all data and operations
- Can create, read, update, delete any user, pet, or medical record
- Can list all users with pagination

### Client (Pet Owner)

- Can manage their own profile
- Can create, read, update, delete their own pets
- Can view medical records for their own pets
- Cannot create or modify medical records

### Veterinarian

- Can manage their own profile
- Can view any pet's details (for medical purposes)
- Can create, read, update, delete medical records for any pet
- Cannot modify pet details or create pets

## API Endpoints

### Authentication

All endpoints require a valid JWT token in the Authorization header:

```bash
Authorization: Bearer <jwt_token>
```

### User Management

#### Create User Profile

```bash
POST /api/v1/users
```

Creates a new user profile (Client or Veterinarian) after Supabase authentication.

**Request Body for Client:**

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+1234567890",
  "address": "123 Main St",
  "role": "client"
}
```

**Request Body for Veterinarian:**

```json
{
  "name": "Dr. Sarah Smith",
  "email": "dr.sarah@vetclinic.com",
  "phone": "+1987654321",
  "role": "veterinarian"
}
```

**Authorization:** Users can only create their own profile type (matching their JWT role), admins can create any type.

#### Get User Profile

```bash
GET /api/v1/users/{id}
```

Retrieves a user profile by ID.

**Authorization:** Users can only access their own profile, admins can access any profile.

#### Update User Profile

```bash
PUT /api/v1/users/{id}
```

Updates a user profile.

**Request Body:**

```json
{
  "name": "John Doe Updated",
  "email": "john.updated@example.com",
  "phone": "+1234567890",
  "address": "456 Oak St"
}
```

**Authorization:** Users can only update their own profile, admins can update any profile.

#### Delete User Profile

```bash
DELETE /api/v1/users/{id}
```

Deletes a user profile.

**Authorization:** Admins only.

#### List Users

```bash
GET /api/v1/users?limit=10&offset=0
```

Lists all users with pagination.

**Query Parameters:**

- `limit` (optional): Number of users to return (default: 10)
- `offset` (optional): Number of users to skip (default: 0)

**Authorization:** Admins only.

#### Get Current User Profile

```bash
GET /api/v1/profile

```

Returns the authenticated user's profile information.

### Pet Management

#### Create Pet

```bash
POST /api/v1/pets

```

Creates a new pet record.

**Request Body:**

```json
{
  "name": "Buddy",
  "type": "Dog",
  "breed": "Golden Retriever",
  "date_of_birth": "2020-01-15T00:00:00Z",
  "weight": 25.5,
  "owner_id": "client-user-id"
}
```

**Authorization:** Clients can only create pets for themselves, admins can create for anyone. Veterinarians cannot create pets.

#### Get Pet

```bash
GET /api/v1/pets/{id}
```

Retrieves a pet by ID.

**Authorization:** Clients can only access their own pets, vets and admins can access any pet.

#### Update Pet

```bash
PUT /api/v1/pets/{id}
```

Updates a pet record.

**Request Body:**

```json
{
  "name": "Buddy Updated",
  "type": "Dog",
  "breed": "Golden Retriever",
  "date_of_birth": "2020-01-15T00:00:00Z",
  "weight": 26.0
}
```

**Authorization:** Clients can only update their own pets, admins can update any pet. Veterinarians cannot update pets.

#### Delete Pet

```bash
DELETE /api/v1/pets/{id}
```

Deletes a pet record.

**Authorization:** Clients can only delete their own pets, admins can delete any pet. Veterinarians cannot delete pets.

#### Get Pets by Client

```bash
GET /api/v1/clients/{clientId}/pets
```

Retrieves all pets for a specific client.

**Authorization:** Clients can only access their own pets, vets and admins can access any client's pets.

### Medical Records

#### Create Medical Record

```bash
POST /api/v1/pets/{petId}/medical-records
```

Creates a new medical record for a pet.

**Request Body:**

```json
{
  "date_of_visit": "2024-01-15T10:00:00Z",
  "reason_for_visit": "Annual checkup",
  "diagnosis": "Healthy",
  "medication_prescribed": ["Vitamin D supplement"],
  "notes": "Pet is in good health"
}
```

**Authorization:** Veterinarians and admins only.

#### Get Medical Records for Pet

```bash
GET /api/v1/pets/{petId}/medical-records
```

Retrieves all medical records for a specific pet.

**Authorization:** Clients can only access their own pet's records, vets and admins can access any pet's records.

#### Get Medical Record

```bash
GET /api/v1/medical-records/{id}
```

Retrieves a specific medical record.

**Authorization:** Clients can only access their own pet's records, vets and admins can access any record.

#### Update Medical Record

```bash
PUT /api/v1/medical-records/{id}
```

Updates a medical record.

**Request Body:**

```json
{
  "date_of_visit": "2024-01-15T10:00:00Z",
  "reason_for_visit": "Annual checkup",
  "diagnosis": "Healthy with minor dental issues",
  "medication_prescribed": ["Vitamin D supplement", "Dental cleaning"],
  "notes": "Pet is in good health, recommend dental cleaning"
}
```

**Authorization:** Only the veterinarian who created the record or admins can update it.

#### Delete Medical Record

```bash
DELETE /api/v1/medical-records/{id}
```

Deletes a medical record.

**Authorization:** Only the veterinarian who created the record or admins can delete it.

## Database Schema

The system uses the following tables in Supabase:

- `clients` - Pet owner profiles
- `veterinarians` - Veterinarian profiles  
- `pets` - Pet records
- `medical_records` - Veterinary visit records

See `database/schema/tables.sql` for the complete schema definition.

## Environment Variables

Create a `.env` file in the backend directory:

```env
PORT=3000
ENV=development
FRONTEND_URL=http://localhost:5173
SUPABASE_URL=your_supabase_url
SUPABASE_SERVICE_KEY=your_supabase_service_key
SUPABASE_JWT_SECRET=your_supabase_jwt_secret
```

## Running the Application

1. **Install dependencies:**

   ```bash
   go mod download
   ```

2. **Set up environment variables:**

   ```bash
   cp .env.example .env
   # Edit .env with your Supabase credentials
   ```

3. **Run the application:**

   ```bash
   go run cmd/api/main.go
   ```

4. **Test the health endpoint:**

   ```bash
   curl http://localhost:3000/ping
   ```

## Error Responses

The API returns consistent error responses:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:

- `200` - Success
- `400` - Bad Request (invalid input)
- `401` - Unauthorized (missing or invalid token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `500` - Internal Server Error

## Success Responses

Successful operations return:

```json
{
  "success": true,
  "data": {
    // Response data
  }
}
```

## Rate Limiting

- Public endpoints: 60 requests per minute per IP
- Protected endpoints: 100 requests per minute per IP

## Security Features

- JWT token validation
- Role-based access control
- Input validation
- Rate limiting
- CORS protection
- No Row Level Security (RLS) - authorization handled in Go backend

## Development

### Adding New Endpoints

1. Create handler functions in the appropriate `handlers/` file
2. Add routes in `routes/routes.go`
3. Update this README with endpoint documentation

### Database Changes

1. Update the schema in `database/schema/tables.sql`
2. Update models in `internal/store/models.go`
3. Update database methods in `internal/store/db.go`

## Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## Deployment

The application can be deployed to any platform that supports Go applications. Ensure all environment variables are properly configured in your deployment environment.
