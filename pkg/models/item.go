package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	ID primitive.ObjectID `json:"_id" form:"id" bson:"_id"`
	// OwnerId is the ObjectID of the owner user
	OwnerId string `json:"owner_id" form:"owner_id" binding:"required" bson:"owner_id"`
	// VaultName contains the owner's vault name that this file is located in
	VaultName string `json:"vault_name" form:"vault_name" binding:"required" bson:"vault_name"`
	// ServerPath is the relative to the user vault file path
	ServerPath string `json:"server_path" form:"server_path" binding:"required" bson:"server_path"`
	// ServerMTime contains the last time the file has had a change
	ServerMTime int `json:"server_m_time" form:"server_m_time" binding:"required" bson:"server_m_time"`
	// SHA256 contains the server caluclated SHA256 of the file
	SHA256 string `json:"sha256" form:"sha256" binding:"required" bson:"sha256"`
	// Size contains the bytes size of the file.
	Size int `json:"size" form:"size" binding:"required" bson:"size"`
}
