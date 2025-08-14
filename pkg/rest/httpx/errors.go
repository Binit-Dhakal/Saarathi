package httpx

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	log "github.com/Binit-Dhakal/Saarathi/pkg/logger"
	"github.com/Binit-Dhakal/Saarathi/pkg/rest/jsonutil"
	validator "github.com/Binit-Dhakal/Saarathi/pkg/rest/validators"
)

type ErrorResponder interface {
	InvalidAuthenticationToken(w http.ResponseWriter, r *http.Request)
	NotFound(w http.ResponseWriter, r *http.Request)
	AuthenticationRequired(w http.ResponseWriter, r *http.Request)
	MethodNotAllowed(w http.ResponseWriter, r *http.Request)
	FailedValidation(w http.ResponseWriter, r *http.Request, v validator.Validator)
	BadRequest(w http.ResponseWriter, r *http.Request, err error)
	ServerError(w http.ResponseWriter, r *http.Request, err error)
}

type errorResponderImpl struct {
	writer *jsonutil.Writer
	logger log.ILogger
}

var _ ErrorResponder = (*errorResponderImpl)(nil)

func NewErrorResponder(jsonWriter *jsonutil.Writer, logger log.ILogger) ErrorResponder {
	return &errorResponderImpl{
		writer: jsonWriter,
		logger: logger,
	}
}

func (e *errorResponderImpl) reportServerError(r *http.Request, err error) {
	var (
		message = err.Error()
		method  = r.Method
		url     = r.URL.String()
		trace   = string(debug.Stack())
	)

	e.logger.Error(message, err, "method", method, "url", url, "trace", trace)
}

func (e *errorResponderImpl) errorMessage(w http.ResponseWriter, r *http.Request, status int, message string, headers http.Header) {
	message = strings.ToUpper(message[:1]) + message[1:]

	err := e.writer.JSONWithHeaders(w, status, map[string]string{"Error": message}, headers)
	if err != nil {
		e.reportServerError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (e *errorResponderImpl) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	e.reportServerError(r, err)

	message := "The server encountered a problem and could not process your request"
	e.errorMessage(w, r, http.StatusInternalServerError, message, nil)
}

func (e *errorResponderImpl) NotFound(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	e.errorMessage(w, r, http.StatusNotFound, message, nil)
}

func (e *errorResponderImpl) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	e.errorMessage(w, r, http.StatusMethodNotAllowed, message, nil)
}

func (e *errorResponderImpl) BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	e.errorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
}

func (e *errorResponderImpl) FailedValidation(w http.ResponseWriter, r *http.Request, v validator.Validator) {
	err := e.writer.JSON(w, http.StatusUnprocessableEntity, v)
	if err != nil {
		e.ServerError(w, r, err)
	}
}

func (e *errorResponderImpl) InvalidAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	headers := make(http.Header)

	http.SetCookie(w, &http.Cookie{
		Name:     "session-token",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		MaxAge:   0,
	})

	e.errorMessage(w, r, http.StatusUnauthorized, "Invalid authentication token", headers)
}

func (e *errorResponderImpl) AuthenticationRequired(w http.ResponseWriter, r *http.Request) {
	e.errorMessage(w, r, http.StatusUnauthorized, "You must be authenticated to access this resource", nil)
}
