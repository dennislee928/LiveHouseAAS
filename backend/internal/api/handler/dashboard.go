package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardHandler struct {
	pool *pgxpool.Pool
}

func NewDashboardHandler(pool *pgxpool.Pool) *DashboardHandler {
	return &DashboardHandler{pool: pool}
}

func (h *DashboardHandler) Stats(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var stats map[string]interface{}
	stats = make(map[string]interface{})

	if role == "venue" {
		var totalRevenue int64
		h.pool.QueryRow(context.Background(),
			`SELECT COALESCE(SUM(o.total_amount), 0)
			 FROM orders o
			 JOIN events e ON o.event_id = e.id
			 JOIN venues v ON e.venue_id = v.id
			 WHERE v.owner_id = $1 AND o.status = 'paid'`, userID).Scan(&totalRevenue)

		var ticketSold int64
		h.pool.QueryRow(context.Background(),
			`SELECT COALESCE(COUNT(*), 0)
			 FROM tickets t
			 JOIN events e ON t.event_id = e.id
			 JOIN venues v ON e.venue_id = v.id
			 WHERE v.owner_id = $1`, userID).Scan(&ticketSold)

		var upcomingEvents int64
		h.pool.QueryRow(context.Background(),
			`SELECT COUNT(*) FROM events WHERE venue_id IN
			 (SELECT id FROM venues WHERE owner_id = $1)
			 AND status = 'published' AND start_at > NOW()`, userID).Scan(&upcomingEvents)

		var pendingBookings int64
		h.pool.QueryRow(context.Background(),
			`SELECT COUNT(*) FROM booking_requests WHERE venue_id IN
			 (SELECT id FROM venues WHERE owner_id = $1)
			 AND status = 'pending'`, userID).Scan(&pendingBookings)

		var kybStatus string
		h.pool.QueryRow(context.Background(),
			`SELECT COALESCE(status, 'not_submitted') FROM business_verifications WHERE user_id = $1`, userID).Scan(&kybStatus)

		stats = map[string]interface{}{
			"total_revenue":     totalRevenue,
			"tickets_sold":      ticketSold,
			"upcoming_events":   upcomingEvents,
			"pending_bookings":  pendingBookings,
			"kyb_status":        kybStatus,
		}

	} else {
		var upcomingEvents int64
		h.pool.QueryRow(context.Background(),
			`SELECT COUNT(*) FROM events WHERE artist_id = $1 AND start_at > NOW()`, userID).Scan(&upcomingEvents)

		var pendingBookings int64
		h.pool.QueryRow(context.Background(),
			`SELECT COUNT(*) FROM booking_requests WHERE artist_id = $1 AND status = 'pending'`, userID).Scan(&pendingBookings)

		var totalTickets int64
		h.pool.QueryRow(context.Background(),
			`SELECT COUNT(*) FROM tickets t JOIN orders o ON t.order_id = o.id WHERE o.user_id = $1`, userID).Scan(&totalTickets)

		var paidOrders int64
		h.pool.QueryRow(context.Background(),
			`SELECT COUNT(*) FROM orders WHERE user_id = $1 AND status = 'paid'`, userID).Scan(&paidOrders)

		stats = map[string]interface{}{
			"upcoming_events":  upcomingEvents,
			"pending_bookings": pendingBookings,
			"total_tickets":    totalTickets,
			"paid_orders":      paidOrders,
		}
	}

	c.JSON(http.StatusOK, stats)
}
