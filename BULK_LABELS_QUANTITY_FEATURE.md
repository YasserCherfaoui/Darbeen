# Bulk Labels with Quantity Support - Enhanced Feature

## Overview
The bulk label printing feature has been significantly enhanced to support variant selection with individual quantities using an intuitive combobox interface. This allows users to specify exactly how many labels they need for each variant.

## Key Enhancements

### 1. Quantity-Based Selection

**Before:**
- Simple checkbox selection (1 label per item)
- No way to specify quantities
- Limited to selecting which items to print

**After:**
- Each variant can have a specific quantity
- Quantities are editable per variant
- Total label count is calculated and displayed
- Labels are generated based on quantities (e.g., 5 qty = 5 identical labels)

### 2. Combobox Search Interface

**Before:**
- Scrollable list with checkboxes
- Difficult to find specific variants
- No search functionality

**After:**
- Searchable combobox powered by `cmdk`
- Search by variant name, product name, or SKU
- Shows parent product name + variant name for clarity
- Visual checkmarks for already-selected variants
- Keyboard navigation support

### 3. Improved Variant Display

**Format:** `[Parent Product Name] - [Variant Name]`

Example displays:
- "Nike Air Max - Size 42 / Red"
- "iPhone 15 - 256GB / Blue"
- "Office Chair - Black / Ergonomic"

This makes it easy to differentiate between variants of different products.

## Backend Changes

### Updated DTOs

```go
// New structures for quantity-based requests
type VariantQuantity struct {
    VariantID uint `json:"variant_id" binding:"required"`
    Quantity  int  `json:"quantity" binding:"required,min=1"`
}

type ProductQuantity struct {
    ProductID uint `json:"product_id" binding:"required"`
    Quantity  int  `json:"quantity" binding:"required,min=1"`
}

type GenerateBulkLabelsRequest struct {
    Products []ProductQuantity `json:"products,omitempty"`
    Variants []VariantQuantity `json:"variants,omitempty"`
    Config   *LabelConfig      `json:"config,omitempty"`
}
```

### Updated Service Logic

The `GenerateBulkLabels` service now:
1. Accepts arrays of products/variants with quantities
2. Repeats label data based on quantity
3. Generates PDF with all labels in sequence

Example: If you select Variant A (qty: 3) and Variant B (qty: 2), the PDF will contain:
- Variant A label
- Variant A label
- Variant A label
- Variant B label
- Variant B label

### API Request Format

```json
{
  "variants": [
    {
      "variant_id": 123,
      "quantity": 5
    },
    {
      "variant_id": 456,
      "quantity": 3
    }
  ],
  "products": [
    {
      "product_id": 789,
      "quantity": 2
    }
  ],
  "config": {
    "width_mm": 80,
    "height_mm": 50,
    "margin_mm": 1,
    "labels_per_row": 1
  }
}
```

## Frontend Implementation

### New Component: GenerateBulkLabelsDialogEnhanced

Features:
- **Combobox Search:** Fast variant lookup
- **Selected Variants List:** Shows all selected variants with quantities
- **Quantity Editor:** Inline number input for each variant
- **Remove Button:** Easy removal of unwanted selections
- **Total Counter:** Displays total number of labels to be generated
- **Parent Product Display:** Shows "Product Name - Variant Name"

### Component Structure

```tsx
<Dialog>
  {/* Header with title and description */}
  
  {/* Alerts for errors/success */}
  
  <Form>
    {/* Combobox for variant selection */}
    <Popover>
      <Command>
        <CommandInput placeholder="Search..." />
        <CommandList>
          {/* Searchable variant list */}
        </CommandList>
      </Command>
    </Popover>
    
    {/* Selected variants with quantities */}
    <div className="selected-variants">
      {selectedVariants.map(variant => (
        <div>
          <span>{variant.name}</span>
          <Input type="number" value={quantity} />
          <Badge>{quantity} labels</Badge>
          <Button onClick={remove} />
        </div>
      ))}
    </div>
    
    {/* Configuration inputs */}
    {/* Submit button */}
  </Form>
</Dialog>
```

### User Flow

1. **Open Dialog:** Click "Bulk Generate Labels"
2. **Search Variant:** Type in combobox to filter
3. **Select Variant:** Click on a variant from dropdown
4. **Set Quantity:** Adjust quantity using number input (default: 1)
5. **Add More:** Repeat steps 2-4 for additional variants
6. **Review:** Check total label count
7. **Generate:** Click "Generate X Labels" button
8. **Download:** PDF automatically downloads

### State Management

```tsx
interface SelectedVariant {
  variant_id: number
  variant: VariantWithProduct
  quantity: number
}

// State
const [selectedVariants, setSelectedVariants] = useState<SelectedVariant[]>([])
const [allVariants, setAllVariants] = useState<VariantWithProduct[]>([])

// Actions
- addVariant(variant) - Add with qty 1
- removeVariant(variantId) - Remove from selection
- updateQuantity(variantId, quantity) - Update qty
- getTotalLabels() - Calculate total
```

## UI Components Added

### 1. Command Component (`command.tsx`)
Provides the command palette interface:
- `Command` - Main container
- `CommandInput` - Search input
- `CommandList` - Results container
- `CommandItem` - Individual result
- `CommandEmpty` - No results message
- `CommandGroup` - Group results

### 2. Popover Component (`popover.tsx`)
Provides dropdown functionality:
- `Popover` - Container
- `PopoverTrigger` - Button to open
- `PopoverContent` - Dropdown content

### Dependencies Installed

```bash
npm install @radix-ui/react-popover cmdk --legacy-peer-deps
```

## Benefits

### For Users
1. **Precise Control:** Specify exact quantities needed
2. **Easy Search:** Find variants quickly by name or SKU
3. **Clear Identification:** See parent product + variant name
4. **Flexible Selection:** Mix different variants with different quantities
5. **Visual Feedback:** See total labels before generating
6. **Efficient Workflow:** No need to generate labels multiple times

### For Operations
1. **Inventory Labeling:** Print multiple labels for stock items
2. **Batch Processing:** Prepare labels for incoming shipments
3. **Variant Management:** Easy differentiation between similar products
4. **Time Savings:** Bulk operations with precise quantities
5. **Error Reduction:** Clear identification prevents mislabeling

## Example Use Cases

### Use Case 1: New Stock Arrival
```
Inventory receives:
- Nike Air Max Size 42 Red (15 units)
- Nike Air Max Size 43 Black (10 units)
- Nike Air Max Size 44 White (8 units)

User actions:
1. Open bulk labels dialog
2. Search "Nike Air Max"
3. Select "Size 42 Red", set qty: 15
4. Select "Size 43 Black", set qty: 10
5. Select "Size 44 White", set qty: 8
6. Generate 33 labels total
```

### Use Case 2: Shelf Restocking
```
Need to add labels for:
- Office Chair Black (3 new units)
- Office Chair White (2 new units)

User actions:
1. Search "Office Chair"
2. Select "Black", qty: 3
3. Select "White", qty: 2
4. Generate 5 labels
```

## Comparison: Old vs New

| Feature | Old Version | New Version |
|---------|-------------|-------------|
| Selection | Checkboxes | Searchable combobox |
| Quantities | 1 per item | Custom per item |
| Search | Scroll only | Full text search |
| Variant Display | Name only | Product + Variant |
| Total Preview | Selected count | Total labels |
| Edit Quantities | Not possible | Inline editing |
| Remove Items | Uncheck box | Remove button |
| Keyboard Nav | Limited | Full support |

## Technical Details

### Variant Fetching
```tsx
// Fetches all variants from all products
const fetchAllVariants = async () => {
  const variantsPromises = products.map(async (product) => {
    const response = await apiClient.productVariants.list(companyId, product.id)
    return response.data.map(variant => ({
      ...variant,
      parent_product_name: product.name,
    }))
  })
  
  const variantsArrays = await Promise.all(variantsPromises)
  const flatVariants = variantsArrays.flat()
  setAllVariants(flatVariants)
}
```

### Search Implementation
Uses `cmdk` library for fuzzy search:
- Searches across: parent product name, variant name, and SKU
- Case-insensitive matching
- Instant results
- Keyboard navigation (arrow keys, enter)

### Quantity Validation
- Minimum: 1 label
- Maximum: Limited by browser/system memory
- Type: Integer only
- Default: 1 when newly added

## Future Enhancements

Potential improvements:
1. **Templates:** Save common selections
2. **Quick Quantities:** Buttons for 5, 10, 25, 50
3. **Import CSV:** Bulk add from file
4. **Recent Selections:** Quick access to recently used variants
5. **Duplicate Detection:** Warn if same variant added twice
6. **Batch Editing:** Update multiple quantities at once
7. **Export Selection:** Save selection for later
8. **Print Preview:** See labels before generating

## Backward Compatibility

✅ The old `GenerateBulkLabelsDialog` component is still available
✅ API accepts both old and new request formats
✅ No breaking changes to existing functionality
✅ Frontend can easily switch between dialogs if needed

## Files Modified

### Backend
- `/internal/application/product/dto.go` - Added quantity structures
- `/internal/application/product/service.go` - Updated bulk generation logic
- `/app/src/lib/api-client.ts` - Updated TypeScript types

### Frontend
- `/app/src/components/products/GenerateBulkLabelsDialogEnhanced.tsx` - New component
- `/app/src/components/products/ProductsTable.tsx` - Use new dialog
- `/app/src/components/ui/command.tsx` - New component
- `/app/src/components/ui/popover.tsx` - New component

## Conclusion

The enhanced bulk label printing feature provides a significantly improved user experience with powerful search, precise quantity control, and clear variant identification. This makes label generation faster, more accurate, and better suited for real-world warehouse and inventory operations.

