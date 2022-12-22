package createroom

const QUEUE_NAME = "room.create"

type Model struct {
	DisplayName string `json:"displayName" binding:"required" bson:"displayName"`
}
