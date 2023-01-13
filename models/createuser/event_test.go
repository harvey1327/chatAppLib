package createuser_test

import (
	"testing"

	"github.com/harvey1327/chatapplib/models/createuser"
)

func Test_ModelConfHasQueueName(t *testing.T) {
	modelConf := createuser.GetModelConf()

	if modelConf.GetQueueName() != "user.create" {
		t.Errorf("modelConf.GetQueueName() does not equal user.create")
	}
}
