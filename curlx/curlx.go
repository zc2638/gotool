package curlx

import "net/http"

/**
 * Created by zc on 2019/12/16.
 */
// 获取文件类型
func GetFileContentType(fileByte []byte) string {
	buffer := make([]byte, 512)
	buffer = append(buffer, fileByte[:512]...)
	return http.DetectContentType(fileByte)
}