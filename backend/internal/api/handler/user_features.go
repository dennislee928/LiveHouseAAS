package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/dennis-lee/LiveHouseAAS/backend/internal/notification"
)

type UserFeaturesHandler struct {
	pool     *pgxpool.Pool
	notifSvc notification.Service
}

func NewUserFeaturesHandler(pool *pgxpool.Pool, notifSvc notification.Service) *UserFeaturesHandler {
	return &UserFeaturesHandler{pool: pool, notifSvc: notifSvc}
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (h *UserFeaturesHandler) RequestEmailVerification(c *gin.Context) {
	userID, _ := c.Get("user_id")
	email, _ := c.Get("email")

	token := generateToken()
	_, err := h.pool.Exec(context.Background(),
		`INSERT INTO user_tokens (id, user_id, token, type, expires_at, created_at)
		 VALUES (gen_random_uuid(), $1, $2, 'email_verify', NOW() + interval '24 hours', NOW())
		 ON CONFLICT (user_id, type) DO UPDATE SET token = $2, expires_at = NOW() + interval '24 hours', created_at = NOW()`,
		userID, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create verification token"})
		return
	}

	verifyURL := "/verify-email?token=" + token
	h.notifSvc.SendEmail(email.(string), "驗證您的 Email", "請點擊以下連結驗證您的 Email：\n\n"+verifyURL)

	c.JSON(http.StatusOK, gin.H{
		"message":    "verification email sent",
		"verify_url": verifyURL,
		"email":      email,
	})
}

func (h *UserFeaturesHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token required"})
		return
	}

	var userID, email string
	err := h.pool.QueryRow(context.Background(),
		`UPDATE users SET updated_at = NOW()
		 WHERE id = (SELECT user_id FROM user_tokens WHERE token = $1 AND type = 'email_verify' AND expires_at > NOW())
		 RETURNING id, email`, token).Scan(&userID, &email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
		return
	}

	h.pool.Exec(context.Background(),
		`DELETE FROM user_tokens WHERE user_id = $1 AND type = 'email_verify'`, userID)

	h.pool.Exec(context.Background(),
		`UPDATE users SET avatar_url = COALESCE(avatar_url, '') WHERE id = $1`, userID)

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}

func (h *UserFeaturesHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID string
	err := h.pool.QueryRow(context.Background(),
		`SELECT id FROM users WHERE email = $1`, req.Email).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "if the email exists, a reset link has been sent"})
		return
	}

	token := generateToken()
	h.pool.Exec(context.Background(),
		`INSERT INTO user_tokens (id, user_id, token, type, expires_at, created_at)
		 VALUES (gen_random_uuid(), $1, $2, 'password_reset', NOW() + interval '1 hour', NOW())
		 ON CONFLICT (user_id, type) DO UPDATE SET token = $2, expires_at = NOW() + interval '1 hour', created_at = NOW()`,
		userID, token)

	resetURL := "/reset-password?token=" + token
	h.notifSvc.SendEmail(req.Email, "重設您的密碼", "請點擊以下連結重設您的密碼（有效期 1 小時）：\n\n"+resetURL)

	c.JSON(http.StatusOK, gin.H{
		"message":   "if the email exists, a reset link has been sent",
		"reset_url": resetURL,
	})
}

func (h *UserFeaturesHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID string
	err := h.pool.QueryRow(context.Background(),
		`SELECT user_id FROM user_tokens
		 WHERE token = $1 AND type = 'password_reset' AND expires_at > NOW()`, req.Token).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	h.pool.Exec(context.Background(),
		`UPDATE users SET password_hash = $2, updated_at = NOW() WHERE id = $1`, userID, string(hashedPassword))
	h.pool.Exec(context.Background(),
		`DELETE FROM user_tokens WHERE user_id = $1 AND type = 'password_reset'`, userID)

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

func (h *UserFeaturesHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		h.pool.Exec(context.Background(),
			`UPDATE users SET name = $2, updated_at = NOW() WHERE id = $1`, userID, req.Name)
	}
	if req.AvatarURL != "" {
		h.pool.Exec(context.Background(),
			`UPDATE users SET avatar_url = $2, updated_at = NOW() WHERE id = $1`, userID, req.AvatarURL)
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}

func (h *UserFeaturesHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var hash string
	h.pool.QueryRow(context.Background(),
		`SELECT password_hash FROM users WHERE id = $1`, userID).Scan(&hash)

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "current password is incorrect"})
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	h.pool.Exec(context.Background(),
		`UPDATE users SET password_hash = $2, updated_at = NOW() WHERE id = $1`, userID, string(newHash))

	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}

func (h *UserFeaturesHandler) UpdateAvatar(c *gin.Context) {
	userID, _ := c.Get("user_id")
	avatarURL := c.PostForm("avatar_url")
	if avatarURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "avatar_url required"})
		return
	}
	h.pool.Exec(context.Background(),
		`UPDATE users SET avatar_url = $2, updated_at = NOW() WHERE id = $1`, userID, avatarURL)
	c.JSON(http.StatusOK, gin.H{"message": "avatar updated"})
}
