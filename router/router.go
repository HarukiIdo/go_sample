package router

import (
	"database/sql"
	"net/http"

	"github.com/HarukiIdo/go_sample/handler"
	"github.com/HarukiIdo/go_sample/health"
	"github.com/HarukiIdo/go_sample/middleware"
	"github.com/HarukiIdo/go_sample/service"
)

func NewRouter(db *sql.DB) *http.ServeMux {

	mux := http.NewServeMux()
	mux.Handle("/comment", handler.NewCommentHandler(service.NewCommentService(db)))
	mux.Handle("healthz", middleware.MiddlewareLogging(http.HandlerFunc(health.Healthz)))
	return mux
}
