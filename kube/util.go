package kube

import (
	"fmt"
	"sort"
	"strings"
)

func ConvertMapToStr(toConvertMap map[string]string) string {
	str := make([]string, len(toConvertMap))
	i := 0
	for name, value := range toConvertMap {
		str[i] = fmt.Sprintf("%s=%s", name, value)
		i++
	}
	sort.Strings(str)
	return strings.Join(str, ",")
}
