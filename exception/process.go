package exception

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"jhr.com/apirelay/model"
	"jhr.com/apirelay/util"
)

type appHandler func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) error

func ErrWrapper(handler appHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
		defer func() {
			r := recover()
			if r != nil {
				zap.S().Error("Unkonw error: ", r)
				util.WriteResponse(rw, model.ResultMap[model.SERVER_ERROR])
			}
		}()
		err := handler(rw, r, params)
		if err != nil {
			if ce, ok := err.(CustomError); ok {
				util.WriteResponse(rw, ce.Result)
				return
			}
			zap.S().Error("Handing request error: ", err.Error())
			util.WriteResponse(rw, model.ResultMap[model.SERVER_ERROR])
		}
	}
}

func Cast(rc model.Result) CustomError {
	return CustomError{rc}
}
