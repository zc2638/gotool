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
type HttpResp struct {
	Err      error
	Response *http.Response
}

func NewHttpResp(request *http.Request, client *http.Client) *HttpResp {
	if client == nil {
		client = &http.Client{}
	}
	resp, err := client.Do(request)
	return &HttpResp{Err: err, Response: resp}
}

func (h *HttpResp) GetBody() ([]byte, error) {
	if h.Err != nil {
		return nil, h.Err
	}
	resp := h.Response
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (h *HttpResp) ParseJSON(data interface{}) error {
	res, err := h.GetBody()
	if err != nil {
		return err
	}
	return json.Unmarshal(res, data)
}

func (h *HttpResp) ParseXMLToMap() (map[string]string, error) {
	res, err := h.GetBody()
	if err != nil {
		return nil, err
	}
	return utilx.XmlToMap(res), nil
}

func (h *HttpResp) ParseXML(data interface{}) error {
	m, err := h.ParseXMLToMap()
	if err != nil {
		return err
	}

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, data)
}
