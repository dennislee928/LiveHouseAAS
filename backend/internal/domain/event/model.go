package event

import "time"

type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusCancelled Status = "cancelled"
	StatusCompleted Status = "completed"
)

type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	VenueID     string    `json:"venue_id"`
	ArtistID    string    `json:"artist_id"`
	BookingID   string    `json:"booking_id,omitempty"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EventDetail struct {
	Event
	ArtistName string `json:"artist_name,omitempty"`
	VenueName  string `json:"venue_name,omitempty"`
	VenueCity  string `json:"venue_city,omitempty"`
}
