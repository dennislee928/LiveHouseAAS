package venue

import "time"

type Status string

const (
	StatusActive      Status = "active"
	StatusInactive    Status = "inactive"
	StatusMaintenance Status = "maintenance"
)

type Venue struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Address       string    `json:"address"`
	City          string    `json:"city"`
	Capacity      int32     `json:"capacity"`
	ContactPhone  string    `json:"contact_phone"`
	ContactEmail  string    `json:"contact_email"`
	OwnerID       string    `json:"owner_id"`
	Status        Status    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateVenueRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	Address      string `json:"address" binding:"required"`
	City         string `json:"city" binding:"required"`
	Capacity     int32  `json:"capacity" binding:"required,min=1"`
	ContactPhone string `json:"contact_phone"`
	ContactEmail string `json:"contact_email"`
}

type UpdateVenueRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	Address      string `json:"address" binding:"required"`
	City         string `json:"city" binding:"required"`
	Capacity     int32  `json:"capacity" binding:"required,min=1"`
	ContactPhone string `json:"contact_phone"`
	ContactEmail string `json:"contact_email"`
	Status       Status `json:"status" binding:"required,oneof=active inactive maintenance"`
}

type VenueSpec struct {
	ID          string    `json:"id"`
	VenueID     string    `json:"venue_id"`
	Category    string    `json:"category"`
	Name        string    `json:"name"`
	Brand       string    `json:"brand"`
	Quantity    int32     `json:"quantity"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateSpecRequest struct {
	Category    string `json:"category" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Brand       string `json:"brand"`
	Quantity    int32  `json:"quantity" binding:"required,min=1"`
	Description string `json:"description"`
}
