package slot

import "time"

type Status string

const (
	StatusAvailable Status = "available"
	StatusBooked    Status = "booked"
	StatusBlocked   Status = "blocked"
)

type Slot struct {
	ID        string    `json:"id"`
	VenueID   string    `json:"venue_id"`
	Date      string    `json:"date"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateSlotRequest struct {
	Date      string `json:"date" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

type BatchCreateSlotsRequest struct {
	Slots []CreateSlotRequest `json:"slots" binding:"required,min=1,dive"`
}
