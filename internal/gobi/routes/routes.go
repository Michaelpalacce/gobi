package routes

import (
	"github.com/Michaelpalacce/gobi/internal/gobi/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the application routes.
func SetupRouter(
	userHandler handlers.UsersHandler,
	itemsHandler handlers.ItemsHandler,
) *gin.Engine {
	r := gin.Default()

	// User routes
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/", userHandler.CreateUser)
        userRoutes.DELETE("/:id", userHandler.DeleteUser)
        userRoutes.GET("/:id", userHandler.GetUser)
	}

	// Items routes
	itemRoutes := r.Group("/items")
	{
		itemRoutes.POST("/", itemsHandler.AddItem)
	}

	return r
}
