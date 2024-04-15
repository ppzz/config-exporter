package excel

import (
	"path"
	"regexp"
	"strings"
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

// 配置表类型
const (
	TypInvalid Typ = iota // 占位
	TypGlobal             // 全局配置表
	TypNormal             // 普通配置表
	TypI18n               // 多语言配置表
)

// 匹配不同类型表名的正则表达式：

var TypGlobalNameMatcher = regexp.MustCompile("global_.+$")    // global
var TypI18nNameMatcher = regexp.MustCompile(`i18n_.+$`)        // i18n表
var TypNormalNameMatcher = regexp.MustCompile("[a-z0-9]+_.+$") // 普通表

// GetConfType 获取配置文件的类型
func GetConfType(filename string) (typ Typ) {
	basename := RemoveDirAndExt(filename)
	lower := strings.ToLower(basename)

	global := TypGlobalNameMatcher.MatchString(lower) // global
	if global {
		return TypGlobal
	}
	i18n := TypI18nNameMatcher.MatchString(lower) // i18n表
	if i18n {
		return TypI18n
	}
	normal := TypNormalNameMatcher.MatchString(lower) // 普通表
	if normal {
		return TypNormal
	}

	return TypInvalid
}

// RemoveDir 去掉目录部分
func RemoveDir(p string) string {
	return path.Base(p)
}

// RemoveExtension 去掉扩展名部分
func RemoveExtension(p string) string {
	ext := path.Ext(p)
	return p[:len(p)-len(ext)]
}

// RemoveDirAndExt 去掉目录部分，去掉扩展名部分
func RemoveDirAndExt(fullPath string) string {
	return RemoveExtension(RemoveDir(fullPath))
}
