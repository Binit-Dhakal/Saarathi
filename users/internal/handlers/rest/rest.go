package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/cookies"
	"github.com/Binit-Dhakal/Saarathi/pkg/rest/httpx"
	"github.com/Binit-Dhakal/Saarathi/pkg/rest/jsonutil"
	"github.com/Binit-Dhakal/Saarathi/users/internal/application"
	"github.com/Binit-Dhakal/Saarathi/users/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/users/internal/dto"
)

type UserHandler struct {
	authApp        application.AuthService
	tokenApp       application.TokenService
	jsonReader     *jsonutil.Reader
	jsonWriter     *jsonutil.Writer
	errorResponder httpx.ErrorResponder
}

func NewUserHandler(authApp application.AuthService, tokenApp application.TokenService, jsonReader *jsonutil.Reader, jsonWriter *jsonutil.Writer, errorResponder httpx.ErrorResponder) *UserHandler {
	return &UserHandler{
		authApp:        authApp,
		tokenApp:       tokenApp,
		jsonReader:     jsonReader,
		jsonWriter:     jsonWriter,
		errorResponder: errorResponder,
	}
}

func (u *UserHandler) CreateRiderHandler(w http.ResponseWriter, r *http.Request) {
	var dst dto.RiderRegistrationDTO
	err := u.jsonReader.DecodeJSONStrict(w, r, &dst)
	if err != nil {
		u.errorResponder.BadRequest(w, r, err)
		return
	}

	userID, err := u.authApp.RegisterRider(&dst)
	if err != nil {
		fmt.Println(err)
		return
	}

	token, err := u.tokenApp.GenerateAccessAndRefreshTokens(userID, domain.RoleRider)
	if err != nil {
		fmt.Println(err)
		return
	}

	cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    token.RefreshToken,
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	}

	err = cookies.Write(w, cookie)
	if err != nil {
		u.errorResponder.ServerError(w, r, err)
		return
	}

	err = u.jsonWriter.JSON(w, 201, map[string]string{"access_token": token.AccessToken})
	if err != nil {
		u.errorResponder.ServerError(w, r, err)
		return
	}
}

func (u *UserHandler) CreateDriverHandler(w http.ResponseWriter, r *http.Request) {
	var dst dto.DriverRegistrationDTO
	err := u.jsonReader.DecodeJSONStrict(w, r, &dst)
	if err != nil {
		u.errorResponder.BadRequest(w, r, err)
		return
	}

	userID, err := u.authApp.RegisterDriver(&dst)
	if err != nil {
		fmt.Println(err)
		return
	}

	token, err := u.tokenApp.GenerateAccessAndRefreshTokens(userID, domain.RoleRider)
	if err != nil {
		fmt.Println(err)
		return
	}

	cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    token.RefreshToken,
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	}

	err = cookies.Write(w, cookie)
	if err != nil {
		u.errorResponder.ServerError(w, r, err)
		return
	}

	err = u.jsonWriter.JSON(w, 201, map[string]string{"access_token": token.AccessToken})
	if err != nil {
		u.errorResponder.ServerError(w, r, err)
		return
	}
}

func (u *UserHandler) CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var userInput dto.LoginRequestDTO
	err := u.jsonReader.DecodeJSONStrict(w, r, &userInput)
	if err != nil {
		u.errorResponder.ServerError(w, r, err)
		return
	}

	userID, err := u.authApp.CreateAuthenticationToken(&userInput)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	token, err := u.tokenApp.GenerateAccessAndRefreshTokens(userID, domain.RoleRider)
	if err != nil {
		fmt.Println(err)
		return
	}

	cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    token.RefreshToken,
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	}

	err = cookies.Write(w, cookie)
	if err != nil {
		u.errorResponder.ServerError(w, r, err)
		return
	}

	err = u.jsonWriter.JSON(w, 201, map[string]string{"access_token": token.AccessToken})
	if err != nil {
		u.errorResponder.ServerError(w, r, err)
		return
	}
}

func (u *UserHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := cookies.Read(r, "refreshToken")
	if err != nil {
		u.errorResponder.BadRequest(w, r, err)
		return
	}

	t, err := u.tokenApp.ValidateRefreshToken(refreshToken)
	if err != nil {
		u.errorResponder.BadRequest(w, r, err)
		return
	}

	err = u.tokenApp.RevokeRefreshToken(refreshToken)
	if err != nil {
		u.errorResponder.ServerError(w, r, err)
		return
	}

	token, err := u.tokenApp.GenerateAccessAndRefreshTokens(t.UserID, t.RoleID)
	if err != nil {
		fmt.Println(err)
		return
	}

	cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    token.RefreshToken,
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	}

	err = cookies.Write(w, cookie)
	if err != nil {
		u.errorResponder.ServerError(w, r, err)
		return
	}

	err = u.jsonWriter.JSON(w, 201, map[string]string{"access_token": token.AccessToken})
	if err != nil {
		u.errorResponder.ServerError(w, r, err)
		return
	}

}
