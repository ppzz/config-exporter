package excel

import (
	"github.com/ppzz/config-exporter/internal/helper"
)

// 与 [][]string 相关的操作

// RemoveEmptyCol 删除空列
func RemoveEmptyCol(grid [][]string) [][]string {
	hasVal := helper.GridDoesColumnHasValue(grid)
	return helper.GridFilterColumn(grid, hasVal)
}
