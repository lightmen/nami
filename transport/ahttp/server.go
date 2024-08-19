package ahttp

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/lightmen/nami/alog"
	"github.com/lightmen/nami/internal/host"
	"github.com/lightmen/nami/pkg/endpoint"
	"github.com/lightmen/nami/transport"
	"github.com/lightmen/nami/transport/ahttp/pprof"
)

var (
	_ transport.Server = (*Server)(nil)
)

var gServer *Server

func GetServer() *Server {
	return gServer
}

type Server struct {
	*http.Server
	lis      net.Listener
	network  string
	address  string
	timeout  time.Duration
	filters  []FilterFunc
	endpoint *url.URL
	router   *mux.Router
	usePprof bool
	method   string
}

func New(address string, opts ...Option) (srv *Server, err error) {
	srv = &Server{
		network:  "tcp",
		address:  address,
		timeout:  3 * time.Second,
		filters:  make([]FilterFunc, 0),
		router:   mux.NewRouter(),
		usePprof: true,
		method:   http.MethodGet,
	}

	srv.router.StrictSlash(true)
	srv.router.NotFoundHandler = http.DefaultServeMux
	srv.router.MethodNotAllowedHandler = http.DefaultServeMux
	srv.router.Use(srv.filter())

	for _, opt := range opts {
		opt(srv)
	}

	srv.Server = &http.Server{
		Handler: FilterChain(srv.filters...)(srv.router),
	}

	err = srv.listen()
	if err != nil {
		return nil, err
	}

	if srv.usePprof {
		srv.HandlePrefix("/debug/pprof", pprof.HandleFunc())
	}

	if gServer == nil {
		gServer = srv
	}

	return
}

func (s *Server) Start(ctx context.Context) error {
	alog.InfoCtx(ctx, "[HTTP] server lintening on: %s", s.lis.Addr().String())

	s.BaseContext = func(net.Listener) context.Context {
		return ctx
	}

	err := s.Serve(s.lis)
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) (err error) {
	alog.InfoCtx(ctx, "[HTTP] server stopping")
	return s.Shutdown(ctx)
}

func (s *Server) listen() error {
	if s.lis != nil {
		return nil
	}

	lis, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}

	s.lis = lis

	return nil
}

func (s *Server) Endpoint() (*url.URL, error) {
	if s.endpoint == nil {
		if err := s.listen(); err != nil {
			return nil, err
		}

		addr, err := host.Extract(s.address, s.lis)
		if err != nil {
			return nil, err
		}

		s.endpoint = endpoint.New(s.Name(), addr)
	}

	return s.endpoint, nil
}

func (s *Server) filter() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				ctx    context.Context
				cancel context.CancelFunc
			)

			if s.timeout > 0 {
				ctx, cancel = context.WithTimeout(r.Context(), s.timeout)
			} else {
				ctx, cancel = context.WithCancel(r.Context())
			}
			defer cancel()

			// pathTemplate := r.URL.Path
			// if route := mux.CurrentRoute(r); route != nil {
			// 	// /path/123 -> /path/{id}
			// 	pathTemplate, _ = route.GetPathTemplate()
			// }

			// tr := &Transport{
			// 	operation:    pathTemplate,
			// 	pathTemplate: pathTemplate,
			// 	reqHeader:    headerCarrier(r.Header),
			// 	replyHeader:  headerCarrier(w.Header()),
			// 	request:      r,
			// }
			// if s.endpoint != nil {
			// 	tr.endpoint = s.endpoint.String()
			// }
			// tr.request = r.WithContext(transport.NewServerContext(ctx, tr))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (s *Server) Name() string {
	return transport.HTTP
}

// Handle registers a new route with a matcher for the URL path.
func (s *Server) Handle(path string, h http.Handler) {
	s.router.Handle(path, h)
}

// HandleFunc registers a new route with a matcher for the URL path.
func (s *Server) HandleFunc(path string, h http.HandlerFunc) {
	s.router.HandleFunc(path, h)
}

// HandlePrefix registers a new route with a matcher for the URL path prefix.
func (s *Server) HandlePrefix(prefix string, h http.Handler) {
	s.router.PathPrefix(prefix).Handler(h)
}

// HandleFuncPrefix registers a new route with a matcher for the URL path prefix.
func (s *Server) HandleFuncPrefix(prefix string, h http.HandlerFunc) {
	s.router.PathPrefix(prefix).HandlerFunc(h)
}
