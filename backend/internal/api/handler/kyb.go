package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type KYBHandler struct {
	pool *pgxpool.Pool
}

func NewKYBHandler(pool *pgxpool.Pool) *KYBHandler {
	return &KYBHandler{pool: pool}
}

func (h *KYBHandler) Submit(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	if role != "venue" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only venues can submit KYB"})
		return
	}

	var req struct {
		BusinessName       string   `json:"business_name" binding:"required"`
		TaxID              string   `json:"tax_id" binding:"required"`
		RegistrationNumber string   `json:"registration_number"`
		Address            string   `json:"address" binding:"required"`
		Phone              string   `json:"phone"`
		DocumentURLs       []string `json:"document_urls"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check existing
	var exists int
	h.pool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM business_verifications WHERE user_id = $1`, userID).Scan(&exists)
	if exists > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "KYB already submitted"})
		return
	}

	docs, _ := json.Marshal(req.DocumentURLs)

	var result struct {
		ID                 string  `json:"id"`
		UserID             string  `json:"user_id"`
		BusinessName       string  `json:"business_name"`
		TaxID              string  `json:"tax_id"`
		RegistrationNumber string  `json:"registration_number"`
		Address            string  `json:"address"`
		Phone              string  `json:"phone"`
		Status             string  `json:"status"`
	}
	err := h.pool.QueryRow(context.Background(),
		`INSERT INTO business_verifications (id, user_id, business_name, tax_id, registration_number, address, phone, documents, status, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, 'pending', NOW(), NOW())
		 RETURNING id, user_id, business_name, tax_id, COALESCE(registration_number,''), address, COALESCE(phone,''), status`,
		userID, req.BusinessName, req.TaxID, req.RegistrationNumber, req.Address, req.Phone, docs).
		Scan(&result.ID, &result.UserID, &result.BusinessName, &result.TaxID, &result.RegistrationNumber, &result.Address, &result.Phone, &result.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit KYB"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *KYBHandler) GetStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var id, businessName, taxID, regNum, address, phone, status, reason string
	var docs []byte
	err := h.pool.QueryRow(context.Background(),
		`SELECT id, COALESCE(business_name,''), COALESCE(tax_id,''), COALESCE(registration_number,''), COALESCE(address,''), COALESCE(phone,''), status, COALESCE(rejection_reason,''), documents
		 FROM business_verifications WHERE user_id = $1`, userID).
		Scan(&id, &businessName, &taxID, &regNum, &address, &phone, &status, &reason, &docs)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no KYB submission found"})
		return
	}

	var docURLs []string
	json.Unmarshal(docs, &docURLs)

	c.JSON(http.StatusOK, gin.H{
		"id": id, "business_name": businessName, "tax_id": taxID,
		"registration_number": regNum, "address": address, "phone": phone,
		"status": status, "rejection_reason": reason, "documents": docURLs,
	})
}

// admin review
func (h *KYBHandler) Review(c *gin.Context) {
	verificationID := c.Param("id")

	var req struct {
		Status          string `json:"status" binding:"required,oneof=verified rejected"`
		RejectionReason string `json:"rejection_reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Status == "rejected" && req.RejectionReason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rejection reason required"})
		return
	}

	var id, userID, bizStatus string
	err := h.pool.QueryRow(context.Background(),
		`UPDATE business_verifications SET status=$2, rejection_reason=$3,
		 verified_at = CASE WHEN $2='verified' THEN NOW() ELSE verified_at END,
		 updated_at = NOW()
		 WHERE id=$1
		 RETURNING id, user_id, status`,
		verificationID, req.Status, req.RejectionReason).
		Scan(&id, &userID, &bizStatus)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "verification not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": bizStatus})
}

func (h *KYBHandler) ListPending(c *gin.Context) {
	rows, err := h.pool.Query(context.Background(),
		`SELECT bv.id, bv.user_id, bv.business_name, bv.status, u.name, u.email
		 FROM business_verifications bv
		 JOIN users u ON bv.user_id = u.id
		 WHERE bv.status = 'pending'
		 ORDER BY bv.created_at`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list"})
		return
	}
	defer rows.Close()

	var list []gin.H
	for rows.Next() {
		var id, userID, bizName, status, userName, email string
		rows.Scan(&id, &userID, &bizName, &status, &userName, &email)
		list = append(list, gin.H{
			"id": id, "user_id": userID, "business_name": bizName,
			"status": status, "user_name": userName, "email": email,
		})
	}
	if list == nil {
		list = []gin.H{}
	}
	c.JSON(http.StatusOK, list)
}
