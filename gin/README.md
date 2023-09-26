# middleware

middleware封装了gin的body返回逻辑。

统一返回结构体：

```go
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
```

统一返回状态码：

- 请求成功：200
- 内部错误：500
- 错误请求：400
- NotFound: 404

# User Guide

示例代码：

```go
package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/huweihuang/golib/gin/middlewares"

	"myrepo/example-api/pkg/services"
)

type ExampleHandler struct {
	service *services.ExampleService
}

func newExampleHandler() *ExampleHandler {
	return &ExampleHandler{
		service: services.NewExampleService(),
	}
}

func (h *ExampleHandler) ListExample(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		middlewares.BadRequestWrapper(c, fmt.Errorf("name is requierd"))
		return
	}

	result, err := h.service.ListExample(name)
	if err != nil {
		middlewares.ErrorWrapper(c, "ListExample", err)
		return
	}
	middlewares.SucceedWrapper(c, "ListExample", result)
}
```

# Logger middler

```go
gin.Use(
	middlewares.RequestIDMiddleware, 
	middlewares.LogMiddleware(log.Logger),
	gin.RecoveryWithWriter(log.Logger.Out),
	cors.Default(),
)
```
