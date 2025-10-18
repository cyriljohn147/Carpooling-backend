# Carpooling Backend API

A complete REST API for carpooling trip management with transaction handling, built with Go, Gin, PostgreSQL, and SQLC.

## Features

### 🚗 Core Functionality
- **User Management**: Complete CRUD operations for users
- **Trip Management**: Create trips with automatic cost splitting
- **Name-based Transactions**: Use participant names instead of IDs
- **Balance Tracking**: Weekly balance tracking per user
- **Database Transactions**: Atomic operations for data consistency

### 🛠 Technical Features
- **RESTful API**: Clean, predictable endpoints
- **Database Migrations**: Version-controlled schema changes
- **Type-safe Queries**: Generated with SQLC
- **Validation**: Comprehensive request validation
- **CORS Support**: Cross-origin resource sharing
- **Health Checks**: API and database health monitoring
- **Graceful Shutdown**: Clean server shutdown

## Quick Start

### Prerequisites
- Go 1.19+
- PostgreSQL
- Docker (optional, for database)
- golang-migrate CLI tool
- SQLC

### 1. Setup Database
```bash
# Using Docker (recommended)
make postgres        # Start PostgreSQL container
make createdb        # Create database
make migrateup       # Run migrations

# Or use existing PostgreSQL
createdb carpooling
migrate -path db/migration -database "postgresql://user:password@localhost:5432/carpooling?sslmode=disable" up
```

### 2. Install Dependencies
```bash
make deps
```

### 3. Configure Environment
Update `app.env` with your database credentials:
```env
DB_DRIVER=postgres
DB_SOURCE=postgresql://root:secret@localhost:5432/carpooling?sslmode=disable
MIGRATION_URL=file://db/migration
SERVER_ADDRESS=0.0.0.0:8080
```

### 4. Run Server
```bash
make run
# or
go run main.go
```

The API will be available at `http://localhost:8080`

## API Endpoints

### System
- `GET /` - API information
- `GET /health` - Health check

### Users
- `POST /users` - Create user
- `GET /users` - List users (with pagination)
- `GET /users/:id` - Get user by ID
- `GET /users/by-name/:name` - Get user by name
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user

### Trip Transactions
- `POST /trips/with-transactions` - Create trip with automatic cost splitting
- `GET /trips/:id/transactions` - Get trip transaction details

## Example Usage

### 1. Create Users
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Cyril",
    "email": "cyril@example.com",
    "phone": "+1234567890"
  }'

curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Abhinav",
    "email": "abhinav@example.com",
    "phone": "+1234567891"
  }'
```

### 2. Create Trip with Automatic Transactions
```bash
curl -X POST http://localhost:8080/trips/with-transactions \
  -H "Content-Type: application/json" \
  -d '{
    "driver_name": "Cyril",
    "date": "2024-01-15",
    "time": "09:00:00",
    "total_cost": 145,
    "week": 3,
    "participants": ["Cyril", "Abhinav"]
  }'
```

This automatically:
- ✅ Validates all participants exist
- ✅ Creates trip record
- ✅ Calculates individual shares (145 ÷ 2 = 72.50 each)
- ✅ Adds participants with share amounts
- ✅ Updates balances (Abhinav owes Cyril 72.50)
- ✅ Returns transaction details

## Project Structure

```
.
├── api/                    # HTTP handlers and server setup
│   ├── server.go          # Main server configuration
│   ├── user_handlers.go   # User CRUD endpoints
│   ├── trip_handlers.go   # Trip transaction endpoints
│   ├── trip_service.go    # Business logic for transactions
│   └── trip_transaction.go # Request/response types
├── db/
│   ├── migration/         # Database migrations
│   ├── query/            # SQL queries for SQLC
│   └── sqlc/             # Generated Go code from SQLC
├── util/
│   ├── config.go         # Configuration management
│   └── random.go         # Utility functions
├── main.go               # Application entry point
├── Makefile              # Build and development commands
├── app.env               # Environment configuration
└── README.md
```

## Development

### Available Make Commands

```bash
# Database
make postgres       # Start PostgreSQL container
make createdb       # Create database
make dropdb         # Drop database
make migrateup      # Run migrations
make migratedown    # Rollback migrations

# Development
make run            # Run server
make build          # Build binary
make dev            # Run in development mode
make sqlc           # Generate SQLC code
make test           # Run tests
make fmt            # Format code
make vet            # Vet code
make check          # Run all checks (fmt, vet, test)

# Dependencies
make deps           # Install/update dependencies
make clean          # Clean build artifacts
```

### Database Schema

The API uses these tables:
- `users`: User information (id, name, email, phone, created_at)
- `trips`: Trip details (id, driver_id, date, time, total_cost)
- `participants`: User-trip relationships (id, trip_id, user_id, share_amount)
- `balances`: Weekly balances (id, user_id, amount_owed, week, created_at)

### Business Logic

#### Cost Splitting Algorithm
1. **Equal Split**: Total cost ÷ Number of participants
2. **Remainder Handling**: Any remainder goes to the last participant
3. **Driver Payment**: Driver pays upfront, others owe their share

#### Balance Management
- Each user has a balance per week
- Balances accumulate over multiple trips in the same week
- Only non-driver participants have their balances updated

#### Transaction Flow
1. Validate all participant names exist as users
2. Create trip record in database
3. Add all participants with calculated share amounts
4. Generate virtual transaction records (who owes whom)
5. Update user balances for the specified week
6. Return comprehensive response with all details

## Testing

```bash
# Setup test database
make setup-test

# Run all tests
make test

# Run specific test categories
make test-users
make test-trips
make test-participants
make test-balances

# Test with coverage
make test-coverage
make test-coverage-html  # Generate HTML coverage report

# Cleanup test database
make teardown-test
```

## Configuration

The application uses environment variables loaded from `app.env`:

| Variable | Description | Example |
|----------|-------------|---------|
| `DB_DRIVER` | Database driver | `postgres` |
| `DB_SOURCE` | Database connection string | `postgresql://user:pass@host:port/db?sslmode=disable` |
| `MIGRATION_URL` | Migration files path | `file://db/migration` |
| `SERVER_ADDRESS` | Server bind address | `0.0.0.0:8080` |

## Error Handling

The API returns consistent error responses:

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

Common HTTP status codes:
- `200` - OK
- `201` - Created
- `400` - Bad Request (validation failed)
- `404` - Not Found
- `409` - Conflict (duplicate resource)
- `500` - Internal Server Error
- `503` - Service Unavailable (database error)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Run checks: `make check`
6. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For detailed API documentation, see [API_DOCUMENTATION.md](./API_DOCUMENTATION.md).

For transaction system details, see [README_trip_transactions.md](./README_trip_transactions.md).

---

**Built with ❤️ using Go, Gin, PostgreSQL, and SQLC**