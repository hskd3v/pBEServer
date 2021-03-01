package docs

import (
	"net/http"
	"strings"
)

// TDocHandler defines the type of the documentation handler structure
type TDocHandler struct {
}

// NewDocHandler initializes a documentation handler with the given params
func NewDocHandler() *TDocHandler {
	return &TDocHandler{}
}

// Get gets the entire dummy list
func (oHandler *TDocHandler) Get(aResponse http.ResponseWriter, aRequest *http.Request) {

	const _wwwPath = "./www/"
	_url := strings.TrimPrefix(aRequest.URL.Path, "/docs")

	// It is necessary because by default it redirects to "/docs/" when request "/docs" or "/docs/index.html"
	if aRequest.URL.Path == "/docs/" {
		http.ServeFile(aResponse, aRequest, _wwwPath+"index.html")
	} else {
		http.ServeFile(aResponse, aRequest, _wwwPath+_url)
	}

}
