package dummy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/harriklein/pBE/pBEServer/log"
	"github.com/harriklein/pBE/pBEServer/utils"
)

// swagger:route GET /dummies Dummy dummyList
// Return a dummy list from the database
// responses:
//	200: dummyListResponse

// Read gets the entire dummy list
func handleRead(aResponse http.ResponseWriter, aRequest *http.Request) {

	log.Log.Debugln("Get dummies.")
	aResponse.Header().Add("Content-Type", "application/json")
	_withBrackets := aRequest.Header.Get("x-id-braces") != ""

	// region GET PARAMS ---------------------------------------
	_params := mux.Vars(aRequest)
	var _id string
	_id = _params["id"]
	// endregion -----------------------------------------------

	var _dummies TDummies
	dbRead(&_dummies, _id, _withBrackets)

	//pResponse.Header().Add("Content-Type", "application/json")
	_error := utils.ToJSON(_dummies, aResponse)
	if _error != nil {
		utils.NewResponseError(http.StatusInternalServerError, "Deserializing Dummy - "+_error.Error()).ToJSON(aResponse)
		return
	}

}

// CreateOrApplyUpdates handles a single or bulk create(post) request
// In case of a bulk request using Action fields, it will be an Apply Update procedure
func handleCreateOrApplyUpdates(aResponse http.ResponseWriter, aRequest *http.Request) {

	log.Log.Debug("CreateOrApplyUpdates dummies")

	// region GET PARAMS ---------------------------------------
	_params := mux.Vars(aRequest)
	var _id string
	_id = _params["id"]
	// endregion -----------------------------------------------

	// Read all body and put it into _body. It is necessary because we have to read only one time,
	_body, _error := ioutil.ReadAll(aRequest.Body)
	if _error != nil {
		utils.NewResponseError(http.StatusInternalServerError, "Reading body - "+_error.Error()).ToJSON(aResponse)
		return
	}

	// log.Log.Debugln(string(_body)) // TODO: REMOVE IT

	// Check if it is an array in order to identify an ApplyUpdates request.
	// Otherwise, it is a single create request
	_isArray := utils.IsArray(&_body)

	// region VALIDATE ID: 2 = Array is not allowed in request with ID in URL
	if _id != "" {
		if _isArray {
			utils.NewResponseError(http.StatusBadRequest, "Invalid content: Array not allowed for this request").ToJSON(aResponse)
			return
		}
	}
	// endregion -----------------------------------------------

	if !_isArray {
		// If it is not an array, attempt to Insert
		log.Log.Debugln("Deserializing Object")

		_dummy := &TDummy{}

		_error := json.Unmarshal(_body, _dummy)
		if _error != nil {
			utils.NewResponseError(http.StatusInternalServerError, "Deserializing Dummy - "+_error.Error()).ToJSON(aResponse)
			return
		}

		// region VALIDATE ID: 3 = Consistency between URL and Body
		if _id != "" {
			if _dummy.ID == "" {
				_error := _dummy.ID.UnmarshalText([]byte(_id))
				if _error != nil {
					utils.NewResponseError(http.StatusBadRequest, "unable to convert ID").ToJSON(aResponse)
					return
				}
			} else if _id != _dummy.ID.String() {
				utils.NewResponseError(http.StatusBadRequest, "Mismatch IDs").ToJSON(aResponse)
				return
			}
		}
		// endregion -----------------------------------------

		_respError := dbCreate(_dummy)
		if _respError != nil {
			_respError.ToJSON(aResponse)
			return
		}

		utils.NewResponse(http.StatusCreated, fmt.Sprintf("Created %s", _dummy.ID), nil).ToJSON(aResponse)

	} else {
		// If it is an array, get the action in Action field and attempt to Insert, Update or Delete
		// If there is any error, rollback all actions
		log.Log.Debugln("Deserializing Array")

		_dummies := &TDummies{}

		_error := json.Unmarshal(_body, _dummies)
		if _error != nil {
			utils.NewResponseError(http.StatusInternalServerError, "Deserializing Dummies - "+_error.Error()).ToJSON(aResponse)
			return
		}

		_cIns, _cUpd, _cDel, _respError := dbApplyUpdates(_dummies, "I", false)
		if _respError != nil {
			_respError.ToJSON(aResponse)
			return
		}

		utils.NewResponse(http.StatusOK, fmt.Sprintf("%d requests >>> %d update(s) applied: %d insert(s), %d update(s) and %d delete(s)", len(*_dummies), (_cIns+_cUpd+_cDel), _cIns, _cUpd, _cDel), nil).ToJSON(aResponse)
	}

}

// Update handles a single or bulk update(put) request
func handleUpdate(aResponse http.ResponseWriter, aRequest *http.Request) {

	log.Log.Debug("Edit dummies")

	// region GET PARAMS ---------------------------------------
	_params := mux.Vars(aRequest)
	var _id string
	_id = _params["id"]
	// endregion -----------------------------------------------

	// Read all body and put it into _body. It is necessary because we have to read only one time,
	_body, _error := ioutil.ReadAll(aRequest.Body)
	if _error != nil {
		utils.NewResponseError(http.StatusInternalServerError, "Reading body - "+_error.Error()).ToJSON(aResponse)
		return
	}

	// Check if it is an array in order to identify an ApplyUpdates request.
	// Otherwise, it is a single request
	_isArray := utils.IsArray(&_body)

	// region VALIDATE ID: 2 = Array is not allowed in request with ID in URL
	if _id != "" {
		if _isArray {
			utils.NewResponseError(http.StatusBadRequest, "Invalid content: Array not allowed for this request").ToJSON(aResponse)
			return
		}
	}
	// endregion -----------------------------------------------

	if !_isArray {

		log.Log.Debugln("Deserializing Object")

		_dummy := &TDummy{}

		_error := json.Unmarshal(_body, _dummy)
		if _error != nil {
			utils.NewResponseError(http.StatusInternalServerError, "Deserializing Dummy - "+_error.Error()).ToJSON(aResponse)
			return
		}

		// region VALIDATE ID: 3 = Consistency between URL and Body
		if _id != "" {
			if _dummy.ID == "" {
				_error := _dummy.ID.UnmarshalText([]byte(_id))
				if _error != nil {
					utils.NewResponseError(http.StatusBadRequest, "unable to convert ID").ToJSON(aResponse)
					return
				}
			} else if _id != _dummy.ID.String() {
				utils.NewResponseError(http.StatusBadRequest, "Mismatch IDs").ToJSON(aResponse)
				return
			}
		}
		// endregion -----------------------------------------

		_respError := dbUpdate(_dummy)
		if _respError != nil {
			_respError.ToJSON(aResponse)
			return
		}

		utils.NewResponse(http.StatusOK, fmt.Sprintf("Edited %s", _dummy.ID), nil).ToJSON(aResponse)

	} else {
		// If it is an array, edit all records
		// If there is any error, rollback all actions
		log.Log.Debugln("Deserializing Array")

		_dummies := &TDummies{}

		_error := json.Unmarshal(_body, _dummies)
		if _error != nil {
			utils.NewResponseError(http.StatusInternalServerError, "Deserializing Dummies - "+_error.Error()).ToJSON(aResponse)
			return
		}

		_cIns, _cUpd, _cDel, _respError := dbApplyUpdates(_dummies, "U", true) // Force UPDATE action in all records
		if _respError != nil {
			_respError.ToJSON(aResponse)
			return
		}

		utils.NewResponse(http.StatusOK, fmt.Sprintf("%d requests >>> %d update(s) applied: %d insert(s), %d update(s) and %d delete(s)", len(*_dummies), (_cIns+_cUpd+_cDel), _cIns, _cUpd, _cDel), nil).ToJSON(aResponse)
	}

}

// Delete handles a single or bulk delete request
func handleDelete(aResponse http.ResponseWriter, aRequest *http.Request) {

	log.Log.Debug("Remove dummies")

	// region GET PARAMS ---------------------------------------
	_params := mux.Vars(aRequest)
	var _id string
	_id = _params["id"]
	// endregion -----------------------------------------------

	// Read all body and put it into _body. It is necessary because we have to read only one time,
	_body, _error := ioutil.ReadAll(aRequest.Body)
	if _error != nil {
		utils.NewResponseError(http.StatusInternalServerError, "Reading body - "+_error.Error()).ToJSON(aResponse)
		return
	}

	// Check if it is an array in order to identify an ApplyUpdates request.
	// Otherwise, it is a single request
	_isArray := utils.IsArray(&_body)

	// region VALIDATE ID: 2 = Array is not allowed in request with ID in URL
	if _id != "" {
		if _isArray {
			utils.NewResponseError(http.StatusBadRequest, "Invalid content: Array not allowed for this request").ToJSON(aResponse)
			return
		}
	}
	// endregion -----------------------------------------------

	if !_isArray {

		log.Log.Debugln("Deserializing Object")

		_dummy := &TDummy{}

		_error := json.Unmarshal(_body, _dummy)
		if _error != nil {
			utils.NewResponseError(http.StatusInternalServerError, "Deserializing Dummy - "+_error.Error()).ToJSON(aResponse)
			return
		}

		// region VALIDATE ID: 3 = Consistency between URL and Body
		if _id != "" {
			if _dummy.ID == "" {
				_error := _dummy.ID.UnmarshalText([]byte(_id))
				if _error != nil {
					utils.NewResponseError(http.StatusBadRequest, "unable to convert ID").ToJSON(aResponse)
					return
				}
			} else if _id != _dummy.ID.String() {
				utils.NewResponseError(http.StatusBadRequest, "Mismatch IDs").ToJSON(aResponse)
				return
			}
		}
		// endregion -----------------------------------------

		_respError := dbDelete(_dummy)
		if _respError != nil {
			_respError.ToJSON(aResponse)
			return
		}

		utils.NewResponse(http.StatusOK, fmt.Sprintf("Deleted %s", _dummy.ID), nil).ToJSON(aResponse)

	} else {
		// If it is an array, delete all records
		// If there is any error, rollback all actions
		log.Log.Debugln("Deserializing Array")

		_dummies := &TDummies{}

		_error := json.Unmarshal(_body, _dummies)
		if _error != nil {
			utils.NewResponseError(http.StatusInternalServerError, "Deserializing Dummies - "+_error.Error()).ToJSON(aResponse)
			return
		}

		_cIns, _cUpd, _cDel, _respError := dbApplyUpdates(_dummies, "D", true) // Force DELETE action in all records
		if _respError != nil {
			_respError.ToJSON(aResponse)
			return
		}

		utils.NewResponse(http.StatusOK, fmt.Sprintf("%d requests >>> %d update(s) applied: %d insert(s), %d update(s) and %d delete(s)", len(*_dummies), (_cIns+_cUpd+_cDel), _cIns, _cUpd, _cDel), nil).ToJSON(aResponse)
	}

}
