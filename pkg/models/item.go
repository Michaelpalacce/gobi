package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	ID          primitive.ObjectID `json:"_id" form:"id" bson:"_id"`
	OwnerId     string             `json:"owner_id" form:"owner_id" binding:"required" bson:"owner_id"`
	ServerPath  string             `json:"server_path" form:"server_path" binding:"required" bson:"server_path"`
	ServerMTime int                `json:"server_m_time" form:"server_m_time" binding:"required" bson:"server_m_time"`
}
