package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"jhr.com/apirelay/model"
	"jhr.com/apirelay/util"
)

func HandleHealth(rw http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	util.WriteResponse(rw, model.ResultMap[model.SUCCESS])
	return nil
}
