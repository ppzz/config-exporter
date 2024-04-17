package cmd

import (
	"github.com/ppzz/golang-csv/internal/gen_code_go/setting"
	"github.com/ppzz/golang-csv/internal/gen_code_go/starter"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// codeGoCmd represents the code command
var codeGoCmd = &cobra.Command{
	Use:   "code_go",
	Short: "gen code",
	Long:  `gen go code via csv file`,
	Run: func(cmd *cobra.Command, args []string) {
		goCmdPrepareSetting(cmd)
		starter.Start()
	},
}

func goCmdPrepareSetting(cmd *cobra.Command) {
	err := viper.BindPFlag("fmtcsv", cmd.Flags().Lookup("fmtcsv"))
	cobra.CheckErr(err)
	err = viper.BindPFlag("code_go", cmd.Flags().Lookup("go"))
	cobra.CheckErr(err)
	setting.Get().Init()
}

func init() {
	rootCmd.AddCommand(codeGoCmd)
	codeGoCmd.Flags().StringP("fmtcsv", "f", "", "specify the csv file path")
	codeGoCmd.Flags().StringP("go", "g", "", "specify the code path")
}