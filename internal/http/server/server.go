package server

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/andreyAKor/otus_project/internal/cache"
	"github.com/andreyAKor/otus_project/internal/http/client"
	"github.com/andreyAKor/otus_project/internal/image"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	limitReadBody int64 = 1024 * 1_000
)

var (
	ErrServerNotInit = errors.New("server not init")

	_ io.Closer = (*Server)(nil)
)

type Server struct {
	client client.Client
	image  image.Image
	cache  cache.Cache
	host   string
	port   int
	server *http.Server
	ctx    context.Context
}

func New(
	client client.Client,
	image image.Image,
	cache cache.Cache,
	host string,
	port int,
) (*Server, error) {
	return &Server{
		client: client,
		image:  image,
		cache:  cache,
		host:   host,
		port:   port,
	}, nil
}

// Running http-server.
func (s *Server) Run(ctx context.Context) error {
	if s.server != nil {
		return nil
	}

	s.ctx = ctx

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.proxy(s.process))

	// middlewares
	handler := s.body(mux)
	handler = s.logger(handler)

	s.server = &http.Server{
		Addr:    net.JoinHostPort(s.host, strconv.Itoa(s.port)),
		Handler: handler,
	}

	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return errors.Wrap(err, "http-server listen fail")
	}

	return nil
}

func (s *Server) Close() error {
	if s.server == nil {
		return ErrServerNotInit
	}

	return s.server.Shutdown(s.ctx)
}

// Middleware logger output log info of request, e.g.: r.Method, r.URL etc.
func (s Server) logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := newAppResponseWriter(w)

		start := time.Now()
		defer func() {
			i := log.Info()

			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				i.Err(err)
			}

			i.Str("ip", host).
				Str("startAt", start.String()).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("proto", r.Proto).
				Int("status", rw.statusCode).
				TimeDiff("latency", time.Now(), start)

			if len(r.UserAgent()) > 0 {
				i.Str("userAgent", r.UserAgent())
			}

			i.Msg("http-request")
		}()

		handler.ServeHTTP(rw, r)
	})
}

// Middleware preparing body request.
func (s Server) body(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, limitReadBody))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		handler.ServeHTTP(w, r)
	})
}

// Proxing requested preview image to response.
func (s Server) proxy(h func(w http.ResponseWriter, r *http.Request) ([]byte, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := h(w, r)
		if err != nil {
			log.Error().Err(err).Msg("proxy fail")
			return
		}

		if len(body) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Msg("body is empty")
			return
		}

		w.Header().Set("Content-Length", strconv.Itoa(len(body)))

		if _, err := w.Write(body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("write fail")
			return
		}
	}
}

var _ http.ResponseWriter = (*appResponseWriter)(nil)

// App wrapper over http.ResponseWriter.
type appResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newAppResponseWriter(w http.ResponseWriter) *appResponseWriter {
	return &appResponseWriter{w, http.StatusOK}
}

func (a *appResponseWriter) WriteHeader(code int) {
	a.statusCode = code
	a.ResponseWriter.WriteHeader(code)
}
