package utilx

import (
	"math/rand"
	"time"
)

/**
 * Created by zc on 2019/12/16.
 */
// 生成随机字符串
func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := []byte(str)
	result := make([]byte, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, b[r.Intn(len(b))])
	}
	return string(result)
}