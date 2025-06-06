package app

import (
	"log/slog"

	"github.com/SkySock/lode/libs/utils/http/middleware"
	"github.com/SkySock/lode/services/user-service/internal/handler/http/v1/auth"
	"github.com/gorilla/mux"
)

type controllers struct {
	SignUp  *auth.SignUp
	SignIn  *auth.SignIn
	SignOut *auth.SignOut
}

func newRouter(log *slog.Logger, controllers controllers) *mux.Router {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	apiV1 := api.PathPrefix("/v1").Subrouter()
	authV1 := apiV1.PathPrefix("/auth").Subrouter()

	authV1.Handle("/sign-in", controllers.SignIn).Methods("POST")
	authV1.Handle("/sign-up", controllers.SignUp).Methods("POST")
	authV1.Handle("/sign-out", controllers.SignOut).Methods("POST")

	r.Use(middleware.Logging(log))
	r.Use(middleware.Error(log))

	return r
}
