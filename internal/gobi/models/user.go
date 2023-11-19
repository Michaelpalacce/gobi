package models

// User model.
// Contains basic information about the user
type User struct {
	InsertedID    interface{} `json:"id" form:"id" bson:"id"`
	Username      string      `json:"username" form:"username" binding:"required" bson:"username"`
	Password      string      `json:"password" form:"password" binding:"required" bson:"password"`
	EncryptionKey string      `json:"encryptionKey" form:"encryptionKey" binding:"required" bson:"encryptionKey"`
}
