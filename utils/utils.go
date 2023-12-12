package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"

func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func PrintObjectJson(obj interface{}) string {
	objByte, err := json.Marshal(obj)
	if err != nil {
		return fmt.Sprintf("json marshal failed, err: %v", err)
	}
	return fmt.Sprintf("%s", objByte)
}

func IsInList(str string, list []string) bool {
	listMap := make(map[string]bool)
	for _, key := range list {
		listMap[key] = true
	}
	if _, ok := listMap[str]; ok {
		return true
	}
	return false
}

func MakeParentDir(fullFilePath string) error {
	dirPath := filepath.Dir(fullFilePath)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func FormatTime(t time.Time) string {
	local, _ := time.LoadLocation("Asia/Shanghai")
	return t.In(local).Format("2006-01-02 15:04:05")
}
