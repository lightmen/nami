package pprof

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/lightmen/nami/alog"
	"github.com/lightmen/nami/pkg/cast"
)

var (
	mainPProfReg = regexp.MustCompile(`<a href=["']([\w?=]*)["']>([\w\s]*)</a>`)
)

const uriPrefix = "/debug/pprof"

type IRouter interface {
	HandleFuncPrefix(prefix string, h http.HandlerFunc)
}

type Transporter struct {
	baseURL string
	dstAddr string
	params  map[string]string
	mux     *http.ServeMux
	router  IRouter
	sync.RWMutex
}

func NewTransporter(dstAddr, baseURL string, opts ...Option) *Transporter {
	trans := &Transporter{
		baseURL: baseURL,
		dstAddr: dstAddr,
		mux:     http.NewServeMux(),
	}

	for _, o := range opts {
		o(trans)
	}

	return trans
}

func (trans *Transporter) buildURL(p *Param, ptype string) string {
	reqURL := trans.dstAddr + uriPrefix
	if ptype != "" {
		reqURL += "/" + ptype
		debugParam := p.Get(keyDebug)
		if debugParam != "" {
			reqURL = fmt.Sprintf("%s?%s=%s", reqURL, keyDebug, debugParam)
		}
	}

	return reqURL
}

func (trans *Transporter) Transport(w http.ResponseWriter, r *http.Request, opts ...Option) {
	param := NewParam(r)

	ptype := param.Get(keyType)
	reqURL := trans.buildURL(param, ptype)

	if param.GetBool(keyOnline) {
		trans.handleOnline(w, r)
		return
	}

	header, body, err := getHttpData(r, reqURL)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	if ptype == "" { //说明是 /debug/pprof 主页面
		body = trans.rebuildMainPage(body)
	}

	copyHeader(w.Header(), header)

	w.Write(body)
}

func (trans *Transporter) rebuildMainPage(body []byte) []byte {
	str := string(body)

	result := mainPProfReg.FindAllStringSubmatch(str, -1)

	for _, vals := range result {
		if len(vals) < 3 {
			continue
		}

		rawHref := vals[1]
		ptype := trans.getPprofType(rawHref)
		newStr := trans.getReplaceContent(ptype, rawHref, vals[2])

		oldStr := vals[0]
		str = strings.ReplaceAll(str, oldStr, newStr)
	}

	return []byte(str)
}

func (trans *Transporter) getPprofType(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		alog.Error("url Parse %s error: %s", rawURL, err.Error())
		return ""
	}

	ptype := u.Path
	return ptype
}

func (trans *Transporter) getReplaceContent(ptype, rawHref, hrefContent string) string {
	u, _ := url.Parse(rawHref)

	query := trans.toParams(map[string]any{
		keyType: ptype,
	})

	content := "<a href='" + trans.baseURL
	if query != "" || u.RawQuery != "" {
		content += "?" + query
		if query != "" && u.RawQuery != "" {
			content += "&" + u.RawQuery
		}
	}

	content += "'>" + hrefContent + "</a>"

	return content + trans.getOnlineHref(ptype, trans.baseURL, query)
}

func (trans *Transporter) getOnlineHref(ptype, baseURL, query string) string {
	if trans.router == nil {
		return ""
	}

	if ptype != "heap" && ptype != "profile" {
		return ""
	}
	content := fmt.Sprintf("&nbsp;<i><a href='%s?%s&online=true' target='_blank' rel='noreferrer'>[online]</a>", GetOnlineURL(baseURL), query)

	return content
}

func (trans *Transporter) toParams(params map[string]any) string {
	u := make(url.Values)
	for key, val := range trans.params {
		u.Add(key, val)
	}

	for key, val := range params {
		u.Add(key, cast.ToString(val))
	}

	return u.Encode()
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func getHttpData(r *http.Request, reqURL string) (header http.Header, body []byte, err error) {
	cli := &http.Client{}

	newReq, err := http.NewRequest(r.Method, reqURL, http.NoBody)
	if err != nil {
		alog.Error("http NewRequest error: %s", err.Error())
		return
	}

	contentType := r.Header.Get("Content-Type")
	newReq.Header.Set("Content-Type", contentType)

	resp, err := cli.Do(newReq)
	if err != nil {
		alog.Error("cli.Do error: %s", err.Error())
		return
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		alog.Error("io.ReadAll error: %s", err.Error())
		return
	}
	header = resp.Header

	return
}
