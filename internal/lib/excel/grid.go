package excel

import (
	"github.com/ppzz/golang-csv/internal/helper"
)

// 与 [][]string 相关的操作

// RemoveEmptyCol 删除空列
func RemoveEmptyCol(grid [][]string) [][]string {
	hasVal := helper.GridDoesColumnHasValue(grid)
	return helper.GridFilterColumn(grid, hasVal)
}
