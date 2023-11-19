package services

import (
	"github.com/Michaelpalacce/gobi/pkg/gobi/database"
)

type ItemsService struct {
	DB *database.Database
}

// NewItemsService will return an instance of the User Service
func NewItemsService(db *database.Database) *ItemsService {
	return &ItemsService{
		DB: db,
	}
}
