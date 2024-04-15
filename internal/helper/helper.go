package helper

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// TypeName 返回类名
func TypeName(val any) string {
	t := reflect.TypeOf(val)
	if len(t.Name()) > 0 { // 这种情况对应的着传值绑定
		t.Name()
	}

	if t.Elem() != nil {
		return t.Elem().Name()
	}
	return ""
}

func AssertNoErrWithFmt(msg string, err error) {
	if err == nil {
		return
	}

	fmt.Println("\nError:", msg)
	panic(err)
}

func PrintJson(m any) {
	buf, _ := json.Marshal(m)
	fmt.Println("")
	fmt.Println("")
	fmt.Println(TypeName(m))
	fmt.Println(string(buf))
	fmt.Println("")
	fmt.Println("")
}

func Json(m any) string {
	buf, _ := json.Marshal(m)
	return TypeName(m) + ": " + string(buf)
}
