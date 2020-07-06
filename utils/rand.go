package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"reflect"
	"time"

	"sipemanager/dao"
)

func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// return len=8  salt
func GetRandomSalt() string {
	return GetRandomString(20)
}

//生成随机字符串
func GetRandomString(lenght int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ%$*&#@"
	bytes := []byte(str)
	bytesLen := len(bytes)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lenght; i++ {
		result = append(result, bytes[r.Intn(bytesLen)])
	}
	return string(result)

}

//结构体数组去重
func RemoveRepByLoop(nodes []dao.InstanceNodes) (result []dao.InstanceNodes) {
	n := len(nodes)
	for i := 0; i < n; i++ {
		state := false
		for j := i + 1; j < n; j++ {
			if j > 0 && reflect.DeepEqual(nodes[i].ChainId, nodes[j].ChainId) {
				state = true
				break
			}
		}
		if !state {
			result = append(result, nodes[i])
		}
	}
	return
}
