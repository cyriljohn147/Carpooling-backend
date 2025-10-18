# Database CRUD Tests

This directory contains comprehensive CRUD operation tests for all database tables in the carpooling application.

## Test Files

- `main_test.go` - Test setup and database connection
- `users_test.go` - User CRUD operations tests
- `trips_test.go` - Trip CRUD operations tests
- `participants_test.go` - Participant CRUD operations tests
- `balances_test.go` - Balance CRUD operations tests

## Prerequisites

1. PostgreSQL database running locally
2. Test database named `carpooling_test`
3. Database schema migrated to the test database

## Database Setup

```bash
# Create test database
createdb carpooling_test

# Run migrations on test database
migrate -path ../migration -database "postgresql://root:secret@localhost:5432/carpooling_test?sslmode=disable" -verbose up
```

## Environment Variables

The tests use the following default database connection:
```
Host: localhost
Port: 5432
User: root
Password: secret
Database: carpooling_test
```

To use different connection parameters, modify the `dbSource` constant in `main_test.go`.

## Running Tests

### Run all tests
```bash
go test -v ./db/sqlc/test
```

### Run specific test file
```bash
go test -v ./db/sqlc/test -run TestUser
go test -v ./db/sqlc/test -run TestTrip
go test -v ./db/sqlc/test -run TestParticipant
go test -v ./db/sqlc/test -run TestBalance
```

### Run with coverage
```bash
go test -v -cover ./db/sqlc/test
```

## Test Coverage

Each test file covers the following operations:

### CREATE operations
- Basic creation with all fields
- Creation with minimal required fields
- Validation of created data

### READ operations
- Get single record by ID
- List records with pagination
- Handle empty results

### UPDATE operations
- Full record updates
- Partial record updates
- Validation of updated data

### DELETE operations
- Delete existing records
- Verify deletion (record not found after delete)

## Random Test Data

The tests use utility functions from `util/random.go` to generate:
- Random names, emails, and phone numbers
- Random locations, distances, and prices
- Random dates and times
- Random student IDs and password hashes

## Dependencies

- `github.com/stretchr/testify` - Test assertions
- `github.com/lib/pq` - PostgreSQL driver
- Standard Go testing package