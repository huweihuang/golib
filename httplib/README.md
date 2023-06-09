# httplib

`httplib`封装了`net/http`包的使用逻辑，便于快速调用http的接口。

# User Guide

调用`httplib.CallURL`方法，参数如下：

**入参**：

- `method` string
- `url` string
- `path` string
- `header` map[string]string
- `request` interface{}
- `response` interface{}

**出参**：

- `statusCode` int
- `body` []byte
- `err` error

示例代码：

```go
import(
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
		return nil, fmt.Errorf("failed to get token by edge api, statusCode :%d, err: %v", statusCode, err)
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("get token request error, %s", body)
	}

	data = (&response).Data
	return data, nil
}
```

