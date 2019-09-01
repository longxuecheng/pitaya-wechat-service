package middle_ware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 是一个应用标准的业务响应模型
type Response struct {
	Data       interface{} `json:"data"`
	ResultDesc string      `json:"resultDesc"`
	ResultCode string      `json:"resultCode"`
}

var okResponse = Response{
	ResultCode: OK,
	ResultDesc: "Success",
}

func newResponse(data interface{}, resultCode string, resultDesc string) Response {
	responseDTO := Response{
		Data:       data,
		ResultCode: resultCode,
		ResultDesc: resultDesc,
	}
	return responseDTO
}

const (
	DATA_KEY      = "data"
	OK            = "Ok"
	UNKNOWN_ERROR = "Unknown"
	UNAUTHORIZED  = "UnAuthorized"
	TOKEN_EXP     = "TokenExpired"
	FORM_INVALID  = "FormInvalid"
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
	response := newResponse(data, OK, "Success")
	c.JSON(http.StatusOK, response)
}

func SetResponseData(c *gin.Context, data interface{}) {
	c.Set(DATA_KEY, data)
}

func SetResponseDataWithStatus(c *gin.Context, data interface{}, code int) {
	c.Status(code)
	c.Set(DATA_KEY, data)
}

func BadRequest(c *gin.Context, reason string) {
	if reason == "" {
		reason = "请求参数有误"
	}
	c.JSON(http.StatusBadRequest, newResponse(nil, reason, FORM_INVALID))
	c.Abort()
}
