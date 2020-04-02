package handlers

import (
	"net/http"
)

// Preflight - prefilight is a handler used to reply to
// OPTIONS request
// gorila mux example
// router.Methods("OPTIONS").HandlerFunc(middlewares.Cors(Preflight))
func Preflight(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}
