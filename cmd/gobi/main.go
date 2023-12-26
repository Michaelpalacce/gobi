package main

import (
	"log"
	"log/slog"

	"github.com/Michaelpalacce/gobi/internal/gobi/handlers"
	"github.com/Michaelpalacce/gobi/internal/gobi/routes"
	internalServices "github.com/Michaelpalacce/gobi/internal/gobi/services"
	"github.com/Michaelpalacce/gobi/pkg/database"
	"github.com/Michaelpalacce/gobi/pkg/gobi/services"
	"github.com/Michaelpalacce/gobi/pkg/logger"
)

func main() {
	logger.ConfigureLogging()

	var (
		db  *database.Database
		err error
	)
	slog.Info("Connecting to Database")

	if db, err = database.NewDatabase(); err != nil {
		log.Fatalf("Error while extablishing connectino to the database: %s", err)
	}

	defer db.Disconnect()

	usersHandler := *handlers.NewUsersHandler(
		services.NewUsersService(db),
	)

	websocketHandler := *handlers.NewWebsocketHandler(
		internalServices.NewWebsocketService(db),
	)

	itemsHandler := *handlers.NewItemsHandler(
		services.NewItemsService(db),
	)

	r := routes.SetupRouter(
		usersHandler,
		websocketHandler,
		itemsHandler,
	)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
