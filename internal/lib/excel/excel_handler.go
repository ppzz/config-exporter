package excel

import (
	"encoding/csv"
	"github.com/ppzz/golang-csv/internal/helper"
	"github.com/samber/lo"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// ExcelHandler 负责处理读取，输出等对接外部的操作
type ExcelHandler struct {
}

func NewExcelHandler() ExcelHandler {
	return ExcelHandler{}
}

// ReadOne 读取一个文件
func (h ExcelHandler) ReadOne(fullPath string) (*Excel, error) {
	stat, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	size := int(stat.Size()) // 读取size信息

	file, err := xlsx.OpenFile(fullPath)
	if err != nil {
		return nil, err
	}

	// 过滤出行数大于 0，列数大于 0 的效 sheet
	list := lo.Filter(file.Sheets, func(item *xlsx.Sheet, index int) bool {
		return item.MaxRow > 0 && item.MaxCol > 0
	})

	sheets := lo.Map(list, func(item *xlsx.Sheet, index int) *Sheet {
		return NewSheet(index, item)
	})

	return &Excel{
		FilePath: fullPath,
		Size:     size,
		Sheets:   sheets,
	}, nil
}

var bomBytes = []byte{0xEF, 0xBB, 0xBF}

// 一些全局参数，
var withBOM bool // Byte Order Mark ， 字节顺序标记， 用来标明这个文件是一个UTF8编码的文件，在某些WIN系统上可能需要

// ExportToCsv 输出csv文件
func (h ExcelHandler) ExportToCsv(dir string, excel *Excel) {
	bareName := helper.FileBareName(excel.FilePath)
	for _, sheet := range excel.Sheets {
		saveCsv(dir, bareName, sheet)
	}
}

func saveCsv(dir string, bareName string, sheet *Sheet) {
	csvPath := path.Join(dir, bareName+"."+strconv.Itoa(sheet.Index)+"."+strings.ToLower(sheet.Name)+".csv")
	file, err := os.OpenFile(csvPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		log.Fatal("open file failed: ", csvPath, err.Error())
	}

	defer func(f *os.File) { _ = f.Close() }(file)

	if withBOM {
		_, err = file.Write(bomBytes)
		if err != nil {
			log.Fatal("write bom failed: ", csvPath, err.Error())
		}
	}

	writer := csv.NewWriter(file)
	writer.UseCRLF = false
	for _, row := range sheet.Grid {
		if err = writer.Write(row); err != nil {
			log.Fatal("write csv failed: ", csvPath, err.Error())
		}
		writer.Flush()
	}
}

func getValDemo(index int, grid [][]string) (string, string, string) {
	first, short, long := "", "", ""

	if len(grid) == 0 {
		return first, short, long
	}

	first = grid[0][index] // 第一个值应该是下标为 0 的行，列为：下标为 index 的值
	short = grid[0][index] // 第一个值应该是下标为 0 的行，列为：下标为 index 的值
	long = grid[0][index]  // 第一个值应该是下标为 0 的行，列为：下标为 index 的值
	for _, row := range grid {
		cell := row[index]
		if len(cell) > len(long) {
			long = cell
			continue
		}
		if len(cell) < len(short) {
			short = cell
		}
	}
	return first, short, long
}

// getMetaVariableName   生成 csv meta 的变量的名字
func getMetaVariableName(filename string) string {
	inputRunes := []rune(filename) // 将字符串转换为 Rune 数组

	// 如果字符串非空且首字母是大写字母，则将首字母转换为小写字母
	if len(inputRunes) > 0 && unicode.IsUpper(inputRunes[0]) {
		inputRunes[0] = unicode.ToLower(inputRunes[0])
	}

	// 将 Rune 数组转换回字符串
	return string(inputRunes) + "CsvMeta"
}

func getGoType(typ string) string {
	typ = strings.ToLower(typ)
	m := map[string]string{
		"other":      "string",
		"i18n":       "string",
		"i18n[]":     "string",
		"i18n[][]":   "string",
		"int":        "int",
		"int[]":      "[]int",
		"int[][]":    "[][]int",
		"string":     "string",
		"string[]":   "[]string",
		"string[][]": "[][]string",
		"intw":       "int",
		"intw[]":     "[]int",
		"intw[][]":   "[][]int",
	}
	t, exist := m[typ]
	if !exist {
		return "NOT_VALID_TYPE"
	}
	return t
}

func UpperFirstLetter(name string) string {
	inputRunes := []rune(name) // 将字符串转换为 Rune 数组

	// 如果字符串非空且首字母是大写字母，则将首字母转换为小写字母
	if len(inputRunes) > 0 && unicode.IsLower(inputRunes[0]) {
		inputRunes[0] = unicode.ToUpper(inputRunes[0])
	}
	return string(inputRunes)

}

// GetBareFileName 返回文件名，返回csv文件名; 规律: 1. 无下划线，2. 无中文字，3. 无扩展名，大小写未改变
func GetBareFileName(filename string) string {
	bareName := RemoveDirAndExt(filename)
	typ := GetConfType(filename)
	Default := "INVALID"

	switch typ {
	case TypGlobal: // 下划线之前的内容
		fallthrough
	case TypNormal: // 下划线之前的内容
		idx := strings.Index(bareName, "_")
		if idx == -1 {
			return bareName
		}
		return bareName[:idx]

	case TypI18n: // 去掉下划线
		return strings.ReplaceAll(bareName, "_", "")
	default:
		return Default
	}
}

// GetNormalConfClassName 输入文件名普通配置文件的类名
func GetNormalConfClassName(filename string) string {
	return GetBareFileName(filename)
}

// GetSnakeCaseFileName 返回蛇形裸文件名
func GetSnakeCaseFileName(filename string) string {
	return camelCaseToSnakeCase(GetBareFileName(filename))
}

var camelCaseToSnakeCaseReg = regexp.MustCompile(`[\p{Han}]`) // 正则表达式，用于匹配中文字符

// camelCaseToSnakeCase 把驼峰转为蛇形命名
func camelCaseToSnakeCase(s string) string {
	s = camelCaseToSnakeCaseReg.ReplaceAllString(s, "") // 去掉中文
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

// IsExist filepath 对应的文件是否存在
func IsExist(filepath string) bool {
	_, err := os.Stat(filepath)
	// 如果 err == nil ，则存在
	return err == nil
}
