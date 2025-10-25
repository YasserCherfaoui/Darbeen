# Darween ERP API - Postman Documentation

This directory contains comprehensive Postman documentation for the Darween ERP API.

## Files

- `darween-api.postman_collection.json` - Complete API collection with all endpoints
- `darween-api.postman_environment.json` - Environment variables for local development

## Quick Start

### 1. Import the Collection

1. Open Postman
2. Click **Import** button
3. Select `darween-api.postman_collection.json`
4. The collection will appear in your Collections sidebar

### 2. Import the Environment (Optional)

1. Click **Import** button
2. Select `darween-api.postman_environment.json`
3. Select the "Darween API - Local" environment from the environment dropdown (top right)

**Note:** The environment file is optional. The collection has default variables configured, so you can use it directly without importing the environment.

### 3. Start Using the API

#### Testing the API

First, verify the API is running:

1. Expand the **Health** folder
2. Select **Health Check** request
3. Click **Send**
4. You should receive: `{"status": "ok"}`

#### Authentication Flow

The API requires JWT authentication for most endpoints. Follow these steps:

1. **Register a New User**
   - Navigate to **Auth > Register**
   - Click **Send**
   - The response will include a JWT token
   - The token is automatically saved to the collection variable

2. **Or Login with Existing User**
   - Navigate to **Auth > Login**
   - Update the email/password in the request body
   - Click **Send**
   - The token is automatically saved to the collection variable

3. **Use Protected Endpoints**
   - All other endpoints automatically use the saved token
   - Try **Users > Get Current User** to verify authentication

## Collection Structure

### 📁 Health
- `GET /api/v1/health` - API health check

### 📁 Auth
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login existing user

**Note:** Registration and login requests automatically save the JWT token to collection variables.

### 📁 Users
- `GET /api/v1/users/me` - Get current user profile
- `PUT /api/v1/users/me` - Update current user
- `GET /api/v1/users?company_id={id}` - List users in a company

### 📁 Companies
- `POST /api/v1/companies` - Create new company
- `GET /api/v1/companies` - List user's companies
- `GET /api/v1/companies/:id` - Get company details
- `PUT /api/v1/companies/:id` - Update company
- `POST /api/v1/companies/:id/users` - Add user to company
- `DELETE /api/v1/companies/:id/users/:userId` - Remove user from company

### 📁 Subscriptions
- `GET /api/v1/companies/:id/subscription` - Get company subscription
- `PUT /api/v1/companies/:id/subscription` - Update subscription

## Configuration

### Collection Variables

The collection includes these variables (accessible via collection settings):

- `base_url` - API base URL (default: `http://localhost:8080`)
- `token` - JWT authentication token (auto-populated after login/register)

### Changing the Base URL

If your API runs on a different URL or port:

1. Click the collection name (Darween ERP API)
2. Select **Variables** tab
3. Update the `base_url` value
4. Save the collection

### Manual Token Configuration

If you need to set the token manually:

1. Click the collection name
2. Go to **Variables** tab
3. Set the `token` variable value
4. Save the collection

## Request Examples

### Register a New User

```http
POST {{base_url}}/api/v1/auth/register
Content-Type: application/json

{
    "email": "john.doe@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
}
```

**Response:**
```json
{
    "success": true,
    "message": "User registered successfully",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
            "id": 1,
            "email": "john.doe@example.com",
            "first_name": "John",
            "last_name": "Doe"
        }
    }
}
```

### Create a Company

```http
POST {{base_url}}/api/v1/companies
Authorization: Bearer {{token}}
Content-Type: application/json

{
    "name": "Acme Corporation",
    "code": "ACME001",
    "description": "A leading provider of innovative solutions"
}
```

**Response:**
```json
{
    "success": true,
    "message": "Company created successfully",
    "data": {
        "id": 1,
        "name": "Acme Corporation",
        "code": "ACME001",
        "description": "A leading provider of innovative solutions",
        "is_active": true
    }
}
```

## Response Format

All API responses follow this structure:

### Success Response
```json
{
    "success": true,
    "message": "Optional success message",
    "data": { /* Response data */ }
}
```

### Error Response
```json
{
    "success": false,
    "error": {
        "code": "ERROR_CODE",
        "message": "Human-readable error message"
    }
}
```

### Common Error Codes

- `VALIDATION_ERROR` - Invalid request data (400)
- `UNAUTHORIZED` - Missing or invalid authentication (401)
- `FORBIDDEN` - Insufficient permissions (403)
- `NOT_FOUND` - Resource not found (404)
- `CONFLICT` - Resource conflict (409)
- `INTERNAL` - Server error (500)

## Testing Workflow

Here's a recommended workflow for testing the complete API:

1. **Health Check** - Verify API is running
2. **Register** - Create a new user account
3. **Get Current User** - Verify authentication works
4. **Create Company** - Create your first company
5. **List Companies** - View all your companies
6. **Get Company Subscription** - Check default subscription
7. **Update Subscription** - Upgrade to a different plan
8. **Add User to Company** - Invite another user (requires second registered user)
9. **List Users** - View company members
10. **Update Company** - Modify company details

## Tips

- **Auto-save tokens**: The collection automatically saves JWT tokens from login/register responses
- **Example responses**: Each request includes example success and error responses
- **Path parameters**: Replace `:id` and `:userId` in URLs with actual values
- **Query parameters**: Modify query parameters in the Params tab
- **Request bodies**: Update request bodies in the Body tab before sending

## Troubleshooting

### 401 Unauthorized Error

- Ensure you've logged in or registered first
- Check that the token variable is set (Collection > Variables)
- Token may have expired - try logging in again

### 404 Not Found Error

- Verify the API server is running
- Check the `base_url` is correct
- Ensure you're using valid resource IDs

### CORS Errors

- The API includes CORS middleware
- If testing from browser, ensure your origin is allowed
- Postman desktop app doesn't have CORS restrictions

### Connection Refused

- Ensure the API server is running: `make run` or `go run cmd/api/main.go`
- Verify the server is listening on the correct port (default: 8080)
- Check `base_url` matches your server configuration

## Environment Setup

For different environments (development, staging, production), create additional environment files:

1. Duplicate the environment
2. Update the `base_url` value
3. Switch between environments using the environment dropdown

Example environments:
- **Local**: `http://localhost:8080`
- **Development**: `https://dev-api.darween.com`
- **Staging**: `https://staging-api.darween.com`
- **Production**: `https://api.darween.com`

## Support

For API documentation updates or issues, please refer to:
- `ARCHITECTURE.md` - System architecture
- `README.md` - Project overview
- `QUICKSTART.md` - Development setup guide

