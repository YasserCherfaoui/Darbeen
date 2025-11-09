package migrations

import (
	"log"

	"github.com/YasserCherfaoui/darween/internal/domain/company"
	"github.com/YasserCherfaoui/darween/internal/domain/franchise"
	"github.com/YasserCherfaoui/darween/internal/domain/inventory"
	"github.com/YasserCherfaoui/darween/internal/domain/product"
	"github.com/YasserCherfaoui/darween/internal/domain/subscription"
	"github.com/YasserCherfaoui/darween/internal/domain/supplier"
	"github.com/YasserCherfaoui/darween/internal/domain/user"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("Running auto-migration...")

	// AutoMigrate will create tables, missing columns, missing indexes, and foreign key constraints
	err := db.AutoMigrate(
		&user.User{},
		&company.Company{},
		&subscription.Subscription{},
		&user.UserCompanyRole{},
		&franchise.Franchise{},
		&user.UserFranchiseRole{},
		&supplier.Supplier{},
		&product.Product{},
		&product.ProductVariant{},
		&inventory.Inventory{},
		&inventory.InventoryMovement{},
		&franchise.FranchisePricing{},
	)

	if err != nil {
		log.Printf("Auto-migration failed: %v", err)
		return err
	}

	log.Println("Auto-migration completed successfully")
	return nil
}
