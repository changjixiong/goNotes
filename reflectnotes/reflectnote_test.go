package reflectnote_test

import (
	"fmt"
	reflectnote "goNotes/reflectnotes"
	"reflect"
	"testing"
)

func TestInvoke(t *testing.T) {

	methodStruct := reflectnote.MethodStruct{make(map[string]*reflectnote.MethodInfo)}

	foo := &reflectnote.Foo{}
	bar := &reflectnote.Bar{}

	methodStruct.RegisterMethod(foo)
	methodStruct.RegisterMethod(bar)

	expectFooFuncZero := fmt.Sprintln("FooFuncZero without arg")
	resultFooFuncZero := methodStruct.Methods["FooFuncZero"].Host.MethodByName("FooFuncZero").Call(nil)
	if expectFooFuncZero != resultFooFuncZero[0].String() {
		t.Errorf("invoke FooFuncZero error")
	}

	argForFooFuncOne := 123
	expectFooFuncOne := fmt.Sprintln("FooFuncOne with arg:", argForFooFuncOne)
	resultFooFuncOne := methodStruct.Methods["FooFuncOne"].Host.MethodByName("FooFuncOne").Call([]reflect.Value{reflect.ValueOf(argForFooFuncOne)})
	if expectFooFuncOne != resultFooFuncOne[0].String() {
		t.Errorf("invoke FooFuncOne error")
	}

	argForFooFuncTwo := []interface{}{"str123", 456}
	expectFooFuncTwo := fmt.Sprintln("FooFuncTwo with argOne:", argForFooFuncTwo[0], "argTwo:", argForFooFuncTwo[1])

	resultFooFuncTwo := methodStruct.Methods["FooFuncTwo"].
		Host.Method(methodStruct.Methods["FooFuncTwo"].Idx).
		Call([]reflect.Value{reflect.ValueOf(argForFooFuncTwo[0]),
			reflect.ValueOf(argForFooFuncTwo[1])})

	if expectFooFuncTwo != resultFooFuncTwo[0].String() {
		t.Errorf("invoke FooFuncTwo error")
	}
}
