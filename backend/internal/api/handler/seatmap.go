package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeatMapHandler struct {
	pool *pgxpool.Pool
}

func NewSeatMapHandler(pool *pgxpool.Pool) *SeatMapHandler {
	return &SeatMapHandler{pool: pool}
}

func (h *SeatMapHandler) GetLayout(c *gin.Context) {
	venueID := c.Param("venueId")

	var id, name string
	var rows, cols int32
	var seatsJSON []byte
	err := h.pool.QueryRow(context.Background(),
		`SELECT id, name, rows, cols, seats FROM seat_layouts WHERE venue_id = $1`,
		venueID).Scan(&id, &name, &rows, &cols, &seatsJSON)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no seat layout found"})
		return
	}

	var seats interface{}
	json.Unmarshal(seatsJSON, &seats)

	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"name":    name,
		"rows":    rows,
		"cols":    cols,
		"seats":   seats,
	})
}

func (h *SeatMapHandler) SaveLayout(c *gin.Context) {
	venueID := c.Param("venueId")

	var req struct {
		Name  string          `json:"name"`
		Rows  int32           `json:"rows" binding:"required,min=1"`
		Cols  int32           `json:"cols" binding:"required,min=1"`
		Seats json.RawMessage `json:"seats"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	seatsJSON, _ := json.Marshal(req.Seats)
	if req.Name == "" {
		req.Name = "Main"
	}

	_, err := h.pool.Exec(context.Background(),
		`INSERT INTO seat_layouts (id, venue_id, name, rows, cols, seats, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, NOW(), NOW())
		 ON CONFLICT (venue_id) DO UPDATE SET name=$2, rows=$3, cols=$4, seats=$5, updated_at=NOW()`,
		venueID, req.Name, req.Rows, req.Cols, seatsJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save layout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "layout saved"})
}

func (h *SeatMapHandler) GetSeatAvailability(c *gin.Context) {
	eventID := c.Param("eventId")

	// get ticket types with seat sections
	rows, err := h.pool.Query(context.Background(),
		`SELECT tt.id, tt.name, tt.price, tt.quantity, tt.seat_section,
		        tt.sale_start_at, tt.sale_end_at, tt.status,
		        COALESCE(
		            (SELECT SUM(oi.quantity) FROM order_items oi
		             JOIN orders o ON oi.order_id = o.id
		             WHERE oi.ticket_type_id = tt.id AND o.status IN ('pending','paid'))
		        , 0) as sold
		 FROM ticket_types tt
		 WHERE tt.event_id = $1 AND tt.status = 'active'
		 ORDER BY tt.price`, eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load seat availability"})
		return
	}
	defer rows.Close()

	var sections []gin.H
	for rows.Next() {
		var id, name, seatSection, ttStatus string
		var price, quantity, sold int32
		var saleStart, saleEnd interface{}
		if err := rows.Scan(&id, &name, &price, &quantity, &seatSection, &saleStart, &saleEnd, &ttStatus, &sold); err != nil {
			continue
		}
		sections = append(sections, gin.H{
			"id": id, "name": name, "price": price, "quantity": quantity,
			"seat_section": seatSection, "sold": sold,
			"remaining": quantity - sold,
			"sale_start_at": saleStart, "sale_end_at": saleEnd, "status": ttStatus,
		})
	}
	if sections == nil {
		sections = []gin.H{}
	}
	c.JSON(http.StatusOK, sections)
}
