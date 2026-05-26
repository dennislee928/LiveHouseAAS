package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UploadHandler struct {
	pool      *pgxpool.Pool
	uploadDir string
	maxSize   int64
}

func NewUploadHandler(pool *pgxpool.Pool, uploadDir string, maxSize int64) *UploadHandler {
	os.MkdirAll(uploadDir, 0755)
	return &UploadHandler{pool: pool, uploadDir: uploadDir, maxSize: maxSize}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	defer file.Close()

	if header.Size > h.maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("file too large (max %d bytes)", h.maxSize)})
		return
	}

	// validate file type
	allowed := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".pdf": true, ".doc": true, ".docx": true,
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed"})
		return
	}

	// generate unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
	dst, err := os.Create(filepath.Join(h.uploadDir, filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	fileURL := fmt.Sprintf("/uploads/%s", filename)
	c.JSON(http.StatusCreated, gin.H{
		"url":      fileURL,
		"filename": filename,
		"size":     header.Size,
	})
}

func (h *UploadHandler) UploadKYBDoc(c *gin.Context) {
	userID, _ := c.Get("user_id")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	defer file.Close()

	if header.Size > h.maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("file too large (max %d bytes)", h.maxSize)})
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	filename := fmt.Sprintf("kyb_%s_%d%s", userID, time.Now().UnixNano(), ext)
	dst, err := os.Create(filepath.Join(h.uploadDir, filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	fileURL := fmt.Sprintf("/uploads/%s", filename)
	c.JSON(http.StatusCreated, gin.H{
		"url":      fileURL,
		"filename": filename,
		"size":     header.Size,
	})
}

func (h *UploadHandler) UploadEventImage(c *gin.Context) {
	eventID := c.Param("id")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	defer file.Close()

	if header.Size > h.maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("file too large (max %d bytes)", h.maxSize)})
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	filename := fmt.Sprintf("event_%s_%d%s", eventID, time.Now().UnixNano(), ext)
	dst, err := os.Create(filepath.Join(h.uploadDir, filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	fileURL := fmt.Sprintf("/uploads/%s", filename)

	// update event with image URL
	h.pool.Exec(context.Background(),
		`UPDATE events SET description = COALESCE(NULLIF(description,''),'') || $2 WHERE id = $1`,
		eventID, fmt.Sprintf(`\n![event image](%s)`, fileURL))

	c.JSON(http.StatusCreated, gin.H{
		"url":      fileURL,
		"filename": filename,
		"size":     header.Size,
	})
}
