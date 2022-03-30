package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"jhr.com/apirelay/exception"
	"jhr.com/apirelay/global"
	"jhr.com/apirelay/model"
	"jhr.com/apirelay/util"
)

func ForwardRequest(rw http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	// 校验请求类型
	if err := checkRequestContentType(r); err != nil {
		return err
	}
	// 校验请求参数
	if err := checkRequestParams(r, params); err != nil {
		return err
	}
	// 转发请求
	if err := forwarding(rw, r); err != nil {
		return err
	}
	return nil
}

func forwarding(rw http.ResponseWriter, r *http.Request) error {
	forwardUrl := global.ServerConfig.Service.Host + r.RequestURI
	remote, err := url.Parse(forwardUrl)
	if err != nil {
		return err
	}
	if !global.ServerConfig.Test {
		r.URL.Path = ""
	}
	r.URL.Scheme = remote.Scheme
	r.URL.Host = remote.Host
	r.Host = remote.Host
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	rp := getReverseProxy(remote)
	rp.ServeHTTP(rw, r)
	return nil
}

func getReverseProxy(remote *url.URL) *httputil.ReverseProxy {
	rp := httputil.NewSingleHostReverseProxy(remote)
	// 添加超时时间
	rp.Transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, time.Millisecond*time.Duration(global.ServerConfig.Service.Timeout))
		},
	}
	// 代理异常处理
	rp.ErrorHandler = func(rw http.ResponseWriter, r *http.Request, e error) {
		zap.S().Error("Proxy error: ", e.Error())
		util.WriteResponse(rw, model.ResultMap[model.BAD_GATEWAY])
	}
	return rp
}

func checkRequestContentType(r *http.Request) error {
	apiConfig := global.ServerConfig.Api
	// 校验请求Content-Type是否正确
	var requestContentTypeConfig = model.RequestContentTypeName[apiConfig.RequestContentType]
	if util.IsNotBlank(requestContentTypeConfig) {
		requestContentType := r.Header.Get("Content-Type")
		if requestContentTypeConfig != requestContentType {
			return exception.Cast(model.ResultMap[model.REQUEST_CONTENT_INCORRECT])
		}
	}
	return nil
}

func checkRequestParams(r *http.Request, params httprouter.Params) error {
	apiConfig := global.ServerConfig.Api
	// 获取body内容
	m, err := getRequestBody(r)
	if err != nil {
		return err
	}
	// 校验参数是否正确
	requestParams := apiConfig.RequestParams
	for _, rp := range requestParams {
		var paramValue string
		switch {
		case rp.Type == model.REQUEST_PARAM_TYPE_QUERY:
			paramValue = r.FormValue(rp.Name)
		case rp.Type == model.REQUEST_PARAM_TYPE_BODY:
			if apiConfig.RequestContentType == model.REQUEST_APPLICATION_FORM {
				paramValue = r.PostFormValue(rp.Name)
			} else {
				paramValue = m[rp.Name]
			}
		case rp.Type == model.REQUEST_PARAM_TYPE_PATH:
			paramValue = params.ByName(rp.Name)
		case rp.Type == model.REQUEST_PARAM_TYPE_HEADER:
			if len(r.Header[rp.Name]) > 0 {
				paramValue = r.Header[rp.Name][0]
			}
		default:
			return exception.Cast(model.ResultMap[model.REQUEST_PARAMETER_TYPE_INCORRECT])
		}
		if util.IsBlank(paramValue) && util.IsNotBlank(rp.DefaultValue) {
			paramValue = rp.DefaultValue
		}
		if rp.Required && util.IsBlank(paramValue) {
			return exception.Cast(model.ResultMap[model.REQUEST_PARAMETER_BLANK])
		}
	}
	return nil
}

func getRequestBody(r *http.Request) (map[string]string, error) {
	apiConfig := global.ServerConfig.Api
	var m map[string]string
	if apiConfig.RequestContentType == model.REQUEST_APPLICATION_JSON || apiConfig.RequestContentType == model.REQUEST_APPLICATION_XML {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return m, err
		}
		if len(b) == 0 {
			return m, nil
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
		if apiConfig.RequestContentType == model.REQUEST_APPLICATION_JSON {
			if err := json.Unmarshal(b, &m); err != nil {
				return m, err
			}
		} else if apiConfig.RequestContentType == model.REQUEST_APPLICATION_XML {
			var sm model.StringMap
			if err := xml.Unmarshal(b, &sm); err != nil {
				return m, err
			}
			m = sm
		}
	}
	return m, nil
}
