package routes

import (
	"github.com/Michaelpalacce/gobi/internal/gobi/handlers"
	"github.com/Michaelpalacce/gobi/internal/gobi/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the application routes.
func SetupRouter(
	userHandler handlers.UsersHandler,
	websocketHandler handlers.WebsocketHandler,
) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	authMiddleware := middleware.Auth(userHandler.Service)

	// User Routes
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/", userHandler.CreateUser)
		userRoutes.DELETE("/", authMiddleware, userHandler.DeleteUser)
	}

	// Websocket Routes
	websocketRoutes := r.Group("/ws")
	websocketRoutes.Use(authMiddleware)
	{
		websocketRoutes.GET("/", websocketHandler.Establish)
	}

	return r
}
