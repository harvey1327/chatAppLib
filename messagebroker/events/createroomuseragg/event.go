package createroomuseragg

const QUEUE_NAME = "roomuseragg.create"

type Model struct {
	RoomID string `json:"roomID" binding:"required" bson:"roomID"`
	UserID string `json:"userID" binding:"required" bson:"userID"`
}
