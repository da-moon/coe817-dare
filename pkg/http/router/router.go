package router

import (
	handlers "github.com/da-moon/dare-cli/pkg/http/handlers"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
)

// JSON2 -
type JSON2 struct {
	Endpoint  string
	Namespace string
	Handler   interface{}
}

// GenerateRPC2Routes ...
func GenerateRPC2Routes(routes []JSON2) *mux.Router {
	router := mux.NewRouter().PathPrefix("/").Subrouter()
	// router.
	// 	Methods(http.MethodOptions).
	// 	HandlerFunc(middlewares.Cors(handlers.Preflight))
	// // healthcheck
	router.
		Methods(http.MethodGet).
		Path("/healthcheck").
		HandlerFunc(handlers.HealthCheck)
	s := rpc.NewServer()
	s.RegisterCodec(json2.NewCodec(), "application/json")
	s.RegisterCodec(json2.NewCodec(), "application/json;charset=UTF-8")
	for _, route := range routes {
		s.RegisterService(route.Handler, route.Namespace)
		router.Handle(route.Endpoint, s)
	}
	return router
}

// Route - gorilla mux route wrapper
type Route struct {
	PathPrefix  string
	Name        string
	Method      string
	Path        string
	HandlerFunc http.HandlerFunc
	Queries     []string
}

// GenerateRoutes - Return a new Gorilla Mux Wrapper
func GenerateRoutes(routes []Route) *mux.Router {

	router := mux.NewRouter().PathPrefix("/").Subrouter()
	// router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		router.
			PathPrefix(route.PathPrefix).
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(route.HandlerFunc).
			Queries(route.Queries...)

	}

	return router
}
