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

func (u *UserHandler) authCookieGenerator(w http.ResponseWriter, r *http.Request, userID string, roleID int, message string) error {
	token, err := u.tokenApp.GenerateAccessAndRefreshTokens(userID, roleID)
	if err != nil {
		return err
	}

	refreshCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    token.RefreshToken,
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Domain:   ".saarathi.com",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Secure:   false,
	}

	err = cookies.Write(w, refreshCookie)
	if err != nil {
		return err
	}

	accessCookie := http.Cookie{
		Name:     "accessToken",
		Value:    token.AccessToken,
		HttpOnly: true,
		Expires:  time.Now().Add(1 * 24 * time.Hour), // dev - temporary
		Domain:   ".saarathi.com",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Secure:   false,
	}

	err = cookies.Write(w, accessCookie)
	if err != nil {
		return err
	}

	err = u.jsonWriter.JSON(w, 201, map[string]string{"message": message})
	if err != nil {
		return err
	}

	return nil
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

	err = u.authCookieGenerator(w, r, userID, domain.RoleRider, "Rider Created successfully")
	if err != nil {
		u.errorResponder.ServerError(w, r, err)
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

	err = u.authCookieGenerator(w, r, userID, domain.RoleDriver, "Driver created successfully")
	if err != nil {
		u.errorResponder.ServerError(w, r, err)
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

	var role int
	switch userInput.Role {
	case "rider":
		role = domain.RoleRider
	case "driver":
		role = domain.RoleDriver
	default:
		u.errorResponder.BadRequest(w, r, fmt.Errorf("Role is not defined(rider/driver)- %v", userInput.Role))
		return
	}

	err = u.authCookieGenerator(w, r, userID, role, "Logged in successfully")
	if err != nil {
		u.errorResponder.ServerError(w, r, err)
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

	err = u.authCookieGenerator(w, r, t.UserID, t.RoleID, "Token Refreshed successfully")
	if err != nil {
		u.errorResponder.ServerError(w, r, err)
	}
}
