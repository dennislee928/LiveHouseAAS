package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VerifyHandler struct {
	pool *pgxpool.Pool
}

func NewVerifyHandler(pool *pgxpool.Pool) *VerifyHandler {
	return &VerifyHandler{pool: pool}
}

type ScanRequest struct {
	Code      string `json:"code" binding:"required"`
	QRSecret  string `json:"qr_secret" binding:"required"`
}

func (h *VerifyHandler) Verify(c *gin.Context) {
	var req ScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, _ := c.Get("role")
	if role != "venue" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only venues can verify tickets"})
		return
	}

	var ticketID, ticketStatus, eventID string
	err := h.pool.QueryRow(context.Background(),
		`SELECT id, status, event_id FROM tickets WHERE code = $1 AND qr_secret = $2`,
		req.Code, req.QRSecret).Scan(&ticketID, &ticketStatus, &eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"valid": false, "error": "ticket not found"})
		return
	}

	if ticketStatus != "active" {
		c.JSON(http.StatusConflict, gin.H{
			"valid":  false,
			"error":  "ticket already used",
			"status": ticketStatus,
		})
		return
	}

	// mark as used
	var usedAt interface{}
	err = h.pool.QueryRow(context.Background(),
		`UPDATE tickets SET status = 'used', used_at = NOW() WHERE id = $1 AND status = 'active'
		 RETURNING used_at`, ticketID).Scan(&usedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify"})
		return
	}

	// get event + ticket info
	var eventTitle, venueName, ticketType, userName string
	h.pool.QueryRow(context.Background(),
		`SELECT e.title, v.name, tt.name, u.name
		 FROM tickets t
		 JOIN events e ON t.event_id = e.id
		 JOIN venues v ON e.venue_id = v.id
		 JOIN ticket_types tt ON t.ticket_type_id = tt.id
		 JOIN orders o ON t.order_id = o.id
		 JOIN users u ON o.user_id = u.id
		 WHERE t.id = $1`, ticketID).
		Scan(&eventTitle, &venueName, &ticketType, &userName)

	c.JSON(http.StatusOK, gin.H{
		"valid":       true,
		"ticket_id":   ticketID,
		"used_at":     usedAt,
		"event_title": eventTitle,
		"venue_name":  venueName,
		"ticket_type": ticketType,
		"holder_name": userName,
	})
}

func (h *VerifyHandler) Lookup(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code required"})
		return
	}

	var ticketID, ticketStatus, eventTitle, venueName, ticketType string
	err := h.pool.QueryRow(context.Background(),
		`SELECT t.id, t.status, COALESCE(e.title,''), COALESCE(v.name,''), COALESCE(tt.name,'')
		 FROM tickets t
		 JOIN events e ON t.event_id = e.id
		 JOIN venues v ON e.venue_id = v.id
		 JOIN ticket_types tt ON t.ticket_type_id = tt.id
		 WHERE t.code = $1`, code).
		Scan(&ticketID, &ticketStatus, &eventTitle, &venueName, &ticketType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"valid": false, "error": "ticket not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":       ticketStatus == "active",
		"ticket_id":   ticketID,
		"status":      ticketStatus,
		"event_title": eventTitle,
		"venue_name":  venueName,
		"ticket_type": ticketType,
	})
}
