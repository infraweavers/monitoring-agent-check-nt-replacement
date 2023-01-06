package httpclient

import (
	"net/http"
	"time"
)

type Interface interface {
	Do(*http.Request) (*http.Response, error)
	SetTimeout(timeout time.Duration)
	SetTransport(transport *http.Transport)
}
