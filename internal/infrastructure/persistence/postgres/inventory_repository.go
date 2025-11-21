package postgres

import (
	"fmt"

	"github.com/YasserCherfaoui/darween/internal/domain/inventory"
	"gorm.io/gorm"
)

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) inventory.Repository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) Create(inv *inventory.Inventory) error {
	return r.db.Create(inv).Error
}

func (r *inventoryRepository) FindByID(id uint) (*inventory.Inventory, error) {
	var inv inventory.Inventory
	err := r.db.Where("id = ?", id).First(&inv).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("inventory not found")
		}
		return nil, err
	}
	return &inv, nil
}

func (r *inventoryRepository) FindByVariantAndCompany(variantID, companyID uint) (*inventory.Inventory, error) {
	var inv inventory.Inventory
	err := r.db.Where("product_variant_id = ? AND company_id = ?", variantID, companyID).First(&inv).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("inventory not found")
		}
		return nil, err
	}
	return &inv, nil
}

func (r *inventoryRepository) FindByVariantAndFranchise(variantID, franchiseID uint) (*inventory.Inventory, error) {
	var inv inventory.Inventory
	err := r.db.Where("product_variant_id = ? AND franchise_id = ?", variantID, franchiseID).First(&inv).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("inventory not found")
		}
		return nil, err
	}
	return &inv, nil
}

func (r *inventoryRepository) FindByCompany(companyID uint) ([]*inventory.Inventory, error) {
	var inventories []*inventory.Inventory
	err := r.db.Where("company_id = ? AND is_active = ?", companyID, true).Find(&inventories).Error
	return inventories, err
}

func (r *inventoryRepository) FindByFranchise(franchiseID uint) ([]*inventory.Inventory, error) {
	var inventories []*inventory.Inventory
	err := r.db.Where("franchise_id = ? AND is_active = ?", franchiseID, true).Find(&inventories).Error
	return inventories, err
}

func (r *inventoryRepository) Update(inv *inventory.Inventory) error {
	return r.db.Save(inv).Error
}

func (r *inventoryRepository) Delete(id uint) error {
	return r.db.Delete(&inventory.Inventory{}, id).Error
}

// Inventory Movements

func (r *inventoryRepository) CreateMovement(movement *inventory.InventoryMovement) error {
	return r.db.Create(movement).Error
}

func (r *inventoryRepository) FindMovementsByInventory(inventoryID uint, limit int) ([]*inventory.InventoryMovement, error) {
	var movements []*inventory.InventoryMovement
	query := r.db.Where("inventory_id = ?", inventoryID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&movements).Error
	return movements, err
}

func (r *inventoryRepository) FindMovementsByInventoryWithFilters(inventoryID uint, movementType *string, startDate *string, endDate *string, page, limit int) ([]*inventory.InventoryMovement, int64, error) {
	var movements []*inventory.InventoryMovement
	var total int64

	query := r.db.Model(&inventory.InventoryMovement{}).Where("inventory_id = ?", inventoryID)

	// Apply movement type filter
	if movementType != nil && *movementType != "" {
		query = query.Where("LOWER(movement_type) = LOWER(?)", *movementType)
	}

	// Apply date range filter
	if startDate != nil && *startDate != "" {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil && *endDate != "" {
		query = query.Where("created_at <= ?", *endDate)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&movements).Error

	return movements, total, err
}

func (r *inventoryRepository) FindMovementsByReference(referenceType string, referenceID string) ([]*inventory.InventoryMovement, error) {
	var movements []*inventory.InventoryMovement
	err := r.db.
		Where("reference_type = ? AND reference_id = ?", referenceType, referenceID).
		Order("created_at DESC").
		Find(&movements).Error
	return movements, err
}




