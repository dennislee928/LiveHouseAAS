package kyb

import "time"

type Status string

const (
	StatusPending  Status = "pending"
	StatusVerified Status = "verified"
	StatusRejected Status = "rejected"
)

type BusinessVerification struct {
	ID                   string    `json:"id"`
	UserID               string    `json:"user_id"`
	BusinessName         string    `json:"business_name"`
	TaxID                string    `json:"tax_id"`
	RegistrationNumber   string    `json:"registration_number"`
	Address              string    `json:"address"`
	Phone                string    `json:"phone"`
	Documents            []string  `json:"documents"`
	Status               Status    `json:"status"`
	RejectionReason      string    `json:"rejection_reason,omitempty"`
	VerifiedAt           *time.Time `json:"verified_at,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type SubmitRequest struct {
	BusinessName       string   `json:"business_name" binding:"required"`
	TaxID              string   `json:"tax_id" binding:"required"`
	RegistrationNumber string   `json:"registration_number"`
	Address            string   `json:"address" binding:"required"`
	Phone              string   `json:"phone"`
	DocumentURLs       []string `json:"document_urls"`
}

type ReviewRequest struct {
	Status          Status `json:"status" binding:"required,oneof=verified rejected"`
	RejectionReason string `json:"rejection_reason"`
}
