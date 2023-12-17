# golib

Golib is a library that encapsulates the go toolkit, making it easy for rapid development.

# logger

Logger encapsulates the two most popular logging libraries, `Zap` and `Logrus`. 
By configuring the following commonly used logging parameters, you can quickly customize the logging function. 
The log parameters include `LogFile`, `LogLevel`, `LogFormat` etc.

## zap

```go
package main

import (
	"time"
	log "github.com/huweihuang/golib/logger/zap"
)

func main() {
	log.InitLogger("log/info.log", "log/error.log", "debug", "text", true)

	log.Logger().Infof("test default log")
	log.Logger().Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", "example.com",
		"attempt", 3,
		"backoff", time.Second)

	log.Logger().Error("test error log")
	log.Logger().Errorw("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", "example.com",
		"attempt", 3,
		"backoff", time.Second)

	log.Logger().With("with_field", map[string]string{"test1": "value1", "test2": "value2"}).Info("test with field")
	log.Logger().With("field1", "value1", "field2", "value2", "field3", "value3").Error("test multi field")
}
```

## logrus

```go
package main

import (
	"github.com/sirupsen/logrus"

	log "github.com/huweihuang/golib/logger/logrus"
)

func main() {
	// init log
	log.InitLogger("./log/logrus.log", "debug", "text", false)

	// Printf
	log.Logger().Debugf("test debugf, %s", "debugf")
	log.Logger().Infof("test infof, %s", "infof")
	log.Logger().Warnf("test warnf, %s", "warnf")
	log.Logger().Errorf("test errorf, %s", "errorf")

	// WithField
	log.Logger().WithField("field1", "debug").Debug("test field, debug")
	log.Logger().WithField("field1", "info").Info("test field, info")
	log.Logger().WithField("field1", "warn").Warn("test field, warn")
	log.Logger().WithField("field1", "error").Error("test field, error")

	// WithFields
	log.Logger().WithFields(logrus.Fields{
		"fields1": "fields1_value",
		"fields2": "fields2_value",
	}).Debug("test fields, debug")

	log.Logger().WithFields(logrus.Fields{
		"fields1": "fields1_value",
		"fields2": "fields2_value",
	}).Info("test fields, info")

	log.Logger().WithFields(logrus.Fields{
		"fields1": "fields1_value",
		"fields2": "fields2_value",
	}).Warn("test fields, warn")

	log.Logger().WithFields(logrus.Fields{
		"fields1": "fields1_value",
		"fields2": "fields2_value",
	}).Error("test fields, error")
}
```

# httplib

The `httplib` encapsulates the usage logic of the `net/http` package, making it easy to quickly call the HTTP interface. 
The input parameters are as follows:

- `method` string
- `url` string
- `path` string
- `header` map[string]string
- `request` interface{}
- `response` interface{}

```go
import(
    "fmt"
    "github.com/huweihuang/golib/httplib"
)

const(
    endpoint = "http://api.example.com"
)

func GetExample(email, role string) (data map[string]string, err error) {
	var response struct {
		Code    int
		Message string
		Data    map[string]string
	}

	path := fmt.Sprintf("/api/v1/token?email=%s&role=%s", email, role)
	statusCode, body, err := httplib.CallURL("GET", endpoint, path, nil, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get token by example api, statusCode :%d, err: %v", statusCode, err)
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("get token request error, %s", body)
	}

	data = (&response).Data
	return data, nil
}
```

# gin

## middleware

gin middleware encapsulates the body return logic of gin.

Unified return structure:

```go
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
```

Unified return status code:

- Request successful: 200
- Internal error: 500
- Invalid request: 400
- NotFound: 404

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

## Logger middler

```go
gin.Use(
	middlewares.RequestIDMiddleware, 
	middlewares.LogMiddleware(log.Logger),
	gin.RecoveryWithWriter(log.Logger.Out),
	cors.Default(),
)
```

# db

The DB library encapsulates the library building and mock operations of Gorm, 
and encapsulates the commonly used `map [string] string` and `[] string` types that support Gorm's format.

```go
import (
    "github.com/huweihuang/golib/db"
    "gorm.io/gorm"
)
// DB is a global db manager
var DB *DBMng

type DBMng struct {
	db *gorm.DB
}

func NewDBMng(db *gorm.DB) *DBMng {
	return &DBMng{
		db: db,
	}
}

func InitDB(dbConf *configs.DBConfig) (*DBMng, error) {
	d, err := db.SetupDB(dbConf.Addr, dbConf.DBName, dbConf.User, dbConf.Password, dbConf.LogLevel)
	if err != nil {
		return nil, err
	}
	DB = NewDBMng(d)
	return DB, nil
}

func Close() error {
	db, err := DB.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
```

# config

`config` encapsulates the use of `viper` and parses the configuration file of the specified path into a structure.

```go
import(
	"github.com/huweihuang/golib/config"
)

func main() {
	err := config.InitConfigObjectByPath(configFile, &configs.GlobalConfig)
	if err != nil {
		panic(err)
	}
}	
```

# kube

The `kube` library encapsulates commonly used k8s operations, such as building clientsets.

```go
	kubeClient, err := kube.NewKubeClient(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize kube client, err: %w", err)
	}
```
