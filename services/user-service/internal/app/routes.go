package app

import (
	"log/slog"

	"github.com/SkySock/lode/libs/utils/http/middleware"
	"github.com/SkySock/lode/services/user-service/docs"
	"github.com/SkySock/lode/services/user-service/internal/handler/http/v1/auth"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
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

	docs.SwaggerInfo.Title = "User Service API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "0.0.0.0:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	apiV1.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.Use(middleware.Logging(log))
	r.Use(middleware.Error(log))

	return r
}
