package endpoint

import (
	"net/url"

	"github.com/lightmen/nami/transport"
)

func New(scheme, host string) *url.URL {
	return &url.URL{
		Scheme: scheme,
		Host:   host,
	}
}

// ParseEndpoint parses an Endpoint URL.
func ParseEndpoint(endpoints []string, scheme string) (string, error) {
	for _, e := range endpoints {
		u, err := url.Parse(e)
		if err != nil {
			return "", err
		}

		if u.Scheme == scheme {
			return u.Host, nil
		}
	}
	return "", nil
}

// GetEndpoint 从endpoints数组中选出等于scheme的元素
func GetEndpoint(endpoints []string, scheme string) string {
	for _, e := range endpoints {
		u, err := url.Parse(e)
		if err != nil {
			return ""
		}

		if u.Scheme == scheme {
			return e
		}
	}
	return ""
}

func GetGrpcEndpoint(endpoints []string) string {
	return GetEndpoint(endpoints, transport.GRPC)
}

func GetHttpEndpoint(endpoints []string) string {
	return GetEndpoint(endpoints, transport.HTTP)
}
