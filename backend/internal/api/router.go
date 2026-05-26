package api

import (
	"github.com/gin-gonic/gin"

	"github.com/dennis-lee/LiveHouseAAS/backend/internal/api/handler"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/api/middleware"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/auth"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/config"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/infra/cache"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/infra/db"
)

func NewRouter(cfg *config.Config, pg *db.Postgres, r *cache.Redis) *gin.Engine {
	gin.SetMode(cfg.GinMode)
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	jwt := auth.NewJWT(cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(pg.Pool, jwt)

	router.GET("/health", handler.HealthCheck)

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		protected := v1.Group("")
		protected.Use(middleware.Auth(jwt))
		{
			protected.GET("/me", authHandler.GetMe)
		}
	}

	return router
}
