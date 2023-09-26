# Usage

## 编译示例

```shell
BUILD_DIR=./_output
BASE_DIR="github.com/huweihuang/gin-api-frame"
VERSION_PACKAGE="github.com/huweihuang/golib/version"

# 构建参数
TARGET_OS=linux
TARGET_ARCH=amd64
VERSION=$(git describe --abbrev=0 --always --tags | sed 's/-/./g')
GIT_COMMIT=$(git rev-parse HEAD)
GIT_TREE_STATE="clean"

GO_LDFLAGS="-X ${VERSION_PACKAGE}.GitVersion=${VERSION} \
	-X ${VERSION_PACKAGE}.GitCommit=${GIT_COMMIT} \
	-X ${VERSION_PACKAGE}.GitTreeState=${GIT_TREE_STATE} \
	-X ${VERSION_PACKAGE}.BuildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ')"

CGO_ENABLED=0 GOOS="$TARGET_OS" GOARCH="$TARGET_ARCH" go build -i -v -ldflags "${GO_LDFLAGS}" \
        -o $BUILD_DIR/bin/"${bin}" ${BASE_DIR}/cmd/server
```

## 代码调用

```go
import (
	"github.com/huweihuang/golib/version/verflag"
)

func main(){
    verflag.PrintAndExitIfRequested()
}
```
