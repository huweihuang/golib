package zap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

func NewFromToml(confPath string) *LogOptions {
	var c *LogOptions
	if _, err := toml.DecodeFile(confPath, &c); err != nil {
		panic(err)
	}
	c.defaultDisplay()
	c.SetCaller(true)
	return c
}

func NewFromYaml(confPath string) *LogOptions {
	var c *LogOptions
	file, err := ioutil.ReadFile(confPath)
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	c.defaultDisplay()
	c.SetCaller(true)
	return c
}

func NewFromJson(confPath string) *LogOptions {
	var c *LogOptions
	file, err := ioutil.ReadFile(confPath)
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
	}
	err = json.Unmarshal(file, &c)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	c.defaultDisplay()
	c.SetCaller(true)
	return c
}
