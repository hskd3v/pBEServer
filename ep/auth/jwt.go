package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/harriklein/pBE/pBEServer/utils"
)

// TInfoClaims struct
type TInfoClaims struct {
	User   string `json:"user,omitempty"`
	System string `json:"system,omitempty"`
}

// TJWTClaim adds email as a claim to the token
type TJWTClaim struct {
	TInfoClaims
	jwt.StandardClaims
}

// GenerateJWT generates token
func GenerateJWT(aInfo TInfoClaims) (string, error) {

	_claims := &TJWTClaim{
		TInfoClaims: aInfo,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(cfgJWTExpirationHours)).Unix(),
			Issuer:    cfgJWTIssuer,
		},
	}

	_token := jwt.NewWithClaims(jwt.SigningMethodHS256, _claims)

	_signedToken, _error := _token.SignedString([]byte(cfgJWTSecretKey))

	return _signedToken, _error
}

// ValidateJWT validates JWT
func ValidateJWT(aSignedToken string) (*TJWTClaim, error) {

	_claims := &TJWTClaim{}

	_splitToken := strings.Split(pSignedToken, "Bearer")
	if len(_splitToken) != 2 {
		return _claims, errors.New("Invalid Token")
	}
	pSignedToken = strings.TrimSpace(_splitToken[1])

	_, _error := jwt.ParseWithClaims(pSignedToken, _claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfgJWTSecretKey), nil
	})

	if _error != nil {
		return _claims, _error
	}

	if _claims.ExpiresAt < time.Now().Local().Unix() {
		return _claims, errors.New("Token is expired")
	}

	return _claims, nil
}

// CheckJWT is a middleware which checks the token
func CheckJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(aResponse http.ResponseWriter, aRequest *http.Request) {
		_claims, _error := ValidateJWT(aRequest.Header.Get("Authorization"))
		if _error != nil {
			utils.NewResponseError(http.StatusUnauthorized, _error.Error()).ToJSON(aResponse)
			return
		}

		// Access context values in handlers like this
		// _jwtclaims := pRequest.Context().Value("jwtclaims")
		_ctx := context.WithValue(aRequest.Context(), "jwtclaims", _claims.TInfoClaims)

		next.ServeHTTP(aResponse, aRequest.WithContext(_ctx))
	})
}
