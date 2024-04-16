package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// formatCmd represents the format command
var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "csv 文件格式化处理",
	Long:  `根据文件名, 文件内容等信息对 csv 文件进行格式化处理`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("format called")
		// 1. csv header name 补全
		// 2. csv type 行 映射 到通用的type 行
		// 3. csv column 过滤(根据c/s标记),
		// 4. csv column 可以加备注行等信息

	},
}

func init() {
	rootCmd.AddCommand(formatCmd)

	// csvCmd.PersistentFlags().StringP("excel", "e", "", "specify the excel file path")
	// csvCmd.PersistentFlags().StringP("output", "o", "", "specify the csv file path")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// formatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// formatCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
