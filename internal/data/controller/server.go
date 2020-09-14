package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/shipa988/fanouter/internal/domain/usecase"
	"github.com/shipa988/fanouter/internal/domain/usecase/fanouter"
	"github.com/shipa988/fanouter/internal/util"
)

const (
	ErrID = "must be id in query"
)
const (
	MainAnswer = `go to /feeds/{id}`
)

type HTTPServer struct {
	logger   usecase.Logger
	server   *http.Server
	fanouter fanouter.Fanouter
}

func NewHttpServer(addr string, logger usecase.Logger, fanouter fanouter.Fanouter) *HTTPServer {
	server := &http.Server{Addr: addr}
	return &HTTPServer{
		server:   server,
		logger:   logger,
		fanouter: fanouter,
	}
}

func (s *HTTPServer) Serve() error {
	s.logger.Log(context.Background(), "starting http server on address [%v]", s.server.Addr)

	router := mux.NewRouter()

	router.HandleFunc("/", s.main).Methods(http.MethodGet)
	router.HandleFunc("/feeds/{id}", s.fanout).Methods(http.MethodGet)

	handler := s.accessLogMiddleware(router)
	handler = s.panicMiddleware(handler)
	s.server.Handler = handler

	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return errors.Wrapf(err, "can't start listen address [%v]", s.server.Addr)
	}
	return nil
}

func (s *HTTPServer) main(w http.ResponseWriter, r *http.Request) {
	s.httpAnswer(w, MainAnswer, http.StatusOK)
}

func (s *HTTPServer) fanout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		s.httpError(r.Context(), w, ErrID, http.StatusBadRequest)
		return
	}
	err := s.fanouter.Fanout(r.Context(), id)
	if err != nil {
		s.httpError(r.Context(), w, err.Error(), http.StatusBadRequest)
		return
	}
	s.httpAnswer(w, "query send", http.StatusOK)
}

func (s *HTTPServer) StopServe() {
	ctx := context.Background()
	s.logger.Log(ctx, "stopping http server")
	defer s.logger.Log(ctx, "http server stopped")
	if s.server == nil {
		s.logger.Log(ctx, "http server is nil")
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Log(ctx, "can't stop http server with error: %v", err)
	}
}

func (s *HTTPServer) accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := util.SetRequestID(r.Context())

		next.ServeHTTP(w, r.WithContext(ctx))

		latency := time.Since(start)
		s.logRequest(ctx, r.RemoteAddr, start.Format(util.LayoutISO), r.Method, r.URL.Path, latency)
	})
}

func (s *HTTPServer) panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.httpError(r.Context(), w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *HTTPServer) logRequest(ctx context.Context, remoteAddr, start, method, path string, latency time.Duration) {
	s.logger.Log(ctx, "%s [%s] %s %s [%s]", remoteAddr, start, method, path, latency)
}

func (s *HTTPServer) httpError(ctx context.Context, w http.ResponseWriter, error string, code int) {
	s.logger.Log(ctx, error)
	http.Error(w, error, code)
}

func (s *HTTPServer) httpAnswer(w http.ResponseWriter, msg interface{}, code int) {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		code = http.StatusInternalServerError
	}
	w.WriteHeader(code)
	w.Write(jmsg) //nolint:errcheck
}
