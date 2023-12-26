package services

import (
	"fmt"

	"github.com/Michaelpalacce/gobi/pkg/database"
	"github.com/Michaelpalacce/gobi/pkg/models"
)

// ItemsService handles the logic for CRUD operations on items
type ItemsService struct {
	DB *database.Database
}

func NewItemsService(db *database.Database) *ItemsService {
	return &ItemsService{
		DB: db,
	}
}

// Upsert will insert or update an item in the database
// @TODO: Implement this
func (s *ItemsService) Upsert(item *models.Item) error {
	return fmt.Errorf("not implemented")
}
