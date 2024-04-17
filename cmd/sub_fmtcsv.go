package cmd

import (
	"github.com/ppzz/golang-csv/internal/format_csv/setting"
	"github.com/ppzz/golang-csv/internal/format_csv/starter"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// formatCmd represents the format command
var formatCmd = &cobra.Command{
	Use:   "fmtcsv",
	Short: "csv 文件格式化处理",
	Long:  `根据文件名, 文件内容等信息对 csv 文件进行格式化处理`,
	Run: func(cmd *cobra.Command, args []string) {
		formatCmdPrepareSetting(cmd)
		starter.Start()
	},
}

func formatCmdPrepareSetting(cmd *cobra.Command) {
	err := viper.BindPFlag("csv", cmd.Flags().Lookup("csv"))
	cobra.CheckErr(err)
	err = viper.BindPFlag("fmtcsv", cmd.Flags().Lookup("fmtcsv"))
	cobra.CheckErr(err)
	setting.Get().Init()
}

func init() {
	rootCmd.AddCommand(formatCmd)
	formatCmd.Flags().StringP("csv", "c", "", "specify the csv file path")
	formatCmd.Flags().StringP("fmtcsv", "f", "", "specify the formatted csv file path")
}
