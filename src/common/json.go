package common

import (
	"encoding/json"
	"os"
)

// 读取Json文件
// file:json文件的完整路径
// conf:config结构实例
func ReadJson(file string, conf interface{}) error {
	r, err := os.Open(file)
	if err != nil {
		println(err)
		return err
	}
	decoder := json.NewDecoder(r)
	err = decoder.Decode(conf)
	return err
}
