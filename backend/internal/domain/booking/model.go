package booking

import "time"

type Status string

const (
	StatusPending   Status = "pending"
	StatusApproved  Status = "approved"
	StatusRejected  Status = "rejected"
	StatusCancelled Status = "cancelled"
	StatusConfirmed Status = "confirmed"
)

type BookingRequest struct {
	ID        string    `json:"id"`
	SlotID    string    `json:"slot_id"`
	VenueID   string    `json:"venue_id"`
	ArtistID  string    `json:"artist_id"`
	Message   string    `json:"message"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateBookingRequest struct {
	SlotID  string `json:"slot_id" binding:"required"`
	Message string `json:"message"`
}

type BookingDetail struct {
	BookingRequest
	Date       string `json:"date"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	ArtistName string `json:"artist_name,omitempty"`
	ArtistEmail string `json:"artist_email,omitempty"`
	VenueName  string `json:"venue_name,omitempty"`
	VenueCity  string `json:"venue_city,omitempty"`
}
