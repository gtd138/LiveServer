// 通用的函数
package common

import (
	"bytes"
	"encoding/gob"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// 获取服务器时间
// @return:返回自1970年1月1日起，到现在的毫秒
func GetTime() int64 {
	// 纳秒级别
	t1 := time.Now().UnixNano()
	// 转成毫秒级别
	t2 := float64(t1) / math.Pow(float64(10), float64(6))
	return int64(t2)
}

// Gob编码
func GobEncode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Gob解码
func GobDecode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}

// 查找两个数组的交集、差集
// @note:uint32版本，其中以b为参考数组，查找出a跟b的交集，差集，以及b的新增集合
// same_set为相同集合，rest_set为a多出的集合，new_set为b多出的集合
func FindMixSet(a, b []uint32) (same_set, rest_set, new_set []uint32) {
	i := 0
	len_a := len(a)
	if len_a == 0 {
		new_set = b
		return
	}
	if len(b) == 0 {
		rest_set = a
		return
	}
	for {
		if i > len_a {
			break
		}
		bIn := false
		for j := i; j < len(b); j++ {
			if b[j] == a[i] {
				bIn = true
				b[j], b[i] = b[i], b[j]
				break
			}
		}

		if !bIn {
			if i >= len_a-1 {
				break
			}
			rest_set = append(rest_set, a[i])
			a[i], a[len_a-1] = a[len_a-1], a[i]
			len_a -= 1
		} else {
			same_set = append(same_set, a[i])
			i++
		}
	}
	new_index := len(same_set)
	if new_index > 0 {
		new_set = b[new_index:]
	} else {
		new_set = b[0:]
	}
	return
}

// 获取当前所在exe所在的目录
func GetDir() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		println(err)
		return ""
	}
	return filepath.Dir(file)
}
