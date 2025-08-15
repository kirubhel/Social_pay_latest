package http

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	log      *log.Logger
	ServeMux *http.ServeMux
	Server   *http.Server
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		origin := r.Header.Get("Origin")
		if origin == "https://dashboard.socialpay.co" {
			w.Header().Set("Access-Control-Allow-Origin", origin)

		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if r.Method == http.MethodPost {
			contentType := r.Header.Get("Content-Type")
			switch {
			case strings.Contains(contentType, "text/xml"):
			case strings.Contains(contentType, "application/json"):
			case strings.Contains(contentType, "text/plain"):
			default:
				http.Error(w, "Invalid Content-Type, expected text/xml, application/json, or text/plain", http.StatusBadRequest)
				return
			}
		}

		if r.Method == http.MethodOptions {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func LimitBodySize(h http.Handler, maxBytes int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		h.ServeHTTP(w, r)
	})
}

func ParseXMLRequest(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return xml.NewDecoder(r.Body).Decode(v)
}

func ParseJSONRequest(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func New(log *log.Logger) Server {
	// sm := &http.ServeMux{}

	// s := &http.Server{
	// 	Handler: accessControl(sm),
	// 	Addr:    "0.0.0.0:8004",
	// }

	sm := &http.ServeMux{}

	var handler http.Handler = sm
	handler = accessControl(handler)
	handler = LimitBodySize(handler, 10<<20)

	s := &http.Server{
		Handler: handler,
		Addr:    "0.0.0.0:8004",
	}

	return Server{log: log, ServeMux: sm, Server: s}
}

func (s Server) Serve() {
	go func() {
		err := s.Server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.log.Println("Server error:", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := s.Server.Shutdown(ctx); err != nil {
		s.log.Println("Server shutdown error:", err)
	}
}
