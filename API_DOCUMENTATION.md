# Carpooling API Documentation

## Overview
This is a REST API for carpooling trip management with simplified balance tracking. The API supports:
- User management with integrated balance tracking
- Trip creation with automatic date/time assignment
- Real-time balance updates and payment processing
- Simple cost splitting among participants

## Base URL
```
http://localhost:8081
```

## Authentication
Currently, the API does not require authentication (can be added later).

## Common Response Formats

### Success Response
```json
{
  "status": "success",
  "data": { ... }
}
```

### Error Response
```json
{
  "error": "Error message",
  "validations": [
    {
      "field": "field_name",
      "message": "Validation error message"
    }
  ]
}
```

## Endpoints

### System Endpoints

#### Health Check
**GET** `/health`

Check if the API and database are running.

**Response (200 OK):**
```json
{
  "status": "healthy",
  "message": "Carpooling API is running",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0"
}
```

#### API Info
**GET** `/`

Get basic API information and available endpoints.

**Response (200 OK):**
```json
{
  "name": "Carpooling API",
  "version": "1.0.0",
  "description": "REST API for carpooling trip management",
  "endpoints": {
    "health": "/health",
    "users": "/users",
    "trips": "/trips",
    "balances": "/balances"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

### User Management

#### Create User
**POST** `/users`

Create a new user.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+1234567890"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+1234567890",
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Error Response (409 Conflict):**
```json
{
  "error": "User with this email already exists"
}
```

#### Get User
**GET** `/users/{id}`

Get a user by ID.

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+1234567890",
  "created_at": "2024-01-15T10:30:00Z"
}
```

#### Get User by Name
**GET** `/users/by-name/{name}`

Get a user by name.

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+1234567890",
  "created_at": "2024-01-15T10:30:00Z"
}
```

#### List Users
**GET** `/users?page=1&limit=10`

List all users with pagination.

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10, max: 100)

**Response (200 OK):**
```json
{
  "users": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "phone": "+1234567890",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 10
}
```

#### Update User
**PUT** `/users/{id}`

Update an existing user.

**Request Body:**
```json
{
  "name": "John Smith",
  "email": "johnsmith@example.com",
  "phone": "+1234567891"
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "John Smith",
  "email": "johnsmith@example.com",
  "phone": "+1234567891",
  "created_at": "2024-01-15T10:30:00Z"
}
```

#### Delete User
**DELETE** `/users/{id}`

Delete a user.

**Response (200 OK):**
```json
{
  "message": "User deleted successfully"
}
```

---

### Trip Management

#### Create Trip
**POST** `/trips`

Create a trip with participants. Date and time are automatically set to the current date/time.

**Request Body:**
```json
{
  "driver_name": "Cyril",
  "total_cost": 145,
  "participants": [
    "Cyril",
    "Abhinav",
    "Vishal"
  ]
}
```

**Response (201 Created):**
```json
{
  "trip_id": 1,
  "driver_id": 1,
  "driver_name": "Cyril",
  "date": "2024-01-15",
  "time": "09:30:45",
  "total_cost": 145,
  "participants": [
    {
      "user_id": 1,
      "name": "Cyril",
      "share_amount": 48.33
    },
    {
      "user_id": 2,
      "name": "Abhinav",
      "share_amount": 48.33
    },
    {
      "user_id": 3,
      "name": "Vishal",
      "share_amount": 48.34
    }
  ],
  "message": "Trip created successfully with 3 participants. Total cost: 145.00, Share per person: 48.33",
  "created_at": "2024-01-15T09:30:45Z"
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "participant 'John' not found"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Validation failed",
  "validations": [
    {
      "field": "driver_name",
      "message": "Driver name is required"
    },
    {
      "field": "participants",
      "message": "Driver must be included in participants list"
    }
  ]
}
```

#### Get Trip
**GET** `/trips/{id}`

Get trip details by ID.

**Response (200 OK):**
```json
{
  "trip_id": 1,
  "driver_id": 1,
  "driver_name": "Cyril",
  "date": "2024-01-15",
  "time": "09:30:45",
  "total_cost": 145,
  "participants": [
    {
      "user_id": 1,
      "name": "Cyril",
      "share_amount": 48.33
    },
    {
      "user_id": 2,
      "name": "Abhinav",
      "share_amount": 48.33
    },
    {
      "user_id": 3,
      "name": "Vishal",
      "share_amount": 48.34
    }
  ],
  "message": "Trip retrieved successfully",
  "created_at": ""
}
```

---

### Balance Management

#### Get All User Balances
**GET** `/balances`

Get current balances for all users who owe money.

**Response (200 OK):**
```json
{
  "balances": [
    {
      "user_id": 2,
      "name": "Abhinav",
      "total_owed": 96.66
    },
    {
      "user_id": 3,
      "name": "Vishal",
      "total_owed": 48.34
    }
  ],
  "message": "User balances retrieved successfully"
}
```

#### Get User Balance
**GET** `/users/{id}/balance`

Get current balance for a specific user.

**Response (200 OK):**
```json
{
  "user_id": 2,
  "name": "Abhinav",
  "total_owed": 96.66
}
```

**Response (200 OK) - No balance:**
```json
{
  "user_id": 1,
  "name": "Cyril",
  "total_owed": 0.0
}
```

#### Pay Off Full Balance
**POST** `/users/{id}/pay-off`

Pay off all outstanding balances for a user.

**Response (200 OK):**
```json
{
  "user_id": 2,
  "name": "Abhinav",
  "amount_paid": 96.66,
  "remaining_owed": 0.0,
  "message": "All balances paid off successfully"
}
```

#### Make Partial Payment
**POST** `/users/{id}/pay`

Make a partial payment towards a user's balance.

**Request Body:**
```json
{
  "user_id": 2,
  "amount": 50.00
}
```

**Response (200 OK):**
```json
{
  "user_id": 2,
  "name": "Abhinav",
  "amount_paid": 50.00,
  "remaining_owed": 46.66,
  "message": "Payment processed successfully"
}
```

**Error Response (400 Bad Request) - No balance:**
```json
{
  "error": "User has no outstanding balance to pay"
}
```

---

## Business Logic

### Cost Splitting
- Total cost is divided evenly among all participants
- Each participant (except driver) owes their share to the driver
- Driver pays upfront and others owe money

### Balance Management
- Each user has a single `balance` field tracking total amount owed
- When trips are created, non-driver participants' balances are increased by their share
- Payments reduce user balances
- Driver's balance is always 0 (they don't owe themselves)

### Trip Creation Flow
1. Validate all participant names exist as users
2. Automatically set current date and time
3. Create trip record with automatic timestamps
4. Add all participants with calculated share amounts
5. Update balances for participants (except driver)

## Example Workflows

### 1. Setup Users
```bash
# Create users
curl -X POST http://localhost:8081/users -H "Content-Type: application/json" -d '{
  "name": "Cyril",
  "email": "cyril@example.com",
  "phone": "+1234567890"
}'

curl -X POST http://localhost:8081/users -H "Content-Type: application/json" -d '{
  "name": "Abhinav",
  "email": "abhinav@example.com",
  "phone": "+1234567891"
}'

curl -X POST http://localhost:8081/users -H "Content-Type: application/json" -d '{
  "name": "Vishal",
  "email": "vishal@example.com",
  "phone": "+1234567892"
}'
```

### 2. Create Trips
```bash
# Create a trip with automatic date/time
curl -X POST http://localhost:8081/trips -H "Content-Type: application/json" -d '{
  "driver_name": "Cyril",
  "total_cost": 145,
  "participants": ["Cyril", "Abhinav", "Vishal"]
}'

# Create another trip
curl -X POST http://localhost:8081/trips -H "Content-Type: application/json" -d '{
  "driver_name": "Cyril",
  "total_cost": 120,
  "participants": ["Cyril", "Vishal"]
}'
```

### 3. Check Results
```bash
# Get trip details
curl http://localhost:8081/trips/1

# List all users
curl http://localhost:8081/users

# Check specific user
curl http://localhost:8081/users/by-name/Abhinav

# Get all user balances
curl http://localhost:8081/balances

# Get specific user's balance
curl http://localhost:8081/users/2/balance
```

### 4. Manage Balances
```bash
# Make a partial payment for user 2
curl -X POST http://localhost:8081/users/2/pay -H "Content-Type: application/json" -d '{
  "user_id": 2,
  "amount": 50.00
}'

# Pay off all balance for user 3
curl -X POST http://localhost:8081/users/3/pay-off

# Check balances after payments
curl http://localhost:8081/balances
```

## Error Codes

| Code | Description |
|------|-------------|
| 200  | OK - Request successful |
| 201  | Created - Resource created successfully |
| 400  | Bad Request - Invalid request format or validation failed |
| 404  | Not Found - Resource not found |
| 409  | Conflict - Resource already exists |
| 500  | Internal Server Error - Server error |
| 503  | Service Unavailable - Database connection failed |

## Data Validation

### User Validation
- `name`: Required, non-empty string
- `email`: Required, valid email format
- `phone`: Required, non-empty string

### Trip Validation  
- `driver_name`: Required, must exist as a user
- `total_cost`: Required, must be > 0
- `participants`: Required, non-empty array, all names must exist as users
- Driver must be included in participants list
- No duplicate participants allowed
- Date and time are automatically set to current timestamp

## Database Schema

The API uses the following tables:
- `users`: User information
- `trips`: Trip details with automatic date/time
- `participants`: Links users to trips with share amounts

The existing schema supports the simplified trip functionality.
