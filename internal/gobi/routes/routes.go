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
	itemsHandler handlers.ItemsHandler,
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

	// V1
	v1Routes := r.Group("/v1/")
	v1Routes.Use(authMiddleware)
	// Item Routes are not in use currently in favor of websocket communication. Leaving this for now
	itemRoutes := v1Routes.Group("/item")
	{
		itemRoutes.GET("/", itemsHandler.GetItem)
	}

	return r
}
