package order

import "time"

type Status string

const (
	StatusPending   Status = "pending"
	StatusPaid      Status = "paid"
	StatusCancelled Status = "cancelled"
	StatusRefunded  Status = "refunded"
)

type Order struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	EventID       string     `json:"event_id"`
	TotalAmount   int32      `json:"total_amount"`
	Status        Status     `json:"status"`
	PaymentMethod string     `json:"payment_method,omitempty"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type OrderItem struct {
	ID           string `json:"id"`
	OrderID      string `json:"order_id"`
	TicketTypeID string `json:"ticket_type_id"`
	Quantity     int32  `json:"quantity"`
	UnitPrice    int32  `json:"unit_price"`
	Subtotal     int32  `json:"subtotal"`
}

type PurchaseRequest struct {
	Items []PurchaseItem `json:"items" binding:"required,min=1,dive"`
}

type PurchaseItem struct {
	TicketTypeID string `json:"ticket_type_id" binding:"required"`
	Quantity     int32  `json:"quantity" binding:"required,min=1"`
}

type PurchaseResponse struct {
	Order   Order        `json:"order"`
	Payment PaymentInfo  `json:"payment"`
	Tickets []TicketInfo `json:"tickets,omitempty"`
}

type PaymentInfo struct {
	Provider      string `json:"provider"`
	RedirectURL   string `json:"redirect_url,omitempty"`
	PaymentURL    string `json:"payment_url,omitempty"`
}

type TicketInfo struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	QRSegret string `json:"qr_secret"`
}
