package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type Record struct {
	ID          int
	Num         int
	IndexString string
	Name        string
}

func FillStruct(structPtr interface{}, data []string) error {
	v := reflect.ValueOf(structPtr)
	//t := reflect.TypeOf(structPtr)
	if reflect.Ptr != v.Type().Kind() {
		return errors.New("must pass a pointer")
	}

	if v.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}

	// baseBataType := t.Elem()
	structValue := v.Elem()
	baseBataType := structValue.Type()

	if len(data) != baseBataType.NumField() {
		return errors.New("field count mismatch")
	}

	for i := 0; i < baseBataType.NumField(); i++ {
		dataTypeField := baseBataType.Field(i)
		field := structValue.Field(i)

		if !field.CanSet() {
			return errors.New("field can not be fill" + dataTypeField.Name)
		}

		switch dataTypeField.Type.Kind() {
		case reflect.Bool:
			var v bool
			v, err := strconv.ParseBool(data[i])
			if err == nil {
				field.SetBool(v)
			}
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			var v int64
			v, err := strconv.ParseInt(data[i], 0, dataTypeField.Type.Bits())
			if err == nil {
				field.SetInt(v)
			}
		case reflect.String:
			field.SetString(data[i])
		}
	}

	return nil
}

func main() {
	datas := [][]string{}
	datas = append(datas, []string{"1", "11", "a", "aaa"})
	datas = append(datas, []string{"2", "22", "b", "bbb"})
	datas = append(datas, []string{"3", "33", "c", "ccc"})
	Records := []*Record{}

	for _, data := range datas {
		record := &Record{}
		if err := FillStruct(record, data); nil == err {
			Records = append(Records, record)
		} else {
			//错误
		}

	}

	for _, record := range Records {
		fmt.Println(record)
	}
}
