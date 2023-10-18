package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// Scan for scanner helper
func Scan(data interface{}, value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, data)
}

// Value for valuer helper
func Value(data interface{}) (interface{}, error) {
	vi := reflect.ValueOf(data)
	// 判断是否为 0 值
	if vi.IsZero() {
		return nil, nil
	}
	return json.Marshal(data)
}

type StringMap map[string]string

// Scan 实现 sql.Scanner 接口，Scan 将字符串变成结构体
func (obj *StringMap) Scan(value interface{}) error {
	return Scan(&obj, value)
}

// Value 实现 driver.Valuer 接口，Value 将结构体变成字符串
func (obj StringMap) Value() (driver.Value, error) {
	return Value(obj)
}

type StringArray []string

func (obj *StringArray) Scan(value interface{}) error {
	return Scan(&obj, value)
}

func (obj StringArray) Value() (driver.Value, error) {
	return Value(obj)
}
