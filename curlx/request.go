package curlx

import (
	"bytes"
	"crypto/tls"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

/**
 * Created by zc on 2019/12/16.
 */
type FileInfo struct {
	Name   string
	Stream io.Reader
}

type Request struct {
	client             *http.Client
	transport          *http.Transport
	Url                string
	Method             string
	Header             map[string]string
	Query              map[string]string
	Params             map[string]string
	FileData           map[string]FileInfo
	Body               []byte
	BodyReader         io.Reader
	InsecureSkipVerify bool
}

func NewRequest(options ...RequestOption) *Request {
	req := &Request{}
	for _, option := range options {
		option(req)
	}
	return req
}

func (h *Request) buildUrl() {
	if h.Query == nil {
		return
	}

	query := url.Values{}
	for k, v := range h.Query {
		query.Set(k, v)
	}

	urlSet := strings.Split(h.Url, "?")
	switch len(urlSet) {
	case 1:
		h.Url += "?" + query.Encode()
	case 2:
		if urlSet[1] != "" {
			urlSet[1] += "&"
		}
		h.Url = urlSet[0] + "?" + url.PathEscape(urlSet[1]+query.Encode())
	}
}

func (h *Request) buildBody() {
	if h.Body != nil || h.BodyReader != nil {
		return
	}

	params := url.Values{}
	for k, v := range h.Params {
		params.Set(k, v)
	}
	h.Body = []byte(params.Encode())
	if h.BodyReader == nil {
		h.BodyReader = bytes.NewReader(h.Body)
	}
}

func (h *Request) initClient() {
	if h.client == nil {
		h.client = &http.Client{}
	}
	if h.transport != nil {
		h.client.Transport = h.transport
	}
	if h.InsecureSkipVerify {
		h.client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
}

func (h *Request) Do() (*Response, error) {

	h.buildUrl()
	h.buildBody()
	h.initClient()

	bReader := h.BodyReader
	if bReader == nil {
		bReader = bytes.NewReader(h.Body)
	}

	req, err := http.NewRequest(h.Method, h.Url, bReader)
	if err != nil {
		return nil, err
	}

	if h.Header != nil {
		for k, v := range h.Header {
			req.Header.Set(k, v)
		}
	}

	res, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	return NewResponse(res)
}

func (h *Request) Get() ([]byte, error) {
	h.Method = MethodGET
	res, err := h.Do()
	if err != nil {
		return nil, err
	}
	return res.ParseBody()
}

func (h *Request) Post() ([]byte, error) {
	h.Method = MethodPOST
	if h.Header == nil {
		h.Header = make(map[string]string)
	}
	if _, ok := h.Header[HeaderContentType]; !ok {
		h.Header[HeaderContentType] = ApplicationFormURLEncoded
	}
	res, err := h.Do()
	if err != nil {
		return nil, err
	}
	return res.ParseBody()
}

func (h *Request) PostForm() ([]byte, error) {
	h.Method = MethodPOST
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	if h.FileData != nil {
		for k, file := range h.FileData {
			part, err := w.CreateFormFile(k, file.Name)
			if err != nil {
				return nil, err
			}
			if _, err = io.Copy(part, file.Stream); err != nil {
				return nil, err
			}
		}
	}

	if h.Params != nil {
		for k, v := range h.Params {
			if err := w.WriteField(k, v); err != nil {
				return nil, err
			}
		}
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	if h.Header == nil {
		h.Header = make(map[string]string)
	}
	h.Header[HeaderContentType] = w.FormDataContentType()
	h.BodyReader = &buf

	res, err := h.Do()
	if err != nil {
		return nil, err
	}
	return res.ParseBody()
}
