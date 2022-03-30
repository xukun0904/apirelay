package web

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"jhr.com/apirelay/global"
	"jhr.com/apirelay/model"
	"jhr.com/apirelay/util"
)

func RunServer(r *httprouter.Router) {
	port := global.ServerConfig.Port
	zap.S().Info("Start the Server. Port: ", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		zap.S().Panic("Failed to initialize server: ", err)
	}
}

func TestServer(r *httprouter.Router) {
	zap.S().Info("Start testing Server.")
	writer := httptest.NewRecorder()
	// 获取请求方法
	apiConfig := global.ServerConfig.Api
	requestMethod := model.RequestMethodName[apiConfig.RequestMethod]
	// 获取请求参数
	rawQuery, formData, headerMap, bodyData, path := getRequestParamter()
	// 封装请求Body
	body := getRequestBody(formData, bodyData)
	req, err := http.NewRequest(requestMethod, path, body)
	if err != nil {
		zap.S().Panic("Failed to test server: ", err)
	}
	// 添加Query参数
	req.URL.RawQuery = rawQuery.Encode()
	// 添加请求头
	req.Header.Set("Content-Type", model.RequestContentTypeName[apiConfig.RequestContentType])
	for k, v := range headerMap {
		req.Header.Set(k, v)
	}
	startTime := time.Now()
	// 测试服务
	r.ServeHTTP(writer, req)
	tc := time.Since(startTime)
	// 总耗时
	zap.S().Named("time").Info(tc)
	// 读取响应值
	b, err2 := io.ReadAll(writer.Body)
	if err2 != nil {
		zap.S().Panic("Failed to read response body: ", err2)
	}
	zap.S().Named("response").Info(string(b))
}

func getRequestBody(formData url.Values, bodyData map[string]string) *strings.Reader {
	apiConfig := global.ServerConfig.Api
	var bodyString string
	if apiConfig.RequestContentType == model.REQUEST_APPLICATION_FORM {
		bodyString = formData.Encode()
	} else if apiConfig.RequestContentType == model.REQUEST_APPLICATION_JSON {
		b, err := json.Marshal(bodyData)
		if err != nil {
			zap.S().Panic("Failed to marshal request body to json: ", err)
		}
		bodyString = string(b)
	} else if apiConfig.RequestContentType == model.REQUEST_APPLICATION_XML {
		b, err := xml.Marshal(model.StringMap(bodyData))
		if err != nil {
			zap.S().Panic("Failed to marshal request body to xml: ", err)
		}
		bodyString = string(b)
	}
	return strings.NewReader(bodyString)
}

func getRequestParamter() (url.Values, url.Values, map[string]string, map[string]string, string) {
	apiConfig := global.ServerConfig.Api
	path := apiConfig.Path
	rawQuery := url.Values{}
	formData := url.Values{}
	bodyData := make(map[string]string, len(apiConfig.RequestParams))
	headerMap := make(map[string]string)
	for _, rp := range apiConfig.RequestParams {
		switch {
		case rp.Type == model.REQUEST_PARAM_TYPE_QUERY:
			rawQuery.Set(rp.Name, rp.Value)
		case rp.Type == model.REQUEST_PARAM_TYPE_BODY:
			if apiConfig.RequestContentType == model.REQUEST_APPLICATION_FORM {
				formData.Set(rp.Name, rp.Value)
			} else {
				bodyData[rp.Name] = rp.Value
			}
		case rp.Type == model.REQUEST_PARAM_TYPE_PATH:
			if util.IsBlank(rp.Value) {
				zap.S().Panic("Path value are not allowed to be empty.")
			} else {
				path = strings.Replace(path, ":"+rp.Name, rp.Value, 1)
			}
		case rp.Type == model.REQUEST_PARAM_TYPE_HEADER:
			headerMap[rp.Name] = rp.Value
		default:
			zap.S().Panic("Request param type incorrect.")
		}
	}
	return rawQuery, formData, headerMap, bodyData, path
}
