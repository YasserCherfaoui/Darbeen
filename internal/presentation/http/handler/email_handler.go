package handler

import (
	"net/http"
	"strconv"

	emailApp "github.com/YasserCherfaoui/darween/internal/application/email"
	"github.com/YasserCherfaoui/darween/internal/presentation/http/middleware"
	"github.com/YasserCherfaoui/darween/internal/presentation/response"
	"github.com/YasserCherfaoui/darween/pkg/errors"
	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	emailService *emailApp.Service
}

func NewEmailHandler(emailService *emailApp.Service) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
	}
}

type SendEmailRequest struct {
	To      []string `json:"to" binding:"required,min=1"`
	Subject string   `json:"subject" binding:"required"`
	Body    string   `json:"body" binding:"required"`
	IsHTML  bool     `json:"is_html"`
}

func (h *EmailHandler) SendEmail(c *gin.Context) {
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

	var req SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	// TODO: Check user has permission to send emails for this company
	_ = userID

	emailReq := &emailApp.SendCustomEmailRequest{
		CompanyID: uint(companyID),
		To:        req.To,
		Subject:   req.Subject,
		Body:      req.Body,
		IsHTML:    req.IsHTML,
	}

	if err := h.emailService.SendCustomEmail(emailReq); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Email queued successfully", nil)
}

