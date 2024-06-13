package helper

import (
	"bytes"
	"fmt"
	"github.com/samber/lo"
	"log"
	"strings"
	"text/template"
)

const TypeB = "bool"
const TypeI = "int"
const TypeStr = "string"
const TypeIArr = "[]int"
const TypeStrArr = "[]string"
const TypeIArr2d = "[][]int"
const TypeStrArr2d = "[][]string"
const TypeMapIToI = "map[int]int"
const TypeMapIToStr = "map[int]string"
const TypeMapIToIArr = "map[int][]int"
const TypeMapIToStrArr = "map[int][]string"
const TypeMapStrToI = "map[string]int"
const TypeMapStrToStr = "map[string]string"
const TypeMapStrToIArr = "map[string][]int"
const TypeMapStrToStrArr = "map[string][]string"

var allTypes = []string{
	// 基础类型
	TypeB, TypeI, TypeStr,

	// 数组类型
	TypeIArr, TypeStrArr, TypeIArr2d, TypeStrArr2d,

	// map 类型
	TypeMapIToI, TypeMapIToStr, TypeMapIToIArr, TypeMapIToStrArr,
	TypeMapStrToI, TypeMapStrToStr, TypeMapStrToIArr, TypeMapStrToStrArr,
}

// TypeAlias 类型别名
var typeAlias = map[string][]string{
	TypeB:              {TypeB},
	TypeI:              {TypeI, "integer", "int8", "int16", "int32", "int64", "intw"},
	TypeStr:            {TypeStr, "str", "other", "i18n"},
	TypeIArr:           {TypeIArr, "intarr", "intw[]", "int[]"},
	TypeStrArr:         {TypeStrArr, "[]str", "string[]", "strarr", "i18n[]"},
	TypeIArr2d:         {TypeIArr2d, "[]intarr", "int[][]"},
	TypeStrArr2d:       {TypeStrArr2d, "[][]str", "[]strarr", "string[][]"},
	TypeMapIToI:        {TypeMapIToI, "int:int"},
	TypeMapIToStr:      {TypeMapIToStr, "int:str"},
	TypeMapIToIArr:     {TypeMapIToIArr, "int:[]int", "int:intarr"},
	TypeMapIToStrArr:   {TypeMapIToStrArr, "int:[]str", "int:strarr"},
	TypeMapStrToI:      {TypeMapStrToI, "str:int"},
	TypeMapStrToStr:    {TypeMapStrToStr, "str:str"},
	TypeMapStrToIArr:   {TypeMapStrToIArr, "str:[]int", "str:intarr"},
	TypeMapStrToStrArr: {TypeMapStrToStrArr, "str:[]str", "str:strarr"},
}

// AllValidTypeNames 返回所有合法的类型名称
func AllValidTypeNames() []string {
	lists := lo.Values(typeAlias)
	return lo.Flatten(lists)
}

// FormatTypeName 格式化类型名称
func FormatTypeName(typeName string) string {
	for _, types := range typeAlias {
		if lo.Contains(types, typeName) {
			return types[0]
		}
	}
	log.Fatal("Error: invalid type name: ", typeName)
	return typeName
}

func GetVariableLiteralValue(typ string, val string) string {
	switch typ {
	case TypeI:
		return GoLiteralValInt(val)
	case TypeIArr:
		return GoLiteralValIntArr(val)
	case TypeIArr2d:
		return GoLiteralValIntArr2d(val)
	case TypeStr:
		return GoLiteralValStr(val)
	case TypeStrArr:
		return "NOT IMPLEMENT"
	case TypeStrArr2d:
		return "NOT IMPLEMENT"
	}
	return ""
}

func GoLiteralValIntArr(val string) string {
	tmp := `[]int{ {{- . -}} }`
	tmpl, err := template.New("").Parse(tmp)
	if err != nil {
		log.Fatal(err)
	}

	b := bytes.Buffer{}
	err = tmpl.Execute(&b, strings.ReplaceAll(val, "#", ", "))
	if err != nil {
		return ""
	}
	return b.String()
}

func GoLiteralValIntArr2d(val string) string {
	tmp := `[][]int{ {{- . -}} }`
	tmpl, err := template.New("").Parse(tmp)
	if err != nil {
		log.Fatal(err)
	}

	list := lo.FilterMap(strings.Split(val, "|"), func(item string, index int) (string, bool) {
		trimmed := strings.TrimSpace(item)
		temp := fmt.Sprintf("{%s}", strings.ReplaceAll(trimmed, "#", ","))
		return temp, len(trimmed) > 0
	})
	if len(list) == 0 {
		return "[][]int"
	}

	b := bytes.Buffer{}
	err = tmpl.Execute(&b, strings.Join(list, ","))
	if err != nil {
		return ""
	}
	return b.String()
}

func GoLiteralValInt(val string) string {
	return val
}

func GoLiteralValStr(val string) string {
	return `"` + val + `"`
}
