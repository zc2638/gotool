package utilx

import (
	"reflect"
	"strconv"
	"strings"
)

/**
 * Created by zc on 2019/12/16.
 */
// 字符串驼峰转下划线
func CamelToUnderline(s string) string {
	num := len(s)
	data := make([]byte, 0, num*2)
	j := false
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// 字符串下划线转驼峰
func UnderlineToCamel(s string) string {
	data := make([]byte, 0, len(s))
	j, k := false, false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if !k && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || !k) {
			d = d - 32
			j, k = false, true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// 判断元素是否存在数组中
func InSlice(val interface{}, array interface{}) (exists bool, index int) {
	index = -1
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				index = i
				exists = true
				return
			}
		}
	}
	return
}

// 通用map转字符串map
func MapToStringMap(data map[string]interface{}) map[string]string {

	var res = make(map[string]string)
	for k, v := range data {
		var val string
		switch v.(type) {
		case int:
			val = strconv.Itoa(v.(int))
		case int8:
			val = strconv.Itoa(int(v.(int8)))
		case int16:
			val = strconv.Itoa(int(v.(int16)))
		case int32:
			val = strconv.Itoa(int(v.(int32)))
		case int64:
			val = strconv.Itoa(int(v.(int64)))
		case float32:
			val = strconv.FormatFloat(float64(v.(float32)), 'f', -1, 64)
		case float64:
			val = strconv.FormatFloat(v.(float64), 'f', -1, 64)
		case uint:
			// 转十进制字符串
			val = strconv.FormatUint(uint64(v.(uint)), 10)
		case uint8:
			val = strconv.FormatUint(uint64(v.(uint8)), 10)
		case uint16:
			val = strconv.FormatUint(uint64(v.(uint16)), 10)
		case uint32:
			val = strconv.FormatUint(uint64(v.(uint32)), 10)
		case uint64:
			val = strconv.FormatUint(v.(uint64), 10)
		case bool:
			val = strconv.FormatBool(v.(bool))
		case string:
			val = v.(string)
		}
		res[k] = val
	}
	return res
}