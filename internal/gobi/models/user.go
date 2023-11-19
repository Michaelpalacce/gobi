package models

// User model.
// Contains basic information about the user
type User struct {
	Username      string `json:"username" form:"username" binding:"required"`
	Password      string `json:"password" form:"password" binding:"required"`
	EncryptionKey string `json:"encryptionKey" form:"encryptionKey" binding:"required"`
}
