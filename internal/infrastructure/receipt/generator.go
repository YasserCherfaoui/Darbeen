package receipt

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/YasserCherfaoui/darween/internal/domain/pos"
)

// ReceiptGenerator handles PDF receipt generation for thermal printers
type ReceiptGenerator struct{}

// NewReceiptGenerator creates a new receipt generator
func NewReceiptGenerator() *ReceiptGenerator {
	return &ReceiptGenerator{}
}

// ReceiptData holds the data needed to generate a receipt
type ReceiptData struct {
	CompanyName   string
	FranchiseName string
	ReceiptNumber string
	Date          time.Time
	CustomerName  string
	CustomerEmail string
	Items         []ReceiptItem
	SubTotal      float64
	TaxAmount     float64
	DiscountAmount float64
	TotalAmount   float64
	Payments      []ReceiptPayment
}

// ReceiptItem represents an item on the receipt
type ReceiptItem struct {
	Name           string
	SKU            string
	Quantity       int
	UnitPrice      float64
	DiscountAmount float64
	TotalAmount    float64
}

// ReceiptPayment represents a payment on the receipt
type ReceiptPayment struct {
	Method string
	Amount float64
}

// GenerateReceipt creates a PDF receipt formatted for thermal printers (80mm width)
func (rg *ReceiptGenerator) GenerateReceipt(data *ReceiptData) ([]byte, error) {
	// Define the initialization options using the InitType struct
	// 80mm width (standard thermal receipt width), 297mm height (A4 length)
	initOptions := gofpdf.InitType{
		OrientationStr: "P",  // "P"ortrait
		UnitStr:        "mm", // millimeters
		FontDirStr:     "",   // Leave empty for default font directory
		Size: gofpdf.SizeType{
			Wd: 80,  // 80mm width
			Ht: 297, // 297mm height
		},
	}

	// Create PDF with custom size
	pdf := gofpdf.NewCustom(&initOptions)
	pdf.SetMargins(5, 5, 5)
	pdf.AddPage()

	// Set font for header
	pdf.SetFont("Helvetica", "B", 10)

	// Get business name
	businessName := data.CompanyName
	if data.FranchiseName != "" {
		businessName = data.FranchiseName
	}

	// Center-aligned header
	pdf.CellFormat(70, 5, businessName, "", 1, "C", false, 0, "")

	// Receipt details
	pdf.SetFont("Helvetica", "", 8)
	saleDate := data.Date.Format("2006-01-02 15:04")
	pdf.CellFormat(70, 4, "Receipt: "+data.ReceiptNumber, "", 1, "C", false, 0, "")
	pdf.CellFormat(70, 4, "Date: "+saleDate, "", 1, "C", false, 0, "")

	// Separator line
	pdf.CellFormat(70, 2, strings.Repeat("-", 35), "", 1, "C", false, 0, "")

	// Items section header
	pdf.SetFont("Courier", "", 8)
	pdf.CellFormat(35, 4, "Item", "B", 0, "L", false, 0, "")
	pdf.CellFormat(10, 4, "Qty", "B", 0, "R", false, 0, "")
	pdf.CellFormat(25, 4, "Price", "B", 1, "R", false, 0, "")

	// Items details
	pdf.SetFont("Helvetica", "", 7)
	for _, item := range data.Items {
		// Item name may need to be truncated
		itemName := item.Name
		if len(itemName) > 20 {
			itemName = itemName[:17] + "..."
		}

		// Print item name
		pdf.CellFormat(35, 4, itemName, "", 0, "L", false, 0, "")
		pdf.CellFormat(10, 4, strconv.Itoa(item.Quantity), "", 0, "R", false, 0, "")

		// Calculate item total with discount
		itemPrice := (item.UnitPrice*float64(item.Quantity) - item.DiscountAmount)
		pdf.CellFormat(25, 4, fmt.Sprintf("%.2f", itemPrice), "", 1, "R", false, 0, "")

		// If SKU is available, print as indented detail
		if item.SKU != "" {
			skuDesc := fmt.Sprintf("  SKU: %s", item.SKU)
			pdf.CellFormat(70, 3, skuDesc, "", 1, "L", false, 0, "")
		}
	}

	// Another separator
	pdf.CellFormat(70, 2, strings.Repeat("-", 35), "", 1, "C", false, 0, "")

	// Summary section
	pdf.SetFont("Courier", "", 8)
	pdf.CellFormat(45, 4, "Subtotal:", "", 0, "R", false, 0, "")
	pdf.CellFormat(25, 4, fmt.Sprintf("%.2f", data.SubTotal), "", 1, "R", false, 0, "")

	if data.DiscountAmount > 0 {
		pdf.CellFormat(45, 4, "Discount:", "", 0, "R", false, 0, "")
		pdf.CellFormat(25, 4, fmt.Sprintf("%.2f", data.DiscountAmount), "", 1, "R", false, 0, "")
	}

	if data.TaxAmount > 0 {
		pdf.CellFormat(45, 4, "Tax:", "", 0, "R", false, 0, "")
		pdf.CellFormat(25, 4, fmt.Sprintf("%.2f", data.TaxAmount), "", 1, "R", false, 0, "")
	}

	// Bold total
	pdf.SetFont("Courier", "B", 9)
	pdf.CellFormat(45, 5, "TOTAL:", "", 0, "R", false, 0, "")
	pdf.CellFormat(25, 5, fmt.Sprintf("%.2f", data.TotalAmount), "", 1, "R", false, 0, "")

	// Payment method
	pdf.SetFont("Courier", "", 8)
	paymentMethod := "Cash"
	if len(data.Payments) > 0 {
		method := strings.ToLower(data.Payments[0].Method)
		if len(method) > 0 {
			paymentMethod = strings.ToUpper(method[:1]) + method[1:]
		}
	}
	pdf.CellFormat(70, 4, "Payment: "+paymentMethod, "", 1, "L", false, 0, "")

	// Footer
	pdf.SetFont("Courier", "", 7)
	pdf.Ln(4)
	pdf.CellFormat(70, 3, "Thank you for your purchase!", "", 1, "C", false, 0, "")

	// Optional: Add receipt ID
	pdf.CellFormat(70, 3, "Receipt ID: "+data.ReceiptNumber, "", 1, "C", false, 0, "")

	// Output PDF to buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// ConvertSaleToReceiptData converts a pos.Sale to ReceiptData
// productVariantMap should map ProductVariantID to a struct with Name and SKU
type ProductVariantInfo struct {
	Name string
	SKU  string
}

func ConvertSaleToReceiptData(sale *pos.Sale, companyName, franchiseName string, productVariantMap map[uint]ProductVariantInfo) *ReceiptData {
	data := &ReceiptData{
		CompanyName:    companyName,
		FranchiseName:  franchiseName,
		ReceiptNumber:  sale.ReceiptNumber,
		Date:           sale.CreatedAt,
		SubTotal:       sale.SubTotal,
		TaxAmount:      sale.TaxAmount,
		DiscountAmount: sale.DiscountAmount,
		TotalAmount:    sale.TotalAmount,
	}
	
	// Customer info
	if sale.Customer != nil {
		data.CustomerName = sale.Customer.Name
		data.CustomerEmail = sale.Customer.Email
	}
	
	// Convert items
	data.Items = make([]ReceiptItem, len(sale.Items))
	for i, item := range sale.Items {
		variantInfo, exists := productVariantMap[item.ProductVariantID]
		itemName := fmt.Sprintf("Item #%d", i+1)
		itemSKU := ""
		if exists {
			itemName = variantInfo.Name
			itemSKU = variantInfo.SKU
		}
		
		data.Items[i] = ReceiptItem{
			Name:           itemName,
			SKU:            itemSKU,
			Quantity:       item.Quantity,
			UnitPrice:      item.UnitPrice,
			DiscountAmount: item.DiscountAmount,
			TotalAmount:    item.TotalAmount,
		}
	}
	
	// Convert payments
	data.Payments = make([]ReceiptPayment, len(sale.Payments))
	for i, payment := range sale.Payments {
		data.Payments[i] = ReceiptPayment{
			Method: string(payment.PaymentMethod),
			Amount: payment.Amount,
		}
	}
	
	return data
}

