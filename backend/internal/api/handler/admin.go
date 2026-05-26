package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminHandler struct {
	pool *pgxpool.Pool
}

func NewAdminHandler(pool *pgxpool.Pool) *AdminHandler {
	return &AdminHandler{pool: pool}
}

func (h *AdminHandler) Stats(c *gin.Context) {
	var totalUsers, totalVenues, totalEvents, totalOrders, totalRevenue int64
	h.pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM users`).Scan(&totalUsers)
	h.pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM venues`).Scan(&totalVenues)
	h.pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM events`).Scan(&totalEvents)
	h.pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM orders`).Scan(&totalOrders)
	h.pool.QueryRow(context.Background(),
		`SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE status = 'paid'`).Scan(&totalRevenue)

	c.JSON(http.StatusOK, gin.H{
		"total_users":   totalUsers,
		"total_venues":  totalVenues,
		"total_events":  totalEvents,
		"total_orders":  totalOrders,
		"total_revenue": totalRevenue,
	})
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	rows, err := h.pool.Query(context.Background(),
		`SELECT id, email, name, role, COALESCE(avatar_url,''), created_at, updated_at
		 FROM users ORDER BY created_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}
	defer rows.Close()

	var list []gin.H
	for rows.Next() {
		var id, email, name, role, avatar string
		var createdAt, updatedAt interface{}
		if err := rows.Scan(&id, &email, &name, &role, &avatar, &createdAt, &updatedAt); err != nil {
			continue
		}
		list = append(list, gin.H{
			"id": id, "email": email, "name": name, "role": role,
			"avatar_url": avatar, "created_at": createdAt, "updated_at": updatedAt,
		})
	}
	if list == nil {
		list = []gin.H{}
	}
	c.JSON(http.StatusOK, list)
}

func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")
	var req struct {
		Role string `json:"role" binding:"required,oneof=admin venue artist"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.pool.Exec(context.Background(),
		`UPDATE users SET role = $2, updated_at = NOW() WHERE id = $1`, userID, req.Role)
	if err != nil || result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "role updated"})
}

func (h *AdminHandler) ListVenues(c *gin.Context) {
	rows, err := h.pool.Query(context.Background(),
		`SELECT v.id, v.name, v.description, v.address, v.city, v.capacity,
		        v.contact_phone, v.contact_email, v.owner_id, v.status, v.created_at, v.updated_at,
		        COALESCE(u.name,'')
		 FROM venues v LEFT JOIN users u ON v.owner_id = u.id
		 ORDER BY v.created_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list venues"})
		return
	}
	defer rows.Close()

	var list []gin.H
	for rows.Next() {
		var id, name, desc, addr, city, phone, email, ownerID, status, ownerName string
		var capacity int32
		var createdAt, updatedAt interface{}
		if err := rows.Scan(&id, &name, &desc, &addr, &city, &capacity, &phone, &email, &ownerID, &status, &createdAt, &updatedAt, &ownerName); err != nil {
			continue
		}
		list = append(list, gin.H{
			"id": id, "name": name, "description": desc, "address": addr,
			"city": city, "capacity": capacity, "contact_phone": phone,
			"contact_email": email, "owner_id": ownerID, "status": status,
			"owner_name": ownerName, "created_at": createdAt, "updated_at": updatedAt,
		})
	}
	if list == nil {
		list = []gin.H{}
	}
	c.JSON(http.StatusOK, list)
}

func (h *AdminHandler) UpdateVenueStatus(c *gin.Context) {
	venueID := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required,oneof=active inactive maintenance"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.pool.Exec(context.Background(),
		`UPDATE venues SET status = $2, updated_at = NOW() WHERE id = $1`, venueID, req.Status)
	if err != nil || result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "venue not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "venue status updated"})
}

func (h *AdminHandler) ListEvents(c *gin.Context) {
	rows, err := h.pool.Query(context.Background(),
		`SELECT e.id, e.title, e.description, e.venue_id, e.artist_id, e.start_at, e.end_at, e.status,
		        e.created_at, e.updated_at, COALESCE(v.name,''), COALESCE(u.name,'')
		 FROM events e
		 LEFT JOIN venues v ON e.venue_id = v.id
		 LEFT JOIN users u ON e.artist_id = u.id
		 ORDER BY e.created_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list events"})
		return
	}
	defer rows.Close()

	var list []gin.H
	for rows.Next() {
		var id, title, desc, venueID, artistID, status, venueName, artistName string
		var startAt, endAt, createdAt, updatedAt interface{}
		if err := rows.Scan(&id, &title, &desc, &venueID, &artistID, &startAt, &endAt, &status, &createdAt, &updatedAt, &venueName, &artistName); err != nil {
			continue
		}
		list = append(list, gin.H{
			"id": id, "title": title, "description": desc,
			"venue_id": venueID, "artist_id": artistID,
			"start_at": startAt, "end_at": endAt, "status": status,
			"venue_name": venueName, "artist_name": artistName,
			"created_at": createdAt, "updated_at": updatedAt,
		})
	}
	if list == nil {
		list = []gin.H{}
	}
	c.JSON(http.StatusOK, list)
}

func (h *AdminHandler) ListBookings(c *gin.Context) {
	rows, err := h.pool.Query(context.Background(),
		`SELECT br.id, br.slot_id, br.venue_id, br.artist_id, br.message, br.status,
		        br.created_at, br.updated_at,
		        COALESCE(v.name,''), COALESCE(u.name,'')
		 FROM booking_requests br
		 LEFT JOIN venues v ON br.venue_id = v.id
		 LEFT JOIN users u ON br.artist_id = u.id
		 ORDER BY br.created_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list bookings"})
		return
	}
	defer rows.Close()

	var list []gin.H
	for rows.Next() {
		var id, slotID, venueID, artistID, message, status, venueName, artistName string
		var createdAt, updatedAt interface{}
		if err := rows.Scan(&id, &slotID, &venueID, &artistID, &message, &status, &createdAt, &updatedAt, &venueName, &artistName); err != nil {
			continue
		}
		list = append(list, gin.H{
			"id": id, "slot_id": slotID, "venue_id": venueID, "artist_id": artistID,
			"message": message, "status": status,
			"venue_name": venueName, "artist_name": artistName,
			"created_at": createdAt, "updated_at": updatedAt,
		})
	}
	if list == nil {
		list = []gin.H{}
	}
	c.JSON(http.StatusOK, list)
}

func (h *AdminHandler) ListOrders(c *gin.Context) {
	rows, err := h.pool.Query(context.Background(),
		`SELECT o.id, o.user_id, o.event_id, o.total_amount, o.status, COALESCE(o.payment_method,''),
		        o.paid_at, o.created_at, o.updated_at,
		        COALESCE(e.title,''), COALESCE(u.name,'')
		 FROM orders o
		 LEFT JOIN events e ON o.event_id = e.id
		 LEFT JOIN users u ON o.user_id = u.id
		 ORDER BY o.created_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list orders"})
		return
	}
	defer rows.Close()

	var list []gin.H
	for rows.Next() {
		var id, userID, eventID, status, paymentMethod, eventTitle, userName string
		var totalAmount int32
		var paidAt, createdAt, updatedAt interface{}
		if err := rows.Scan(&id, &userID, &eventID, &totalAmount, &status, &paymentMethod, &paidAt, &createdAt, &updatedAt, &eventTitle, &userName); err != nil {
			continue
		}
		list = append(list, gin.H{
			"id": id, "user_id": userID, "event_id": eventID,
			"total_amount": totalAmount, "status": status,
			"payment_method": paymentMethod, "paid_at": paidAt,
			"event_title": eventTitle, "user_name": userName,
			"created_at": createdAt, "updated_at": updatedAt,
		})
	}
	if list == nil {
		list = []gin.H{}
	}
	c.JSON(http.StatusOK, list)
}
