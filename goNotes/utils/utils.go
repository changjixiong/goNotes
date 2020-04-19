package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func Struct2MapString(st interface{}) map[string]string {
	v := reflect.ValueOf(st)

	ma := map[string]string{}

	if v.Kind() != reflect.Ptr {
		return ma
	}

	e := v.Elem()
	t := e.Type()

	for i := 0; i < e.NumField(); i++ {
		field := t.Field(i)
		key := field.Tag.Get("json")
		if len(key) <= 0 {
			key = field.Name
		}
		ma[key] = formatAtom(e.Field(i))
	}

	return ma

}

// Any formats any value as a string.
func Any(value interface{}) string {
	return formatAtom(reflect.ValueOf(value))
}

// formatAtom formats a value without inspecting its internal structure.
func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return v.String()
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	case reflect.Struct:
		switch v.Type().String() {
		case "time.Time":
			return v.Interface().(time.Time).Format(time.RFC3339)
		default:
			return fmt.Sprintf("%v", v)
		}
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}

}

func MapString2Struct(ma map[string]string, st interface{}) bool {
	v := reflect.ValueOf(st)

	if v.Kind() != reflect.Ptr {
		return false
	}

	e := v.Elem()
	t := e.Type()

	matchField := false

	for i := 0; i < e.NumField(); i++ {
		field := t.Field(i)

		fe := e.Field(i)
		key := field.Tag.Get("json")
		if len(key) <= 0 {
			key = field.Name
		}
		if v, found := ma[key]; found {
			matchField = true
			setField(&fe, v)
		}
	}

	return matchField
}

func setField(refV *reflect.Value, value string) {

	switch refV.Kind() {
	case reflect.Bool:
		var v bool
		v, err := strconv.ParseBool(value)
		if err == nil {
			refV.SetBool(v)
		}
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		var v int64
		v, err := strconv.ParseInt(value, 0, refV.Type().Bits())
		if err == nil {
			refV.SetInt(v)
		}

	case reflect.Float64:
		v, err := strconv.ParseFloat(value, refV.Type().Bits())
		if err == nil {
			refV.SetFloat(v)
		}
	case reflect.String:
		refV.SetString(value)

	case reflect.Struct:
		switch refV.Type().String() {
		case "time.Time":
			v, err := time.Parse(time.RFC3339, value)

			if nil == err {
				refV.Set(reflect.ValueOf(v))
			}

		default:

		}
	}
}
