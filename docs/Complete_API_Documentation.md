# Darween ERP - Complete API Documentation

## Overview

Complete API documentation for the Darween ERP platform including Users, Companies, Subscriptions, and Products management.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

All protected endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## Table of Contents

1. [Authentication API](#authentication-api)
2. [Users API](#users-api)
3. [Companies API](#companies-api)
4. [Subscriptions API](#subscriptions-api)
5. [Products API](#products-api)
6. [Product Variants API](#product-variants-api)
7. [Stock Management API](#stock-management-api)
8. [Error Responses](#error-responses)
9. [Authorization Rules](#authorization-rules)

---

## Authentication API

### Register User

**POST** `/auth/register`

Creates a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe"
    }
  }
}
```

### Login User

**POST** `/auth/login`

Authenticates a user and returns a JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe"
    }
  }
}
```

---

## Users API

### Get Current User Profile

**GET** `/users/me`

Retrieves the current authenticated user's profile.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "is_active": true
  }
}
```

### Update Current User Profile

**PUT** `/users/me`

Updates the current authenticated user's profile.

**Request Body:**
```json
{
  "first_name": "John Updated",
  "last_name": "Doe Updated"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John Updated",
    "last_name": "Doe Updated",
    "is_active": true
  }
}
```

### List Users in Company

**GET** `/users?company_id={companyId}`

Lists all users in a specific company with their roles.

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "is_active": true,
      "role": "owner"
    }
  ]
}
```

---

## Companies API

### Create Company

**POST** `/companies`

Creates a new company. The user becomes the owner.

**Request Body:**
```json
{
  "name": "My Company",
  "code": "MYCO",
  "description": "My company description"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Company created successfully",
  "data": {
    "id": 1,
    "name": "My Company",
    "code": "MYCO",
    "description": "My company description",
    "is_active": true
  }
}
```

### List User's Companies

**GET** `/companies`

Lists all companies the current user belongs to.

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "My Company",
      "code": "MYCO",
      "description": "My company description",
      "is_active": true
    }
  ]
}
```

### Get Company Details

**GET** `/companies/{id}`

Retrieves details of a specific company.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "My Company",
    "code": "MYCO",
    "description": "My company description",
    "is_active": true
  }
}
```

### Update Company

**PUT** `/companies/{id}`

Updates company details. Requires Admin or Owner role.

**Request Body:**
```json
{
  "name": "Updated Company Name",
  "description": "Updated description",
  "is_active": true
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Company updated successfully",
  "data": {
    "id": 1,
    "name": "Updated Company Name",
    "code": "MYCO",
    "description": "Updated description",
    "is_active": true
  }
}
```

### Add User to Company

**POST** `/companies/{id}/users`

Adds a user to the company with a specific role. Requires Admin or Owner role.

**Request Body:**
```json
{
  "email": "newuser@example.com",
  "role": "employee"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User added to company successfully",
  "data": null
}
```

### Remove User from Company

**DELETE** `/companies/{id}/users/{userId}`

Removes a user from the company. Requires Admin or Owner role. Cannot remove company owner.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User removed from company successfully",
  "data": null
}
```

---

## Subscriptions API

### Get Company Subscription

**GET** `/companies/{id}/subscription`

Retrieves the subscription details for a company.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "company_id": 1,
    "plan_type": "free",
    "status": "active",
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": null,
    "max_users": 5
  }
}
```

### Update Subscription

**PUT** `/companies/{id}/subscription`

Updates the subscription plan. Requires Owner role.

**Request Body:**
```json
{
  "plan_type": "basic",
  "status": "active"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Subscription updated successfully",
  "data": {
    "id": 1,
    "company_id": 1,
    "plan_type": "basic",
    "status": "active",
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": null,
    "max_users": 20
  }
}
```

---

## Products API

### Create Product

**POST** `/companies/{companyId}/products`

Creates a new product for the specified company.

**Request Body:**
```json
{
  "name": "T-Shirt",
  "description": "Comfortable cotton t-shirt",
  "sku": "TSHIRT-001",
  "base_price": 20.00
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Product created successfully",
  "data": {
    "id": 1,
    "company_id": 1,
    "name": "T-Shirt",
    "description": "Comfortable cotton t-shirt",
    "sku": "TSHIRT-001",
    "base_price": 20.00,
    "is_active": true,
    "variants": []
  }
}
```

### List Products

**GET** `/companies/{companyId}/products`

Retrieves a paginated list of products for the specified company.

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20, max: 100)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "data": [
      {
        "id": 1,
        "company_id": 1,
        "name": "T-Shirt",
        "description": "Comfortable cotton t-shirt",
        "sku": "TSHIRT-001",
        "base_price": 20.00,
        "is_active": true,
        "variants": []
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 20,
    "total_pages": 1
  }
}
```

### Get Product

**GET** `/companies/{companyId}/products/{id}`

Retrieves a specific product with its variants.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "company_id": 1,
    "name": "T-Shirt",
    "description": "Comfortable cotton t-shirt",
    "sku": "TSHIRT-001",
    "base_price": 20.00,
    "is_active": true,
    "variants": [
      {
        "id": 1,
        "product_id": 1,
        "name": "Small Blue",
        "sku": "TSHIRT-001-S-BLU",
        "price": 20.00,
        "stock": 50,
        "attributes": {
          "size": "S",
          "color": "Blue"
        },
        "is_active": true
      }
    ]
  }
}
```

### Update Product

**PUT** `/companies/{companyId}/products/{id}`

Updates an existing product. Requires Admin or Owner role.

**Request Body:**
```json
{
  "name": "Updated T-Shirt",
  "description": "Updated description",
  "sku": "TSHIRT-001-UPDATED",
  "base_price": 25.00,
  "is_active": true
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Product updated successfully",
  "data": {
    "id": 1,
    "company_id": 1,
    "name": "Updated T-Shirt",
    "description": "Updated description",
    "sku": "TSHIRT-001-UPDATED",
    "base_price": 25.00,
    "is_active": true,
    "variants": []
  }
}
```

### Delete Product

**DELETE** `/companies/{companyId}/products/{id}`

Soft deletes a product (sets is_active = false). Requires Admin or Owner role.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Product deleted successfully",
  "data": null
}
```

---

## Product Variants API

### Create Product Variant

**POST** `/companies/{companyId}/products/{productId}/variants`

Creates a new variant for the specified product.

**Request Body:**
```json
{
  "name": "Small Blue",
  "sku": "TSHIRT-001-S-BLU",
  "price": 20.00,
  "stock": 50,
  "attributes": {
    "size": "S",
    "color": "Blue"
  }
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Product variant created successfully",
  "data": {
    "id": 1,
    "product_id": 1,
    "name": "Small Blue",
    "sku": "TSHIRT-001-S-BLU",
    "price": 20.00,
    "stock": 50,
    "attributes": {
      "size": "S",
      "color": "Blue"
    },
    "is_active": true
  }
}
```

### List Product Variants

**GET** `/companies/{companyId}/products/{productId}/variants`

Retrieves all variants for the specified product.

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "product_id": 1,
      "name": "Small Blue",
      "sku": "TSHIRT-001-S-BLU",
      "price": 20.00,
      "stock": 50,
      "attributes": {
        "size": "S",
        "color": "Blue"
      },
      "is_active": true
    }
  ]
}
```

### Get Product Variant

**GET** `/companies/{companyId}/products/{productId}/variants/{id}`

Retrieves a specific product variant.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "product_id": 1,
    "name": "Small Blue",
    "sku": "TSHIRT-001-S-BLU",
    "price": 20.00,
    "stock": 50,
    "attributes": {
      "size": "S",
      "color": "Blue"
    },
    "is_active": true
  }
}
```

### Update Product Variant

**PUT** `/companies/{companyId}/products/{productId}/variants/{id}`

Updates an existing product variant.

**Request Body:**
```json
{
  "name": "Small Blue Updated",
  "sku": "TSHIRT-001-S-BLU-UPDATED",
  "price": 22.00,
  "stock": 75,
  "attributes": {
    "size": "S",
    "color": "Blue",
    "material": "Cotton"
  },
  "is_active": true
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Product variant updated successfully",
  "data": {
    "id": 1,
    "product_id": 1,
    "name": "Small Blue Updated",
    "sku": "TSHIRT-001-S-BLU-UPDATED",
    "price": 22.00,
    "stock": 75,
    "attributes": {
      "size": "S",
      "color": "Blue",
      "material": "Cotton"
    },
    "is_active": true
  }
}
```

### Delete Product Variant

**DELETE** `/companies/{companyId}/products/{productId}/variants/{id}`

Soft deletes a product variant (sets is_active = false).

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Product variant deleted successfully",
  "data": null
}
```

---

## Stock Management API

### Update Variant Stock

**PUT** `/companies/{companyId}/products/{productId}/variants/{id}/stock`

Sets the exact stock level for a product variant. Requires Manager, Admin, or Owner role.

**Request Body:**
```json
{
  "stock": 100
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Stock updated successfully",
  "data": {
    "id": 1,
    "product_id": 1,
    "name": "Small Blue",
    "sku": "TSHIRT-001-S-BLU",
    "price": 20.00,
    "stock": 100,
    "attributes": {
      "size": "S",
      "color": "Blue"
    },
    "is_active": true
  }
}
```

### Adjust Variant Stock

**POST** `/companies/{companyId}/products/{productId}/variants/{id}/stock/adjust`

Adjusts the stock level by adding or subtracting from the current amount. Requires Manager, Admin, or Owner role.

**Request Body:**
```json
{
  "amount": 10
}
```

**Notes:**
- Positive `amount`: Adds to current stock
- Negative `amount`: Subtracts from current stock
- Will fail if trying to subtract more than available stock

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Stock adjusted successfully",
  "data": {
    "id": 1,
    "product_id": 1,
    "name": "Small Blue",
    "sku": "TSHIRT-001-S-BLU",
    "price": 20.00,
    "stock": 110,
    "attributes": {
      "size": "S",
      "color": "Blue"
    },
    "is_active": true
  }
}
```

---

## Error Responses

All endpoints return standardized error responses:

### Validation Error (400 Bad Request)
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "invalid product data"
  }
}
```

### Unauthorized (401 Unauthorized)
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "invalid or expired token"
  }
}
```

### Forbidden (403 Forbidden)
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "insufficient permissions for this operation"
  }
}
```

### Not Found (404 Not Found)
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "product not found"
  }
}
```

### Conflict (409 Conflict)
```json
{
  "success": false,
  "error": {
    "code": "CONFLICT",
    "message": "product with this SKU already exists in the company"
  }
}
```

---

## Authorization Rules

### Role-Based Access Control

- **Owner**: Full control over company (cannot be removed)
- **Admin**: Can manage company settings, users, products, and variants
- **Manager**: Can view all, update stock levels, create/update variants
- **Employee**: Read-only access to products and variants

### Company Scoping

- All operations are scoped to the user's company
- Users can only access data from companies they belong to
- SKU uniqueness is enforced per company, not globally

### Subscription Plans

- **Free**: Up to 5 users
- **Basic**: Up to 20 users
- **Premium**: Up to 100 users
- **Enterprise**: Up to 1000 users

---

## Health Check

### API Health Check

**GET** `/health`

Checks if the API is running.

**Response (200 OK):**
```json
{
  "status": "ok"
}
```

---

## Complete Example Workflow

### 1. Register and Login
```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123",
    "first_name": "Admin",
    "last_name": "User"
  }'

# Login (save the token)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123"
  }'
```

### 2. Create Company
```bash
curl -X POST http://localhost:8080/api/v1/companies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "My Company",
    "code": "MYCO",
    "description": "My company description"
  }'
```

### 3. Create Product
```bash
curl -X POST http://localhost:8080/api/v1/companies/1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "T-Shirt",
    "description": "Comfortable cotton t-shirt",
    "sku": "TSHIRT-001",
    "base_price": 20.00
  }'
```

### 4. Create Product Variants
```bash
# Small Blue variant
curl -X POST http://localhost:8080/api/v1/companies/1/products/1/variants \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Small Blue",
    "sku": "TSHIRT-001-S-BLU",
    "price": 20.00,
    "stock": 50,
    "attributes": {
      "size": "S",
      "color": "Blue"
    }
  }'

# Large Red variant
curl -X POST http://localhost:8080/api/v1/companies/1/products/1/variants \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Large Red",
    "sku": "TSHIRT-001-L-RED",
    "price": 22.00,
    "stock": 30,
    "attributes": {
      "size": "L",
      "color": "Red"
    }
  }'
```

### 5. List Products with Variants
```bash
curl -X GET http://localhost:8080/api/v1/companies/1/products \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 6. Update Stock
```bash
curl -X PUT http://localhost:8080/api/v1/companies/1/products/1/variants/1/stock \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "stock": 75
  }'
```

---

## Postman Collection

A complete Postman collection is available at:
`postman/Darween_ERP_Products_API.postman_collection.json`

Import this collection into Postman to test all endpoints with pre-configured requests and environment variables.

---

## Support

For API support or questions, please refer to the main project documentation or create an issue in the project repository.
