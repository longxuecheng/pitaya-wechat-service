package wechat

import (
	"gotrue/service/tencloud"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
)

func TestGetParameterizedMPCode(t *testing.T) {
	tm := &TokenManager{}
	tm.at = "28_DOWO4KcusZtqLqd6_XGUAp8B2jHdhjlmNrjsXvE_yzWuzTyaGhv23iJ54K7gU31ZQD-xTRoQGx9GfDAbCYuPIV8m9_iVlESiG8axBXKXNNbghoCMqSNefxfBDdmg4uvszu-9deQboHzxxnT0EMWhAIAMUE"
	ws := wechatService{
		TokenManager: tm,
	}
	dataBytes, err := ws.GetParameterizedMPCode("id=1", "pages/goods/goods", 430, true)
	assert.NoError(t, err)
	cosService := tencloud.NewCosService()
	name := uuid.New().String() + ".jpeg"
	err = cosService.PushData(name, dataBytes)
	assert.NoError(t, err)
}
