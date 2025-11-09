package inventory

// Request DTOs

type CreateInventoryRequest struct {
	ProductVariantID uint  `json:"product_variant_id" binding:"required"`
	CompanyID        *uint `json:"company_id,omitempty"`
	FranchiseID      *uint `json:"franchise_id,omitempty"`
	Stock            int   `json:"stock" binding:"required"`
}

type UpdateInventoryStockRequest struct {
	Stock int `json:"stock" binding:"required,min=0"`
}

type AdjustInventoryStockRequest struct {
	Adjustment int    `json:"adjustment" binding:"required"`
	Notes      string `json:"notes,omitempty"`
}

type ReserveStockRequest struct {
	Quantity      int    `json:"quantity" binding:"required,min=1"`
	ReferenceType string `json:"reference_type,omitempty"`
	ReferenceID   string `json:"reference_id,omitempty"`
	Notes         string `json:"notes,omitempty"`
}

type ReleaseStockRequest struct {
	Quantity int    `json:"quantity" binding:"required,min=1"`
	Notes    string `json:"notes,omitempty"`
}

// Response DTOs

type InventoryResponse struct {
	ID               uint   `json:"id"`
	ProductVariantID uint   `json:"product_variant_id"`
	VariantName      string `json:"variant_name,omitempty"`
	VariantSKU       string `json:"variant_sku,omitempty"`
	CompanyID        *uint  `json:"company_id,omitempty"`
	FranchiseID      *uint  `json:"franchise_id,omitempty"`
	FranchiseName    string `json:"franchise_name,omitempty"`
	Stock            int    `json:"stock"`
	ReservedStock    int    `json:"reserved_stock"`
	AvailableStock   int    `json:"available_stock"`
	IsActive         bool   `json:"is_active"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

type InventoryListResponse struct {
	Inventories []*InventoryResponse `json:"inventories"`
	Total       int                  `json:"total"`
}

type InventoryMovementResponse struct {
	ID            uint   `json:"id"`
	InventoryID   uint   `json:"inventory_id"`
	MovementType  string `json:"movement_type"`
	Quantity      int    `json:"quantity"`
	PreviousStock int    `json:"previous_stock"`
	NewStock      int    `json:"new_stock"`
	ReferenceType string `json:"reference_type,omitempty"`
	ReferenceID   string `json:"reference_id,omitempty"`
	Notes         string `json:"notes,omitempty"`
	CreatedByID   uint   `json:"created_by_id"`
	CreatedAt     string `json:"created_at"`
}

type InventoryMovementListResponse struct {
	Movements []*InventoryMovementResponse `json:"movements"`
	Total     int                          `json:"total"`
}




