package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var descLong = `游戏配置表的处理程序,包括: Excel转Csv, Csv 过滤, 根据 Csv 生成 Go 代码`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:              "golang-csv",
	Short:            "关于csv的配置表处理",
	Long:             descLong,
	TraverseChildren: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./csv-gen.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	filename := "csv-gen"
	fileExt := "yaml"
	filepath := "."
	envPrefix := "APP"
	envReplacer := strings.NewReplacer(".", "_")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		// home, err := os.UserHomeDir()
		// cobra.CheckErr(err)

		viper.SetConfigName(filename)
		viper.SetConfigType(fileExt)
		viper.AddConfigPath(filepath)
	}
	viper.SetEnvPrefix(envPrefix)        // 设置环境变量前缀
	viper.SetEnvKeyReplacer(envReplacer) // 环境变量的分割符号
	viper.AutomaticEnv()                 // read in environment variables that match

	err := viper.ReadInConfig()
	cobra.CheckErr(err)
}
