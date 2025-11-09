<!-- a4575041-5921-4563-b061-326c2f5f6534 7f9eb968-c0c7-4501-a50a-bc2589270203 -->
# ERP Platform with DDD Architecture - Extended with Products

## Products & Variants Feature Specification

### Product Entity

- **Fields**: ID, CompanyID, Name, Description, SKU (unique per company), BasePrice, IsActive, CreatedAt, UpdatedAt
- **Relationships**: Belongs to Company, has many ProductVariants
- **Business Rules**: SKU must be unique within a company

### ProductVariant Entity

- **Fields**: ID, ProductID, Name, SKU (unique per company), Price, Stock, Attributes (JSON), IsActive, CreatedAt, UpdatedAt
- **Relationships**: Belongs to Product
- **Business Rules**: 
  - Variant SKU must be unique within a company
  - Price can override base product price
  - Attributes stored as JSON for flexibility (e.g., `{"size": "L", "color": "Blue"}`)

## Product Management Endpoints

### Product Routes (Protected)

- `POST /api/v1/companies/:companyId/products` - Create product
- `GET /api/v1/companies/:companyId/products` - List company products (with pagination)
- `GET /api/v1/companies/:companyId/products/:id` - Get product details with variants
- `PUT /api/v1/companies/:companyId/products/:id` - Update product (admin/owner only)
- `DELETE /api/v1/companies/:companyId/products/:id` - Soft delete product (admin/owner only)

### Product Variant Routes (Protected)

- `POST /api/v1/companies/:companyId/products/:productId/variants` - Create variant
- `GET /api/v1/companies/:companyId/products/:productId/variants` - List product variants
- `GET /api/v1/companies/:companyId/products/:productId/variants/:id` - Get variant details
- `PUT /api/v1/companies/:companyId/products/:productId/variants/:id` - Update variant
- `DELETE /api/v1/companies/:companyId/products/:productId/variants/:id` - Soft delete variant

## Implementation Structure

```
internal/
├── domain/
│   └── product/
│       ├── entity.go             # Product & ProductVariant entities
│       └── repository.go         # Repository interfaces
├── application/
│   └── product/
│       ├── dto.go                # Product & Variant DTOs
│       └── service.go            # Product management service
├── infrastructure/
│   └── persistence/
│       └── postgres/
│           └── product_repository.go  # Repository implementation
└── presentation/
    └── http/
        └── handler/
            └── product_handler.go     # Product & Variant handlers
```

## Database Schema

### products Table

```sql
- id (PK)
- company_id (FK → companies.id, indexed)
- name (varchar, not null)
- description (text)
- sku (varchar, not null)
- base_price (decimal)
- is_active (boolean, default true)
- created_at, updated_at (timestamps)
- UNIQUE INDEX (company_id, sku)
```

### product_variants Table

```sql
- id (PK)
- product_id (FK → products.id, indexed)
- name (varchar, not null)
- sku (varchar, not null)
- price (decimal)
- stock (integer, default 0)
- attributes (jsonb)
- is_active (boolean, default true)
- created_at, updated_at (timestamps)
- UNIQUE INDEX (product_id, sku)
```

## Authorization Rules

- Users must belong to the company to access products
- **Owner/Admin**: Full CRUD access to products and variants
- **Manager**: Can view all, update stock levels
- **Employee**: Read-only access

## Key Implementation Details

### 1. Company Scoping

- All product queries filtered by company_id
- Validate user has access to company before any operation
- Prevent cross-company data access

### 2. SKU Uniqueness

- Enforce unique SKU per company (not globally)
- Composite unique constraint: (company_id, sku)
- Validate on create and update

### 3. Soft Deletes

- Set is_active = false instead of hard delete
- Filter inactive products in list queries
- Allow reactivation by admins

### 4. Variant Attributes

- Use PostgreSQL JSONB for flexible attributes
- Support any key-value pairs
- Enable querying by attributes (future enhancement)

### 5. Pagination

- List endpoints support page and limit query params
- Default: page=1, limit=20
- Return total count in response

## Example API Flows

### Create Product with Variants

```
1. POST /api/v1/companies/1/products
   Body: {name: "T-Shirt", sku: "TSHIRT-001", base_price: 20.00}
   
2. POST /api/v1/companies/1/products/1/variants
   Body: {name: "Small Blue", sku: "TSHIRT-001-S-BLU", price: 20.00, 
          attributes: {"size": "S", "color": "Blue"}}
   
3. POST /api/v1/companies/1/products/1/variants
   Body: {name: "Large Red", sku: "TSHIRT-001-L-RED", price: 22.00,
          attributes: {"size": "L", "color": "Red"}}
```

### List Products with Variants

```
GET /api/v1/companies/1/products?page=1&limit=20
Response includes products with their variants embedded
```

### To-dos

- [ ] Create DDD folder structure and initialize Go modules with dependencies
- [ ] Implement configuration management and environment loading
- [ ] Create domain entities (User, Company, Subscription, UserCompanyRole) and repository interfaces
- [ ] Implement database connection, GORM setup, and auto-migration
- [ ] Implement repository pattern with PostgreSQL for all entities
- [ ] Implement JWT token generation and validation utilities
- [ ] Create authentication application service (register, login)
- [ ] Implement application services for User, Company, and Subscription
- [ ] Create HTTP handlers and middleware (auth, CORS)
- [ ] Set up Gin router with routes and create main.go entry point
- [ ] Create .env.example file with all required environment variables