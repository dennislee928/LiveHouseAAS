package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dennis-lee/LiveHouseAAS/backend/internal/blockchain"
)

type NFTHandler struct {
	pool    *pgxpool.Pool
	svc     blockchain.Service
}

func NewNFTHandler(pool *pgxpool.Pool, svc blockchain.Service) *NFTHandler {
	return &NFTHandler{pool: pool, svc: svc}
}

func (h *NFTHandler) Claim(c *gin.Context) {
	ticketID := c.Param("ticketId")
	userID, _ := c.Get("user_id")

	var req struct {
		OwnerAddress string `json:"owner_address" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// verify ticket ownership
	var orderUserID, ticketStatus, eventID string
	err := h.pool.QueryRow(context.Background(),
		`SELECT o.user_id, t.status, t.event_id
		 FROM tickets t JOIN orders o ON t.order_id = o.id WHERE t.id = $1`, ticketID).
		Scan(&orderUserID, &ticketStatus, &eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}
	if orderUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your ticket"})
		return
	}
	if ticketStatus != "active" && ticketStatus != "used" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket cannot be minted as NFT"})
		return
	}

	// check if NFT already exists
	var nftCount int
	h.pool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM nft_tickets WHERE ticket_id = $1`, ticketID).Scan(&nftCount)
	if nftCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "NFT already claimed for this ticket"})
		return
	}

	// get event + ticket info for metadata
	var eventTitle, venueName, ticketTypeName, eventDate string
	var artistID string
	h.pool.QueryRow(context.Background(),
		`SELECT COALESCE(e.title,''), COALESCE(v.name,''), COALESCE(tt.name,''),
		        e.start_at::text, e.artist_id
		 FROM events e
		 JOIN venues v ON e.venue_id = v.id
		 JOIN ticket_types tt ON tt.event_id = e.id
		 WHERE e.id = $1 LIMIT 1`, eventID).
		Scan(&eventTitle, &venueName, &ticketTypeName, &eventDate, &artistID)

	// get ticket code
	var code string
	h.pool.QueryRow(context.Background(),
		`SELECT code FROM tickets WHERE id = $1`, ticketID).Scan(&code)

	meta := blockchain.TicketMetadata{
		Name:        eventTitle + " Ticket",
		Description: "LiveHouseAAS NFT Ticket for " + eventTitle,
		Image:       "https://api.livehouseaas.com/nft/ticket.png",
		EventName:   eventTitle,
		VenueName:   venueName,
		EventDate:   eventDate,
		TicketType:  ticketTypeName,
		IsPOAP:      false,
	}

	result, err := h.svc.MintTicket(context.Background(), req.OwnerAddress, code, meta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "NFT minting failed: " + err.Error()})
		return
	}

	// save NFT record
	var nftID, nftStatus string
	h.pool.QueryRow(context.Background(),
		`INSERT INTO nft_tickets (id, ticket_id, token_id, contract_address, tx_hash, network, token_uri, owner_address, is_poap, status, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, false, 'minted', NOW(), NOW())
		 RETURNING id, status`,
		ticketID, result.TokenID, result.ContractAddress, result.TxHash, result.Network, result.TokenURI, req.OwnerAddress).
		Scan(&nftID, &nftStatus)

	c.JSON(http.StatusCreated, gin.H{
		"nft_id":    nftID,
		"status":    nftStatus,
		"token_id":  result.TokenID,
		"tx_hash":   result.TxHash,
		"network":   result.Network,
		"token_uri": result.TokenURI,
	})
}

func (h *NFTHandler) ClaimPOAP(c *gin.Context) {
	ticketID := c.Param("ticketId")
	userID, _ := c.Get("user_id")

	// verify ownership
	var orderUserID string
	h.pool.QueryRow(context.Background(),
		`SELECT o.user_id FROM tickets t JOIN orders o ON t.order_id = o.id WHERE t.id = $1`, ticketID).
		Scan(&orderUserID)
	if orderUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your ticket"})
		return
	}

	// get existing NFT
	var nftID, nftStatus, ownerAddress string
	var tokenID int64
	err := h.pool.QueryRow(context.Background(),
		`SELECT id, status, token_id, owner_address FROM nft_tickets WHERE ticket_id = $1`, ticketID).
		Scan(&nftID, &nftStatus, &tokenID, &ownerAddress)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "NFT not minted yet, claim NFT first"})
		return
	}
	if nftStatus != "minted" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NFT not eligible for POAP"})
		return
	}

	result, err := h.svc.ClaimPOAP(context.Background(), uint64(tokenID), ownerAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "POAP claim failed: " + err.Error()})
		return
	}

	var poapStatus string
	h.pool.QueryRow(context.Background(),
		`UPDATE nft_tickets SET is_poap=true, poap_claimed_at=NOW(), status='claimed', tx_hash=$2, updated_at=NOW()
		 WHERE id=$1 RETURNING status`, nftID, result.TxHash).Scan(&poapStatus)

	c.JSON(http.StatusOK, gin.H{
		"nft_id":  nftID,
		"status":  poapStatus,
		"tx_hash": result.TxHash,
		"poap":    true,
	})
}

func (h *NFTHandler) GetStatus(c *gin.Context) {
	ticketID := c.Param("ticketId")

	var nftID, status, network, contractAddr, tokenURI string
	var tokenID int64
	var isPOAP bool
	err := h.pool.QueryRow(context.Background(),
		`SELECT id, status, token_id, COALESCE(contract_address,''), COALESCE(token_uri,''), network, is_poap
		 FROM nft_tickets WHERE ticket_id = $1`, ticketID).
		Scan(&nftID, &status, &tokenID, &contractAddr, &tokenURI, &network, &isPOAP)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"nft_claimed": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"nft_claimed": true,
		"nft_id":      nftID,
		"status":      status,
		"token_id":    tokenID,
		"contract":    contractAddr,
		"network":     network,
		"token_uri":   tokenURI,
		"is_poap":     isPOAP,
	})
}
