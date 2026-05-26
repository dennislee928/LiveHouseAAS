package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalyticsHandler struct {
	pool *pgxpool.Pool
}

func NewAnalyticsHandler(pool *pgxpool.Pool) *AnalyticsHandler {
	return &AnalyticsHandler{pool: pool}
}

func (h *AnalyticsHandler) RevenueOverTime(c *gin.Context) {
	period := c.DefaultQuery("period", "daily") // daily, weekly, monthly
	days := c.DefaultQuery("days", "30")

	var rows interface{ Scan(...interface{}) error; Close() }
	var err error

	switch period {
	case "weekly":
		rows, err = h.pool.Query(context.Background(),
			`SELECT date_trunc('week', paid_at)::date as date, COALESCE(SUM(total_amount), 0)
			 FROM orders WHERE status = 'paid' AND paid_at > NOW() - $1::interval
			 GROUP BY 1 ORDER BY 1`, days+" days")
	case "monthly":
		rows, err = h.pool.Query(context.Background(),
			`SELECT date_trunc('month', paid_at)::date as date, COALESCE(SUM(total_amount), 0)
			 FROM orders WHERE status = 'paid' AND paid_at > NOW() - $1::interval
			 GROUP BY 1 ORDER BY 1`, days+" days")
	default:
		rows, err = h.pool.Query(context.Background(),
			`SELECT paid_at::date as date, COALESCE(SUM(total_amount), 0)
			 FROM orders WHERE status = 'paid' AND paid_at > NOW() - $1::interval
			 GROUP BY 1 ORDER BY 1`, days+" days")
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query revenue"})
		return
	}
	defer rows.Close()

	type DataPoint struct {
		Date   string `json:"date"`
		Amount int64  `json:"amount"`
	}
	var data []DataPoint
	for rows.Next() {
		var date time.Time
		var amount int64
		if err := rows.Scan(&date, &amount); err != nil {
			continue
		}
		data = append(data, DataPoint{Date: date.Format("2006-01-02"), Amount: amount})
	}
	if data == nil {
		data = []DataPoint{}
	}
	c.JSON(http.StatusOK, data)
}

func (h *AnalyticsHandler) TopVenues(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	rows, err := h.pool.Query(context.Background(),
		`SELECT v.id, v.name, v.city, COUNT(DISTINCT o.id) as order_count,
		        COALESCE(SUM(o.total_amount), 0) as revenue
		 FROM venues v
		 LEFT JOIN events e ON v.id = e.venue_id
		 LEFT JOIN orders o ON e.id = o.event_id AND o.status = 'paid'
		 GROUP BY v.id, v.name, v.city
		 ORDER BY revenue DESC LIMIT $1`, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query top venues"})
		return
	}
	defer rows.Close()

	var venues []gin.H
	for rows.Next() {
		var id, name, city string
		var orderCount int64
		var revenue int64
		if err := rows.Scan(&id, &name, &city, &orderCount, &revenue); err != nil {
			continue
		}
		venues = append(venues, gin.H{
			"id": id, "name": name, "city": city,
			"order_count": orderCount, "revenue": revenue,
		})
	}
	if venues == nil {
		venues = []gin.H{}
	}
	c.JSON(http.StatusOK, venues)
}

func (h *AnalyticsHandler) TopEvents(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	rows, err := h.pool.Query(context.Background(),
		`SELECT e.id, e.title, e.start_at, v.name as venue_name,
		        COUNT(DISTINCT o.id) as order_count,
		        COUNT(t.id) as tickets_sold,
		        COALESCE(SUM(o.total_amount), 0) as revenue
		 FROM events e
		 JOIN venues v ON e.venue_id = v.id
		 LEFT JOIN orders o ON e.id = o.event_id AND o.status = 'paid'
		 LEFT JOIN tickets t ON e.id = t.event_id
		 GROUP BY e.id, e.title, e.start_at, v.name
		 ORDER BY revenue DESC LIMIT $1`, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query top events"})
		return
	}
	defer rows.Close()

	var events []gin.H
	for rows.Next() {
		var id, title, venueName string
		var startAt time.Time
		var orderCount, ticketsSold, revenue int64
		if err := rows.Scan(&id, &title, &startAt, &venueName, &orderCount, &ticketsSold, &revenue); err != nil {
			continue
		}
		events = append(events, gin.H{
			"id": id, "title": title, "start_at": startAt,
			"venue_name": venueName, "order_count": orderCount,
			"tickets_sold": ticketsSold, "revenue": revenue,
		})
	}
	if events == nil {
		events = []gin.H{}
	}
	c.JSON(http.StatusOK, events)
}

func (h *AnalyticsHandler) BookingTrends(c *gin.Context) {
	days := c.DefaultQuery("days", "30")
	rows, err := h.pool.Query(context.Background(),
		`SELECT created_at::date as date,
		        COUNT(*) as total_bookings,
		        COUNT(*) FILTER (WHERE status = 'approved') as approved,
		        COUNT(*) FILTER (WHERE status = 'rejected') as rejected
		 FROM booking_requests
		 WHERE created_at > NOW() - $1::interval
		 GROUP BY 1 ORDER BY 1`, days+" days")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query booking trends"})
		return
	}
	defer rows.Close()

	type BookingPoint struct {
		Date           string `json:"date"`
		TotalBookings  int64  `json:"total_bookings"`
		Approved       int64  `json:"approved"`
		Rejected       int64  `json:"rejected"`
	}
	var data []BookingPoint
	for rows.Next() {
		var date time.Time
		var total, approved, rejected int64
		if err := rows.Scan(&date, &total, &approved, &rejected); err != nil {
			continue
		}
		data = append(data, BookingPoint{
			Date: date.Format("2006-01-02"), TotalBookings: total,
			Approved: approved, Rejected: rejected,
		})
	}
	if data == nil {
		data = []BookingPoint{}
	}
	c.JSON(http.StatusOK, data)
}

func (h *AnalyticsHandler) Summary(c *gin.Context) {
	var totalRevenue, totalOrders, totalTicketsSold, totalBookings, totalUsers, totalVenues int64

	h.pool.QueryRow(context.Background(),
		`SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE status = 'paid'`).Scan(&totalRevenue)
	h.pool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM orders`).Scan(&totalOrders)
	h.pool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM tickets WHERE status IN ('active', 'used')`).Scan(&totalTicketsSold)
	h.pool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM booking_requests`).Scan(&totalBookings)
	h.pool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM users`).Scan(&totalUsers)
	h.pool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM venues`).Scan(&totalVenues)

	// current period (this month)
	var thisMonthRevenue int64
	h.pool.QueryRow(context.Background(),
		`SELECT COALESCE(SUM(total_amount), 0) FROM orders
		 WHERE status = 'paid' AND date_trunc('month', paid_at) = date_trunc('month', NOW())`).Scan(&thisMonthRevenue)

	// last month
	var lastMonthRevenue int64
	h.pool.QueryRow(context.Background(),
		`SELECT COALESCE(SUM(total_amount), 0) FROM orders
		 WHERE status = 'paid' AND date_trunc('month', paid_at) = date_trunc('month', NOW() - interval '1 month')`).Scan(&lastMonthRevenue)

	var growth float64
	if lastMonthRevenue > 0 {
		growth = float64(thisMonthRevenue-lastMonthRevenue) / float64(lastMonthRevenue) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"total_revenue":        totalRevenue,
		"total_orders":         totalOrders,
		"total_tickets_sold":   totalTicketsSold,
		"total_bookings":       totalBookings,
		"total_users":          totalUsers,
		"total_venues":         totalVenues,
		"this_month_revenue":   thisMonthRevenue,
		"last_month_revenue":   lastMonthRevenue,
		"revenue_growth_pct":   growth,
	})
}

func (h *AnalyticsHandler) VenuePerformance(c *gin.Context) {
	userID, _ := c.Get("user_id")
	venueID := c.Query("venue_id")

	query := `SELECT v.id, v.name,
		COUNT(DISTINCT e.id) as total_events,
		COUNT(DISTINCT br.id) as total_bookings,
		COUNT(DISTINCT o.id) as total_orders,
		COALESCE(SUM(o.total_amount) FILTER (WHERE o.status = 'paid'), 0) as revenue
		FROM venues v
		LEFT JOIN events e ON v.id = e.venue_id
		LEFT JOIN booking_requests br ON v.id = br.venue_id
		LEFT JOIN orders o ON e.id = o.event_id
		WHERE v.owner_id = $1`

	args := []interface{}{userID}
	argIdx := 2

	if venueID != "" {
		query += ` AND v.id = $` + string(rune('0'+argIdx))
		args = append(args, venueID)
	}

	query += ` GROUP BY v.id, v.name ORDER BY revenue DESC`

	rows, err := h.pool.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query venue performance"})
		return
	}
	defer rows.Close()

	var performances []gin.H
	for rows.Next() {
		var id, name string
		var totalEvents, totalBookings, totalOrders, revenue int64
		if err := rows.Scan(&id, &name, &totalEvents, &totalBookings, &totalOrders, &revenue); err != nil {
			continue
		}
		performances = append(performances, gin.H{
			"id": id, "name": name,
			"total_events": totalEvents, "total_bookings": totalBookings,
			"total_orders": totalOrders, "revenue": revenue,
		})
	}
	if performances == nil {
		performances = []gin.H{}
	}
	c.JSON(http.StatusOK, performances)
}
