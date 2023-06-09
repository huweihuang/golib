package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/huweihuang/golib/gin/types"
)

// 封装请求成功的处理逻辑，状态码 200
func SucceedWrapper(c *gin.Context, msg string, data interface{}) {
	resp := types.Response{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("%s succeed", msg),
		Data:    data,
	}
	log.WithField("resp", resp).Info(msg)
	c.JSON(http.StatusOK, resp)
}

// 封装请求失败的处理逻辑，状态码 500
func ErrorWrapper(c *gin.Context, msg string, err error) {
	resp := types.Response{
		Code:    http.StatusInternalServerError,
		Message: fmt.Sprintf("%s failed", msg),
		Data:    map[string]interface{}{"error": err.Error()},
	}
	log.WithField("resp", resp).Error(msg)
	c.JSON(http.StatusInternalServerError, resp)
}

// 封装NotFound的处理逻辑，状态码 404
func NotFoundWrapper(c *gin.Context, msg string, data interface{}) {
	resp := types.Response{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%s not found", msg),
		Data:    data,
	}
	log.WithField("resp", resp).Error(msg)
	c.JSON(http.StatusNotFound, resp)
}

// 封装非法请求的处理逻辑，状态码 400
func BadRequestWrapper(c *gin.Context, err error) {
	resp := types.Response{
		Code:    http.StatusBadRequest,
		Message: "invalid request",
		Data:    map[string]interface{}{"error": err.Error()},
	}
	log.WithField("resp", resp).Error("invalid request")
	c.JSON(http.StatusBadRequest, resp)
}
