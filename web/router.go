package web

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"jhr.com/apirelay/exception"
	"jhr.com/apirelay/global"
	"jhr.com/apirelay/handler"
	"jhr.com/apirelay/model"
	"jhr.com/apirelay/util"
)

func RegisterRouter() *httprouter.Router {
	zap.S().Info("Start registering routes")
	router := httprouter.New()
	router.GET("/health", exception.ErrWrapper(handler.HandleHealth))
	router.GET("/ready", exception.ErrWrapper(handler.HandleReady))
	// 获取请求方法
	apiConfig := global.ServerConfig.Api
	requestMethod := model.RequestMethodName[apiConfig.RequestMethod]
	// 绑定自定义路径
	router.Handle(requestMethod, apiConfig.Path, exception.ErrWrapper(handler.ForwardRequest))
	zap.S().Infof("Registered route successfully. Api info: %+v", &apiConfig)
	// 404处理
	router.NotFound = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		util.WriteResponse(rw, model.ResultMap[model.NOT_FOUND])
	})
	// 405处理
	router.MethodNotAllowed = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		util.WriteResponse(rw, model.ResultMap[model.METHOD_NOT_ALLOWED])
	})
	return router
}
