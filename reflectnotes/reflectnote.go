package reflectnote

import (
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"
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

type MethodStruct struct {
	Methods map[string]*MethodInfo
}

type MethodInfo struct {
	Method reflect.Method
	Host   reflect.Value
	Idx    int
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

		methodStruct.Methods[m.Name] = &MethodInfo{Method: m, Host: host, Idx: i}
	}
}
