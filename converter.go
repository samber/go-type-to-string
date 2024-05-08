package typetostring

import (
	"fmt"
	"reflect"
)

func GetType[T any]() string {
	var t T
	return GetValueType(t)
}

func GetValueType[T any](t T) string {
	return getType(reflect.TypeOf(&t).Elem())
}

func GetReflectType(typeOfT reflect.Type) string {
	return getType(typeOfT)
}

func GetReflectValueType(v reflect.Value) string {
	return GetReflectType(v.Type())
}

// getType generates a service name from a type.
// It returns a string like "*[]*github.com/samber/example.test".
func getType(typeOfT reflect.Type) string {
	if typeOfT.Kind() == reflect.Pointer {
		return "*" + getType(typeOfT.Elem())
	} else if typeOfT.Kind() == reflect.Slice || typeOfT.Kind() == reflect.Array {
		return "[]" + getType(typeOfT.Elem())
	} else if typeOfT.Kind() == reflect.Map {
		key := getType(typeOfT.Key())
		value := getType(typeOfT.Elem())
		return fmt.Sprintf("map[%s]%s", key, value)
	} else if typeOfT.Kind() == reflect.Chan {
		var prefix string

		switch typeOfT.ChanDir() {
		case reflect.RecvDir:
			prefix = "<-chan"
		case reflect.SendDir:
			prefix = "chan<-"
		case reflect.BothDir:
			prefix = "chan"
		}

		return fmt.Sprintf("%s %s", prefix, getType(typeOfT.Elem()))
	} else if typeOfT.Kind() == reflect.Func {
		// @TODO: handle arguments and returned types recursively
		return typeOfT.String()
	}

	pkgPath := typeOfT.PkgPath()
	if pkgPath == "" {
		// anonymous type
		return typeOfT.String()
	}

	if typeOfT.Name() == "" {
		// any + interface{} + anonymous type
		return typeOfT.String()
	}

	return pkgPath + "." + typeOfT.Name()
}
