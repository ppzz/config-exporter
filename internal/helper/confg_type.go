package helper

import (
	"log"
	"regexp"
	"strings"
	"unicode"
)

// 配置表类型
const (
	TypInvalid Typ = iota // 占位
	TypGlobal             // 枚举配置表
	TypNormal             // 普通配置表
	TypI18n               // 多语言配置表
)

type Typ int

func (t Typ) String() string {
	switch t {
	case TypGlobal:
		return "global"
	case TypNormal:
		return "normal"
	case TypI18n:
		return "i18n"
	default:
		return "invalid"
	}
}

var TypGlobalNameMatcher = regexp.MustCompile("global_.+$")    // global
var TypI18nNameMatcher = regexp.MustCompile(`i18n_.+$`)        // i18n表
var TypNormalNameMatcher = regexp.MustCompile("[a-z0-9]+_.+$") // 普通表

// GetConfType 获取配置文件的类型
func GetConfType(filename string) (typ Typ) {
	lower := strings.ToLower(FileBareName(filename))

	// global
	if TypGlobalNameMatcher.MatchString(lower) {
		return TypGlobal
	}
	// i18n表
	if TypI18nNameMatcher.MatchString(lower) {
		return TypI18n
	}
	// 普通表
	if TypNormalNameMatcher.MatchString(lower) {
		return TypNormal
	}

	return TypInvalid
}

// FilenameByType 返回文件名
func FilenameByType(bareName string) string {
	bareName = strings.Split(bareName, ".")[0]
	typ := GetConfType(bareName)
	switch typ {
	case TypGlobal, TypNormal: // 下划线之前的内容
		idx := strings.Index(bareName, "_")
		if idx == -1 {
			return bareName
		}
		return bareName[:idx]

	case TypI18n: // 去掉下划线
		return strings.ReplaceAll(bareName, "_", "")
	default:
		log.Fatal("invalid config file type: ", bareName)
		return "INVALID"
	}
}

// CamelCaseToSnakeCase 把驼峰转为蛇形命名
func CamelCaseToSnakeCase(s string) string {
	s = regexp.MustCompile(`\p{Han}`).ReplaceAllString(s, "") // 去掉中文
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}
