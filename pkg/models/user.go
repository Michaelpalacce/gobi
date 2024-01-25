package models

import (
	"unicode"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User model.
// Contains basic information about the user
type User struct {
	ID       primitive.ObjectID `json:"_id" form:"id" bson:"_id"`
	Username string             `json:"username" form:"username" binding:"required" bson:"username"`
	Password string             `json:"password" form:"password" binding:"required" bson:"password"`
}

// ValidateUsername checks if the username is valid. A valid username contains only
// alphanumeric characters and underscores.
func ValidateUsername(username string) bool {
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' && char != '-' && char != '@' {
			return false
		}
	}

	return true
}
