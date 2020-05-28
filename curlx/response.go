package curlx

import (
	"encoding/json"
	"github.com/zc2638/gotool/utilx"
	"io/ioutil"
	"net/http"
)

/**
 * Created by zc on 2019/12/16.
 */
type Response struct {
	*http.Response
	Result []byte
}

func NewResponse(res *http.Response) (*Response, error) {
	return &Response{Response: res}, nil
}

func (res *Response) ParseBody() ([]byte, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := res.Body.Close(); err != nil {
		return nil, err
	}
	res.Result = body
	return body, nil
}

func (res *Response) ParseJSON(data interface{}) error {
	body, err := res.ParseBody()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, data)
}

func (res *Response) ParseXMLToMap() (map[string]string, error) {
	body, err := res.ParseBody()
	if err != nil {
		return nil, err
	}
	return utilx.XmlToMap(body), nil
}

func (res *Response) ParseXML(data interface{}) error {
	m, err := res.ParseXMLToMap()
	if err != nil {
		return err
	}
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, data)
}
