package pprof

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/google/pprof/driver"
	"github.com/lightmen/nami/alog"
	"github.com/lightmen/nami/pkg/cast"
)

func (trans *Transporter) UI(w http.ResponseWriter, r *http.Request) {
	trans.mux.ServeHTTP(w, r)
}

func (trans *Transporter) handleOnline(w http.ResponseWriter, r *http.Request) {
	if trans.router == nil {
		return
	}

	param := NewParam(r)
	ptype := param.Get(keyType)
	reqURL := trans.buildURL(param, ptype)

	_, body, err := getHttpData(r, reqURL)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	id, err := trans.registerOnline(ptype, body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	query := r.URL.RawQuery

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	u := url.URL{
		Scheme:   scheme,
		Host:     r.Host,
		Path:     path.Join(GetOnlineURL(trans.baseURL), id),
		RawQuery: query,
	}
	target := u.String()

	trans.router.HandleFuncPrefix(path.Join(GetOnlineURL(trans.baseURL), id), trans.UI)

	http.Redirect(w, r, target, http.StatusSeeOther)

	return
}

func (trans *Transporter) registerOnline(ptype string, data []byte) (id string, err error) {
	id = cast.ToString(time.Now().UnixMilli())

	file := path.Join(os.TempDir(), fmt.Sprintf("%s_%s", ptype, id))
	if err = os.WriteFile(file, data, 0600); err != nil {
		alog.Error("write file %s error: %s", file, err.Error())
		return
	}

	flags := &flags{
		args: []string{"-http=localhost:0", "-no_browser", file},
	}

	curPath := path.Join(GetOnlineURL(trans.baseURL), id) + "/"

	options := &driver.Options{
		Flagset: flags,
		HTTPServer: func(args *driver.HTTPServerArgs) error {
			for pattern, handler := range args.Handlers {
				var joinedPattern string
				if pattern == "/" {
					joinedPattern = curPath
				} else {
					joinedPattern = path.Join(curPath, pattern)
				}
				trans.mux.Handle(joinedPattern, handler)

				alog.Info("register online url: %s", joinedPattern)
			}
			return nil
		},
	}

	if err = driver.PProf(options); err != nil {
		return
	}

	return
}

func GetOnlineURL(baseURL string) string {
	return baseURL
}

func hexEscapeNonASCII(s string) string {
	newLen := 0
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			newLen += 3
		} else {
			newLen++
		}
	}
	if newLen == len(s) {
		return s
	}
	b := make([]byte, 0, newLen)
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			b = append(b, '%')
			b = strconv.AppendInt(b, int64(s[i]), 16)
		} else {
			b = append(b, s[i])
		}
	}
	return string(b)
}
