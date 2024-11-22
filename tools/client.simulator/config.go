package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var configJsonPath string                  // 配置文件 路径
var configJson *ConfigJson                 // 配置文件
var ignoreMsgIDMap = map[uint32]struct{}{} // 忽略的消息ID

type ConfigJson struct {
	Addr string `json:"addr"`
}

// 解析配置文件
func parseConfigJson(path string) error {
	// Open the config file
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return err
	}
	defer file.Close()
	// Read the file content
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return err
	}
	configJson = &ConfigJson{}
	// Parse the JSON content
	if err := json.Unmarshal(bytes, configJson); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return err
	}
	// Print the parsed config
	fmt.Printf("Parsed Config: %+v\n", configJson)
	return nil
}
