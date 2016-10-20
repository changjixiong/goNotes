package main

import (
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"
)

type Foo struct {
}

func (foo *Foo) FooFuncZero() {
	fmt.Println("FooFuncZero without arg")
}

func (foo *Foo) FooFuncOne(arg int) {
	fmt.Println("FooFuncOne with arg:", arg)
}

func (foo *Foo) FooFuncTwo(argStr string, argInt int) {
	fmt.Println("FooFuncTwo with argOne:", argStr, "argTwo:", argInt)
}

type Bar struct {
}

func (bar *Bar) BarFuncZero() {
	fmt.Println("BarFuncZero without arg")
}

func (bar *Bar) BarFuncOne(arg float64) {
	fmt.Println("BarFuncOne with arg:", arg)
}

func (bar *Bar) BarFuncTwo(argStr bool, argInt int) {
	fmt.Println("BarFuncTwo with argOne:", argStr, "argTwo:", argInt)
}

type MethodStruct struct {
	Methods map[string]struct {
		Method reflect.Method
		host   reflect.Value
		idx    int
	}
}

func (methodStruct *MethodStruct) RegisterMethod(v interface{}) {

	reflectType := reflect.TypeOf(v)
	host := reflect.ValueOf(v)

	for i := 0; i < reflectType.NumMethod(); i++ {
		m := reflectType.Method(i)

		char, _ := utf8.DecodeRuneInString(m.Name)
		//非导出函数不注册
		if !unicode.IsUpper(char) {
			continue
		}

		methodStruct.Methods[m.Name] = struct {
			Method reflect.Method
			host   reflect.Value
			idx    int
		}{Method: m, host: host, idx: i}
	}
}

func main() {
	methodStruct := MethodStruct{Methods: make(map[string]struct {
		Method reflect.Method
		host   reflect.Value
		idx    int
	})}

	foo := &Foo{}
	bar := &Bar{}

	methodStruct.RegisterMethod(foo)
	methodStruct.RegisterMethod(bar)

	methodStruct.Methods["FooFuncZero"].host.MethodByName("FooFuncZero").Call(nil)
	methodStruct.Methods["FooFuncOne"].host.MethodByName("FooFuncOne").Call([]reflect.Value{reflect.ValueOf(123)})

	methodStruct.Methods["FooFuncTwo"].host.Method(methodStruct.Methods["FooFuncTwo"].idx).Call([]reflect.Value{reflect.ValueOf("str123"), reflect.ValueOf(456)})
}
