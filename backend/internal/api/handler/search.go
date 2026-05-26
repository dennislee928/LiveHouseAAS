package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SearchHandler struct {
	pool *pgxpool.Pool
}

func NewSearchHandler(pool *pgxpool.Pool) *SearchHandler {
	return &SearchHandler{pool: pool}
}

func (h *SearchHandler) Events(c *gin.Context) {
	q := c.Query("q")
	city := c.Query("city")
	venueID := c.Query("venue_id")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0")

	query := `SELECT e.id, e.title, e.description, e.start_at, e.end_at, e.status,
		v.name as venue_name, v.city as venue_city, u.name as artist_name
		FROM events e
		JOIN venues v ON e.venue_id = v.id
		JOIN users u ON e.artist_id = u.id
		WHERE e.status = 'published'`

	args := []interface{}{}
	argIdx := 1

	if q != "" {
		query += ` AND (e.title ILIKE $` + itoa(argIdx) + ` OR e.description ILIKE $` + itoa(argIdx) + `)`
		args = append(args, "%"+q+"%")
		argIdx++
	}
	if city != "" {
		query += ` AND v.city = $` + itoa(argIdx)
		args = append(args, city)
		argIdx++
	}
	if venueID != "" {
		query += ` AND e.venue_id = $` + itoa(argIdx)
		args = append(args, venueID)
		argIdx++
	}
	if dateFrom != "" {
		query += ` AND e.start_at >= $` + itoa(argIdx)
		args = append(args, dateFrom)
		argIdx++
	}
	if dateTo != "" {
		query += ` AND e.start_at <= $` + itoa(argIdx)
		args = append(args, dateTo)
		argIdx++
	}

	query += ` ORDER BY e.start_at LIMIT $` + itoa(argIdx) + ` OFFSET $` + itoa(argIdx+1)
	args = append(args, limit, offset)

	rows, err := h.pool.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}
	defer rows.Close()

	var results []gin.H
	for rows.Next() {
		var id, title, desc, status, venueName, venueCity, artistName string
		var startAt, endAt interface{}
		if err := rows.Scan(&id, &title, &desc, &startAt, &endAt, &status, &venueName, &venueCity, &artistName); err != nil {
			continue
		}
		results = append(results, gin.H{
			"id": id, "title": title, "description": desc,
			"start_at": startAt, "end_at": endAt, "status": status,
			"venue_name": venueName, "venue_city": venueCity, "artist_name": artistName,
		})
	}
	if results == nil {
		results = []gin.H{}
	}
	c.JSON(http.StatusOK, results)
}

func (h *SearchHandler) Venues(c *gin.Context) {
	q := c.Query("q")
	city := c.Query("city")
	minCap := c.Query("min_capacity")
	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0")

	query := `SELECT v.id, v.name, v.description, v.address, v.city, v.capacity,
		v.contact_phone, v.contact_email
		FROM venues v WHERE v.status = 'active'`

	args := []interface{}{}
	argIdx := 1

	if q != "" {
		query += ` AND (v.name ILIKE $` + itoa(argIdx) + ` OR v.description ILIKE $` + itoa(argIdx) + ` OR v.city ILIKE $` + itoa(argIdx) + `)`
		args = append(args, "%"+q+"%")
		argIdx++
	}
	if city != "" {
		query += ` AND v.city = $` + itoa(argIdx)
		args = append(args, city)
		argIdx++
	}
	if minCap != "" {
		query += ` AND v.capacity >= $` + itoa(argIdx)
		args = append(args, minCap)
		argIdx++
	}

	query += ` ORDER BY v.name LIMIT $` + itoa(argIdx) + ` OFFSET $` + itoa(argIdx+1)
	args = append(args, limit, offset)

	rows, err := h.pool.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}
	defer rows.Close()

	var results []gin.H
	for rows.Next() {
		var id, name, desc, addr, city, phone, email string
		var capacity int32
		if err := rows.Scan(&id, &name, &desc, &addr, &city, &capacity, &phone, &email); err != nil {
			continue
		}
		results = append(results, gin.H{
			"id": id, "name": name, "description": desc,
			"address": addr, "city": city, "capacity": capacity,
			"contact_phone": phone, "contact_email": email,
		})
	}
	if results == nil {
		results = []gin.H{}
	}
	c.JSON(http.StatusOK, results)
}

func itoa(i int) string {
	return string(rune('0' + i%10)) + func() string {
		if i >= 10 {
			return string(rune('0' + i/10))
		}
		return ""
	}()
}
