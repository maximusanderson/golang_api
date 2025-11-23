# Deployment Guide for Dokploy

This guide explains how to deploy the Golang authentication API to Dokploy alongside your existing PostgreSQL and SuperTokens services.

## Prerequisites

- Dokploy instance running
- PostgreSQL deployed on Dokploy: `mynew-postgres-jgcyi2`
- SuperTokens container (from docker-compose.yml)

## Deployment Options

### Option 1: Deploy with Docker Compose (Recommended)

The `docker-compose.yml` file includes both SuperTokens and the Golang API configured to work on Dokploy.

**Steps:**

1. **Push your code to a Git repository** (GitHub, GitLab, etc.)

2. **In Dokploy:**
   - Create a new application
   - Select "Docker Compose" as the deployment type
   - Connect your Git repository
   - Point to the `docker-compose.yml` file
   - Deploy

3. **Environment Variables** (already configured in docker-compose.yml):
   ```yaml
   PORT: 8080
   DB_HOST: mynew-postgres-jgcyi2
   DB_PORT: 5432
   DB_USER: db-user
   DB_PASSWORD: db-password
   DB_NAME: db-name
   DB_SSLMODE: disable
   SUPERTOKENS_CONNECTION_URI: http://supertokens:3567
   ```

4. **Access your API:**
   - The API will be available on port 8080
   - Configure your domain/reverse proxy in Dokploy

### Option 2: Deploy Only the Golang API

If you already have SuperTokens running separately:

1. **Create a new Dockerfile-only application in Dokploy**

2. **Set environment variables:**
   ```
   PORT=8080
   DB_HOST=mynew-postgres-jgcyi2
   DB_PORT=5432
   DB_USER=db-user
   DB_PASSWORD=db-password
   DB_NAME=db-name
   DB_SSLMODE=disable
   SUPERTOKENS_CONNECTION_URI=http://your-supertokens-service:3567
   ```

3. **Deploy**

## Network Configuration

The application is configured to use the `dokploy-network` external network, which allows:
- Communication with PostgreSQL (`mynew-postgres-jgcyi2`)
- Communication with SuperTokens service
- All services can discover each other by service name

## Database Initialization

The application automatically creates the `users` table on first startup:

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

No manual database setup required!

## Verifying Deployment

1. **Check service health:**
   ```bash
   curl http://your-domain:8080/api/health
   ```

   Expected response:
   ```json
   {"status":"healthy"}
   ```

2. **Test registration:**
   ```bash
   curl -X POST http://your-domain:8080/api/register \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'
   ```

3. **Check logs in Dokploy:**
   - View application logs for any errors
   - Look for "Successfully connected to database"
   - Look for "Database schema initialized"

## Troubleshooting

### Cannot connect to PostgreSQL

**Error:** `error connecting to database`

**Solution:**
- Verify PostgreSQL service name: `mynew-postgres-jgcyi2`
- Ensure both services are on the same Docker network: `dokploy-network`
- Check PostgreSQL credentials match your deployment

### Cannot connect to SuperTokens

**Error:** `Failed to initialize SuperTokens`

**Solution:**
- Ensure SuperTokens is running on the same network
- Verify connection URI: `http://supertokens:3567`
- Check SuperTokens logs for startup issues

### Port conflicts

**Error:** `bind: address already in use`

**Solution:**
- Change the `PORT` environment variable
- Update port mappings in `docker-compose.yml`

## Production Considerations

1. **Security:**
   - Use strong database passwords
   - Enable SSL/TLS for PostgreSQL (`DB_SSLMODE=require`)
   - Set up HTTPS/reverse proxy in Dokploy
   - Consider adding rate limiting

2. **SuperTokens API Key:**
   - Generate an API key in SuperTokens dashboard
   - Add to `SUPERTOKENS_API_KEY` environment variable

3. **CORS Configuration:**
   - Update allowed origins in `cmd/api/main.go` (corsMiddleware)
   - Replace `*` with your actual frontend domain

4. **Database Connection Pooling:**
   - Consider adding connection pool settings for production load

## Updating the Application

1. Push changes to your Git repository
2. In Dokploy, trigger a redeploy
3. Monitor logs for successful startup

## Rolling Back

If deployment fails:
1. In Dokploy, select previous successful deployment
2. Click "Redeploy"

## Support

For issues:
- Check Dokploy application logs
- Verify environment variables are set correctly
- Ensure all services are on the same Docker network
