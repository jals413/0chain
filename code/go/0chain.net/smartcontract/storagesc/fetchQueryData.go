package storagesc

import (
	"net/http"

	"0chain.net/core/common"
	"0chain.net/smartcontract/dbs/event"
)

func (srh *StorageRestHandler) getQueryData(w http.ResponseWriter, r *http.Request) {
	// read all data from query_data table and return
	edb := srh.GetQueryStateContext().GetEventDB()
	if edb == nil {
		common.Respond(w, r, nil, common.NewErrInternal("no db connection"))
		return
	}
	entity := r.URL.Query().Get("entity")
	fields := r.URL.Query().Get("fields")
	preload := r.URL.Query().Get("preload")
	var table interface{}
	switch entity {
	case "blobber":
		table = &event.Blobber{}
	case "Sharder":
		table = &event.Sharder{}
	case "miner":
		table = &event.Miner{}
	case "authorizer":
		table = &event.Authorizer{}
	case "validator":
		table = &event.Validator{}
	}

	result, err := edb.GetQueryData(preload, fields, table)
	if err != nil {
		common.Respond(w, r, nil, err)
		return
	}

	common.Respond(w, r, result, nil)
}
