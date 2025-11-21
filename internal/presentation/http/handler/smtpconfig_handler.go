package handler

import (
	"net/http"
	"strconv"

	smtpconfigApp "github.com/YasserCherfaoui/darween/internal/application/smtpconfig"
	"github.com/YasserCherfaoui/darween/internal/presentation/http/middleware"
	"github.com/YasserCherfaoui/darween/internal/presentation/response"
	"github.com/YasserCherfaoui/darween/pkg/errors"
	"github.com/gin-gonic/gin"
)

type SMTPConfigHandler struct {
	smtpConfigService *smtpconfigApp.Service
}

func NewSMTPConfigHandler(smtpConfigService *smtpconfigApp.Service) *SMTPConfigHandler {
	return &SMTPConfigHandler{
		smtpConfigService: smtpConfigService,
	}
}

func (h *SMTPConfigHandler) CreateSMTPConfig(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	companyID, err := strconv.ParseUint(c.Param("companyId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid company id"))
		return
	}

	var req smtpconfigApp.CreateSMTPConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.smtpConfigService.CreateSMTPConfig(userID, uint(companyID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusCreated, "SMTP config created successfully", result)
}

func (h *SMTPConfigHandler) ListSMTPConfigs(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	companyID, err := strconv.ParseUint(c.Param("companyId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid company id"))
		return
	}

	result, err := h.smtpConfigService.GetSMTPConfigs(userID, uint(companyID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *SMTPConfigHandler) GetSMTPConfig(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	companyID, err := strconv.ParseUint(c.Param("companyId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid company id"))
		return
	}

	configID, err := strconv.ParseUint(c.Param("configId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid config id"))
		return
	}

	result, err := h.smtpConfigService.GetSMTPConfig(userID, uint(companyID), uint(configID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *SMTPConfigHandler) UpdateSMTPConfig(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	companyID, err := strconv.ParseUint(c.Param("companyId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid company id"))
		return
	}

	configID, err := strconv.ParseUint(c.Param("configId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid config id"))
		return
	}

	var req smtpconfigApp.UpdateSMTPConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.smtpConfigService.UpdateSMTPConfig(userID, uint(companyID), uint(configID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "SMTP config updated successfully", result)
}

func (h *SMTPConfigHandler) DeleteSMTPConfig(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	companyID, err := strconv.ParseUint(c.Param("companyId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid company id"))
		return
	}

	configID, err := strconv.ParseUint(c.Param("configId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid config id"))
		return
	}

	err = h.smtpConfigService.DeleteSMTPConfig(userID, uint(companyID), uint(configID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "SMTP config deleted successfully", nil)
}

func (h *SMTPConfigHandler) SetDefaultSMTPConfig(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	companyID, err := strconv.ParseUint(c.Param("companyId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid company id"))
		return
	}

	configID, err := strconv.ParseUint(c.Param("configId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid config id"))
		return
	}

	err = h.smtpConfigService.SetDefaultSMTPConfig(userID, uint(companyID), uint(configID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Default SMTP config set successfully", nil)
}

