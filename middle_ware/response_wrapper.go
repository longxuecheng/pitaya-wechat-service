package middle_ware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 是一个应用标准的业务响应模型
type Response struct {
	Data     interface{} `json:"data"`
	ErrorMsg string      `json:"errorMsg"`
	ErrorNo  ErrorCode   `json:"errno"`
}

type ErrorCode int

const (
	DATA_KEY                 = "data"
	STATUS_SUCCESS ErrorCode = 0
	STATUS_ERROR   ErrorCode = 999
)

// WrapResponse 在请求controller的handler完成之后获取上下文数据进行统一的业务封装
func WrapResponse(c *gin.Context) {
	c.Next()
	data, _ := c.Get(DATA_KEY)
	responseDTO := newResponseDTO(data, "", STATUS_SUCCESS)
	c.JSON(http.StatusOK, responseDTO)
}

func newResponseDTO(data interface{}, errorMsg string, errorCode ErrorCode) Response {
	responseDTO := Response{
		Data:     data,
		ErrorMsg: errorMsg,
		ErrorNo:  errorCode,
	}
	return responseDTO
}
