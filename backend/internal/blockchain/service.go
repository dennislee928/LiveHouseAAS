package blockchain

import (
	"context"
	"errors"
	"fmt"
	"math/big"
)

var (
	ErrNotImplemented = errors.New("blockchain: not implemented")
	ErrMockOnly       = errors.New("blockchain: running in mock mode")
)

type Network string

const (
	NetworkPolygon Network = "polygon"
	NetworkMock    Network = "mock"
)

type TicketMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	EventName   string `json:"event_name"`
	VenueName   string `json:"venue_name"`
	EventDate   string `json:"event_date"`
	TicketType  string `json:"ticket_type"`
	IsPOAP      bool   `json:"is_poap"`
}

type MintResult struct {
	TokenID         uint64 `json:"token_id"`
	ContractAddress string `json:"contract_address"`
	TxHash          string `json:"tx_hash"`
	Network         string `json:"network"`
	TokenURI        string `json:"token_uri"`
}

type Service interface {
	MintTicket(ctx context.Context, ownerAddress, code string, metadata TicketMetadata) (*MintResult, error)
	ClaimPOAP(ctx context.Context, tokenID uint64, claimerAddress string) (*MintResult, error)
	Network() Network
}

// --- Mock Service ---

type MockService struct {
	nextTokenID    uint64
	contractAddr   string
}

func NewMockService() *MockService {
	return &MockService{
		nextTokenID:  1,
		contractAddr: "0x0000000000000000000000000000000000000000",
	}
}

func (s *MockService) Network() Network { return NetworkMock }

func (s *MockService) MintTicket(ctx context.Context, ownerAddress, code string, metadata TicketMetadata) (*MintResult, error) {
	tokenID := s.nextTokenID
	s.nextTokenID++

	return &MintResult{
		TokenID:         tokenID,
		ContractAddress: s.contractAddr,
		TxHash:          fmt.Sprintf("0x%x_mock_tx_%d", tokenID, tokenID),
		Network:         string(NetworkMock),
		TokenURI:        fmt.Sprintf("https://api.livehouseaas.com/metadata/%d", tokenID),
	}, nil
}

func (s *MockService) ClaimPOAP(ctx context.Context, tokenID uint64, claimerAddress string) (*MintResult, error) {
	return &MintResult{
		TokenID:         tokenID,
		ContractAddress: s.contractAddr,
		TxHash:          fmt.Sprintf("0x%x_mock_poap_%d", tokenID, tokenID),
		Network:         string(NetworkMock),
	}, nil
}

// --- Polygon Service (stub) ---

type PolygonService struct {
	rpcURL         string
	contractAddr   string
	privateKey     string
	nextTokenID    uint64
}

func NewPolygonService(rpcURL, contractAddr, privateKey string) *PolygonService {
	return &PolygonService{
		rpcURL:       rpcURL,
		contractAddr: contractAddr,
		privateKey:   privateKey,
		nextTokenID:  1,
	}
}

func (s *PolygonService) Network() Network { return NetworkPolygon }

func (s *PolygonService) MintTicket(ctx context.Context, ownerAddress, code string, metadata TicketMetadata) (*MintResult, error) {
	return nil, fmt.Errorf("%w: polygon mint not implemented yet", ErrNotImplemented)
}

func (s *PolygonService) ClaimPOAP(ctx context.Context, tokenID uint64, claimerAddress string) (*MintResult, error) {
	return nil, fmt.Errorf("%w: polygon claim not implemented yet", ErrNotImplemented)
}

// --- Token URI Generator ---

func GenerateTokenURI(metadata TicketMetadata) string {
	return fmt.Sprintf(
		`{"name":"%s","description":"%s","image":"%s","attributes":[{"trait_type":"Event","value":"%s"},{"trait_type":"Venue","value":"%s"},{"trait_type":"Date","value":"%s"},{"trait_type":"Ticket Type","value":"%s"},{"trait_type":"POAP","value":"%v"}]}`,
		metadata.Name, metadata.Description, metadata.Image,
		metadata.EventName, metadata.VenueName, metadata.EventDate,
		metadata.TicketType, metadata.IsPOAP,
	)
}

func NewBigInt(v uint64) *big.Int {
	return new(big.Int).SetUint64(v)
}
