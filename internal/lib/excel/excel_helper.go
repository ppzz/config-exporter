package excel

import (
	"github.com/samber/lo"
	"regexp"
	"strconv"
)

var MetaDataRowCount = 4 // 最小的行数，小于这个行数的表可能无意义

func fixMetaEmptyName(grid [][]string) {
	for i := 0; i < len(grid[0]); i++ {
		if grid[0][i] == "" {
			grid[0][i] = "_col" + strconv.Itoa(i+1)
		}
	}
}

func IsAutoGenColumnName(name string) bool {
	e := regexp.MustCompile(`_col[\d]+`)
	return e.MatchString(name)
}

// GenColumnMeta 根据meta信息返回column对象
func GenColumnMeta(grid [][]string) []ColumnMeta {
	names, descriptions, types, flags := grid[0], grid[1], grid[2], grid[3]
	return lo.Map(names, func(name string, index int) ColumnMeta {
		return NewColumnMeta(index, name, descriptions[index], types[index], flags[index])
	})
}
