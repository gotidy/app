package cli

import (
	"errors"
	"fmt"
	"reflect"
)

func CombineStructs(structs ...interface{}) (interface{}, error) {
	fields := []reflect.StructField{}
	for _, s := range structs {
		v := reflect.Indirect(reflect.ValueOf(s))
		if v.Kind() != reflect.Struct {
			return nil, fmt.Errorf("invalid value type, must be struct: %T", s)
		}
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			fields = append(fields, t.Field(i))
		}
	}
	result := reflect.New(reflect.StructOf(fields))
	// Copy data
	resultIndirect := reflect.Indirect(result)
	for _, s := range structs {
		v := reflect.Indirect(reflect.ValueOf(s))
		for i := 0; i < v.NumField(); i++ {
			src := v.Field(i)
			dest := resultIndirect.FieldByName(v.Type().Field(i).Name)
			dest.Set(src)
		}
	}
	return result.Interface(), nil
}

func CopyStruct(src, dest interface{}) error {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	if !destValue.CanSet() {
		return errors.New("unable to set destination, must be pointer")
	}
	if destValue.Kind() != reflect.Struct {
		return fmt.Errorf("invalid value type, must be struct: %T", dest)
	}
	srcValue := reflect.Indirect(reflect.ValueOf(src))
	if srcValue.Kind() != reflect.Struct {
		return fmt.Errorf("invalid value type, must be struct: %T", src)
	}
	destType := destValue.Type()
	for i := 0; i < destValue.NumField(); i++ {
		dest := destValue.Field(i)
		src := srcValue.FieldByName(destType.Field(i).Name)
		if !src.IsValid() {
			continue
		}
		if src.Type().AssignableTo(dest.Type()) {
			dest.Set(src)
		}
	}
	return nil
}
