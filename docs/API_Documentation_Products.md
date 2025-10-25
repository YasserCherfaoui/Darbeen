# Darween ERP - Products & Variants API Documentation

## Overview

This document provides comprehensive API documentation for the Products and Product Variants management endpoints in the Darween ERP platform.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

All endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## Products API

### Create Product

**POST** `/companies/{companyId}/products`

Creates a new product for the specified company.

**Headers:**
- `Content-Type: application/json`
- `Authorization: Bearer <token>`

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

**Validation Rules:**
- `name`: Required, string
- `sku`: Required, string, unique within company
- `base_price`: Optional, number >= 0

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

**Validation Rules:**
- `name`: Required, string
- `sku`: Required, string, unique within product
- `price`: Optional, number >= 0
- `stock`: Optional, integer >= 0
- `attributes`: Optional, JSON object

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
    },
    {
      "id": 2,
      "product_id": 1,
      "name": "Large Red",
      "sku": "TSHIRT-001-L-RED",
      "price": 22.00,
      "stock": 30,
      "attributes": {
        "size": "L",
        "color": "Red"
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

## Authorization Rules

### Role-Based Access Control

- **Owner**: Full CRUD access to all products and variants
- **Admin**: Full CRUD access to all products and variants
- **Manager**: Can view all, update stock levels, create/update variants
- **Employee**: Read-only access to products and variants

### Company Scoping

- All operations are scoped to the user's company
- Users can only access products from companies they belong to
- SKU uniqueness is enforced per company, not globally

## Example Workflows

### Complete Product Setup

1. **Create Company** (if not exists)
2. **Create Product**
3. **Create Multiple Variants**
4. **Update Stock Levels**

### Sample cURL Commands

```bash
# 1. Create a product
curl -X POST http://localhost:8080/api/v1/companies/1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "T-Shirt",
    "description": "Comfortable cotton t-shirt",
    "sku": "TSHIRT-001",
    "base_price": 20.00
  }'

# 2. Create a variant
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

# 3. List products with variants
curl -X GET http://localhost:8080/api/v1/companies/1/products \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 4. Update stock
curl -X PUT http://localhost:8080/api/v1/companies/1/products/1/variants/1/stock \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "stock": 75
  }'
```

## Rate Limiting

Currently no rate limiting is implemented. Consider implementing rate limiting for production use.

## Versioning

This API uses URL versioning (`/api/v1/`). Future versions will use `/api/v2/`, etc.

## Support

For API support or questions, please refer to the main project documentation or create an issue in the project repository.
