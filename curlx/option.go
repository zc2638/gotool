/**
 * Created by zc on 2020/5/2.
 */
package curlx

import (
	"crypto/tls"
	"errors"
	"net/http"
)

type RequestOption func(*Request)

func SetClient(client *http.Client) RequestOption {
	return func(request *Request) {
		request.client = client
	}
}

func SetTransport(certFilePath, keyFilePath string) (RequestOption, error) {
	if certFilePath == "" {
		return nil, errors.New("certFilePath is empty")
	}
	cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil {
		return nil, err
	}
	return func(request *Request) {
		request.transport = &http.Transport{
			DisableCompression: true,
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		}
	}, nil
}
