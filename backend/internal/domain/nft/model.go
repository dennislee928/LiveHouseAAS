package nft

import "time"

type Status string

const (
	StatusPending  Status = "pending"
	StatusMinted   Status = "minted"
	StatusFailed   Status = "failed"
	StatusClaimed  Status = "claimed"
)

type TicketNFT struct {
	ID              string    `json:"id"`
	TicketID        string    `json:"ticket_id"`
	TokenID         int64     `json:"token_id"`
	ContractAddress string    `json:"contract_address"`
	TxHash          string    `json:"tx_hash,omitempty"`
	Network         string    `json:"network"`
	TokenURI        string    `json:"token_uri,omitempty"`
	OwnerAddress    string    `json:"owner_address"`
	IsPOAP          bool      `json:"is_poap"`
	POAPClaimedAt   *time.Time `json:"poap_claimed_at,omitempty"`
	Status          Status    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ClaimNFTRequest struct {
	OwnerAddress string `json:"owner_address" binding:"required"`
}

type NFTMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	EventName   string `json:"event_name"`
	VenueName   string `json:"venue_name"`
	EventDate   string `json:"event_date"`
	TicketType  string `json:"ticket_type"`
	IsPOAP      bool   `json:"is_poap"`
}
