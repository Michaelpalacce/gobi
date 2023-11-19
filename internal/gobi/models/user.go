package models

// User model.
// Contains basic information about the user
type User struct {
	ID            interface{} `json:"_id" form:"id" bson:"_id,omitempty"`
	Username      string      `json:"username" form:"username" binding:"required" bson:"username"`
	Password      string      `json:"password" form:"password" binding:"required" bson:"password"`
	EncryptionKey string      `json:"encryptionKey" form:"encryptionKey" binding:"required" bson:"encryptionKey"`
}
