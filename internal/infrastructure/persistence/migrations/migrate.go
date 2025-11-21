package migrations

import (
	"log"

	"github.com/YasserCherfaoui/darween/internal/domain/company"
	"github.com/YasserCherfaoui/darween/internal/domain/emailqueue"
	"github.com/YasserCherfaoui/darween/internal/domain/franchise"
	"github.com/YasserCherfaoui/darween/internal/domain/inventory"
	"github.com/YasserCherfaoui/darween/internal/domain/pos"
	"github.com/YasserCherfaoui/darween/internal/domain/product"
	"github.com/YasserCherfaoui/darween/internal/domain/smtpconfig"
	"github.com/YasserCherfaoui/darween/internal/domain/subscription"
	"github.com/YasserCherfaoui/darween/internal/domain/supplier"
	"github.com/YasserCherfaoui/darween/internal/domain/user"
	"github.com/YasserCherfaoui/darween/internal/domain/warehousebill"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("Running auto-migration...")

	// AutoMigrate will create tables, missing columns, missing indexes, and foreign key constraints
	// migrate only when gin_mode is release
	if gin.Mode() == gin.ReleaseMode {
		log.Println("Running auto-migration in release mode")
		err := db.AutoMigrate(
			&user.User{},
			&company.Company{},
			&subscription.Subscription{},
			&user.UserCompanyRole{},
			&franchise.Franchise{},
			&user.UserFranchiseRole{},
			&supplier.Supplier{},
			&supplier.SupplierBill{},
			&supplier.SupplierBillItem{},
			&supplier.SupplierPayment{},
			&supplier.SupplierPaymentDistribution{},
			&product.Product{},
			&product.ProductVariant{},
			&inventory.Inventory{},
			&inventory.InventoryMovement{},
			&franchise.FranchisePricing{},
			&pos.Customer{},
			&pos.Sale{},
			&pos.SaleItem{},
			&pos.Payment{},
			&pos.CashDrawer{},
			&pos.CashDrawerTransaction{},
			&pos.Refund{},
			&warehousebill.WarehouseBill{},
			&warehousebill.WarehouseBillItem{},
			&smtpconfig.SMTPConfig{},
			&emailqueue.EmailQueue{},
		)

		if err != nil {
			log.Printf("Auto-migration failed: %v", err)
			return err
		}
	}

	log.Println("Auto-migration completed successfully")
	return nil
}
