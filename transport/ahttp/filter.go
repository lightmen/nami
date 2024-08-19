package ahttp

import "net/http"

type FilterFunc func(http.Handler) http.Handler

func FilterChain(filters ...FilterFunc) FilterFunc {
	return func(next http.Handler) http.Handler {
		for i := len(filters) - 1; i >= 0; i-- {
			next = filters[i](next)
		}

		return next
	}
}
