# Golang User Authentication API

A simple Golang REST API for user authentication (login/logout) using PostgreSQL and SuperTokens.

## Features

- User registration
- User login with session management
- User logout
- PostgreSQL database integration
- SuperTokens for session handling
- Password hashing with bcrypt

## Prerequisites

- Go 1.19+
- PostgreSQL (deployed on Dokploy)
- SuperTokens (deployed on Dokploy)

## Project Structure

```
golang_api/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── config/
│   └── config.go            # Configuration management
├── database/
│   └── database.go          # Database connection and schema
├── handlers/
│   └── auth.go              # Authentication handlers
├── models/
│   └── user.go              # User models
├── .env.example             # Environment variables template
├── go.mod
└── README.md
```

## Setup

### Local Development

1. **Clone and navigate to the project:**
   ```bash
   cd /Users/leo/Developer/golang_api
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Configure environment variables:**
   ```bash
   cp .env.example .env
   ```

   Edit `.env` with your configuration:
   ```env
   PORT=8080

   # PostgreSQL from Dokploy
   DB_HOST=mynew-postgres-jgcyi2
   DB_PORT=5432
   DB_USER=db-user
   DB_PASSWORD=db-password
   DB_NAME=db-name
   DB_SSLMODE=disable

   # SuperTokens - use localhost:3567 if port forwarding from Dokploy
   SUPERTOKENS_CONNECTION_URI=http://localhost:3567
   SUPERTOKENS_API_KEY=
   ```

4. **Load environment variables:**
   ```bash
   export $(cat .env | xargs)
   ```

5. **Run the application:**
   ```bash
   go run cmd/api/main.go
   ```

### Docker Deployment on Dokploy

1. **Build and run with Docker Compose:**
   ```bash
   docker-compose up -d
   ```

   The `docker-compose.yml` includes both the Golang API and SuperTokens services configured to use your existing PostgreSQL database on Dokploy.

2. **Check logs:**
   ```bash
   docker-compose logs -f golang-api
   ```

3. **Stop services:**
   ```bash
   docker-compose down
   ```

**Note:** The Docker Compose file is configured to work within the Dokploy network (`dokploy-network`) and connects to your existing PostgreSQL instance.

## API Endpoints

### Health Check
```bash
GET /api/health
```

Response:
```json
{
  "status": "healthy"
}
```

### Register User
```bash
POST /api/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```

Response:
```json
{
  "message": "Registration successful",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "created_at": "2025-11-23T14:00:00Z",
    "updated_at": "2025-11-23T14:00:00Z"
  }
}
```

### Login
```bash
POST /api/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```

Response:
```json
{
  "message": "Login successful",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "created_at": "2025-11-23T14:00:00Z",
    "updated_at": "2025-11-23T14:00:00Z"
  }
}
```

**Note:** Session cookies will be set automatically by SuperTokens.

### Logout
```bash
POST /api/logout
```

**Note:** Requires valid session (include cookies from login).

Response:
```json
{
  "message": "Logout successful"
}
```

## Testing with cURL

1. **Register a user:**
   ```bash
   curl -X POST http://localhost:8080/api/register \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'
   ```

2. **Login:**
   ```bash
   curl -X POST http://localhost:8080/api/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}' \
     -c cookies.txt
   ```

3. **Logout:**
   ```bash
   curl -X POST http://localhost:8080/api/logout \
     -b cookies.txt
   ```

## Database Schema

The application automatically creates the following schema on startup:

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | - |
| `DB_NAME` | Database name | `userdb` |
| `DB_SSLMODE` | SSL mode | `disable` |
| `SUPERTOKENS_CONNECTION_URI` | SuperTokens URI | `http://localhost:3567` |
| `SUPERTOKENS_API_KEY` | SuperTokens API key | - |

## Security Features

- Passwords are hashed using bcrypt
- Session management via SuperTokens
- CORS middleware for cross-origin requests
- SQL injection prevention using parameterized queries

## Troubleshooting

**Database connection issues:**
- Verify PostgreSQL is running on Dokploy
- Check connection details in `.env`
- Ensure database exists or create it: `CREATE DATABASE userdb;`

**SuperTokens connection issues:**
- Verify SuperTokens is running on Dokploy
- Check the connection URI
- Review SuperTokens logs

## License

MIT
