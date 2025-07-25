package routes

import (
	"backend/handlers"
	"backend/middleware"
	"backend/services"
	"backend/websocket"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	authService *services.AuthService,
	wgService *services.WireGuardService,
	hub *websocket.Hub,
) {
	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	wgHandler := handlers.NewWireGuardHandler(wgService)

	// Public routes
	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
		}
	}

	// WebSocket endpoint (requires authentication)
	router.GET("/ws", middleware.AuthMiddleware(authService), hub.HandleWebSocket)

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		// Auth protected routes
		auth := protected.Group("/auth")
		{
			auth.GET("/profile", authHandler.GetProfile)
			auth.POST("/change-password", authHandler.ChangePassword)
		}

		// WireGuard routes
		wg := protected.Group("/wireguard")
		{
			// Interface routes
			interfaces := wg.Group("/interfaces")
			{
				interfaces.GET("", wgHandler.GetInterfaces)
				interfaces.POST("", wgHandler.CreateInterface)
				interfaces.GET("/:id", wgHandler.GetInterface)
				interfaces.PUT("/:id", wgHandler.UpdateInterface)
				interfaces.DELETE("/:id", wgHandler.DeleteInterface)
				interfaces.POST("/:id/start", wgHandler.StartInterface)
				interfaces.POST("/:id/stop", wgHandler.StopInterface)
				interfaces.GET("/:id/config", wgHandler.GetInterfaceConfig)
			}

			// Peer routes
			peers := wg.Group("/peers")
			{
				peers.GET("", wgHandler.GetPeers)
				peers.POST("", wgHandler.CreatePeer)
				peers.GET("/:id", wgHandler.GetPeer)
				peers.PUT("/:id", wgHandler.UpdatePeer)
				peers.DELETE("/:id", wgHandler.DeletePeer)
				peers.GET("/:id/config", wgHandler.GetPeerConfig)
			}
		}
	}
}
