package response

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// Error represents the structure of an error message
type Error struct {
	Error   bool   `json:"error"`
	Code    int    `json:"statusCode"`
	Message string `json:"message"`
}

// LogErrorResponse - logs an errornous server server at level debug
// on standard logger
func LogErrorResponse(r *http.Request, err error, code int, message string) {
	logrus.WithFields(logrus.Fields{
		"host":       r.Host,
		"address":    r.RemoteAddr,
		"method":     r.Method,
		"requestURI": r.RequestURI,
		"proto":      r.Proto,
		"useragent":  r.UserAgent(),
	}).WithError(err).Debug(message)
}
