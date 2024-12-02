package excel

import (
	"github.com/ppzz/config-exporter/internal/helper"
)

// 与 [][]string 相关的操作

// RemoveEmptyCol 删除空列
func RemoveEmptyCol(grid [][]string) [][]string {
	if len(grid) == 0 {
		return grid
	}
	hasVal := helper.GridDoesColumnHasValue(grid)
	return helper.GridFilterColumn(grid, hasVal)
}
