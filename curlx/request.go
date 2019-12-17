package curlx

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

/**
 * Created by zc on 2019/12/16.
 */
const (
	METHOD_POST = "POST"
	METHOD_GET  = "GET"
)

const (
	HEADER_CONTENT_TYPE = "Content-Type"
)

const (
	CT_APPLICATION_FORM_URLENCODED = "application/x-www-form-urlencoded"
	CT_APPLICATION_JSON            = "application/json"
	CT_APPLICATION_JSON_UTF8       = "application/json; charset=UTF-8"
	CT_APPLICATION_OCTET_STREAM    = "application/octet-stream"
	CT_APPLICATION_PDF             = "application/pdf"
	CT_APPLICATION_XML             = "application/xml"
	CT_IMAGE_GIF                   = "image/gif"
	CT_IMAGE_JPEG                  = "image/jpeg"
	CT_IMAGE_PNG                   = "image/png"
	CT_TEXT_HTML                   = "text/html"
	CT_TEXT_MARKDOWN               = "text/markdown"
	CT_TEXT_PLAIN                  = "text/plain"
	CT_TEXT_XML                    = "text/xml"
	CT_TEXT_XML_UTF8               = "text/xml; charset=UTF-8"
)

type FileInfo struct {
	Name   string
	Stream io.Reader
}

type FormData struct {
	File   map[string]FileInfo
	Params map[string]string
}

type HttpReq struct {
	client     *http.Client
	transport  *http.Transport
	resp       *HttpResp
	Url        string
	Method     string
	Header     map[string]string
	Query      map[string]string
	Params     map[string]string
	FormData   FormData
	Body       []byte
	BodyReader io.Reader
	CertFile   string
	KeyFile    string
	Timeout    time.Duration
}

func NewTransportCert(certFilePath, keyFilePath string) (tls.Certificate, error) {
	return tls.LoadX509KeyPair(certFilePath, keyFilePath)
}

func NewTransport(certFilePath, keyFilePath string, InsecureSkipVerify bool) (*http.Transport, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: InsecureSkipVerify},
	}
	if certFilePath == "" {
		return nil, errors.New("certFilePath is empty")
	}
	cert, err := NewTransportCert(certFilePath, keyFilePath)
	if err != nil {
		return nil, err
	}
	tr.DisableCompression = true
	tr.TLSClientConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	return tr, nil
}

func (h *HttpReq) buildUrl() {
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

func (h *HttpReq) buildBody() {
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

func (h *HttpReq) SetClient(client *http.Client) {
	h.client = client
}

func (h *HttpReq) SetTransport(transport *http.Transport) {
	h.transport = transport
}

func (h *HttpReq) initClient() error {
	if h.client != nil {
		return nil
	}
	h.client = &http.Client{
		Timeout: h.Timeout,
	}
	if h.CertFile != "" {
		tr, err := NewTransport(h.CertFile, h.KeyFile, false)
		if err != nil {
			return err
		}
		h.client.Transport = tr
	}
	return nil
}

func (h *HttpReq) Do() *HttpResp {

	h.buildUrl()
	h.buildBody()
	response := &HttpResp{}
	if err := h.initClient(); err != nil {
		response.Err = err
		return response
	}

	var bReader = h.BodyReader
	if h.BodyReader == nil {
		bReader = bytes.NewReader(h.Body)
	}

	req, err := http.NewRequest(h.Method, h.Url, bReader)
	if err != nil {
		response.Err = err
		return response
	}

	if h.Header != nil {
		for k, v := range h.Header {
			req.Header.Set(k, v)
		}
	}
	return NewHttpResp(req, h.client)
}

func (h *HttpReq) Get() ([]byte, error) {
	h.Method = METHOD_GET
	return h.Do().GetBody()
}

func (h *HttpReq) Post() ([]byte, error) {

	h.Method = METHOD_POST
	if h.Header == nil {
		h.Header = make(map[string]string)
	}
	if _, ok := h.Header[HEADER_CONTENT_TYPE]; !ok {
		h.Header[HEADER_CONTENT_TYPE] = CT_APPLICATION_FORM_URLENCODED
	}
	return h.Do().GetBody()
}

func (h *HttpReq) PostForm() ([]byte, error) {

	h.Method = METHOD_POST

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	if h.FormData.File != nil {
		for k, file := range h.FormData.File {
			part, err := w.CreateFormFile(k, file.Name)
			if err != nil {
				return nil, err
			}
			if _, err = io.Copy(part, file.Stream); err != nil {
				return nil, err
			}
		}
	}

	if h.FormData.Params != nil {
		for k, v := range h.FormData.Params {
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
	h.Header[HEADER_CONTENT_TYPE] = w.FormDataContentType()
	h.BodyReader = &buf
	return h.Do().GetBody()
}
