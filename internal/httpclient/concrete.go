package httpclient

import (
	"net/http"
	"time"
)

type concrete struct {
	// DoFunc will be executed whenever Do function is executed
	// so we'll be able to create a custom response
	DoFunc        func(*http.Request) (*http.Response, error)
	wrappedClient *http.Client
}

func (H concrete) Do(r *http.Request) (*http.Response, error) {
	//fmt.Printf("DOING REQUEST %#v", r)
	return H.wrappedClient.Do(r)
}

func (H concrete) SetTimeout(timeout time.Duration) {
	H.wrappedClient.Timeout = timeout
	//fmt.Printf("SETTING TIMEOUT%#v", timeout)
}

func (H concrete) SetTransport(transport *http.Transport) {
	H.wrappedClient.Transport = transport
	//fmt.Printf("SETTING TRANSPORT%#v", transport)
}

func NewHTTPClient() *concrete {
	client := new(concrete)
	client.wrappedClient = &http.Client{}
	return client
}
