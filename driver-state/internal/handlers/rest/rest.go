package rest

import (
	"errors"
	"net/http"

	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/application"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/dto"
	"github.com/Binit-Dhakal/Saarathi/pkg/rest/httpx"
	"github.com/Binit-Dhakal/Saarathi/pkg/rest/jsonutil"
)

type LocationHandler struct {
	locationSvc    application.LocationService
	jsonReader     *jsonutil.Reader
	jsonWriter     *jsonutil.Writer
	errorResponder httpx.ErrorResponder
}

func NewLocationHandler(locationSvc application.LocationService, jsonReader *jsonutil.Reader, jsonWriter *jsonutil.Writer, errorResponder httpx.ErrorResponder) *LocationHandler {
	return &LocationHandler{
		locationSvc:    locationSvc,
		jsonReader:     jsonReader,
		jsonWriter:     jsonWriter,
		errorResponder: errorResponder,
	}
}

func (l *LocationHandler) UpdateDriverLocation(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateDriverLocationRequest
	err := l.jsonReader.DecodeJSON(w, r, &req)
	if err != nil {
		l.errorResponder.BadRequest(w, r, err)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		l.errorResponder.BadRequest(w, r, errors.New("Authentication failed"))
		return
	}

	err = l.locationSvc.UpsertDriverLocation(&req, userID)
	if err != nil {
		l.errorResponder.ServerError(w, r, err)
		return
	}
}
