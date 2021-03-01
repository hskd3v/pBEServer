package auth

import (
	"net/http"
	"strconv"

	"github.com/harriklein/pBE/pBEServer/log"
	"github.com/harriklein/pBE/pBEServer/utils"
)

// Read gets the entire dummy list
func handleSignUp(aResponse http.ResponseWriter, aRequest *http.Request) {
	log.Log.Debugln("try handleSignUp")
	utils.NewResponseError(http.StatusNotImplemented, aRequest.URL.String()+" not Implemented").ToJSON(aResponse)
	return
}

// Read gets the entire dummy list
func handleLoginTest(aResponse http.ResponseWriter, aRequest *http.Request) {
	_jwtclaims := aRequest.Context().Value("jwtclaims")

	utils.NewResponse(http.StatusOK, "OK!", _jwtclaims).ToJSON(aResponse)
}

// Read gets the entire dummy list
func handleLoginBasic(aResponse http.ResponseWriter, aRequest *http.Request) {

	aResponse.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

	_auth := aRequest.Header.Get("Authorization")
	_username, _password, _authOK := aRequest.BasicAuth()

	log.Log.Debugf("Auth Type=%s, User=%s, Pwd=%s, OK=%s\n", _auth, _username, _password, strconv.FormatBool(_authOK))

	if _authOK == false {
		utils.NewResponseError(http.StatusUnauthorized, "Not authorized").ToJSON(aResponse)
		return
	}

	if _username != "api" || _password != _username {
		utils.NewResponseError(http.StatusUnauthorized, "Not authorized").ToJSON(aResponse)
		return
	}

	claims := TInfoClaims{
		User:   "api",
		System: "system",
	}

	_token, _error := GenerateJWT(claims)
	if _error != nil {
		utils.NewResponseError(http.StatusUnauthorized, "Not authorized").ToJSON(aResponse)
		return
	}

	var ResponseLoginBasicData struct {
		Token string `json:"token"`
	}
	ResponseLoginBasicData.Token = _token

	utils.NewResponse(http.StatusOK, "Login", ResponseLoginBasicData).ToJSON(aResponse)
}

// Read gets the entire dummy list
func handleLogout(aResponse http.ResponseWriter, aRequest *http.Request) {
	utils.NewResponseError(http.StatusNotImplemented, aRequest.URL.String()+" not Implemented").ToJSON(aResponse)
	return
}
