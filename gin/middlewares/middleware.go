package middlerwares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/huweihuang/golib/logger/logrus"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/huweihuang/golib/gin/types"
)

// SucceedWrapper 封装请求成功的处理逻辑，状态码 200
func SucceedWrapper(c *gin.Context, msg string, data interface{}) {
	resp := types.Response{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("%s succeed", msg),
		Data:    data,
	}
	log.Log().WithField("resp", resp).Info(msg)
	c.JSON(http.StatusOK, resp)
}

// ErrorWrapper 封装请求失败的处理逻辑，状态码 500
func ErrorWrapper(c *gin.Context, msg string, err error) {
	resp := types.Response{
		Code:    http.StatusInternalServerError,
		Message: fmt.Sprintf("%s failed", msg),
		Data:    map[string]interface{}{"error": err.Error()},
	}
	log.Log().WithField("resp", resp).Error(msg)
	c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
}

// NotFoundWrapper 封装NotFound的处理逻辑，状态码 404
func NotFoundWrapper(c *gin.Context, msg string, data interface{}) {
	resp := types.Response{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%s not found", msg),
		Data:    data,
	}
	log.Log().WithField("resp", resp).Error(msg)
	c.AbortWithStatusJSON(http.StatusNotFound, resp)
}

// BadRequestWrapper 封装非法请求的处理逻辑，状态码 400
func BadRequestWrapper(c *gin.Context, err error) {
	resp := types.Response{
		Code:    http.StatusBadRequest,
		Message: "invalid request",
		Data:    map[string]interface{}{"error": err.Error()},
	}
	log.Log().WithField("resp", resp).Error("invalid request")
	c.AbortWithStatusJSON(http.StatusBadRequest, resp)
}

// ValidateBadRequestWrapper 封装多项校验非法请求的处理逻辑，状态码 400
func ValidateBadRequestWrapper(c *gin.Context, errs field.ErrorList) {
	resp := types.Response{
		Code:    http.StatusBadRequest,
		Message: "invalid request",
		Data:    map[string]interface{}{"error": errs},
	}
	log.Log().WithField("resp", resp).Error("invalid request")
	c.AbortWithStatusJSON(http.StatusBadRequest, resp)
}

func ParseRequest(c *gin.Context, request interface{}) {
	if err := c.BindJSON(request); err != nil {
		resp := types.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
			Data:    map[string]interface{}{"error": err},
		}
		log.Log().WithField("err", err).Warn("invalid request body")
		c.AbortWithStatusJSON(http.StatusBadRequest, resp)
	}
}
