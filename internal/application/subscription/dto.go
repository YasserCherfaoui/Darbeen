package subscription

import "time"

type SubscriptionResponse struct {
	ID        uint       `json:"id"`
	CompanyID uint       `json:"company_id"`
	PlanType  string     `json:"plan_type"`
	Status    string     `json:"status"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	MaxUsers  int        `json:"max_users"`
}

type UpdateSubscriptionRequest struct {
	PlanType string `json:"plan_type" binding:"required"`
	Status   string `json:"status"`
}


