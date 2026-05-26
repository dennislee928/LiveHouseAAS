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
	authH := handler.NewAuthHandler(pg.Pool, jwt)
	venueH := handler.NewVenueHandler(pg.Pool)
	slotH := handler.NewSlotHandler(pg.Pool)
	bookingH := handler.NewBookingHandler(pg.Pool)

	router.GET("/health", handler.HealthCheck)

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authH.Register)
			auth.POST("/login", authH.Login)
		}

		protected := v1.Group("")
		protected.Use(middleware.Auth(jwt))
		{
			protected.GET("/me", authH.GetMe)

			// --- Venues ---
			venues := protected.Group("/venues")
			{
				venues.GET("", venueH.List)
				venues.GET("/all", venueH.ListAll)
				venues.POST("", venueH.Create)
				venues.GET("/:id", venueH.Get)
				venues.PUT("/:id", venueH.Update)
				venues.DELETE("/:id", venueH.Delete)

				// Venue Specs
				venues.GET("/:id/specs", venueH.ListSpecs)
				venues.POST("/:id/specs", venueH.CreateSpec)
				venues.DELETE("/:id/specs/:specId", venueH.DeleteSpec)

				// Slots
				venues.GET("/:venueId/slots", slotH.List)
				venues.POST("/:venueId/slots", slotH.Create)
				venues.POST("/:venueId/slots/batch", slotH.BatchCreate)
				venues.DELETE("/slots/:id", slotH.Delete)
			}

			// --- Bookings ---
			protected.GET("/bookings/artist", bookingH.ListByArtist)
			protected.GET("/venues/:venueId/bookings", bookingH.ListByVenue)
			protected.POST("/bookings", bookingH.Create)
			protected.PUT("/bookings/:id/status", bookingH.UpdateStatus)
		}
	}

	return router
}
