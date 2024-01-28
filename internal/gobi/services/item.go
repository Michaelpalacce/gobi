package services

type ItemService struct{}

// NewItemService will instantiate a new ItemService given the database
func NewItemService() *ItemService {
	return &ItemService{}
}
