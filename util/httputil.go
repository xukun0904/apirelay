package util

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"

	"jhr.com/apirelay/global"
	"jhr.com/apirelay/model"
)

func GetHttpStatusCode(code string) int {
	if strings.Contains(code, "A") {
		return http.StatusBadRequest
	} else if strings.ContainsAny(code, "BC") {
		return http.StatusInternalServerError
	} else {
		return http.StatusOK
	}
}

func WriteResponse(rw http.ResponseWriter, rc model.Result) error {
	var b []byte
	var err error
	if global.ServerConfig.Api.ResponseContentType == model.RESPONSE_APPLICATION_XML {
		// 转换成xml字符串
		if b, err = xml.Marshal(&rc); err != nil {
			return err
		}
	} else {
		// 默认转换成json字符串
		if b, err = json.Marshal(&rc); err != nil {
			return err
		}
	}
	var responseContentType = model.ResponseContentTypeName[global.ServerConfig.Api.ResponseContentType]
	rw.Header().Add("Content-Type", responseContentType+"; charset=utf-8")
	rw.WriteHeader(GetHttpStatusCode(rc.Code))
	// 返回响应值
	rw.Write(b)
	return err
}
