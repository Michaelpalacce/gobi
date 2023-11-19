package main

import (
	"log"
	"net/http"

	"github.com/Michaelpalacce/gobi/pkg/gobi/database"
	"github.com/gin-gonic/gin"
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

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
