package wechat

import (
	"fmt"
	"gotrue/service/tencloud"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
)

func TestGetParameterizedMPCode(t *testing.T) {
	tm := &TokenManager{}
	tm.at = "27_-AIkrZebUuXVaRfkDX1ml8SXVs_UhvFOD0qE5gbJv1UkFsXwxjGQ1QBTFHGiHuIojA8NOGgYn4uo0rsvA2-OnBcVbKiK9CPV_nglM4Vyl52I8vurPaWnizVziv029Tn68cX5Q5cK1fHw6PpyFJMeABAVZC"
	ws := wechatService{
		TokenManager: tm,
	}
	dataBytes, err := ws.GetParameterizedMPCode()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	cosService := tencloud.NewCosService()
	name := uuid.New().String() + ".jpeg"
	err = cosService.PushData(name, dataBytes)
	assert.NoError(t, err)
}
