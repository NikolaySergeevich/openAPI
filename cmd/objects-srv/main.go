package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	h "openapi/internal/handler"
	"openapi/internal/memstore"
	"openapi/pkg/api/objapi"
)

func main() {
	router := chi.NewRouter()
	memStore := memstore.New()
	handler := h.NewHandler(memStore)
	router.Mount(
		"/api", objapi.HandlerWithOptions(
			handler, objapi.ChiServerOptions{
				BaseURL: "/v1",
				ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
					slog.Error("handle error", slog.String("err", err.Error()))
				},
			},
		),
	)

	srv := http.Server{
		Addr:              ":8181",
		Handler:           router,
		ReadTimeout:       20 * time.Second,
		ReadHeaderTimeout: 20 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       20 * time.Second,
		MaxHeaderBytes:    10 * 1024 * 1024, // 10mib
	}

	slog.Info(fmt.Sprintf("http server was started %s", ":8181"))
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("http.Server ListenAndServe", slog.String("err", err.Error()))
	}
}
