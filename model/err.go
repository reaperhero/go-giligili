package model

import (
	"encoding/json"
	"fmt"
	validator "gopkg.in/go-playground/validator.v8"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

//
//const ErrorCodeOK = 1000
//
//var (
//	ErrorParamsInvalid     = errors.New("json数值错误")
//	ErrorStructBind        = errors.New("json字段错误")
//	ErrorDomainExisting    = errors.New("域名已经存在")
//	ErrorDomainNotExisting = errors.New("域名不存在存在")
//	ErrorDomainId          = errors.New("域名配置不存在")
//)
//
//func ErrorToCode(err error) (code int) {
//	switch err {
//	case ErrorParamsInvalid:
//		return 1001
//	case ErrorStructBind:
//		return 1002
//	case ErrorDomainExisting:
//		return 1003
//	case ErrorDomainNotExisting:
//		return 1004
//	case ErrorDomainId:
//		return 1005
//	case nil:
//		return ErrorCodeOK
//	default:
//		return 0
//	}
//}
//
//// GetErrorMap ...
//func GetErrorMap(err error) map[string]interface{} {
//	var msg = "OK"
//	if err != nil {
//		msg = err.Error()
//	}
//
//	return map[string]interface{}{
//		"errcode": ErrorToCode(err),
//		"msg":     msg,
//	}
//}


var Dictinary *map[interface{}]interface{}

// LoadLocales 读取国际化文件
func LoadLocales(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		return err
	}

	Dictinary = &m

	return nil
}

// T 翻译
func T(key string) string {
	dic := *Dictinary
	keys := strings.Split(key, ".")
	for index, path := range keys {
		// 如果到达了最后一层，寻找目标翻译
		if len(keys) == (index + 1) {
			for k, v := range dic {
				if k, ok := k.(string); ok {
					if k == path {
						if value, ok := v.(string); ok {
							return value
						}
					}
				}
			}
			return path
		}
		// 如果还有下一层，继续寻找
		for k, v := range dic {
			if ks, ok := k.(string); ok {
				if ks == path {
					if dic, ok = v.(map[interface{}]interface{}); ok == false {
						return path
					}
				}
			} else {
				return ""
			}
		}
	}

	return ""
}

func ErrorResponse(err error) Response {
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			field := T(fmt.Sprintf("Field.%s", e.Field))
			tag := T(fmt.Sprintf("Tag.Valid.%s", e.Tag))
			return Response{
				Status: 40001,
				Msg:    fmt.Sprintf("%s%s", field, tag),
				Error:  fmt.Sprint(err),
			}
		}
	}
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return Response{
			Status: 40002,
			Msg:    "JSON类型不匹配",
			Error:  fmt.Sprint(err),
		}
	}

	return Response{
		Status: 40003,
		Msg:    "参数错误",
		Error:  fmt.Sprint(err),
	}
}
