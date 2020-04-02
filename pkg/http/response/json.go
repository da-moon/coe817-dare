package response

import (
	"net/http"
	"strconv"

	"github.com/da-moon/coe817-dare/pkg/jsonutil"
)

// WriteErrorJSON - logs and sends a json response to the client
// showing the error message
func WriteErrorJSON(w *http.ResponseWriter, r *http.Request, code int, message string) {
	// LogErrorResponse(r, err, code, message)
	response, err := jsonutil.EncodeJSON(Error{
		Error:   true,
		Code:    code,
		Message: message,
	})
	if err != nil {
		panic(err)
	}
	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")
	(*w).Header().Set("Content-Length", strconv.Itoa(len(response)))
	(*w).WriteHeader(code)
	if response != nil {
		(*w).Write(response)
		(*w).(http.Flusher).Flush()
	}
	return
}

// WriteSuccessfulJSON - logs and sends a new json response to the client
func WriteSuccessfulJSON(w *http.ResponseWriter, r *http.Request, data interface{}) {
	LogSuccessfulResponse(r, data)
	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")
	(*w).WriteHeader(http.StatusOK)
	response, err := jsonutil.EncodeJSON(data)
	if err != nil {
		panic(err)
	}
	(*w).Header().Set("Content-Length", strconv.Itoa(len(response)))
	if response != nil {
		(*w).Write(response)
		(*w).(http.Flusher).Flush()
	}
	return
}
