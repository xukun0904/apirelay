package model

import (
	"encoding/xml"
	"io"
)

type Result struct {
	XMLName xml.Name    `json:"-" xml:"result,omitempty"`
	Success bool        `json:"success" xml:"success"`
	Code    string      `json:"code,omitempty" xml:"code,omitempty"`
	Message string      `json:"message,omitempty" xml:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"  xml:"data,omitempty"`
}

type StringMap map[string]string

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (sm StringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(sm) == 0 {
		return nil
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for k, v := range sm {
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
	}
	return e.EncodeToken(start.End())
}

func (sm *StringMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*sm = StringMap{}
	for {
		var e xmlMapEntry
		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		(*sm)[e.XMLName.Local] = e.Value
	}
	return nil
}

const (
	SUCCESS uint16 = iota
	SERVER_ERROR
	NOT_FOUND_HEALTHY_INSTANCE
	BAD_GATEWAY
	NOT_FOUND
	METHOD_NOT_ALLOWED
	RESPONSE_CONTENT_INCORRECT
	REQUEST_CONTENT_INCORRECT
	REQUEST_PARAMETER_TYPE_INCORRECT
	REQUEST_PARAMETER_BLANK
)

var ResultMap = map[uint16]Result{
	SUCCESS: {
		Success: true,
		Code:    "00000",
		Message: "操作成功！",
	},
	SERVER_ERROR: {
		Success: false,
		Code:    "B0001",
		Message: "抱歉，系统繁忙，请稍后重试！",
	},
	NOT_FOUND_HEALTHY_INSTANCE: {
		Success: false,
		Code:    "B0002",
		Message: "没有可用的数据源客户端实例！",
	},
	BAD_GATEWAY: {
		Success: false,
		Code:    "B0003",
		Message: "后端服务错误！",
	},
	NOT_FOUND: {
		Success: false,
		Code:    "A0001",
		Message: "请求接口未找到！",
	},
	METHOD_NOT_ALLOWED: {
		Success: false,
		Code:    "A0002",
		Message: "请求方式不允许！",
	},
	RESPONSE_CONTENT_INCORRECT: {
		Success: false,
		Code:    "A0101",
		Message: "响应类型不支持！",
	},
	REQUEST_CONTENT_INCORRECT: {
		Success: false,
		Code:    "A0102",
		Message: "请求类型不正确！",
	},
	REQUEST_PARAMETER_TYPE_INCORRECT: {
		Success: false,
		Code:    "A0103",
		Message: "请求参数类型不正确！",
	},
	REQUEST_PARAMETER_BLANK: {
		Success: false,
		Code:    "A0104",
		Message: "请求参数不能为空！",
	},
}
