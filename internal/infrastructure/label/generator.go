package label

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

// LabelConfig holds configuration for label generation
type LabelConfig struct {
	WidthMM    float64 // Label width in millimeters
	HeightMM   float64 // Label height in millimeters
	MarginMM   float64 // Margin in millimeters
	QRSizeMM   float64 // QR code size in millimeters
	FontSize   float64 // Font size for text
	LabelsPerRow int   // Number of labels per row for bulk generation
}

// DefaultLabelConfig returns default configuration
func DefaultLabelConfig() LabelConfig {
	return LabelConfig{
		WidthMM:      80.0,  // 80mm width
		HeightMM:     50.0,  // 50mm height
		MarginMM:     1.0,   // 1mm margin
		QRSizeMM:     40.0,  // QR code size (adjusted for larger label)
		FontSize:     24.0,  // Font size (larger for better readability)
		LabelsPerRow: 1,     // 1 label per row
	}
}

// LabelGenerator handles PDF label generation with QR codes
type LabelGenerator struct {
	config LabelConfig
}

// NewLabelGenerator creates a new label generator with the given config
func NewLabelGenerator(config LabelConfig) *LabelGenerator {
	return &LabelGenerator{
		config: config,
	}
}

// GenerateQRCode generates a QR code image from the given text and saves to temp file
func (lg *LabelGenerator) GenerateQRCode(text string) (string, error) {
	// Create temporary file path
	tempFile := "/tmp/qrcode_" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".png"
	
	// Generate QR code and save to file
	err := qrcode.WriteFile(text, qrcode.Medium, 256, tempFile)
	if err != nil {
		return "", fmt.Errorf("failed to create QR code: %w", err)
	}

	return tempFile, nil
}

// GenerateSingleLabel creates a PDF with a single label containing QR code and SKU text
func (lg *LabelGenerator) GenerateSingleLabel(sku string) ([]byte, error) {
	// Create new PDF
	pdf := gofpdf.New("P", "mm", "", "")
	pdf.SetMargins(0, 0, 0)
	pdf.SetAutoPageBreak(false, 0)
	
	// Add page with custom size (label dimensions)
	pdf.AddPageFormat("P", gofpdf.SizeType{
		Wd: lg.config.WidthMM,
		Ht: lg.config.HeightMM,
	})

	// Generate QR code to temp file
	qrCodeFile, err := lg.GenerateQRCode(sku)
	if err != nil {
		return nil, err
	}
	defer os.Remove(qrCodeFile) // Clean up temp file

	// Calculate QR code size to fit the label
	qrCodeSize := lg.config.QRSizeMM
	if qrCodeSize > lg.config.WidthMM*0.5 {
		qrCodeSize = lg.config.WidthMM * 0.5
	}
	if qrCodeSize > lg.config.HeightMM*0.8 {
		qrCodeSize = lg.config.HeightMM * 0.8
	}

	// Calculate positions - QR code on left, text on right
	imgX := lg.config.MarginMM
	imgY := (lg.config.HeightMM - qrCodeSize) / 2 // Center vertically

	// Draw QR code image
	pdf.Image(qrCodeFile, imgX, imgY, qrCodeSize, qrCodeSize, false, "", 0, "")

	// Add SKU text next to QR code
	pdf.SetFont("Arial", "", lg.config.FontSize)
	textX := imgX + qrCodeSize + 2
	textY := lg.config.HeightMM / 2
	pdf.SetXY(textX, textY)
	
	// Calculate remaining width for text
	textWidth := lg.config.WidthMM - qrCodeSize - lg.config.MarginMM - 4
	pdf.CellFormat(textWidth, 5, sku, "", 0, "L", false, 0, "")

	// Output PDF to buffer
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// LabelData represents data for a single label
type LabelData struct {
	SKU  string
	Name string // Optional product name
}

// GenerateBulkLabels creates a PDF with multiple labels in a grid layout
func (lg *LabelGenerator) GenerateBulkLabels(labels []LabelData) ([]byte, error) {
	if len(labels) == 0 {
		return nil, fmt.Errorf("no labels to generate")
	}

	// Create new PDF with A4 page size
	pdf := gofpdf.New("P", "mm", "A4", "")
	margin := lg.config.MarginMM
	pdf.SetMargins(margin, margin, margin)
	pdf.SetAutoPageBreak(true, margin)

	// Calculate sticker dimensions based on A4 page and desired grid
	pageWidth := 210.0  // A4 width in mm
	pageHeight := 297.0 // A4 height in mm
	
	columnsPerPage := lg.config.LabelsPerRow
	
	// Calculate sticker dimensions to fit the page
	stickerWidth := (pageWidth - (2 * margin)) / float64(columnsPerPage)
	
	// Calculate rows that fit on the page based on label height
	rowsPerPage := int((pageHeight - (2 * margin)) / lg.config.HeightMM)
	if rowsPerPage < 1 {
		rowsPerPage = 1
	}
	
	stickerHeight := (pageHeight - (2 * margin)) / float64(rowsPerPage)

	// Calculate QR code dimensions to fit the sticker
	qrCodeSize := stickerWidth * 0.5 // Using half of sticker width for QR code
	if stickerHeight < qrCodeSize {
		qrCodeSize = stickerHeight * 0.8
	}

	// Track position and page management
	col := 0
	row := 0
	labelCount := 0
	maxLabelsPerPage := columnsPerPage * rowsPerPage
	totalLabels := len(labels)

	// Calculate total pages needed
	totalPages := int(math.Ceil(float64(totalLabels) / float64(maxLabelsPerPage)))
	currentPage := 1

	// Track temp files for cleanup
	var tempFiles []string
	defer func() {
		// Clean up all temp files
		for _, file := range tempFiles {
			os.Remove(file)
		}
	}()

	// Add the first page
	pdf.AddPage()

	for _, labelData := range labels {
		// Calculate sticker position
		x := margin + (float64(col) * stickerWidth)
		y := margin + (float64(row) * stickerHeight)

		// Generate QR code to temp file
		qrCodeFile, err := lg.GenerateQRCode(labelData.SKU)
		if err != nil {
			return nil, fmt.Errorf("failed to generate QR code for %s: %w", labelData.SKU, err)
		}
		tempFiles = append(tempFiles, qrCodeFile)

		// Add QR code image to PDF, positioned on left side of sticker
		imgX := x + 2
		imgY := y + (stickerHeight-qrCodeSize)/2 // Center vertically
		pdf.Image(qrCodeFile, imgX, imgY, qrCodeSize, qrCodeSize, false, "", 0, "")

		// Add SKU text next to QR code
		pdf.SetFont("Arial", "", lg.config.FontSize)
		textX := imgX + qrCodeSize + 2
		textY := y + (stickerHeight / 2)
		pdf.SetXY(textX, textY)
		pdf.CellFormat(stickerWidth-qrCodeSize-4, 5, labelData.SKU, "", 0, "L", false, 0, "")

		// Increment label count
		labelCount++

		// Move to next position
		col++
		if col >= columnsPerPage {
			col = 0
			row++
			if row >= rowsPerPage && labelCount < totalLabels {
				// Reset row and add new page
				row = 0
				currentPage++
				if currentPage <= totalPages {
					pdf.AddPage()
				}
			}
		}
	}

	// Output PDF to buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate bulk PDF: %w", err)
	}

	return buf.Bytes(), nil
}

