package utilx

import (
	"bytes"
	"encoding/xml"
)

/**
 * Created by zc on 2019/12/16.
 */
// map转xml
func MapToXml(data map[string]string) []byte {

	var buf bytes.Buffer
	buf.WriteString(`<xml>`)
	for k, v := range data {
		buf.WriteString(`<`)
		buf.WriteString(k)
		buf.WriteString(`><![CDATA[`)
		buf.WriteString(v)
		buf.WriteString(`]]></`)
		buf.WriteString(k)
		buf.WriteString(`>`)
	}
	buf.WriteString(`</xml>`)

	return buf.Bytes()
}

// xml转map
func XmlToMap(b []byte) map[string]string {

	params := make(map[string]string)
	decoder := xml.NewDecoder(bytes.NewReader(b))

	var key, value string
	for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement: // 开始标签
			key = token.Name.Local
		case xml.CharData: // 标签内容
			content := string([]byte(token))
			value = content
		}
		if key != "xml" {
			if value != "\n" {
				params[key] = value
			}
		}
	}
	return params
}