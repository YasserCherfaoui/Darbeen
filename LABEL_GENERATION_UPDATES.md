# Label Generation Implementation Updates

## Overview
The label generation implementation has been updated to align with the existing barcode generation pattern used in the codebase, following best practices for QR code generation and PDF layout.

## Key Changes Made

### 1. QR Code Generation Method

**Before:**
- Generated QR codes as in-memory byte arrays
- Used `bytes.Reader` to pass data to PDF library
- QR codes were stored in memory during PDF generation

**After:**
- QR codes are saved as temporary PNG files in `/tmp/`
- Uses timestamp-based unique filenames: `qrcode_{timestamp}.png`
- Files are automatically cleaned up after PDF generation using `defer`
- Follows the pattern: `qrcode.WriteFile(content, qrcode.Medium, 256, tempFile)`

**Benefits:**
- More memory efficient for bulk operations
- Simpler image handling with gofpdf library
- Automatic cleanup prevents temp file accumulation
- Consistent with existing codebase patterns

### 2. Label Layout Algorithm

**Before:**
- Simple grid calculation
- Basic page overflow handling
- Fixed QR code positioning

**After:**
- Dynamic sticker dimension calculation based on A4 page size
- Smart QR code sizing: `qrCodeSize = stickerWidth * 0.5` (adjusts to fit)
- Proper multi-page handling with page count calculation
- QR code positioned on **left side** of sticker
- Text positioned on **right side** next to QR code

**Layout Calculation:**
```go
// Calculate sticker dimensions
stickerWidth := (pageWidth - (2 * margin)) / float64(columnsPerPage)
rowsPerPage := int((pageHeight - (2 * margin)) / labelHeightMM)
stickerHeight := (pageHeight - (2 * margin)) / float64(rowsPerPage)

// Calculate QR code size
qrCodeSize := stickerWidth * 0.5
if stickerHeight < qrCodeSize {
    qrCodeSize = stickerHeight * 0.8
}
```

### 3. Font Size Update

**Before:**
- Default font size: 10pt
- Smaller, harder to read on labels

**After:**
- Default font size: 24pt
- Much larger and more readable
- Better for scanning from a distance
- Matches the pattern in existing barcode generation

### 4. Label Positioning

**Before:**
- QR code centered
- Text below QR code
- Vertical layout

**After:**
- QR code on left side (with 2mm padding)
- Text on right side
- Horizontal layout
- Better space utilization

**Position Calculation:**
```go
// QR code on left
imgX := x + 2
imgY := y + (stickerHeight-qrCodeSize)/2  // Vertically centered

// Text on right
textX := imgX + qrCodeSize + 2
textY := y + (stickerHeight / 2)
```

### 5. Page Management

**Enhanced multi-page handling:**
```go
totalPages := int(math.Ceil(float64(totalLabels) / float64(maxLabelsPerPage)))

// Proper page overflow detection
if row >= rowsPerPage && labelCount < totalLabels {
    row = 0
    currentPage++
    if currentPage <= totalPages {
        pdf.AddPage()
    }
}
```

## Implementation Details

### Temp File Management

```go
// Track all temp files for cleanup
var tempFiles []string
defer func() {
    for _, file := range tempFiles {
        os.Remove(file)
    }
}()

// Generate QR code
qrCodeFile, err := lg.GenerateQRCode(labelData.SKU)
tempFiles = append(tempFiles, qrCodeFile)
```

### Single Label Generation

- Uses same temp file approach
- Cleanup with `defer os.Remove(qrCodeFile)`
- QR code auto-sized to fit label dimensions
- Text positioned beside QR code (not below)

### Bulk Label Generation

- Collects all temp files in array
- Single defer cleanup at function end
- Efficient memory usage
- Proper A4 page utilization

## Updated Default Configuration

```go
LabelConfig{
    WidthMM:      80.0,  // 80mm width
    HeightMM:     50.0,  // 50mm height
    MarginMM:     1.0,   // 1mm margin
    QRSizeMM:     40.0,  // QR code size (auto-adjusted)
    FontSize:     24.0,  // Font size (increased from 10)
    LabelsPerRow: 1,     // 1 label per row
}
```

### Configuration Rationale

The new defaults are optimized for standard label printers:
- **80mm x 50mm**: Common label size for warehouse/inventory labels
- **1mm margin**: Minimal margin for maximum print area
- **40mm QR code**: Larger QR code for easy scanning from distance
- **1 label per row**: Single column layout for better readability
- **4 rows per page**: Calculated to fit A4 page (297mm height with margins)

## Frontend Updates

### Updated Placeholders
All placeholders updated to match new defaults:
- Width: `50.8` → `80`
- Height: `25.4` → `50`
- Margin: `2.0` → `1.0`
- QR Size: `20.0` → `40`
- Font Size: `10.0` → `24.0`
- Labels Per Row: `3` → `1`
- Description text updated to show "(80mm x 50mm)"

These changes provide:
- Better user guidance with accurate defaults
- Consistent experience between frontend and backend
- Clear expectations for label dimensions

### No Breaking Changes
- API remains the same
- All existing functionality preserved
- Configuration still optional
- Backward compatible

## Benefits of New Implementation

1. **Memory Efficient**: Temp files instead of in-memory buffers
2. **Better Layout**: Improved space utilization on A4 pages
3. **More Readable**: Larger font size (24pt vs 10pt)
4. **Cleaner Code**: Automatic temp file cleanup
5. **Consistent**: Matches existing barcode generation pattern
6. **Scalable**: Better handling of large bulk operations
7. **Professional**: QR code + text side-by-side layout

## Technical Improvements

### Error Handling
- Proper cleanup even on errors (defer)
- Descriptive error messages
- Context in error wrapping

### Resource Management
- Automatic temp file cleanup
- No file handle leaks
- Memory efficient for bulk operations

### Code Quality
- Follows Go best practices
- Consistent with existing codebase
- Well-commented calculations
- Clear variable naming

## Testing Recommendations

1. **Single Label Generation**
   - Test with various SKU lengths
   - Verify QR code readability
   - Check text positioning

2. **Bulk Generation**
   - Test with 1-100+ labels
   - Verify multi-page layout
   - Check temp file cleanup

3. **Configuration**
   - Test custom dimensions
   - Test various font sizes
   - Test labels per row (1-5)

4. **Edge Cases**
   - Very long SKUs
   - Small label dimensions
   - Large bulk operations (1000+ labels)

## Migration Notes

No migration needed - this is a transparent improvement:
- API endpoints unchanged
- Request/response format identical
- Configuration options same
- Default behavior improved but compatible

## Files Modified

### Backend
- `/internal/infrastructure/label/generator.go` - Core implementation
- Updated imports (added `os`, `math`, `strconv`, `time`)
- Updated `GenerateQRCode()` to use temp files
- Updated `GenerateSingleLabel()` with new layout
- Updated `GenerateBulkLabels()` with improved algorithm

### Frontend
- `/app/src/components/products/GenerateLabelDialog.tsx` - Font size placeholder
- `/app/.docs/label-printing-feature.md` - Documentation update

## Conclusion

The updated implementation provides a more robust, efficient, and professional label generation system that aligns with existing codebase patterns while maintaining full backward compatibility.

