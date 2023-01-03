package createuser

import "github.com/harvey1327/chatapplib/models"

type Model struct {
	DisplayName string `json:"displayName" binding:"required" bson:"displayName"`
}

func GetModelConf() models.ModelConfiguration[Model] {
	return new(Model)
}

func (m Model) GetQueueName() string {
	return "user.create"
}
