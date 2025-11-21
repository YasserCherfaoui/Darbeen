package inventory

type Repository interface {
	// Inventory CRUD
	Create(inventory *Inventory) error
	FindByID(id uint) (*Inventory, error)
	FindByVariantAndCompany(variantID, companyID uint) (*Inventory, error)
	FindByVariantAndFranchise(variantID, franchiseID uint) (*Inventory, error)
	FindByCompany(companyID uint) ([]*Inventory, error)
	FindByFranchise(franchiseID uint) ([]*Inventory, error)
	Update(inventory *Inventory) error
	Delete(id uint) error

	// Inventory Movements (Audit Trail)
	CreateMovement(movement *InventoryMovement) error
	FindMovementsByInventory(inventoryID uint, limit int) ([]*InventoryMovement, error)
	FindMovementsByInventoryWithFilters(inventoryID uint, movementType *string, startDate *string, endDate *string, page, limit int) ([]*InventoryMovement, int64, error)
	FindMovementsByReference(referenceType string, referenceID string) ([]*InventoryMovement, error)
}




