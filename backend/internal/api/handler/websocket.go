package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dennis-lee/LiveHouseAAS/backend/internal/auth"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WSClient struct {
	Conn   *websocket.Conn
	UserID string
	Send   chan []byte
}

type WSHub struct {
	mu      sync.RWMutex
	clients map[string][]*WSClient
}

var Hub = &WSHub{clients: make(map[string][]*WSClient)}

func NewWebSocketHandler(pool *pgxpool.Pool, jwt *auth.JWT) *WebSocketHandler {
	return &WebSocketHandler{pool: pool, jwt: jwt}
}

type WebSocketHandler struct {
	pool *pgxpool.Pool
	jwt  *auth.JWT
}

func (h *WebSocketHandler) Serve(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	claims, err := h.jwt.ValidateToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("ws upgrade error: %v", err)
		return
	}

	client := &WSClient{
		Conn:   conn,
		UserID: claims.UserID,
		Send:   make(chan []byte, 256),
	}

	Hub.mu.Lock()
	Hub.clients[client.UserID] = append(Hub.clients[client.UserID], client)
	Hub.mu.Unlock()

	go client.writePump()
	client.readPump()
}

func (c *WSClient) writePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

func (c *WSClient) readPump() {
	defer func() {
		c.Conn.Close()
		Hub.mu.Lock()
		clients := Hub.clients[c.UserID]
		for i, cl := range clients {
			if cl == c {
				Hub.clients[c.UserID] = append(clients[:i], clients[i+1:]...)
				break
			}
		}
		Hub.mu.Unlock()
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func SendToUser(userID string, msg interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	Hub.mu.RLock()
	clients := Hub.clients[userID]
	Hub.mu.RUnlock()
	for _, c := range clients {
		select {
		case c.Send <- data:
		default:
		}
	}
}

type NotificationHandler struct {
	pool *pgxpool.Pool
}

func NewNotificationHandler(pool *pgxpool.Pool) *NotificationHandler {
	return &NotificationHandler{pool: pool}
}

func (h *NotificationHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")

	rows, err := h.pool.Query(context.Background(),
		`SELECT id, type, title, COALESCE(body,''), read, created_at
		 FROM notifications WHERE user_id = $1
		 ORDER BY created_at DESC LIMIT 50`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list notifications"})
		return
	}
	defer rows.Close()

	var list []gin.H
	for rows.Next() {
		var id, ntype, title, body string
		var read bool
		var createdAt interface{}
		if err := rows.Scan(&id, &ntype, &title, &body, &read, &createdAt); err != nil {
			continue
		}
		list = append(list, gin.H{
			"id": id, "type": ntype, "title": title, "body": body,
			"read": read, "created_at": createdAt,
		})
	}
	if list == nil {
		list = []gin.H{}
	}
	c.JSON(http.StatusOK, list)
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	userID, _ := c.Get("user_id")
	notifID := c.Param("id")

	result, err := h.pool.Exec(context.Background(),
		`UPDATE notifications SET read = true WHERE id = $1 AND user_id = $2`, notifID, userID)
	if err != nil || result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var count int32
	h.pool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND read = false`, userID).Scan(&count)
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func CreateNotification(ctx context.Context, pool *pgxpool.Pool, userID, ntype, title, body string, data map[string]interface{}) {
	d, _ := json.Marshal(data)
	pool.Exec(ctx,
		`INSERT INTO notifications (id, user_id, type, title, body, data)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)`,
		userID, ntype, title, body, d)
}

func (h *NotificationHandler) ListAll(c *gin.Context) {
	rows, err := h.pool.Query(context.Background(),
		`SELECT n.id, n.user_id, n.type, n.title, COALESCE(n.body,''), n.read, n.created_at, COALESCE(u.name,'')
		 FROM notifications n
		 JOIN users u ON n.user_id = u.id
		 ORDER BY n.created_at DESC LIMIT 100`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list notifications"})
		return
	}
	defer rows.Close()

	var list []gin.H
	for rows.Next() {
		var id, userID, ntype, title, body, userName string
		var read bool
		var createdAt interface{}
		if err := rows.Scan(&id, &userID, &ntype, &title, &body, &read, &createdAt, &userName); err != nil {
			continue
		}
		list = append(list, gin.H{
			"id": id, "user_id": userID, "type": ntype, "title": title,
			"body": body, "read": read, "user_name": userName,
			"created_at": createdAt,
		})
	}
	if list == nil {
		list = []gin.H{}
	}
	c.JSON(http.StatusOK, list)
}
