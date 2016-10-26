package reflectnote

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"
)

type MethodMap struct {
	Methods map[string]*MethodInfo
}

type MethodInfo struct {
	Method reflect.Method
	Host   reflect.Value
	Idx    int
}

type Request struct {
	FuncName string        `json:"func_name"`
	Params   []interface{} `json:"params"`
}

type Response struct {
	FuncName string        `json:"func_name"`
	Data     []interface{} `json:"data"`
}

var methodStruct MethodMap = MethodMap{make(map[string]*MethodInfo)}

func RegisterMethod(v interface{}) {

	reflectType := reflect.TypeOf(v)
	host := reflect.ValueOf(v)

	for i := 0; i < reflectType.NumMethod(); i++ {
		m := reflectType.Method(i)

		char, _ := utf8.DecodeRuneInString(m.Name)
		//非导出函数不注册
		if !unicode.IsUpper(char) {
			continue
		}

		methodStruct.Methods[m.Name] = &MethodInfo{Method: m, Host: host, Idx: i}
	}
}

func convertParamType(v interface{}, targetType reflect.Type) (targetValue reflect.Value, ok bool) {
	defer func() {
		if re := recover(); re != nil {
			ok = false
			fmt.Println(re)
		}
	}()

	ok = true

	if targetType.Kind() == reflect.Interface ||
		targetType.Kind() == reflect.TypeOf(v).Kind() {

		targetValue = reflect.ValueOf(v)

	} else if reflect.TypeOf(v).Kind() == reflect.Float64 {
		f := v.(float64)
		switch targetType.Kind() {
		case reflect.Int:
			targetValue = reflect.ValueOf(int(f))
		case reflect.Uint8:
			targetValue = reflect.ValueOf(uint8(f))
		case reflect.Uint16:
			targetValue = reflect.ValueOf(uint16(f))
		case reflect.Uint32:
			targetValue = reflect.ValueOf(uint32(f))
		case reflect.Uint64:
			targetValue = reflect.ValueOf(uint64(f))
		case reflect.Int8:
			targetValue = reflect.ValueOf(int8(f))
		case reflect.Int16:
			targetValue = reflect.ValueOf(int16(f))
		case reflect.Int32:
			targetValue = reflect.ValueOf(int32(f))
		case reflect.Int64:
			targetValue = reflect.ValueOf(int64(f))
		case reflect.Float32:
			targetValue = reflect.ValueOf(float32(f))
		default:
			ok = false
		}
	} else {
		ok = false
	}

	return
}

func convertParam(funcName string, Params []interface{}) []reflect.Value {

	paramsValue := make([]reflect.Value, 0, len(Params))
	//跳过 receiver
	for i := 1; i < methodStruct.Methods[funcName].Method.Type.NumIn(); i++ {
		inParaType := methodStruct.Methods[funcName].Method.Type.In(i)
		value, _ := convertParamType(Params[i-1], inParaType)
		paramsValue = append(paramsValue, value)
	}

	return paramsValue
}

func InvokeByArgs(funcName string, par []reflect.Value) []reflect.Value {

	return methodStruct.Methods[funcName].Host.MethodByName(funcName).Call(par)
}

func InvokeByParams(funcName string, Params []interface{}) []reflect.Value {

	paramsValue := convertParam(funcName, Params)

	return methodStruct.Methods[funcName].Host.MethodByName(funcName).Call(paramsValue)
}

func InvokeByString(strData string) []reflect.Value {

	req := &Request{}
	err := json.Unmarshal([]byte(strData), req)

	if err != nil {
		return nil
	}

	return InvokeByParams(req.FuncName, req.Params)
}
