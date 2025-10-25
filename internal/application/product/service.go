package product

import (
	"encoding/json"

	"github.com/YasserCherfaoui/darween/internal/domain/product"
	"github.com/YasserCherfaoui/darween/internal/domain/user"
	"github.com/YasserCherfaoui/darween/pkg/errors"
)

type Service struct {
	productRepo product.Repository
	userRepo    user.Repository
}

func NewService(productRepo product.Repository, userRepo user.Repository) *Service {
	return &Service{
		productRepo: productRepo,
		userRepo:    userRepo,
	}
}

// Product operations
func (s *Service) CreateProduct(userID, companyID uint, req *CreateProductRequest) (*ProductResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	// Check if SKU already exists in company
	existingProduct, _ := s.productRepo.FindProductBySKUAndCompany(req.SKU, companyID)
	if existingProduct != nil {
		return nil, errors.NewConflictError("product with this SKU already exists in the company")
	}

	// Create product
	newProduct := req.ToProduct(companyID)
	if !newProduct.IsValid() {
		return nil, errors.NewValidationError("invalid product data")
	}

	if err := s.productRepo.CreateProduct(newProduct); err != nil {
		return nil, errors.NewInternalError("failed to create product", err)
	}

	return ToProductResponse(newProduct), nil
}

func (s *Service) GetProductsByCompanyID(userID, companyID uint, page, limit int) (*PaginatedResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	products, total, err := s.productRepo.FindProductsByCompanyID(companyID, page, limit)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch products", err)
	}

	// Convert to response DTOs
	productResponses := make([]*ProductResponse, len(products))
	for i, p := range products {
		productResponses[i] = ToProductResponse(p)
	}

	return NewPaginatedResponse(productResponses, total, page, limit), nil
}

func (s *Service) GetProductByID(userID, companyID, productID uint) (*ProductResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	product, err := s.productRepo.FindProductByIDAndCompany(productID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("product not found")
	}

	return ToProductResponse(product), nil
}

func (s *Service) UpdateProduct(userID, companyID, productID uint, req *UpdateProductRequest) (*ProductResponse, error) {
	// Check user authorization (only admin/owner can update)
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleManager); err != nil {
		return nil, err
	}

	// Get existing product
	existingProduct, err := s.productRepo.FindProductByIDAndCompany(productID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("product not found")
	}

	// Check SKU uniqueness if SKU is being changed
	if req.SKU != "" && req.SKU != existingProduct.SKU {
		skuProduct, _ := s.productRepo.FindProductBySKUAndCompany(req.SKU, companyID)
		if skuProduct != nil && skuProduct.ID != productID {
			return nil, errors.NewConflictError("product with this SKU already exists in the company")
		}
		existingProduct.SKU = req.SKU
	}

	// Update fields
	if req.Name != "" {
		existingProduct.Name = req.Name
	}
	if req.Description != "" {
		existingProduct.Description = req.Description
	}
	if req.BasePrice >= 0 {
		existingProduct.BasePrice = req.BasePrice
	}
	if req.IsActive != nil {
		existingProduct.IsActive = *req.IsActive
	}

	if err := s.productRepo.UpdateProduct(existingProduct); err != nil {
		return nil, errors.NewInternalError("failed to update product", err)
	}

	return ToProductResponse(existingProduct), nil
}

func (s *Service) DeleteProduct(userID, companyID, productID uint) error {
	// Check user authorization (only admin/owner can delete)
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleManager); err != nil {
		return err
	}

	// Check if product exists
	_, err := s.productRepo.FindProductByIDAndCompany(productID, companyID)
	if err != nil {
		return errors.NewNotFoundError("product not found")
	}

	if err := s.productRepo.SoftDeleteProduct(productID); err != nil {
		return errors.NewInternalError("failed to delete product", err)
	}

	return nil
}

// Product variant operations
func (s *Service) CreateProductVariant(userID, companyID, productID uint, req *CreateProductVariantRequest) (*ProductVariantResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	// Verify product exists and belongs to company
	_, err := s.productRepo.FindProductByIDAndCompany(productID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("product not found")
	}

	// Check if variant SKU already exists for this product
	existingVariant, _ := s.productRepo.FindProductVariantBySKUAndProduct(req.SKU, productID)
	if existingVariant != nil {
		return nil, errors.NewConflictError("product variant with this SKU already exists")
	}

	// Create variant
	newVariant := req.ToProductVariant(productID)
	if !newVariant.IsValid() {
		return nil, errors.NewValidationError("invalid product variant data")
	}

	if err := s.productRepo.CreateProductVariant(newVariant); err != nil {
		return nil, errors.NewInternalError("failed to create product variant", err)
	}

	return ToProductVariantResponse(newVariant), nil
}

func (s *Service) GetProductVariantsByProductID(userID, companyID, productID uint) ([]*ProductVariantResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	// Verify product exists and belongs to company
	_, err := s.productRepo.FindProductByIDAndCompany(productID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("product not found")
	}

	variants, err := s.productRepo.FindProductVariantsByProductID(productID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch product variants", err)
	}

	// Convert to response DTOs
	variantResponses := make([]*ProductVariantResponse, len(variants))
	for i, v := range variants {
		variantResponses[i] = ToProductVariantResponse(v)
	}

	return variantResponses, nil
}

func (s *Service) GetProductVariantByID(userID, companyID, productID, variantID uint) (*ProductVariantResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	// Verify product exists and belongs to company
	_, err := s.productRepo.FindProductByIDAndCompany(productID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("product not found")
	}

	variant, err := s.productRepo.FindProductVariantByIDAndProduct(variantID, productID)
	if err != nil {
		return nil, errors.NewNotFoundError("product variant not found")
	}

	return ToProductVariantResponse(variant), nil
}

func (s *Service) UpdateProductVariant(userID, companyID, productID, variantID uint, req *UpdateProductVariantRequest) (*ProductVariantResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	// Verify product exists and belongs to company
	_, err := s.productRepo.FindProductByIDAndCompany(productID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("product not found")
	}

	// Get existing variant
	existingVariant, err := s.productRepo.FindProductVariantByIDAndProduct(variantID, productID)
	if err != nil {
		return nil, errors.NewNotFoundError("product variant not found")
	}

	// Check SKU uniqueness if SKU is being changed
	if req.SKU != "" && req.SKU != existingVariant.SKU {
		skuVariant, _ := s.productRepo.FindProductVariantBySKUAndProduct(req.SKU, productID)
		if skuVariant != nil && skuVariant.ID != variantID {
			return nil, errors.NewConflictError("product variant with this SKU already exists")
		}
		existingVariant.SKU = req.SKU
	}

	// Update fields
	if req.Name != "" {
		existingVariant.Name = req.Name
	}
	if req.Price >= 0 {
		existingVariant.Price = req.Price
	}
	if req.Stock >= 0 {
		existingVariant.Stock = req.Stock
	}
	if req.Attributes != nil {
		// Convert attributes to JSON
		attributesJSON, _ := json.Marshal(req.Attributes)
		existingVariant.Attributes = attributesJSON
	}
	if req.IsActive != nil {
		existingVariant.IsActive = *req.IsActive
	}

	if err := s.productRepo.UpdateProductVariant(existingVariant); err != nil {
		return nil, errors.NewInternalError("failed to update product variant", err)
	}

	return ToProductVariantResponse(existingVariant), nil
}

func (s *Service) DeleteProductVariant(userID, companyID, productID, variantID uint) error {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return err
	}

	// Verify product exists and belongs to company
	_, err := s.productRepo.FindProductByIDAndCompany(productID, companyID)
	if err != nil {
		return errors.NewNotFoundError("product not found")
	}

	// Check if variant exists
	_, err = s.productRepo.FindProductVariantByIDAndProduct(variantID, productID)
	if err != nil {
		return errors.NewNotFoundError("product variant not found")
	}

	if err := s.productRepo.SoftDeleteProductVariant(variantID); err != nil {
		return errors.NewInternalError("failed to delete product variant", err)
	}

	return nil
}

// Stock management operations
func (s *Service) UpdateVariantStock(userID, companyID, productID, variantID uint, req *UpdateStockRequest) (*ProductVariantResponse, error) {
	// Check user authorization (managers can update stock)
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleManager); err != nil {
		return nil, err
	}

	// Verify product exists and belongs to company
	_, err := s.productRepo.FindProductByIDAndCompany(productID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("product not found")
	}

	// Get existing variant
	variant, err := s.productRepo.FindProductVariantByIDAndProduct(variantID, productID)
	if err != nil {
		return nil, errors.NewNotFoundError("product variant not found")
	}

	if err := s.productRepo.UpdateVariantStock(variantID, req.Stock); err != nil {
		return nil, errors.NewInternalError("failed to update stock", err)
	}

	// Update local variant object
	variant.Stock = req.Stock
	return ToProductVariantResponse(variant), nil
}

func (s *Service) AdjustVariantStock(userID, companyID, productID, variantID uint, req *AdjustStockRequest) (*ProductVariantResponse, error) {
	// Check user authorization (managers can adjust stock)
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleManager); err != nil {
		return nil, err
	}

	// Verify product exists and belongs to company
	_, err := s.productRepo.FindProductByIDAndCompany(productID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("product not found")
	}

	// Get existing variant
	variant, err := s.productRepo.FindProductVariantByIDAndProduct(variantID, productID)
	if err != nil {
		return nil, errors.NewNotFoundError("product variant not found")
	}

	// Adjust stock
	if req.Amount > 0 {
		if err := s.productRepo.AddVariantStock(variantID, req.Amount); err != nil {
			return nil, errors.NewInternalError("failed to add stock", err)
		}
		variant.Stock += req.Amount
	} else if req.Amount < 0 {
		amount := -req.Amount
		if err := s.productRepo.RemoveVariantStock(variantID, amount); err != nil {
			return nil, errors.NewInternalError("insufficient stock or failed to remove stock", err)
		}
		variant.Stock -= amount
	}

	return ToProductVariantResponse(variant), nil
}

// Helper function to check user access to company
func (s *Service) checkUserCompanyAccess(userID, companyID uint, minRole user.Role) error {
	userRole, err := s.userRepo.FindUserRoleInCompany(userID, companyID)
	if err != nil {
		return errors.NewForbiddenError("you don't have access to this company")
	}

	// Check if user has sufficient role
	if !s.hasSufficientRole(userRole.Role, minRole) {
		return errors.NewForbiddenError("insufficient permissions for this operation")
	}

	return nil
}

// Helper function to check role hierarchy
func (s *Service) hasSufficientRole(userRole, requiredRole user.Role) bool {
	roleHierarchy := map[user.Role]int{
		user.RoleEmployee: 1,
		user.RoleManager:  2,
		user.RoleAdmin:    3,
		user.RoleOwner:    4,
	}

	return roleHierarchy[userRole] >= roleHierarchy[requiredRole]
}
