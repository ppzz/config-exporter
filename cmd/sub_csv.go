package cmd

import (
	"github.com/ppzz/golang-csv/internal/excel_to_csv/setting"
	"github.com/ppzz/golang-csv/internal/excel_to_csv/starter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// csvCmd represents the csv command
var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "生成 csv 文件",
	Long:  `根据输入的 Excel 文件生成 csv 文件`,
	Run: func(cmd *cobra.Command, args []string) {
		prepareSetting(cmd)
		starter.Start()
	},
}

func prepareSetting(cmd *cobra.Command) {
	err := viper.BindPFlag("excel", cmd.PersistentFlags().Lookup("excel"))
	cobra.CheckErr(err)
	err = viper.BindPFlag("output", cmd.PersistentFlags().Lookup("output"))
	cobra.CheckErr(err)
	setting.Get().Init()
}

func init() {
	rootCmd.AddCommand(csvCmd)

	// rootCmd.MarkPersistentFlagRequired("excel") // Optional: makes the excel flag required
	csvCmd.PersistentFlags().StringP("excel", "e", "", "specify the excel file path")
	csvCmd.PersistentFlags().StringP("output", "o", "", "specify the csv file path")
}
