package rest

import (
	"fmt"
	"net/http"

	"github.com/Binit-Dhakal/Saarathi/pkg/rest/httpx"
	"github.com/Binit-Dhakal/Saarathi/pkg/rest/jsonutil"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/application"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/domain"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/dto"
)

type TripHandler struct {
	rideSvc        application.RideService
	routeSvc       application.RouteService
	jsonReader     *jsonutil.Reader
	jsonWriter     *jsonutil.Writer
	errorResponder httpx.ErrorResponder
}

func NewTripHandler(rideSvc application.RideService, routeSvc application.RouteService, jsonReader *jsonutil.Reader, jsonWriter *jsonutil.Writer, errorResponder httpx.ErrorResponder) *TripHandler {
	return &TripHandler{
		rideSvc:        rideSvc,
		routeSvc:       routeSvc,
		jsonReader:     jsonReader,
		jsonWriter:     jsonWriter,
		errorResponder: errorResponder,
	}
}

func (t *TripHandler) PreviewFare(w http.ResponseWriter, r *http.Request) {
	var locInfo *dto.FareEstimateRequest
	err := t.jsonReader.DecodeJSONStrict(w, r, &locInfo)
	if err != nil {
		t.errorResponder.BadRequest(w, r, err)
		return
	}

	source := domain.Coordinate{
		Lon: locInfo.PickUpLocation[0],
		Lat: locInfo.PickUpLocation[1],
	}
	destination := domain.Coordinate{
		Lon: locInfo.DropOffLocation[0],
		Lat: locInfo.DropOffLocation[1],
	}

	// get location distance and other info from OSRM backend
	route, err := t.routeSvc.GetRouteDetailFromOSRM(
		source, destination,
	)
	if err != nil {
		t.errorResponder.ServerError(w, r, err)
		return
	}

	var fareResponse dto.FareEstimateResponse
	fareResponse.Geometry = route.Geometry

	fares, fareID, err := t.rideSvc.EstimateFare(route)
	if err != nil {
		t.errorResponder.ServerError(w, r, err)
		return
	}

	fareResponse.FareID = fareID
	fareResponse.Fares = fares

	err = t.jsonWriter.JSON(w, http.StatusOK, fareResponse)
	if err != nil {
		t.errorResponder.ServerError(w, r, err)
		return
	}
}

func (t *TripHandler) ConfirmFare(w http.ResponseWriter, r *http.Request) {
	var confirmRequest dto.FareConfirmRequest
	err := t.jsonReader.DecodeJSON(w, r, &confirmRequest)
	if err != nil {
		t.errorResponder.BadRequest(w, r, err)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		t.errorResponder.BadRequest(w, r, fmt.Errorf("missing X-User-ID"))
		return
	}

	rideID, err := t.rideSvc.FareAcceptByRider(&confirmRequest, userID)
	if err != nil {
		t.errorResponder.ServerError(w, r, err)
		return
	}

	resp := dto.FareConfirmResponse{
		RideID: rideID,
	}

	err = t.jsonWriter.JSON(w, http.StatusOK, resp)
	if err != nil {
		t.errorResponder.ServerError(w, r, err)
		return
	}
}
