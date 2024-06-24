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
