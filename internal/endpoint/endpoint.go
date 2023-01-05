package endpoint

import "net/url"

func New(scheme, host string) *url.URL {
	return &url.URL{
		Scheme: scheme,
		Host:   host,
	}
}
