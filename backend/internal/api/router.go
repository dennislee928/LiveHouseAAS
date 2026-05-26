package api

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/dennis-lee/LiveHouseAAS/backend/internal/api/handler"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/api/middleware"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/auth"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/blockchain"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/config"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/infra/cache"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/infra/db"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/payment"
)

func NewRouter(cfg *config.Config, pg *db.Postgres, r *cache.Redis) *gin.Engine {
	gin.SetMode(cfg.GinMode)
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimit(100, time.Minute))

	router.Static("/uploads", cfg.UploadDir)

	jwt := auth.NewJWT(cfg.JWTSecret)
	authH := handler.NewAuthHandler(pg.Pool, jwt)
	venueH := handler.NewVenueHandler(pg.Pool)
	slotH := handler.NewSlotHandler(pg.Pool)
	bookingH := handler.NewBookingHandler(pg.Pool)
	eventH := handler.NewEventHandler(pg.Pool)
	kybH := handler.NewKYBHandler(pg.Pool)
	dashH := handler.NewDashboardHandler(pg.Pool)
	verifyH := handler.NewVerifyHandler(pg.Pool)
	payRouter := payment.NewRouter()
	ticketH := handler.NewTicketHandler(pg.Pool, payRouter)
	nftSvc := blockchain.NewMockService()
	nftH := handler.NewNFTHandler(pg.Pool, nftSvc)
	adminH := handler.NewAdminHandler(pg.Pool)
	wsH := handler.NewWebSocketHandler(pg.Pool, jwt)
	uploadH := handler.NewUploadHandler(pg.Pool, cfg.UploadDir, cfg.MaxUploadSize)
	seatH := handler.NewSeatMapHandler(pg.Pool)
	notifH := handler.NewNotificationHandler(pg.Pool)
	analyticsH := handler.NewAnalyticsHandler(pg.Pool)
	searchH := handler.NewSearchHandler(pg.Pool)
	userFeaturesH := handler.NewUserFeaturesHandler(pg.Pool)
	callbackH := handler.NewCallbackHandler(pg.Pool, payRouter)

	router.GET("/health", handler.HealthCheck)
	router.GET("/api/v1/ws", wsH.Serve)

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authH.Register)
			auth.POST("/login", authH.Login)
			auth.POST("/forgot-password", userFeaturesH.ForgotPassword)
			auth.POST("/reset-password", userFeaturesH.ResetPassword)
		}

		// Public search endpoints (rate limited but no auth)
		v1.GET("/search/events", searchH.Events)
		v1.GET("/search/venues", searchH.Venues)
		v1.GET("/verify-email", userFeaturesH.VerifyEmail)

		// Payment callbacks (no auth, provider calls these)
		v1.POST("/payments/callback", callbackH.PaymentCallback)
		v1.POST("/payments/ecpay/notify", callbackH.ECPayNotify)
		v1.POST("/payments/newebpay/notify", callbackH.NewebPayNotify)

		protected := v1.Group("")
		protected.Use(middleware.Auth(jwt))
		{
			protected.GET("/me", authH.GetMe)
			protected.GET("/dashboard/stats", dashH.Stats)

			// --- User Features ---
			protected.PUT("/me/profile", userFeaturesH.UpdateProfile)
			protected.POST("/me/change-password", userFeaturesH.ChangePassword)
			protected.POST("/me/avatar", userFeaturesH.UpdateAvatar)
			protected.POST("/me/verify-email", userFeaturesH.RequestEmailVerification)

			// --- Notifications ---
			protected.GET("/notifications", notifH.List)
			protected.GET("/notifications/unread", notifH.UnreadCount)
			protected.PUT("/notifications/:id/read", notifH.MarkRead)

			// --- File Upload ---
			protected.POST("/upload", uploadH.Upload)
			protected.POST("/upload/kyb", uploadH.UploadKYBDoc)
			protected.POST("/events/:id/upload", uploadH.UploadEventImage)

			// --- Venues ---
			venues := protected.Group("/venues")
			{
				venues.GET("", venueH.List)
				venues.GET("/all", venueH.ListAll)
				venues.POST("", venueH.Create)
				venues.GET("/:id", venueH.Get)
				venues.PUT("/:id", venueH.Update)
				venues.DELETE("/:id", venueH.Delete)

				venues.GET("/:id/specs", venueH.ListSpecs)
				venues.POST("/:id/specs", venueH.CreateSpec)
				venues.DELETE("/:id/specs/:specId", venueH.DeleteSpec)

				venues.GET("/:venueId/slots", slotH.List)
				venues.POST("/:venueId/slots", slotH.Create)
				venues.POST("/:venueId/slots/batch", slotH.BatchCreate)
				venues.DELETE("/slots/:id", slotH.Delete)

				venues.GET("/:venueId/events", eventH.ListByVenue)

				venues.GET("/:venueId/seats", seatH.GetLayout)
				venues.PUT("/:venueId/seats", seatH.SaveLayout)
			}

			// --- Bookings ---
			protected.GET("/bookings/artist", bookingH.ListByArtist)
			protected.GET("/bookings/owner", bookingH.ListByOwner)
			protected.GET("/venues/:venueId/bookings", bookingH.ListByVenue)
			protected.POST("/bookings", bookingH.Create)
			protected.PUT("/bookings/:id/status", bookingH.UpdateStatus)

			// --- Events ---
			events := protected.Group("/events")
			{
				events.GET("/published", eventH.ListPublished)
				events.GET("/artist", eventH.ListByArtist)
				events.GET("/:id", eventH.Get)
				events.PUT("/:id", eventH.Update)
				events.POST("/:id/publish", eventH.Publish)

				events.GET("/:id/ticket-types", eventH.ListTicketTypes)
				events.POST("/:id/ticket-types", eventH.CreateTicketType)

				events.GET("/:eventId/seats", seatH.GetSeatAvailability)

				events.POST("/:id/purchase", ticketH.Purchase)
			}

			// --- Orders & Tickets ---
			protected.GET("/orders", ticketH.ListOrders)
			protected.POST("/orders/:orderId/refund", callbackH.Refund)
			protected.GET("/tickets", ticketH.ListTickets)

			// --- Ticket Verification ---
			protected.POST("/tickets/verify", verifyH.Verify)
			protected.GET("/tickets/lookup", verifyH.Lookup)

			// --- NFT ---
			protected.POST("/tickets/:ticketId/nft/claim", nftH.Claim)
			protected.POST("/tickets/:ticketId/nft/poap", nftH.ClaimPOAP)
			protected.GET("/tickets/:ticketId/nft", nftH.GetStatus)

			// --- KYB ---
			protected.POST("/kyb", kybH.Submit)
			protected.GET("/kyb", kybH.GetStatus)

			// --- Analytics ---
			protected.GET("/analytics/summary", analyticsH.Summary)
			protected.GET("/analytics/revenue", analyticsH.RevenueOverTime)
			protected.GET("/analytics/top-venues", analyticsH.TopVenues)
			protected.GET("/analytics/top-events", analyticsH.TopEvents)
			protected.GET("/analytics/booking-trends", analyticsH.BookingTrends)
			protected.GET("/analytics/venue-performance", analyticsH.VenuePerformance)

			// --- Admin ---
			admin := protected.Group("/admin")
			admin.Use(middleware.RoleCheck("admin"))
			{
				admin.GET("/stats", adminH.Stats)
				admin.GET("/users", adminH.ListUsers)
				admin.PUT("/users/:id/role", adminH.UpdateUserRole)
				admin.GET("/venues", adminH.ListVenues)
				admin.PUT("/venues/:id/status", adminH.UpdateVenueStatus)
				admin.GET("/events", adminH.ListEvents)
				admin.GET("/bookings", adminH.ListBookings)
				admin.GET("/orders", adminH.ListOrders)
				admin.GET("/kyb/pending", kybH.ListPending)
				admin.PUT("/kyb/:id/review", kybH.Review)
				admin.GET("/notifications", notifH.ListAll)
			}
		}
	}

	return router
}
