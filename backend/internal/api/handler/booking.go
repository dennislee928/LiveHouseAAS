package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dennis-lee/LiveHouseAAS/backend/internal/domain/booking"
)

type BookingHandler struct {
	pool *pgxpool.Pool
}

func NewBookingHandler(pool *pgxpool.Pool) *BookingHandler {
	return &BookingHandler{pool: pool}
}

func (h *BookingHandler) Create(c *gin.Context) {
	var req booking.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	artistID, _ := c.Get("user_id")

	// get venue_id and verify slot exists & is available
	var venueID, slotDate, startTime, endTime string
	var slotStatus string
	err := h.pool.QueryRow(context.Background(),
		`SELECT venue_id, date::text, start_time::text, end_time::text, status FROM slots WHERE id = $1`,
		req.SlotID).Scan(&venueID, &slotDate, &startTime, &endTime, &slotStatus)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "slot not found"})
		return
	}

	if slotStatus != "available" {
		c.JSON(http.StatusConflict, gin.H{"error": "slot is not available"})
		return
	}

	var b booking.BookingRequest
	err = h.pool.QueryRow(context.Background(),
		`INSERT INTO booking_requests (id, slot_id, venue_id, artist_id, message, status, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, 'pending', NOW(), NOW())
		 RETURNING id, slot_id, venue_id, artist_id, message, status, created_at, updated_at`,
		req.SlotID, venueID, artistID, req.Message).
		Scan(&b.ID, &b.SlotID, &b.VenueID, &b.ArtistID, &b.Message, &b.Status, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create booking request"})
		return
	}

	// mark slot as booked
	h.pool.Exec(context.Background(),
		`UPDATE slots SET status = 'booked', updated_at = NOW() WHERE id = $1`, req.SlotID)

	// notify venue owner
	ctx := context.Background()
	var venueOwnerID string
	h.pool.QueryRow(ctx,
		`SELECT owner_id FROM venues WHERE id = $1`, venueID).Scan(&venueOwnerID)
	if venueOwnerID != "" {
		CreateNotification(ctx, h.pool, venueOwnerID, "booking_created",
			"新的演出申請", "您收到了一筆新的演出申請", gin.H{
				"booking_id": b.ID,
				"venue_id":   venueID,
				"date":       slotDate,
				"start_time": startTime,
				"end_time":   endTime,
			})
		SendToUser(venueOwnerID, gin.H{
			"type":       "booking_created",
			"title":      "新的演出申請",
			"body":       "您收到了一筆新的演出申請",
			"booking_id": b.ID,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"booking": b,
		"slot": gin.H{
			"date":       slotDate,
			"start_time": startTime,
			"end_time":   endTime,
		},
	})
}

func (h *BookingHandler) ListByVenue(c *gin.Context) {
	venueID := c.Param("venueId")

	rows, err := h.pool.Query(context.Background(),
		`SELECT br.id, br.slot_id, br.venue_id, br.artist_id, br.message, br.status, br.created_at, br.updated_at,
		        s.date::text, s.start_time::text, s.end_time::text,
		        u.name, u.email
		 FROM booking_requests br
		 JOIN slots s ON br.slot_id = s.id
		 JOIN users u ON br.artist_id = u.id
		 WHERE br.venue_id = $1
		 ORDER BY br.created_at DESC`, venueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list bookings"})
		return
	}
	defer rows.Close()

	var bookings []booking.BookingDetail
	for rows.Next() {
		var b booking.BookingDetail
		if err := rows.Scan(&b.ID, &b.SlotID, &b.VenueID, &b.ArtistID, &b.Message, &b.Status, &b.CreatedAt, &b.UpdatedAt,
			&b.Date, &b.StartTime, &b.EndTime, &b.ArtistName, &b.ArtistEmail); err != nil {
			continue
		}
		bookings = append(bookings, b)
	}

	if bookings == nil {
		bookings = []booking.BookingDetail{}
	}

	c.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) ListByArtist(c *gin.Context) {
	artistID, _ := c.Get("user_id")

	rows, err := h.pool.Query(context.Background(),
		`SELECT br.id, br.slot_id, br.venue_id, br.artist_id, br.message, br.status, br.created_at, br.updated_at,
		        s.date::text, s.start_time::text, s.end_time::text,
		        v.name, v.city
		 FROM booking_requests br
		 JOIN slots s ON br.slot_id = s.id
		 JOIN venues v ON br.venue_id = v.id
		 WHERE br.artist_id = $1
		 ORDER BY br.created_at DESC`, artistID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list bookings"})
		return
	}
	defer rows.Close()

	var bookings []booking.BookingDetail
	for rows.Next() {
		var b booking.BookingDetail
		if err := rows.Scan(&b.ID, &b.SlotID, &b.VenueID, &b.ArtistID, &b.Message, &b.Status, &b.CreatedAt, &b.UpdatedAt,
			&b.Date, &b.StartTime, &b.EndTime, &b.VenueName, &b.VenueCity); err != nil {
			continue
		}
		bookings = append(bookings, b)
	}

	if bookings == nil {
		bookings = []booking.BookingDetail{}
	}

	c.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status booking.Status `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// start a transaction to update booking and slot status atomically
	tx, err := h.pool.Begin(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	defer tx.Rollback(context.Background())

	var b booking.BookingRequest
	err = tx.QueryRow(context.Background(),
		`UPDATE booking_requests SET status = $2, updated_at = NOW() WHERE id = $1
		 RETURNING id, slot_id, venue_id, artist_id, message, status, created_at, updated_at`,
		id, req.Status).
		Scan(&b.ID, &b.SlotID, &b.VenueID, &b.ArtistID, &b.Message, &b.Status, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	switch req.Status {
	case booking.StatusApproved, booking.StatusConfirmed:
		_, _ = tx.Exec(context.Background(),
			`UPDATE slots SET status = 'booked', updated_at = NOW() WHERE id = $1`, b.SlotID)

		// create event
		var slotDate, startTime, endTime string
		_ = tx.QueryRow(context.Background(),
			`SELECT date::text, start_time::text, end_time::text FROM slots WHERE id = $1`, b.SlotID).
			Scan(&slotDate, &startTime, &endTime)

		_, _ = tx.Exec(context.Background(),
			`INSERT INTO events (id, title, description, venue_id, artist_id, booking_id, start_at, end_at, status, created_at, updated_at)
			 VALUES (gen_random_uuid(), '', '', $1, $2, $3,
			         ($4::date + $5::time)::timestamptz, ($4::date + $6::time)::timestamptz, 'draft', NOW(), NOW())`,
			b.VenueID, b.ArtistID, b.ID, slotDate, startTime, endTime)

	case booking.StatusRejected, booking.StatusCancelled:
		// release the slot back to available
		_, _ = tx.Exec(context.Background(),
			`UPDATE slots SET status = 'available', updated_at = NOW() WHERE id = $1`, b.SlotID)
	}

	if err := tx.Commit(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update booking"})
		return
	}

	c.JSON(http.StatusOK, b)
}

func (h *BookingHandler) ListByOwner(c *gin.Context) {
	ownerID, _ := c.Get("user_id")

	rows, err := h.pool.Query(context.Background(),
		`SELECT br.id, br.slot_id, br.venue_id, br.artist_id, br.message, br.status, br.created_at, br.updated_at,
		        s.date::text, s.start_time::text, s.end_time::text,
		        u.name, u.email,
		        v.name, v.city
		 FROM booking_requests br
		 JOIN slots s ON br.slot_id = s.id
		 JOIN users u ON br.artist_id = u.id
		 JOIN venues v ON br.venue_id = v.id
		 WHERE br.venue_id IN (SELECT id FROM venues WHERE owner_id = $1)
		 ORDER BY br.created_at DESC`, ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list bookings"})
		return
	}
	defer rows.Close()

	var bookings []booking.BookingDetail
	for rows.Next() {
		var b booking.BookingDetail
		if err := rows.Scan(&b.ID, &b.SlotID, &b.VenueID, &b.ArtistID, &b.Message, &b.Status, &b.CreatedAt, &b.UpdatedAt,
			&b.Date, &b.StartTime, &b.EndTime, &b.ArtistName, &b.ArtistEmail, &b.VenueName, &b.VenueCity); err != nil {
			continue
		}
		bookings = append(bookings, b)
	}

	if bookings == nil {
		bookings = []booking.BookingDetail{}
	}

	c.JSON(http.StatusOK, bookings)
}
