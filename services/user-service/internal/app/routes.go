package app

import (
	"log/slog"

	"github.com/SkySock/lode/libs/utils/http/middleware"
	v1 "github.com/SkySock/lode/services/user-service/internal/handler/http/v1"
	"github.com/gorilla/mux"
)

type controllers struct {
	SignUp *v1.SignUp
	SignIn *v1.SignIn
}

func newRouter(log *slog.Logger, controllers controllers) *mux.Router {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	apiV1 := api.PathPrefix("/v1").Subrouter()

	apiV1.Handle("/sign-in", controllers.SignIn).Methods("POST")
	apiV1.Handle("/sign-up", controllers.SignUp).Methods("POST")

	r.Use(middleware.Logging(log))
	r.Use(middleware.Error(log))

	return r
}
