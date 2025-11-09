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
	if req.BaseRetailPrice != nil && *req.BaseRetailPrice >= 0 {
		existingProduct.BaseRetailPrice = *req.BaseRetailPrice
	}
	if req.BaseWholesalePrice != nil && *req.BaseWholesalePrice >= 0 {
		existingProduct.BaseWholesalePrice = *req.BaseWholesalePrice
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
	if req.RetailPrice != nil && *req.RetailPrice >= 0 {
		existingVariant.RetailPrice = req.RetailPrice
		existingVariant.UseParentPricing = false
	}
	if req.WholesalePrice != nil && *req.WholesalePrice >= 0 {
		existingVariant.WholesalePrice = req.WholesalePrice
		existingVariant.UseParentPricing = false
	}
	if req.UseParentPricing != nil && *req.UseParentPricing {
		existingVariant.UseParentPricing = true
		existingVariant.RetailPrice = nil
		existingVariant.WholesalePrice = nil
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

// BulkCreateProductVariants creates multiple product variants from attribute combinations
func (s *Service) BulkCreateProductVariants(userID, companyID, productID uint, req *BulkCreateProductVariantsRequest) (*BulkCreateProductVariantsResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	// Verify product exists and belongs to company
	parentProduct, err := s.productRepo.FindProductByIDAndCompany(productID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("product not found")
	}

	// Validate that we have at least one attribute
	if len(req.Attributes) == 0 {
		return nil, errors.NewValidationError("at least one attribute is required")
	}

	// Generate all combinations
	combinations := s.generateAttributeCombinations(req.Attributes)

	// Create variants for each combination
	variants := make([]*product.ProductVariant, 0, len(combinations))
	variantResponses := make([]ProductVariantResponse, 0, len(combinations))

	for _, combo := range combinations {
		// Generate variant name: "Red - 39"
		variantName := s.generateVariantName(combo)

		// Generate SKU: "PARENT-SKU>RED>39"
		variantSKU := s.generateVariantSKU(parentProduct.SKU, combo)

		// Check if variant SKU already exists
		existingVariant, _ := s.productRepo.FindProductVariantBySKUAndProduct(variantSKU, productID)
		if existingVariant != nil {
			return nil, errors.NewConflictError("variant with SKU " + variantSKU + " already exists")
		}

		// Convert combination to attributes JSON
		attributesMap := make(map[string]interface{})
		for _, attr := range combo {
			attributesMap[attr.Name] = attr.Value
		}
		attributesJSON, _ := json.Marshal(attributesMap)

		// Create variant entity
		newVariant := &product.ProductVariant{
			ProductID:        productID,
			Name:             variantName,
			SKU:              variantSKU,
			UseParentPricing: req.UseParentPricing,
			Attributes:       attributesJSON,
			IsActive:         true,
		}

		// Set pricing based on strategy
		if !req.UseParentPricing {
			// If not using parent pricing, set explicit prices
			newVariant.RetailPrice = &parentProduct.BaseRetailPrice
			newVariant.WholesalePrice = &parentProduct.BaseWholesalePrice
		}

		if !newVariant.IsValid() {
			return nil, errors.NewValidationError("invalid variant data for SKU: " + variantSKU)
		}

		variants = append(variants, newVariant)
	}

	// Batch create all variants
	for _, variant := range variants {
		if err := s.productRepo.CreateProductVariant(variant); err != nil {
			return nil, errors.NewInternalError("failed to create product variant", err)
		}
		variantResponses = append(variantResponses, *ToProductVariantResponse(variant))
	}

	return &BulkCreateProductVariantsResponse{
		CreatedCount: len(variantResponses),
		Variants:     variantResponses,
	}, nil
}

// generateAttributeCombinations generates all possible combinations (cartesian product)
func (s *Service) generateAttributeCombinations(attributes []AttributeDefinition) [][]AttributeValue {
	if len(attributes) == 0 {
		return [][]AttributeValue{}
	}

	// Start with the first attribute
	var combinations [][]AttributeValue
	for _, value := range attributes[0].Values {
		combinations = append(combinations, []AttributeValue{{Name: attributes[0].Name, Value: value}})
	}

	// For each subsequent attribute, multiply the combinations
	for i := 1; i < len(attributes); i++ {
		var newCombinations [][]AttributeValue
		for _, combo := range combinations {
			for _, value := range attributes[i].Values {
				newCombo := make([]AttributeValue, len(combo))
				copy(newCombo, combo)
				newCombo = append(newCombo, AttributeValue{Name: attributes[i].Name, Value: value})
				newCombinations = append(newCombinations, newCombo)
			}
		}
		combinations = newCombinations
	}

	return combinations
}

// generateVariantName creates a human-readable name from attributes
func (s *Service) generateVariantName(combo []AttributeValue) string {
	if len(combo) == 0 {
		return ""
	}

	names := make([]string, len(combo))
	for i, attr := range combo {
		names[i] = attr.Value
	}
	return join(names, " - ")
}

// generateVariantSKU creates SKU with pattern: PARENT>VALUE1>VALUE2
func (s *Service) generateVariantSKU(parentSKU string, combo []AttributeValue) string {
	sku := parentSKU
	for _, attr := range combo {
		// Normalize value: uppercase, remove spaces
		normalizedValue := normalizeForSKU(attr.Value)
		sku += ">" + normalizedValue
	}
	return sku
}

// Helper type for attribute combinations
type AttributeValue struct {
	Name  string
	Value string
}

// Helper function to join strings
func join(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// Helper function to normalize values for SKU
func normalizeForSKU(value string) string {
	// Convert to uppercase and remove spaces
	normalized := ""
	for _, char := range value {
		if char != ' ' {
			if char >= 'a' && char <= 'z' {
				normalized += string(char - 32) // Convert to uppercase
			} else {
				normalized += string(char)
			}
		}
	}
	return normalized
}
