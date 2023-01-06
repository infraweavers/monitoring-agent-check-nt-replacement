package httpclient

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type mock struct {
	Timeout            time.Duration
	Transport          *http.Transport
	RequestBodyContent string
	RequestVerb        string
	RequestHeaders     map[string][]string
	RequestURI         *url.URL
	RequestHost        string

	DoFunc         func(*http.Request) (*http.Response, error)
	DoSetTimeout   func(time.Duration)
	DoSetTransport func(*http.Transport)
}

func (H mock) Do(r *http.Request) (*http.Response, error) {
	return H.DoFunc(r)
}

func (H mock) SetTimeout(timeout time.Duration) {
	H.DoSetTimeout(timeout)
}

func (H mock) SetTransport(transport *http.Transport) {
	H.DoSetTransport(transport)
}

func NewMockHTTPClient(jsonResponse string, httpResponseCode int) *mock {
	client := new(mock)

	client.DoFunc = func(r *http.Request) (*http.Response, error) {
		bodyContent, _ := ioutil.ReadAll(r.Body)
		client.RequestBodyContent = string(bodyContent)
		client.RequestVerb = r.Method
		client.RequestURI = r.URL
		client.RequestHost = r.Host
		client.RequestHeaders = r.Header

		return &http.Response{
			Body:       io.NopCloser(strings.NewReader(jsonResponse)),
			StatusCode: httpResponseCode,
		}, nil
	}

	client.DoSetTimeout = func(timeout time.Duration) {
		client.Timeout = timeout
	}
	client.DoSetTransport = func(transport *http.Transport) {
		client.Transport = transport
	}

	return client
}
