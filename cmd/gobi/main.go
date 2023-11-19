package main

import (
	"log"

	"github.com/Michaelpalacce/gobi/internal/gobi/handlers"
	"github.com/Michaelpalacce/gobi/internal/gobi/routes"
	"github.com/Michaelpalacce/gobi/internal/gobi/services"
	"github.com/Michaelpalacce/gobi/pkg/gobi/database"
)

func main() {
	var (
		db  *database.Database
		err error
	)

	if db, err = database.NewDatabase(); err != nil {
		log.Fatalf("Error while creating a new Database: %s", err)
	}

	defer db.Disconnect()

	usersHandler := *handlers.NewUsersHandler(
		services.NewUserService(db),
	)

	r := routes.SetupRouter(
		usersHandler,
	)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
