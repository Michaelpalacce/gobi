package routes

import (
	"github.com/Michaelpalacce/gobi/internal/gobi/handlers"
	"github.com/Michaelpalacce/gobi/internal/gobi/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the application routes.
func SetupRouter(
	userHandler handlers.UsersHandler,
	itemsHandler handlers.ItemsHandler,
	websocketHandler handlers.WebsocketHandler,
) *gin.Engine {
	r := gin.Default()

	authMiddleware := middleware.Auth(userHandler.Service)

	// User routes
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/", userHandler.CreateUser)
		userRoutes.DELETE("/", authMiddleware, userHandler.DeleteUser)
	}

	websocketRoutes := r.Group("/ws")
	websocketRoutes.Use(authMiddleware)
	{
		websocketRoutes.GET("", websocketHandler.Establish)
	}

	// Items routes
	itemRoutes := r.Group("/items")
	itemRoutes.Use(authMiddleware)
	{
		itemRoutes.POST("/", itemsHandler.AddItem)
	}

	return r
}
