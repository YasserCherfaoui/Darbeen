package supplier

import (
	"strings"

	"github.com/YasserCherfaoui/darween/internal/domain/supplier"
	"github.com/YasserCherfaoui/darween/internal/domain/user"
	"github.com/YasserCherfaoui/darween/pkg/errors"
)

type Service struct {
	supplierRepo supplier.Repository
	userRepo     user.Repository
}

func NewService(supplierRepo supplier.Repository, userRepo user.Repository) *Service {
	return &Service{
		supplierRepo: supplierRepo,
		userRepo:     userRepo,
	}
}

// Supplier operations
func (s *Service) CreateSupplier(userID, companyID uint, req *CreateSupplierRequest) (*SupplierResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleManager); err != nil {
		return nil, err
	}

	// Validate email format if provided
	if req.Email != "" && !isValidEmail(req.Email) {
		return nil, errors.NewValidationError("invalid email format")
	}

	// Create supplier
	newSupplier := req.ToSupplier(companyID)
	if !newSupplier.IsValid() {
		return nil, errors.NewValidationError("invalid supplier data")
	}

	if err := s.supplierRepo.CreateSupplier(newSupplier); err != nil {
		return nil, errors.NewInternalError("failed to create supplier", err)
	}

	return ToSupplierResponse(newSupplier, 0), nil
}

func (s *Service) GetSuppliersByCompanyID(userID, companyID uint, page, limit int) (*PaginatedResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	suppliers, total, err := s.supplierRepo.FindSuppliersByCompanyID(companyID, page, limit)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch suppliers", err)
	}

	// Convert to response DTOs
	supplierResponses := make([]*SupplierResponse, len(suppliers))
	for i, sup := range suppliers {
		// Get product count for each supplier
		products, _ := s.supplierRepo.FindProductsBySupplier(sup.ID)
		supplierResponses[i] = ToSupplierResponse(sup, len(products))
	}

	return NewPaginatedResponse(supplierResponses, total, page, limit), nil
}

func (s *Service) GetSupplierByID(userID, companyID, supplierID uint) (*SupplierResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	sup, err := s.supplierRepo.FindSupplierByIDAndCompany(supplierID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("supplier not found")
	}

	// Get product count
	products, _ := s.supplierRepo.FindProductsBySupplier(sup.ID)
	return ToSupplierResponse(sup, len(products)), nil
}

func (s *Service) UpdateSupplier(userID, companyID, supplierID uint, req *UpdateSupplierRequest) (*SupplierResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleManager); err != nil {
		return nil, err
	}

	// Get existing supplier
	existingSupplier, err := s.supplierRepo.FindSupplierByIDAndCompany(supplierID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("supplier not found")
	}

	// Update fields
	if req.Name != "" {
		existingSupplier.Name = req.Name
	}
	if req.ContactPerson != "" {
		existingSupplier.ContactPerson = req.ContactPerson
	}
	if req.Email != "" {
		if !isValidEmail(req.Email) {
			return nil, errors.NewValidationError("invalid email format")
		}
		existingSupplier.Email = req.Email
	}
	if req.Phone != "" {
		existingSupplier.Phone = req.Phone
	}
	if req.Address != "" {
		existingSupplier.Address = req.Address
	}
	if req.IsActive != nil {
		existingSupplier.IsActive = *req.IsActive
	}

	if err := s.supplierRepo.UpdateSupplier(existingSupplier); err != nil {
		return nil, errors.NewInternalError("failed to update supplier", err)
	}

	// Get product count
	products, _ := s.supplierRepo.FindProductsBySupplier(existingSupplier.ID)
	return ToSupplierResponse(existingSupplier, len(products)), nil
}

func (s *Service) DeleteSupplier(userID, companyID, supplierID uint) error {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleManager); err != nil {
		return err
	}

	// Check if supplier exists
	_, err := s.supplierRepo.FindSupplierByIDAndCompany(supplierID, companyID)
	if err != nil {
		return errors.NewNotFoundError("supplier not found")
	}

	if err := s.supplierRepo.SoftDeleteSupplier(supplierID); err != nil {
		return errors.NewInternalError("failed to delete supplier", err)
	}

	return nil
}

func (s *Service) GetSupplierProducts(userID, companyID, supplierID uint) (*SupplierWithProductsResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	// Get supplier
	sup, err := s.supplierRepo.FindSupplierByIDAndCompany(supplierID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("supplier not found")
	}

	// Get products for this supplier
	products, err := s.supplierRepo.FindProductsBySupplier(supplierID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch supplier products", err)
	}

	// Convert to response
	productInfos := make([]SupplierProductInfo, len(products))
	for i, p := range products {
		productInfos[i] = SupplierProductInfo{
			ID:           p.ID,
			Name:         p.Name,
			SKU:          p.SKU,
			SupplierCost: p.SupplierCost,
		}
	}

	return &SupplierWithProductsResponse{
		ID:            sup.ID,
		CompanyID:     sup.CompanyID,
		Name:          sup.Name,
		ContactPerson: sup.ContactPerson,
		Email:         sup.Email,
		Phone:         sup.Phone,
		Address:       sup.Address,
		IsActive:      sup.IsActive,
		Products:      productInfos,
	}, nil
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

// Helper function to validate email format
func isValidEmail(email string) bool {
	// Basic email validation
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

