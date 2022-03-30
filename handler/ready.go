package handler

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/julienschmidt/httprouter"
	"jhr.com/apirelay/global"
	"jhr.com/apirelay/model"
	"jhr.com/apirelay/util"
)

func HandleReady(rw http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	u, err := url.Parse(global.ServerConfig.Service.Host)
	if err != nil {
		return err
	}
	if c, err := net.DialTimeout("tcp", u.Host, time.Second*3); err != nil {
		return err
	} else {
		defer c.Close()
	}
	util.WriteResponse(rw, model.ResultMap[model.SUCCESS])
	return nil
}
