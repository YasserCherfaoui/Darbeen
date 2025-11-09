# Complete POS System - User Guide

## 🎉 System Overview

Your ERP now has a fully functional Point of Sale (POS) system with all essential features for managing sales, customers, and cash operations.

## 🔐 Accessing the POS System

### Method 1: Sidebar Navigation (Recommended)
1. **Login** to your ERP account
2. **Select a company** from the company selector
3. Look for the **"POS"** menu item in the sidebar (with a 🛒 shopping cart icon)
4. Click on **"POS"** to access the POS dashboard

### Method 2: Direct URL
Navigate directly to: `/companies/{companyId}/pos/`

### User Permissions
- **Minimum Role**: Employee
- Users must have access to the company/franchise to use POS features

## 📋 Available Features

### 1. **POS Dashboard** (`/companies/{companyId}/pos/`)
The main hub with quick access to all POS features:
- Create new sales
- View sales history
- Manage customers
- Open/close cash drawer
- Quick start guide

### 2. **Create Sale** (`/companies/{companyId}/pos/sales/new`)

**Features:**
- **Product Search**: Search products by name or SKU
- **Shopping Cart**: Add multiple items with quantity controls
- **Individual Discounts**: Apply discounts to specific items
- **Customer Selection**: Link sales to customers (optional)
- **Walk-in Sales**: Process sales without customer info
- **Sale Notes**: Add notes to transactions
- **Real-time Calculations**: Automatic subtotal and total calculations
- **Payment Processing**: Complete checkout with multiple payment methods
- **Receipt Generation**: Automatic receipt preview after sale

**Workflow:**
1. Search and add products to cart
2. Adjust quantities and discounts as needed
3. Optionally select a customer
4. Click "Checkout"
5. Choose payment method (cash/card/other)
6. For cash: Enter amount received, get automatic change calculation
7. Complete payment
8. View/print receipt

### 3. **Sales History** (`/companies/{companyId}/pos/sales/`)

**Features:**
- View all sales transactions
- See receipt numbers, dates, customers
- Check payment status (paid/partially paid/unpaid/refunded)
- View sale details including items and payments
- Color-coded status badges

**Information Displayed:**
- Receipt number
- Transaction date
- Customer name (or "Walk-in")
- Number of items
- Total amount
- Payment status
- Detailed view with full breakdown

### 4. **Customer Management** (`/companies/{companyId}/pos/customers/`)

**Features:**
- Create new customers
- Edit customer information
- View customer purchase history
- Search customers
- Track total purchases per customer

**Customer Information:**
- Name (required)
- Email
- Phone
- Address
- Total purchases (automatic tracking)

### 5. **Cash Drawer Management** (`/companies/{companyId}/pos/cash-drawer/`)

**Features:**
- Open cash drawer at start of shift
- View current drawer status
- See expected balance vs actual
- Track all cash transactions
- Close drawer with reconciliation
- View drawer history
- Identify overages/shortages

**Cash Drawer Workflow:**

**Opening:**
1. Click "Open Cash Drawer"
2. Enter starting cash amount
3. Add optional notes
4. Confirm to open

**During Shift:**
- View opening balance
- See expected balance (based on transactions)
- Track number of transactions
- Monitor all cash sales

**Closing:**
1. Click "Close Cash Drawer"
2. Count actual cash
3. Enter closing balance
4. System calculates expected vs actual
5. Shows difference (overage in green, shortage in red)
6. Add closing notes
7. Confirm to close

### 6. **Receipt Preview**
- Professional receipt layout
- Company and franchise information
- Customer details (if linked)
- Itemized list with prices
- Discounts and taxes
- Payment information
- Unique receipt number
- Timestamp
- Print functionality

## 🔄 Complete Sale Workflow

### Standard Sale Process:
1. **Open Cash Drawer** (if using cash)
2. **Create New Sale**
   - Search and add products
   - Adjust quantities
   - Apply discounts if needed
3. **Select Customer** (optional)
4. **Add Sale Notes** (optional)
5. **Checkout**
6. **Process Payment**
   - Choose payment method
   - Enter amount (for cash, get change calculation)
   - Complete payment
7. **View Receipt**
   - Preview receipt
   - Print if needed
8. **Continue or Close**

### End of Day:
1. Go to Cash Drawer Management
2. Count physical cash
3. Close cash drawer
4. Review reconciliation report
5. Note any discrepancies

## 💡 Key Features Explained

### Automatic Inventory Integration
- ✅ Real-time inventory checking before sale
- ✅ Automatic stock deduction on completed sales
- ✅ Inventory movements logged with sale references
- ✅ Low stock prevention (can't sell more than available)

### Payment Processing
- ✅ Multiple payment methods (cash, card, other)
- ✅ Change calculation for cash payments
- ✅ Partial payment support
- ✅ Payment status tracking
- ✅ Automatic cash drawer updates for cash payments

### Customer Management
- ✅ Quick customer search and selection
- ✅ Create customers on-the-fly during checkout
- ✅ Track total purchase history
- ✅ Walk-in customer support (no customer required)

### Discounts
- ✅ Item-level discounts
- ✅ Sale-level discounts
- ✅ Real-time total updates

### Reporting
- ✅ Sales history with filters
- ✅ Cash drawer reconciliation
- ✅ Transaction tracking
- ✅ Payment method breakdown

## 🎯 Components Created

### UI Components (`/app/src/components/pos/`):
1. **CustomerForm.tsx** - Create/edit customers
2. **PaymentDialog.tsx** - Process payments with change calculation
3. **CashDrawerDialog.tsx** - Open/close cash drawer
4. **ProductSearch.tsx** - Search and select products
5. **Cart.tsx** - Shopping cart with quantity and discount controls
6. **CustomerSelector.tsx** - Search and select customers
7. **ReceiptPreview.tsx** - Preview and print receipts

### Route Pages (`/app/src/routes/`):
1. **companies.$companyId.pos.index.tsx** - POS dashboard
2. **companies.$companyId.pos.sales.new.tsx** - Create sale interface
3. **companies.$companyId.pos.sales.index.tsx** - Sales history
4. **companies.$companyId.pos.customers.index.tsx** - Customer management
5. **companies.$companyId.pos.cash-drawer.index.tsx** - Cash drawer management
6. **franchises.$franchiseId.pos.index.tsx** - Franchise POS dashboard

## 🔧 Technical Details

### Backend API Endpoints
All endpoints require authentication:

**Customers:**
- `POST /api/v1/companies/:companyId/pos/customers` - Create customer
- `GET /api/v1/companies/:companyId/pos/customers` - List customers
- `GET /api/v1/companies/:companyId/pos/customers/:id` - Get customer
- `PUT /api/v1/companies/:companyId/pos/customers/:id` - Update customer
- `DELETE /api/v1/companies/:companyId/pos/customers/:id` - Delete customer

**Sales:**
- `POST /api/v1/companies/:companyId/pos/sales` - Create sale
- `GET /api/v1/companies/:companyId/pos/sales` - List sales
- `GET /api/v1/companies/:companyId/pos/sales/:id` - Get sale details
- `POST /api/v1/companies/:companyId/pos/sales/:id/payments` - Add payment
- `POST /api/v1/companies/:companyId/pos/sales/:id/refund` - Process refund

**Cash Drawer:**
- `POST /api/v1/companies/:companyId/pos/cash-drawer/open` - Open drawer
- `GET /api/v1/companies/:companyId/pos/cash-drawer/active` - Get active drawer
- `PUT /api/v1/companies/:companyId/pos/cash-drawer/:id/close` - Close drawer
- `GET /api/v1/companies/:companyId/pos/cash-drawer` - List drawer history

**Reports:**
- `POST /api/v1/companies/:companyId/pos/reports/sales` - Sales report

### Database Tables
- `customers` - Customer information
- `sales` - Sale transactions
- `sale_items` - Line items
- `payments` - Payment records
- `cash_drawers` - Cash drawer sessions
- `cash_drawer_transactions` - Transaction log
- `refunds` - Refund records

## 🚀 Getting Started

### First Time Setup:
1. **Login** to the system
2. **Select your company**
3. Navigate to **POS** from the sidebar
4. **Create some customers** (optional but recommended)
5. Ensure you have **products with inventory**
6. **Open the cash drawer** if accepting cash
7. **Create your first sale**!

### Daily Operations:
1. **Morning**: Open cash drawer
2. **Throughout Day**: Process sales as needed
3. **Evening**: Close cash drawer and reconcile

## 📝 Tips & Best Practices

1. **Always open cash drawer** before accepting cash payments
2. **Link customers** to sales for better tracking and marketing
3. **Use sale notes** for special instructions or details
4. **Reconcile cash drawer daily** to identify discrepancies quickly
5. **Check inventory** before promising products to customers
6. **Keep receipts** for returns and exchanges
7. **Review sales history** regularly for insights

## 🐛 Troubleshooting

**Can't see POS in sidebar?**
- Ensure you have at least Employee role in the company
- Refresh the page
- Check that you've selected a company

**Product search not working?**
- Ensure products exist in the system
- Check that product variants are active
- Verify products have proper SKUs and names

**Can't complete sale?**
- Verify inventory is available
- Check that all required fields are filled
- Ensure payment amount is sufficient
- Confirm cash drawer is open (for cash payments)

**Cash drawer won't open?**
- Check if a drawer is already open
- Verify you have proper permissions
- Ensure you're in the correct company/franchise

## ✨ Future Enhancements

Possible future additions:
- Barcode scanning
- Receipt email functionality
- Advanced reporting dashboard
- Customer loyalty programs
- Product quick access buttons
- Multi-tender payments
- Gift cards and store credit
- Returns and exchange processing
- Shift reports
- Employee performance tracking

## 📞 Support

For issues or questions:
1. Check this guide first
2. Review error messages carefully
3. Verify your permissions
4. Check backend logs for API errors
5. Contact your system administrator

---

**Congratulations!** Your POS system is fully operational and ready to process sales. Happy selling! 🎉

