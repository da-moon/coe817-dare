package middlewares

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// Log - uses logrus library to handles logging of common
// request fields
func Log(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"host":       r.Host,
			"address":    r.RemoteAddr,
			"method":     r.Method,
			"requestURI": r.RequestURI,
			"proto":      r.Proto,
			"useragent":  r.UserAgent(),
		}).Info("HTTP request information")
		next.ServeHTTP(w, r)
	})
}
