package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/dennis-lee/LiveHouseAAS/backend/internal/auth"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/domain/user"
)

type AuthHandler struct {
	pool *pgxpool.Pool
	jwt  *auth.JWT
}

func NewAuthHandler(pool *pgxpool.Pool, jwt *auth.JWT) *AuthHandler {
	return &AuthHandler{pool: pool, jwt: jwt}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req user.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exists bool
	err := h.pool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	id := auth.NewID()
	_, err = h.pool.Exec(context.Background(),
		`INSERT INTO users (id, email, password_hash, name, role, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`,
		id, req.Email, string(hashedPassword), req.Name, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	token, err := h.jwt.GenerateToken(id, req.Email, string(req.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, user.AuthResponse{
		Token: token,
		User: user.User{
			ID:    id,
			Email: req.Email,
			Name:  req.Name,
			Role:  req.Role,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req user.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var u user.User
	err := h.pool.QueryRow(context.Background(),
		`SELECT id, email, password_hash, name, role, COALESCE(avatar_url, ''), created_at, updated_at
		 FROM users WHERE email = $1`, req.Email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	token, err := h.jwt.GenerateToken(u.ID, u.Email, string(u.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, user.AuthResponse{
		Token: token,
		User:  u,
	})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var u user.User
	err := h.pool.QueryRow(context.Background(),
		`SELECT id, email, name, role, COALESCE(avatar_url, ''), created_at, updated_at
		 FROM users WHERE id = $1`, userID).Scan(
		&u.ID, &u.Email, &u.Name, &u.Role, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, u)
}
