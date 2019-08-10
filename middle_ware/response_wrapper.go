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
	DATA_KEY                      = "data"
	STATUS_SUCCESS      ErrorCode = 0
	STATUS_ERROR        ErrorCode = 999
	STATUS_UNAUTHORIZED ErrorCode = 401
	STATUS_TOKEN_EXP    ErrorCode = 401
	FORM_INVALID        ErrorCode = 400
)

// WrapResponse 在请求controller的handler完成之后获取上下文数据进行统一的业务封装
// 因此应该将此中间件放在请求到达主handler函数之前的最后一个
func WrapResponse(c *gin.Context) {
	c.Next()
	// 若跳过handler则不要再对数据进行加工,因为其他的handler在调用abort时应该对响应数据以及状态提前处理完毕了，
	// 如果不判断的话会导致数据重复
	if c.IsAborted() {
		return
	}
	data, _ := c.Get(DATA_KEY)
	responseDTO := newResponse(data, "Ok", STATUS_SUCCESS)
	c.JSON(http.StatusOK, responseDTO)
}

func newResponse(data interface{}, errorMsg string, errorCode ErrorCode) Response {
	responseDTO := Response{
		Data:     data,
		ErrorMsg: errorMsg,
		ErrorNo:  errorCode,
	}
	return responseDTO
}

func SetResponseData(c *gin.Context, data interface{}) {
	c.Set(DATA_KEY, data)
}

func SetResponseDataWithStatus(c *gin.Context, data interface{}, code int) {
	c.Status(code)
	c.Set(DATA_KEY, data)
}

func BadRequet(c *gin.Context, reason string) {
	c.JSON(http.StatusBadRequest, newResponse(nil, reason, FORM_INVALID))
	c.Abort()
}
