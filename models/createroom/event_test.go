package createroom_test

import (
	"testing"

	"github.com/harvey1327/chatapplib/models/createroom"
)

func Test_ModelConfHasQueueName(t *testing.T) {
	modelConf := createroom.GetModelConf()

	if modelConf.GetQueueName() != "room.create" {
		t.Errorf("modelConf.GetQueueName() does not equal room.create")
	}
}
