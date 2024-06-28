package models

type Item struct {
	// ID primitive.ObjectID `json:"_id" form:"id" bson:"_id"`
	// OwnerId is the ObjectID of the owner user
	OwnerId string `json:"owner_id" form:"owner_id" binding:"required"`
	// ServerPath is the relative to the user vault file path
	ServerPath string `json:"server_path" form:"server_path" binding:"required"`
	// ServerMTime contains the last time the file has had a change
	ServerMTime int64 `json:"server_m_time" form:"server_m_time" binding:"required"`
	// SHA256 contains the server caluclated SHA256 of the file
	SHA256 string `json:"sha256" form:"sha256" binding:"required"`
	// Size contains the bytes size of the file.
	Size int `json:"size" form:"size" binding:"required"`
}
