package model

import (
	"net/http"
)

type ServerConfig struct {
	Port    uint16        `mapstructure:"port" json:"port" validate:"required"`
	Test    bool          `mapstructure:"test" json:"test"`
	Api     ApiConfig     `mapstructure:"api" json:"api"`
	Service ServiceConfig `mapstructure:"service" json:"service"`
}

type ApiConfig struct {
	Protocol            Protocol            `mapstructure:"protocol" json:"protocol" validate:"required"`
	RequestMethod       RequestMethod       `mapstructure:"request_method" json:"request_method" validate:"required"`
	RequestContentType  RequestContentType  `mapstructure:"request_content_type" json:"request_content_type"`
	RequestParams       []RequestParam      `mapstructure:"request_params" json:"request_params"`
	Path                string              `mapstructure:"path" json:"path" validate:"required"`
	ResponseContentType ResponseContentType `mapstructure:"response_content_type" json:"response_content_type" validate:"required"`
}

type RequestParam struct {
	Name         string           `mapstructure:"name" json:"name" validate:"required"`
	Type         RequestParamType `mapstructure:"type" json:"type" validate:"required"`
	Required     bool             `mapstructure:"required" json:"required"`
	DefaultValue string           `mapstructure:"default_value" json:"default_value"`
	Value        string           `mapstructure:"value" json:"value"`
}

type ServiceConfig struct {
	Host    string `mapstructure:"host" json:"host" validate:"url"`
	Timeout uint8  `mapstructure:"timeout" json:"timeout"`
}

type RequestContentType uint8

const (
	REQUEST_CONTENT_TYPE_NOT_SPECIFIED RequestContentType = iota
	REQUEST_APPLICATION_JSON
	REQUEST_APPLICATION_XML
	REQUEST_APPLICATION_FORM
)

var (
	RequestContentTypeName = map[RequestContentType]string{
		REQUEST_CONTENT_TYPE_NOT_SPECIFIED: "",
		REQUEST_APPLICATION_JSON:           "application/json",
		REQUEST_APPLICATION_XML:            "application/xml",
		REQUEST_APPLICATION_FORM:           "application/x-www-form-urlencoded",
	}
)

type ResponseContentType uint8

const (
	RESPONSE_CONTENT_TYPE_NOT_SPECIFIED ResponseContentType = iota
	RESPONSE_APPLICATION_JSON
	RESPONSE_APPLICATION_XML
)

var (
	ResponseContentTypeName = map[ResponseContentType]string{
		RESPONSE_CONTENT_TYPE_NOT_SPECIFIED: "application/json",
		RESPONSE_APPLICATION_JSON:           "application/json",
		RESPONSE_APPLICATION_XML:            "application/xml",
	}
)

type Protocol uint8

const (
	PROTOCOL_NOT_SPECIFIED Protocol = iota
	PROTOCOL_HTTP
	PROTOCOL_HTTPS
)

type RequestMethod uint8

const (
	REQUEST_METHOD_NOT_SPECIFIED RequestMethod = iota
	REQUEST_METHOD_GET
	REQUEST_METHOD_POST
)

var (
	RequestMethodName = map[RequestMethod]string{
		REQUEST_METHOD_NOT_SPECIFIED: http.MethodGet,
		REQUEST_METHOD_GET:           http.MethodGet,
		REQUEST_METHOD_POST:          http.MethodPost,
	}
)

type RequestParamType uint8

const (
	REQUEST_PARAM_TYPE_NOT_SPECIFIED RequestParamType = iota
	REQUEST_PARAM_TYPE_QUERY
	REQUEST_PARAM_TYPE_BODY
	REQUEST_PARAM_TYPE_PATH
	REQUEST_PARAM_TYPE_HEADER
)
