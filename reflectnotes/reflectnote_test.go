package reflectnote_test

import (
	"fmt"
	reflectnote "goNotes/reflectnotes"
	"reflect"
	"testing"
)

type Foo struct {
}

func (foo *Foo) FooFuncZero() string {
	return fmt.Sprintln("FooFuncZero without arg")
}

func (foo *Foo) FooFuncOne(arg int) string {

	return fmt.Sprintln("FooFuncOne with arg:", arg)
}

func (foo *Foo) FooFuncTwo(argStr string, argInt int) string {
	return fmt.Sprintln("FooFuncTwo with argOne:", argStr, "argTwo:", argInt)
}

type Bar struct {
}

func (bar *Bar) BarFuncZero() string {
	return fmt.Sprintln("BarFuncZero without arg")
}

func (bar *Bar) BarFuncOne(arg float64) string {
	return fmt.Sprintln("BarFuncOne with arg:", arg)
}

func (bar *Bar) BarFuncTwo(argStr bool, argInt int) string {
	return fmt.Sprintln("BarFuncTwo with argOne:", argStr, "argTwo:", argInt)
}

func init() {
	foo := &Foo{}
	bar := &Bar{}

	reflectnote.RegisterMethod(foo)
	reflectnote.RegisterMethod(bar)
}

func TestInvoke(t *testing.T) {

	expectFooFuncZero := fmt.Sprintln("FooFuncZero without arg")
	resultFooFuncZero := reflectnote.InvokeByArgs("FooFuncZero", nil)

	if expectFooFuncZero != resultFooFuncZero[0].String() {
		t.Errorf("invoke FooFuncZero error")
	}

	argForFooFuncOne := 123
	expectFooFuncOne := fmt.Sprintln("FooFuncOne with arg:", argForFooFuncOne)
	resultFooFuncOne := reflectnote.InvokeByArgs("FooFuncOne",
		[]reflect.Value{reflect.ValueOf(argForFooFuncOne)})

	if expectFooFuncOne != resultFooFuncOne[0].String() {
		t.Errorf("invoke FooFuncOne error")
	}

	argForFooFuncTwo := []reflect.Value{reflect.ValueOf("str123"), reflect.ValueOf(456)}
	expectFooFuncTwo := fmt.Sprintln("FooFuncTwo with argOne:", argForFooFuncTwo[0],
		"argTwo:", argForFooFuncTwo[1])

	resultFooFuncTwo := reflectnote.InvokeByArgs("FooFuncTwo", argForFooFuncTwo)

	if expectFooFuncTwo != resultFooFuncTwo[0].String() {
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

	expectFooFuncTwo := fmt.Sprintln("FooFuncTwo"+" with argOne:", "str123", "argTwo:", "456")
	fmt.Println(expectFooFuncTwo)

	resultFooFuncTwo := reflectnote.InvokeByString(jsonData)

	if expectFooFuncTwo != resultFooFuncTwo[0].String() {
		t.Errorf("invoke FooFuncTwo error")
	}

}
