package dummy

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/harriklein/pBE/pBEServer/db"
	"github.com/harriklein/pBE/pBEServer/log"
	"github.com/harriklein/pBE/pBEServer/utils"
)

// ErrDummyNotFound is an error raised when a dummy can not be found in the database
var ErrDummyNotFound = fmt.Errorf("Dummy not found")

const (
	cmdDummyInsert       = "INSERT INTO Dummy (Description, _ID) VALUES (?,    UUID_TO_BIN(?));"
	cmdDummyDelete       = "DELETE FROM Dummy                   WHERE  (   ID=UUID_TO_BIN(?));"
	cmdDummyUpdate       = "UPDATE      Dummy SET Description=? WHERE  (   ID=UUID_TO_BIN(?));"
	cmdDummySelect       = "SELECT BIN_TO_UUID(ID) AS ID, Description FROM Dummy"
	cmdDummySelectFilter = " WHERE ID=UUID_TO_BIN(?)"
)

// Get gets rows
var mysqlRead = func(oDummies *TDummies, aID string, aWithCurlyBraces bool) *utils.TResponseError {

	log.Log.Debugln(" SELECT:")
	log.Log.Debugln("   ID  : ", aID)

	_cmdDummySelect := cmdDummySelect
	if aID != "" {
		_cmdDummySelect = cmdDummySelect + cmdDummySelectFilter
	}

	_stmtSelect, _error := db.ConnPBE.Prepare(_cmdDummySelect)
	if _error != nil {
		return utils.NewResponseError(http.StatusInternalServerError, _error.Error())
	}
	defer _stmtSelect.Close()

	var _rows *sql.Rows
	var _errorQuery error

	if aID == "" {
		_rows, _errorQuery = _stmtSelect.Query()
	} else {
		_rows, _errorQuery = _stmtSelect.Query(aID)
	}
	if _errorQuery != nil {
		return utils.NewResponseError(http.StatusInternalServerError, _errorQuery.Error())
	}

	for _rows.Next() {
		_dummy := &TDummy{}

		_error := _rows.Scan(&_dummy.ID, &_dummy.Description)

		// TODO: make it better... maybe
		if aWithCurlyBraces {
			if !utils.UUIDHasCurlyBraces(_dummy.ID.String()) {
				_dummy.ID = "{" + _dummy.ID + "}"
			}
		}

		if _error != nil {
			return utils.NewResponseError(http.StatusInternalServerError, _error.Error())
		}

		*oDummies = append(*oDummies, _dummy)
	}

	return nil
}

// Create into tha database
var mysqlCreate = func(oDummy *TDummy) *utils.TResponseError {

	log.Log.Debugln(" INSERT:")
	log.Log.Debugln("   ID  : ", oDummy.ID)
	log.Log.Debugln("   Desc: ", oDummy.Description)

	_stmtInsert, _error := db.ConnPBE.Prepare(cmdDummyInsert)
	if _error != nil {
		return utils.NewResponseError(http.StatusInternalServerError, _error.Error())
	}
	defer _stmtInsert.Close()

	_result, _error := _stmtInsert.Exec(oDummy.Description, oDummy.ID)
	if _error != nil {
		return utils.NewResponseError(http.StatusInternalServerError, _error.Error())
	}

	_rowsAffected, _error := _result.RowsAffected()
	if _error != nil {
		return utils.NewResponseError(http.StatusInternalServerError, _error.Error())
	}

	if _rowsAffected == 0 {
		return utils.NewResponseError(http.StatusBadRequest, fmt.Sprintf("%s not found", oDummy.ID))
	}

	log.Log.Debugln(" INSERTED ", _rowsAffected, " row(s)")

	return nil
}

// Update the database
var mysqlUpate = func(oDummy *TDummy) *utils.TResponseError {

	log.Log.Debugln(" UPDATE:")
	log.Log.Debugln("   ID  : ", oDummy.ID)
	log.Log.Debugln("   Desc: ", oDummy.Description)

	_stmtUpdate, _error := db.ConnPBE.Prepare(cmdDummyUpdate)
	if _error != nil {
		return utils.NewResponseError(http.StatusInternalServerError, _error.Error())
	}
	defer _stmtUpdate.Close()

	_result, _error := _stmtUpdate.Exec(oDummy.Description, oDummy.ID)
	if _error != nil {
		return utils.NewResponseError(http.StatusInternalServerError, _error.Error())
	}

	_rowsAffected, _error := _result.RowsAffected()
	if _error != nil {
		return utils.NewResponseError(http.StatusInternalServerError, _error.Error())
	}

	//if _rowsAffected == 0 {
	//	return utils.NewResponseError(http.StatusBadRequest, fmt.Sprintf("%s not found", oDummy.ID))
	//}

	log.Log.Debugln(" UPDATED ", _rowsAffected, " row(s)")

	return nil
}

// Delete from the database
var mysqlDelete = func(oDummy *TDummy) *utils.TResponseError {

	log.Log.Debugln(" DELETE:")
	log.Log.Debugln("   ID  : ", oDummy.ID)
	log.Log.Debugln("   Desc: ", oDummy.Description)

	_stmtDelete, _error := db.ConnPBE.Prepare(cmdDummyDelete)
	if _error != nil {
		return utils.NewResponseError(http.StatusInternalServerError, _error.Error())
	}
	defer _stmtDelete.Close()

	_result, _error := _stmtDelete.Exec(oDummy.ID)
	if _error != nil {
		return utils.NewResponseError(http.StatusInternalServerError, _error.Error())
	}

	_rowsAffected, _error := _result.RowsAffected()
	if _error != nil {
		return utils.NewResponseError(http.StatusInternalServerError, _error.Error())
	}

	//if _rowsAffected == 0 {
	//	return utils.NewResponseError(http.StatusBadRequest, fmt.Sprintf("%s not found", oDummy.ID))
	//}

	log.Log.Debugln(" DELETED ", _rowsAffected, " row(s)")

	return nil

}

// ApplyUpdates applies all updates defined in the Action field
var mysqlApplyUpdates = func(oDummies *TDummies, aDefaultAction string, aReplaceActionsByDefaultAction bool) (int64, int64, int64, *utils.TResponseError) {

	var _errorReason string
	var _countInserted, _countUpdated, _countDeleted int64

	_countInserted = 0
	_countUpdated = 0
	_countDeleted = 0

	// region START TRANSACTION =================================
	_tx, _errorBegin := db.ConnPBE.Begin()
	if _errorBegin != nil {
		return 0, 0, 0, utils.NewResponseError(http.StatusInternalServerError, _errorBegin.Error())
	}
	// endregion ================================================

	// region PREPARE QUERIES -----------------------------------
	// Prepare the queries in the transaction context
	_stmtInsert, _errorPrepare := _tx.Prepare(cmdDummyInsert)
	if _errorPrepare != nil {
		_errorReason = _errorPrepare.Error()
		if _errorRollback := _tx.Rollback(); _errorRollback != nil {
			_errorReason = _errorReason + " && " + _errorRollback.Error()
		}
		return 0, 0, 0, utils.NewResponseError(http.StatusInternalServerError, _errorReason)
	}
	defer _stmtInsert.Close()

	_stmtUpdate, _errorPrepare := _tx.Prepare(cmdDummyUpdate)
	if _errorPrepare != nil {
		_errorReason = _errorPrepare.Error()
		if _errorRollback := _tx.Rollback(); _errorRollback != nil {
			_errorReason = _errorReason + " && " + _errorRollback.Error()
		}
		return 0, 0, 0, utils.NewResponseError(http.StatusInternalServerError, _errorReason)
	}
	defer _stmtUpdate.Close()

	_stmtDelete, _errorPrepare := _tx.Prepare(cmdDummyDelete)
	if _errorPrepare != nil {
		_errorReason = _errorPrepare.Error()
		if _errorRollback := _tx.Rollback(); _errorRollback != nil {
			_errorReason = _errorReason + " && " + _errorRollback.Error()
		}
		return 0, 0, 0, utils.NewResponseError(http.StatusInternalServerError, _errorReason)
	}
	defer _stmtDelete.Close()
	// endregion -----------------------------------------------

	for _, _dummy := range *oDummies {

		// Get the action from Action field
		_action := aDefaultAction
		if !aReplaceActionsByDefaultAction {
			if len(_dummy.Action) > 0 {
				_action = strings.ToUpper(_dummy.Action[:1])
			}
		}

		// region EXECUTE QUERIES ------------------------------
		var _errorExec error
		var _result sql.Result
		switch _action {
		case "I":
			_result, _errorExec = _stmtInsert.Exec(_dummy.Description, _dummy.ID)
		case "U":
			_result, _errorExec = _stmtUpdate.Exec(_dummy.Description, _dummy.ID)
		case "D":
			_result, _errorExec = _stmtDelete.Exec(_dummy.ID)
		default:
			_errorExec = errors.New("Action not found")
		}

		if _errorExec != nil {
			_errorReason = _errorExec.Error()
			if _errorRollback := _tx.Rollback(); _errorRollback != nil {
				_errorReason = _errorReason + " && " + _errorRollback.Error()
			}
			return 0, 0, 0, utils.NewResponseError(http.StatusInternalServerError, _errorReason)
		}

		_rowsAffected, _error := _result.RowsAffected()
		if _error != nil {
			_errorReason = _errorExec.Error()
			if _errorRollback := _tx.Rollback(); _errorRollback != nil {
				_errorReason = _errorReason + " && " + _errorRollback.Error()
			}
			return 0, 0, 0, utils.NewResponseError(http.StatusInternalServerError, _errorReason)
		}

		switch _action {
		case "I":
			_countInserted += _rowsAffected
		case "U":
			_countUpdated += _rowsAffected
		case "D":
			_countDeleted += _rowsAffected
		}

		// DISABLED!!! Changed to "< 0" instead of "== 0"
		if _rowsAffected < 0 {
			_errorReason = fmt.Sprintf("%s not found", _dummy.ID)
			if _errorRollback := _tx.Rollback(); _errorRollback != nil {
				_errorReason = _errorReason + " && " + _errorRollback.Error()
			}
			return 0, 0, 0, utils.NewResponseError(http.StatusBadRequest, _errorReason)
		}

		// endregion -------------------------------------------

	}

	// region COMMIT TRANSACTION =================================
	_errorCommit := _tx.Commit()
	if _errorCommit != nil {
		_errorReason = _errorCommit.Error()
		if _errorRollback := _tx.Rollback(); _errorRollback != nil {
			_errorReason = _errorReason + " && " + _errorRollback.Error()
		}
		return 0, 0, 0, utils.NewResponseError(http.StatusInternalServerError, _errorReason)
	}
	// endregion ===============================================

	return _countInserted, _countUpdated, _countDeleted, nil
}
