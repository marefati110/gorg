package utils

import (
	"reflect"
	"runtime"
	"strings"
)

func GetFunctionName(fn any) string {
	handlerValue := reflect.ValueOf(fn)
	handlerNameWithPackageName := runtime.FuncForPC(handlerValue.Pointer()).Name()
	lastDotIndex := strings.LastIndex(handlerNameWithPackageName, ".")
	return handlerNameWithPackageName[lastDotIndex+1:]
}
