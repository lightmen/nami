package transport

import "net/url"

type Endpointer interface {
	Endpoint() (*url.URL, error)
}
