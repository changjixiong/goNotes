package reflectnote_test

import (
	"encoding/json"
	"errors"
	"fmt"
	reflectnote "goNotes/reflectnotes"
	"reflect"
	"strconv"
	"testing"
)

type Foo struct {
}

func (foo *Foo) FooFuncZero() bool {
	fmt.Sprintln("FooFuncZero without arg")

	return true
}

func (foo *Foo) FooFuncOne(arg int) string {

	fmt.Sprintln("FooFuncOne with arg:", arg)

	return strconv.Itoa(arg)
}

func (foo *Foo) FooFuncTwo(argStr string, argInt int) string {
	fmt.Sprintln("FooFuncTwo with argOne:", argStr, "argTwo:", argInt)

	return argStr + strconv.Itoa(argInt)
}

type Bar struct {
}

func (bar *Bar) BarFuncZero() string {
	return fmt.Sprintln("BarFuncZero")
}

func (bar *Bar) BarFuncOne(arg float64) int {
	fmt.Sprintln("BarFuncOne with arg:", arg)

	return int(arg)
}

func (bar *Bar) BarFuncTwo(argStr bool, argInt int) int {
	fmt.Sprintln("BarFuncTwo with argOne:", argStr, "argTwo:", argInt)

	if argStr {
		return argInt
	} else {
		return -argInt
	}
}

func init() {
	foo := &Foo{}
	bar := &Bar{}

	reflectnote.RegisterMethod(foo)
	reflectnote.RegisterMethod(bar)
}

func TestInvoke(t *testing.T) {

	resultFooFuncZero := reflectnote.InvokeByReflectArgs("FooFuncZero", nil)
	if false == resultFooFuncZero[0].Bool() {
		t.Errorf("invoke FooFuncZero error")
	}

	resultFooFuncOne := reflectnote.InvokeByReflectArgs("FooFuncOne",
		[]reflect.Value{reflect.ValueOf(123)})

	if "123" != resultFooFuncOne[0].String() {
		t.Errorf("invoke FooFuncOne error")
	}

	argForFooFuncTwo := []reflect.Value{reflect.ValueOf("str123"), reflect.ValueOf(456)}
	resultFooFuncTwo := reflectnote.InvokeByReflectArgs("FooFuncTwo", argForFooFuncTwo)

	if "str123456" != resultFooFuncTwo[0].String() {
		t.Errorf("invoke FooFuncTwo error")
	}

	resultBarFuncZero := reflectnote.InvokeByReflectArgs("BarFuncZero", nil)
	if "BarFuncZero" == resultBarFuncZero[0].String() {
		t.Errorf("invoke BarFuncZero error")
	}

	resultBarFuncOne := reflectnote.InvokeByReflectArgs("BarFuncOne",
		[]reflect.Value{reflect.ValueOf(123.0)})

	if 123 != resultBarFuncOne[0].Int() {
		t.Errorf("invoke BarFuncOne error")
	}

	argForBarFuncTwo := []reflect.Value{reflect.ValueOf(false), reflect.ValueOf(456)}
	resultBarFuncTwo := reflectnote.InvokeByReflectArgs("BarFuncTwo", argForBarFuncTwo)

	if -456 != resultBarFuncTwo[0].Int() {
		t.Errorf("invoke BarFuncTwo error")
	}

}

func TestInvokeByInterfaceArgs(t *testing.T) {
	resultFooFuncZero := reflectnote.InvokeByInterfaceArgs("FooFuncZero", nil)
	if false == resultFooFuncZero[0].Bool() {
		t.Errorf("invoke FooFuncZero error")
	}

	resultFooFuncOne := reflectnote.InvokeByInterfaceArgs("FooFuncOne", []interface{}{123})

	if "123" != resultFooFuncOne[0].String() {
		t.Errorf("invoke FooFuncOne error")
	}

	resultFooFuncTwo := reflectnote.InvokeByInterfaceArgs("FooFuncTwo",
		[]interface{}{"str123", 456})

	if "str123456" != resultFooFuncTwo[0].String() {
		t.Errorf("invoke FooFuncTwo error")
	}

	resultBarFuncZero := reflectnote.InvokeByInterfaceArgs("BarFuncZero", nil)
	if "BarFuncZero" == resultBarFuncZero[0].String() {
		t.Errorf("invoke BarFuncZero error")
	}

	resultBarFuncOne := reflectnote.InvokeByInterfaceArgs("BarFuncOne", []interface{}{123.1})
	if 123 != resultBarFuncOne[0].Int() {
		t.Errorf("invoke BarFuncOne error")
	}

	resultBarFuncTwo := reflectnote.InvokeByInterfaceArgs("BarFuncTwo", []interface{}{false, 456})

	if -456 != resultBarFuncTwo[0].Int() {
		t.Errorf("invoke BarFuncTwo error")
	}
}

func testInvokeByJson(jsonStr, funcName string, expectResult interface{}) error {

	result := make(map[string]interface{})

	err := json.Unmarshal(reflectnote.InvokeByJson([]byte(jsonStr)), &result)

	if err != nil {
		return err
	}

	if resultData, ok := result[funcName]; !ok {

		return errors.New("invoke " + funcName + " error: result not found")

	} else {

		var resultDataConvert interface{}
		switch resultData.(type) {
		case float64:
			resultDataConvert = int(resultData.(float64))
		default:
			resultDataConvert = resultData
		}
		if resultDataConvert != expectResult {
			return errors.New("invoke " + funcName + " error: result not equal")
		}
	}

	return nil
}

func TestInvokeByJson(t *testing.T) {

	var err error
	jsonDataFooFuncZero := `
				{
				    "func_name":"FooFuncZero",
				    "params":[
				    ]
				}
				`
	err = testInvokeByJson(jsonDataFooFuncZero, "FooFuncZero", true)
	if err != nil {
		t.Error(err)
	}

	jsonDataFooFuncOne := `
				{
				    "func_name":"FooFuncOne",
				    "params":[
				        456
				    ]
				}
				`
	err = testInvokeByJson(jsonDataFooFuncOne, "FooFuncOne", "456")
	if err != nil {
		t.Error(err)
	}

	jsonDataFooFuncTwo := `
				{
				    "func_name":"FooFuncTwo",
				    "params":[
				        "str123",
				        456
				    ]
				}
				`
	err = testInvokeByJson(jsonDataFooFuncTwo, "FooFuncTwo", "str123456")
	if err != nil {
		t.Error(err)
	}

	jsonDataBarFuncOne := `
				{
				    "func_name":"BarFuncOne",
				    "params":[
				        456.0
				    ]
				}
				`
	err = testInvokeByJson(jsonDataBarFuncOne, "BarFuncOne", 456)
	if err != nil {
		t.Error(err)
	}

	jsonDataBarFuncTwo := `
					{
					    "func_name":"BarFuncTwo",
					    "params":[
					        false,
					        456
					    ]
					}
					`
	err = testInvokeByJson(jsonDataBarFuncTwo, "BarFuncTwo", -456)
	if err != nil {
		t.Error(err)
	}

}
