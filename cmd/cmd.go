package cmd

import (
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var descLong = `游戏配置表的处理程序,包括: Excel转Csv, Csv 过滤, 根据 Csv 生成 Go 代码`

var rootCmd = &cobra.Command{
	Use:              "golang-csv",
	Short:            "关于csv的配置表处理",
	Long:             descLong,
	TraverseChildren: true,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config-exporter.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	envPrefix := "APP"
	envReplacer := strings.NewReplacer(".", "_")

	defCfgName := "config-exporter"
	defCfgExt := "yaml"
	defCfgDir := "."

	filePath := path.Join(defCfgDir, defCfgName+"."+defCfgExt)
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		err := viper.ReadInConfig()
		cobra.CheckErr(err)
	} else if ExistFile(filePath) {
		viper.SetConfigName(defCfgName)
		viper.SetConfigType(defCfgExt)
		viper.AddConfigPath(defCfgDir)
		err := viper.ReadInConfig()
		cobra.CheckErr(err)
	}
	viper.SetEnvPrefix(envPrefix)        // 设置环境变量前缀
	viper.SetEnvKeyReplacer(envReplacer) // 环境变量的分割符号
	viper.AutomaticEnv()                 // read in environment variables that match

}

// ExistFile returns true if a file or directory exists.
func ExistFile(name string) bool {
	info, err := os.Stat(name)
	return err == nil && info != nil && !info.IsDir()
}
