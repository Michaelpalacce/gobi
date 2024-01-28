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
	itemHandler handlers.ItemHandler,
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

	// Items Routes
	itemsRoutes := r.Group("/items")
	{
		itemsRoutes.GET("/", itemHandler.GetItem)
		itemsRoutes.POST("/", itemHandler.CreateItem)
		itemsRoutes.DELETE("/", itemHandler.DeleteItem)
	}

	return r
}
