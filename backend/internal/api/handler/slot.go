package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dennis-lee/LiveHouseAAS/backend/internal/domain/slot"
)

type SlotHandler struct {
	pool *pgxpool.Pool
}

func NewSlotHandler(pool *pgxpool.Pool) *SlotHandler {
	return &SlotHandler{pool: pool}
}

func (h *SlotHandler) Create(c *gin.Context) {
	venueID := c.Param("venueId")

	var req slot.CreateSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if conflict, err := h.checkOverlap(venueID, req.Date, req.StartTime, req.EndTime, ""); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check slot conflicts"})
		return
	} else if conflict {
		c.JSON(http.StatusConflict, gin.H{"error": "slot overlaps with an existing slot"})
		return
	}

	var s slot.Slot
	err := h.pool.QueryRow(context.Background(),
		`INSERT INTO slots (id, venue_id, date, start_time, end_time, status, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, 'available', NOW(), NOW())
		 RETURNING id, venue_id, date, start_time, end_time, status, created_at, updated_at`,
		venueID, req.Date, req.StartTime, req.EndTime).
		Scan(&s.ID, &s.VenueID, &s.Date, &s.StartTime, &s.EndTime, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create slot"})
		return
	}

	c.JSON(http.StatusCreated, s)
}

func (h *SlotHandler) BatchCreate(c *gin.Context) {
	venueID := c.Param("venueId")

	var req slot.BatchCreateSlotsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var created []slot.Slot
	for _, sr := range req.Slots {
		conflict, err := h.checkOverlap(venueID, sr.Date, sr.StartTime, sr.EndTime, "")
		if err != nil {
			continue
		}
		if conflict {
			continue
		}

		var s slot.Slot
		err = h.pool.QueryRow(context.Background(),
			`INSERT INTO slots (id, venue_id, date, start_time, end_time, status, created_at, updated_at)
			 VALUES (gen_random_uuid(), $1, $2, $3, $4, 'available', NOW(), NOW())
			 RETURNING id, venue_id, date, start_time, end_time, status, created_at, updated_at`,
			venueID, sr.Date, sr.StartTime, sr.EndTime).
			Scan(&s.ID, &s.VenueID, &s.Date, &s.StartTime, &s.EndTime, &s.Status, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			continue
		}
		created = append(created, s)
	}

	if created == nil {
		created = []slot.Slot{}
	}

	c.JSON(http.StatusCreated, gin.H{
		"created": len(created),
		"slots":   created,
	})
}

func (h *SlotHandler) List(c *gin.Context) {
	venueID := c.Param("venueId")
	statusFilter := c.Query("status")

	query := `SELECT id, venue_id, date, start_time, end_time, status, created_at, updated_at
			  FROM slots WHERE venue_id = $1`
	args := []interface{}{venueID}

	if statusFilter != "" {
		query += fmt.Sprintf(" AND status = $2")
		args = append(args, statusFilter)
	}
	query += " ORDER BY date, start_time"

	rows, err := h.pool.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list slots"})
		return
	}
	defer rows.Close()

	var slots []slot.Slot
	for rows.Next() {
		var s slot.Slot
		if err := rows.Scan(&s.ID, &s.VenueID, &s.Date, &s.StartTime, &s.EndTime, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			continue
		}
		slots = append(slots, s)
	}

	if slots == nil {
		slots = []slot.Slot{}
	}

	c.JSON(http.StatusOK, slots)
}

func (h *SlotHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	_, err := h.pool.Exec(context.Background(), `DELETE FROM slots WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete slot"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *SlotHandler) checkOverlap(venueID, date, startTime, endTime, excludeID string) (bool, error) {
	var count int64
	err := h.pool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM slots
		 WHERE venue_id = $1 AND date = $2
		 AND start_time < $4::time AND end_time > $3::time
		 AND ($5 = '' OR id != $5::uuid)`,
		venueID, date, startTime, endTime, excludeID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
