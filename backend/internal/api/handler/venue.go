package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dennis-lee/LiveHouseAAS/backend/internal/domain/venue"
)

type VenueHandler struct {
	pool *pgxpool.Pool
}

func NewVenueHandler(pool *pgxpool.Pool) *VenueHandler {
	return &VenueHandler{pool: pool}
}

func (h *VenueHandler) Create(c *gin.Context) {
	var req venue.CreateVenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ownerID, _ := c.Get("user_id")

	var v venue.Venue
	err := h.pool.QueryRow(context.Background(),
		`INSERT INTO venues (id, name, description, address, city, capacity, contact_phone, contact_email, owner_id, status, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, 'active', NOW(), NOW())
		 RETURNING id, name, description, address, city, capacity, contact_phone, contact_email, owner_id, status, created_at, updated_at`,
		req.Name, req.Description, req.Address, req.City, req.Capacity, req.ContactPhone, req.ContactEmail, ownerID).
		Scan(&v.ID, &v.Name, &v.Description, &v.Address, &v.City, &v.Capacity, &v.ContactPhone, &v.ContactEmail, &v.OwnerID, &v.Status, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create venue"})
		return
	}

	c.JSON(http.StatusCreated, v)
}

func (h *VenueHandler) List(c *gin.Context) {
	ownerID, _ := c.Get("user_id")

	rows, err := h.pool.Query(context.Background(),
		`SELECT id, name, description, address, city, capacity, contact_phone, contact_email, owner_id, status, created_at, updated_at
		 FROM venues WHERE owner_id = $1 ORDER BY created_at DESC`, ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list venues"})
		return
	}
	defer rows.Close()

	var venues []venue.Venue
	for rows.Next() {
		var v venue.Venue
		if err := rows.Scan(&v.ID, &v.Name, &v.Description, &v.Address, &v.City, &v.Capacity, &v.ContactPhone, &v.ContactEmail, &v.OwnerID, &v.Status, &v.CreatedAt, &v.UpdatedAt); err != nil {
			continue
		}
		venues = append(venues, v)
	}

	if venues == nil {
		venues = []venue.Venue{}
	}

	c.JSON(http.StatusOK, venues)
}

func (h *VenueHandler) ListAll(c *gin.Context) {
	rows, err := h.pool.Query(context.Background(),
		`SELECT id, name, description, address, city, capacity, contact_phone, contact_email, owner_id, status, created_at, updated_at
		 FROM venues WHERE status = 'active' ORDER BY created_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list venues"})
		return
	}
	defer rows.Close()

	var venues []venue.Venue
	for rows.Next() {
		var v venue.Venue
		if err := rows.Scan(&v.ID, &v.Name, &v.Description, &v.Address, &v.City, &v.Capacity, &v.ContactPhone, &v.ContactEmail, &v.OwnerID, &v.Status, &v.CreatedAt, &v.UpdatedAt); err != nil {
			continue
		}
		venues = append(venues, v)
	}

	if venues == nil {
		venues = []venue.Venue{}
	}

	c.JSON(http.StatusOK, venues)
}

func (h *VenueHandler) Get(c *gin.Context) {
	id := c.Param("id")

	var v venue.Venue
	err := h.pool.QueryRow(context.Background(),
		`SELECT id, name, description, address, city, capacity, contact_phone, contact_email, owner_id, status, created_at, updated_at
		 FROM venues WHERE id = $1`, id).
		Scan(&v.ID, &v.Name, &v.Description, &v.Address, &v.City, &v.Capacity, &v.ContactPhone, &v.ContactEmail, &v.OwnerID, &v.Status, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "venue not found"})
		return
	}

	c.JSON(http.StatusOK, v)
}

func (h *VenueHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req venue.UpdateVenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var v venue.Venue
	err := h.pool.QueryRow(context.Background(),
		`UPDATE venues SET name=$2, description=$3, address=$4, city=$5, capacity=$6, contact_phone=$7, contact_email=$8, status=$9, updated_at=NOW()
		 WHERE id=$1
		 RETURNING id, name, description, address, city, capacity, contact_phone, contact_email, owner_id, status, created_at, updated_at`,
		id, req.Name, req.Description, req.Address, req.City, req.Capacity, req.ContactPhone, req.ContactEmail, req.Status).
		Scan(&v.ID, &v.Name, &v.Description, &v.Address, &v.City, &v.Capacity, &v.ContactPhone, &v.ContactEmail, &v.OwnerID, &v.Status, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "venue not found"})
		return
	}

	c.JSON(http.StatusOK, v)
}

func (h *VenueHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	_, err := h.pool.Exec(context.Background(), `DELETE FROM venues WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete venue"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// --- Venue Specs ---

func (h *VenueHandler) CreateSpec(c *gin.Context) {
	venueID := c.Param("id")

	var req venue.CreateSpecRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var s venue.VenueSpec
	err := h.pool.QueryRow(context.Background(),
		`INSERT INTO venue_specs (id, venue_id, category, name, brand, quantity, description, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, NOW(), NOW())
		 RETURNING id, venue_id, category, name, brand, quantity, description, created_at, updated_at`,
		venueID, req.Category, req.Name, req.Brand, req.Quantity, req.Description).
		Scan(&s.ID, &s.VenueID, &s.Category, &s.Name, &s.Brand, &s.Quantity, &s.Description, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create spec"})
		return
	}

	c.JSON(http.StatusCreated, s)
}

func (h *VenueHandler) ListSpecs(c *gin.Context) {
	venueID := c.Param("id")

	rows, err := h.pool.Query(context.Background(),
		`SELECT id, venue_id, category, name, brand, quantity, description, created_at, updated_at
		 FROM venue_specs WHERE venue_id = $1 ORDER BY category, name`, venueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list specs"})
		return
	}
	defer rows.Close()

	var specs []venue.VenueSpec
	for rows.Next() {
		var s venue.VenueSpec
		if err := rows.Scan(&s.ID, &s.VenueID, &s.Category, &s.Name, &s.Brand, &s.Quantity, &s.Description, &s.CreatedAt, &s.UpdatedAt); err != nil {
			continue
		}
		specs = append(specs, s)
	}

	if specs == nil {
		specs = []venue.VenueSpec{}
	}

	c.JSON(http.StatusOK, specs)
}

func (h *VenueHandler) DeleteSpec(c *gin.Context) {
	specID := c.Param("specId")

	_, err := h.pool.Exec(context.Background(), `DELETE FROM venue_specs WHERE id = $1`, specID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete spec"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
