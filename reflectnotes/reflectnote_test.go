package reflectnote_test

import (
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
	return fmt.Sprintln("BarFuncZero without arg")
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

	resultFooFuncZero := reflectnote.InvokeByArgs("FooFuncZero", nil)
	if false == resultFooFuncZero[0].Bool() {
		t.Errorf("invoke FooFuncZero error")
	}

	resultFooFuncOne := reflectnote.InvokeByArgs("FooFuncOne",
		[]reflect.Value{reflect.ValueOf(123)})

	if "123" != resultFooFuncOne[0].String() {
		t.Errorf("invoke FooFuncOne error")
	}

	argForFooFuncTwo := []reflect.Value{reflect.ValueOf("str123"), reflect.ValueOf(456)}
	resultFooFuncTwo := reflectnote.InvokeByArgs("FooFuncTwo", argForFooFuncTwo)

	if "str123456" != resultFooFuncTwo[0].String() {
		t.Errorf("invoke FooFuncTwo error")
	}

}

func TestInvokeByJson(t *testing.T) {

	jsonData := `
				{
				    "func_name":"FooFuncTwo",
				    "params":[
				        "str123",
				        456
				    ]
				}
				`

	resultFooFuncTwo := reflectnote.InvokeByString(jsonData)

	if "str123456" != resultFooFuncTwo[0].String() {
		t.Errorf("invoke FooFuncTwo error")
	}

}
