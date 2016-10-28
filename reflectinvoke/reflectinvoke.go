package reflectinvoke

import (
	"encoding/json"
	"errors"
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

func convertParamType(v interface{}, targetType reflect.Type) (
	targetValue reflect.Value, ok bool) {
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

func convertParam(methodInfo *MethodInfo, Params []interface{}) ([]reflect.Value, error) {

	if len(Params) != methodInfo.Method.Type.NumIn()-1 {
		return nil, errors.New("convertParam number error")
	}

	paramsValue := make([]reflect.Value, 0, len(Params))
	//跳过 receiver
	for i := 1; i < methodInfo.Method.Type.NumIn(); i++ {
		inParaType := methodInfo.Method.Type.In(i)
		value, ok := convertParamType(Params[i-1], inParaType)
		if !ok {
			return nil, errors.New("convertParamType error")
		}
		paramsValue = append(paramsValue, value)
	}

	return paramsValue, nil
}

func InvokeByReflectArgs(funcName string, par []reflect.Value) []reflect.Value {

	return methodStruct.Methods[funcName].Host.MethodByName(funcName).Call(par)
}

func InvokeByInterfaceArgs(funcName string, Params []interface{}) []reflect.Value {

	paramsValue, err := convertParam(methodStruct.Methods[funcName], Params)

	if err != nil {
		return nil
	}

	return methodStruct.Methods[funcName].Host.MethodByName(funcName).Call(paramsValue)
}

func InvokeByValues(methodInfo *MethodInfo, params []reflect.Value) (
	data map[string]interface{}) {

	data = make(map[string]interface{})
	result := methodInfo.Host.Method(methodInfo.Idx).Call(params)
	if len(result) > 0 {
		data[methodInfo.Method.Name] = result[0].Interface()
	} else {
		data[methodInfo.Method.Name] = nil
	}

	return
}

func InvokeByJson(byteData []byte) []byte {

	req := &Request{}
	err := json.Unmarshal(byteData, req)

	resultData := make(map[string]interface{})

	if err != nil {
		resultData = map[string]interface{}{"Error": err.Error()}
	} else {

		methodInfo := methodStruct.Methods[req.FuncName]
		paramsValue, err := convertParam(methodInfo, req.Params)

		if err != nil {
			resultData["error"] = err.Error()
		} else {
			resultData = InvokeByValues(methodInfo, paramsValue)
		}

	}

	data, _ := json.Marshal(resultData)

	return data

}
