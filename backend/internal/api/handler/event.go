package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dennis-lee/LiveHouseAAS/backend/internal/domain/event"
)

type EventHandler struct {
	pool *pgxpool.Pool
}

func NewEventHandler(pool *pgxpool.Pool) *EventHandler {
	return &EventHandler{pool: pool}
}

func (h *EventHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var e event.EventDetail
	err := h.pool.QueryRow(context.Background(),
		`SELECT e.id, e.title, e.description, e.venue_id, e.artist_id, COALESCE(e.booking_id, '00000000-0000-0000-0000-000000000000'), e.start_at, e.end_at, e.status, e.created_at, e.updated_at,
		        v.name, u.name
		 FROM events e
		 JOIN venues v ON e.venue_id = v.id
		 JOIN users u ON e.artist_id = u.id
		 WHERE e.id = $1`, id).
		Scan(&e.ID, &e.Title, &e.Description, &e.VenueID, &e.ArtistID, &e.BookingID, &e.StartAt, &e.EndAt, &e.Status, &e.CreatedAt, &e.UpdatedAt,
			&e.VenueName, &e.ArtistName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	c.JSON(http.StatusOK, e)
}

func (h *EventHandler) ListPublished(c *gin.Context) {
	rows, err := h.pool.Query(context.Background(),
		`SELECT e.id, e.title, e.description, e.venue_id, e.artist_id, COALESCE(e.booking_id, '00000000-0000-0000-0000-000000000000'), e.start_at, e.end_at, e.status, e.created_at, e.updated_at,
		        v.name, v.city, u.name
		 FROM events e
		 JOIN venues v ON e.venue_id = v.id
		 JOIN users u ON e.artist_id = u.id
		 WHERE e.status = 'published' AND e.start_at > NOW()
		 ORDER BY e.start_at`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list events"})
		return
	}
	defer rows.Close()

	var events []event.EventDetail
	for rows.Next() {
		var e event.EventDetail
		if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.VenueID, &e.ArtistID, &e.BookingID, &e.StartAt, &e.EndAt, &e.Status, &e.CreatedAt, &e.UpdatedAt,
			&e.VenueName, &e.VenueCity, &e.ArtistName); err != nil {
			continue
		}
		events = append(events, e)
	}
	if events == nil {
		events = []event.EventDetail{}
	}
	c.JSON(http.StatusOK, events)
}

func (h *EventHandler) ListByVenue(c *gin.Context) {
	venueID := c.Param("venueId")
	h.listEvents(c, `SELECT e.id, e.title, e.description, e.venue_id, e.artist_id, COALESCE(e.booking_id, '00000000-0000-0000-0000-000000000000'), e.start_at, e.end_at, e.status, e.created_at, e.updated_at,
		v.name, u.name FROM events e JOIN venues v ON e.venue_id = v.id JOIN users u ON e.artist_id = u.id WHERE e.venue_id = $1 ORDER BY e.start_at DESC`, venueID)
}

func (h *EventHandler) ListByArtist(c *gin.Context) {
	userID, _ := c.Get("user_id")
	h.listEvents(c, `SELECT e.id, e.title, e.description, e.venue_id, e.artist_id, COALESCE(e.booking_id, '00000000-0000-0000-0000-000000000000'), e.start_at, e.end_at, e.status, e.created_at, e.updated_at,
		v.name, u.name FROM events e JOIN venues v ON e.venue_id = v.id JOIN users u ON e.artist_id = u.id WHERE e.artist_id = $1 OR e.artist_id IN (SELECT id FROM users WHERE id = $1) ORDER BY e.start_at DESC`, userID)
}

func (h *EventHandler) listEvents(c *gin.Context, query string, arg interface{}) {
	rows, err := h.pool.Query(context.Background(), query, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list events"})
		return
	}
	defer rows.Close()

	var events []event.EventDetail
	for rows.Next() {
		var e event.EventDetail
		if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.VenueID, &e.ArtistID, &e.BookingID, &e.StartAt, &e.EndAt, &e.Status, &e.CreatedAt, &e.UpdatedAt,
			&e.VenueName, &e.ArtistName); err != nil {
			continue
		}
		e.VenueCity = ""
		events = append(events, e)
	}
	if events == nil {
		events = []event.EventDetail{}
	}
	c.JSON(http.StatusOK, events)
}

func (h *EventHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var e event.Event
	err := h.pool.QueryRow(context.Background(),
		`UPDATE events SET title = $2, description = $3, updated_at = NOW() WHERE id = $1
		 RETURNING id, title, description, venue_id, artist_id, COALESCE(booking_id, '00000000-0000-0000-0000-000000000000'), start_at, end_at, status, created_at, updated_at`,
		id, req.Title, req.Description).
		Scan(&e.ID, &e.Title, &e.Description, &e.VenueID, &e.ArtistID, &e.BookingID, &e.StartAt, &e.EndAt, &e.Status, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	c.JSON(http.StatusOK, e)
}

func (h *EventHandler) Publish(c *gin.Context) {
	id := c.Param("id")

	// check that the event has at least one ticket type
	var count int32
	h.pool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM ticket_types WHERE event_id = $1 AND status = 'active'`, id).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event must have at least one ticket type before publishing"})
		return
	}

	var e event.Event
	err := h.pool.QueryRow(context.Background(),
		`UPDATE events SET status = 'published', updated_at = NOW() WHERE id = $1 AND status = 'draft'
		 RETURNING id, title, description, venue_id, artist_id, COALESCE(booking_id, '00000000-0000-0000-0000-000000000000'), start_at, end_at, status, created_at, updated_at`,
		id).
		Scan(&e.ID, &e.Title, &e.Description, &e.VenueID, &e.ArtistID, &e.BookingID, &e.StartAt, &e.EndAt, &e.Status, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event not found or already published"})
		return
	}
	c.JSON(http.StatusOK, e)
}

func (h *EventHandler) CreateTicketType(c *gin.Context) {
	eventID := c.Param("id")
	var req struct {
		Name        string     `json:"name" binding:"required"`
		Description string     `json:"description"`
		Price       int32      `json:"price" binding:"required,min=0"`
		Quantity    int32      `json:"quantity" binding:"required,min=1"`
		MaxPerOrder int32      `json:"max_per_order" binding:"required,min=1"`
		SaleStartAt *time.Time `json:"sale_start_at"`
		SaleEndAt   *time.Time `json:"sale_end_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tt struct {
		ID          string     `json:"id"`
		EventID     string     `json:"event_id"`
		Name        string     `json:"name"`
		Description string     `json:"description"`
		Price       int32      `json:"price"`
		Quantity    int32      `json:"quantity"`
		MaxPerOrder int32      `json:"max_per_order"`
		SaleStartAt *time.Time `json:"sale_start_at"`
		SaleEndAt   *time.Time `json:"sale_end_at"`
		Status      string     `json:"status"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   time.Time  `json:"updated_at"`
	}
	err := h.pool.QueryRow(context.Background(),
		`INSERT INTO ticket_types (id, event_id, name, description, price, quantity, max_per_order, sale_start_at, sale_end_at, status, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, 'active', NOW(), NOW())
		 RETURNING id, event_id, name, description, price, quantity, max_per_order, sale_start_at, sale_end_at, status, created_at, updated_at`,
		eventID, req.Name, req.Description, req.Price, req.Quantity, req.MaxPerOrder, req.SaleStartAt, req.SaleEndAt).
		Scan(&tt.ID, &tt.EventID, &tt.Name, &tt.Description, &tt.Price, &tt.Quantity, &tt.MaxPerOrder, &tt.SaleStartAt, &tt.SaleEndAt, &tt.Status, &tt.CreatedAt, &tt.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create ticket type"})
		return
	}
	c.JSON(http.StatusCreated, tt)
}

func (h *EventHandler) ListTicketTypes(c *gin.Context) {
	eventID := c.Param("id")

	rows, err := h.pool.Query(context.Background(),
		`SELECT id, event_id, name, description, price, quantity, max_per_order, sale_start_at, sale_end_at, status, created_at, updated_at
		 FROM ticket_types WHERE event_id = $1 ORDER BY price`, eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list ticket types"})
		return
	}
	defer rows.Close()

	var types []map[string]interface{}
	for rows.Next() {
		var id, eventID2, name, description, status string
		var price, quantity, maxPerOrder int32
		var saleStartAt, saleEndAt, createdAt, updatedAt *time.Time
		if err := rows.Scan(&id, &eventID2, &name, &description, &price, &quantity, &maxPerOrder, &saleStartAt, &saleEndAt, &status, &createdAt, &updatedAt); err != nil {
			continue
		}
		types = append(types, map[string]interface{}{
			"id": id, "event_id": eventID2, "name": name, "description": description,
			"price": price, "quantity": quantity, "max_per_order": maxPerOrder,
			"sale_start_at": saleStartAt, "sale_end_at": saleEndAt,
			"status": status, "created_at": createdAt, "updated_at": updatedAt,
		})
	}
	if types == nil {
		types = []map[string]interface{}{}
	}
	c.JSON(http.StatusOK, types)
}
