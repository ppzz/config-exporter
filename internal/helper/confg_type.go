package helper

import (
	"log"
	"regexp"
	"strings"
	"unicode"
)

// 配置表类型
const (
	ConfigTypInvalid ConfigTyp = iota // 占位
	ConfigTypGlobal                   // 枚举配置表
	ConfigTypNormal                   // 普通配置表
	ConfigTypI18n                     // 多语言配置表
)

type ConfigTyp int

func (t ConfigTyp) String() string {
	switch t {
	case ConfigTypGlobal:
		return "global"
	case ConfigTypNormal:
		return "normal"
	case ConfigTypI18n:
		return "i18n"
	default:
		return "invalid"
	}
}

var TypGlobalNameMatcher = regexp.MustCompile("^[gG]lobal.*$") // global
var TypI18nNameMatcher = regexp.MustCompile(`^i18n.*$`)        // i18n表
var TypNormalNameMatcher = regexp.MustCompile("^[a-z0-9]+.*$") // 普通表

// GetConfType 获取配置文件的类型
func GetConfType(bareName string) (typ ConfigTyp) {
	lower := strings.ToLower(bareName)

	// global
	if TypGlobalNameMatcher.MatchString(lower) {
		return ConfigTypGlobal
	}
	// i18n表
	if TypI18nNameMatcher.MatchString(lower) {
		return ConfigTypI18n
	}
	// 普通表
	if TypNormalNameMatcher.MatchString(lower) {
		return ConfigTypNormal
	}

	return ConfigTypInvalid
}

// FilenameByType 返回文件名
func FilenameByType(bareName string) string {
	bareName = strings.Split(bareName, ".")[0]
	typ := GetConfType(bareName)
	switch typ {
	case ConfigTypGlobal, ConfigTypNormal: // 下划线之前的内容
		idx := strings.Index(bareName, "_")
		if idx == -1 {
			return bareName
		}
		return bareName[:idx]

	case ConfigTypI18n: // 去掉下划线
		return strings.ReplaceAll(bareName, "_", "")
	default:
		log.Fatal("invalid config file type: ", bareName)
		return "INVALID"
	}
}

func NameRemoveHan(s string) string {
	return regexp.MustCompile(`\p{Han}`).ReplaceAllString(s, "") // 去掉中文
}

// NameToSnakeCase converts CamelCase to snake_case, avoiding consecutive underscores.
func NameToSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 && result[len(result)-1] != '_' {
				result = append(result, '_')
			}
			r = unicode.ToLower(r)
		}
		result = append(result, r)
	}
	return string(result)
}

// NameToCamelCase converts snake_case strings to camelCase.
func NameToCamelCase(s string) string {
	s = strings.ToLower(s)

	words := strings.Split(s, "_")
	for i := range words {
		if i > 0 && len(words[i]) > 0 {
			words[i] = strings.ToUpper(words[i][:1]) + words[i][1:]
		} else {
			words[i] = strings.ToLower(words[i])
		}
	}
	return strings.Join(words, "")
}

func NameUpperFirstLetter(name string) string {
	inputRunes := []rune(name) // 将字符串转换为 Rune 数组

	// 如果字符串非空且首字母是大写字母，则将首字母转换为小写字母
	if len(inputRunes) > 0 && unicode.IsLower(inputRunes[0]) {
		inputRunes[0] = unicode.ToUpper(inputRunes[0])
	}
	return string(inputRunes)
}
