package pprof

import (
	"net/http"

	"github.com/lightmen/nami/pkg/cast"
)

const (
	keyType   = "type"
	keyDebug  = "debug"
	keyID     = "id"
	keyOnline = "online"
)

type Param struct {
	r *http.Request
}

func NewParam(r *http.Request) *Param {
	r.ParseForm()

	return &Param{
		r: r,
	}
}

func (p *Param) Get(key string) string {
	return p.r.FormValue(key)
}

func (p *Param) GetBool(key string) bool {
	val := p.Get(key)
	return cast.ToBool(val)
}
