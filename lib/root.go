package sman

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfgFile string

//RootCmd for cobra
var RootCmd = &cobra.Command{
	Use:   "sman",
	Short: "CLI Snippet Manager",
	Long:  ``,
}

//Execute for cobra
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.SetOutput(os.Stderr)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}
	viper.SetConfigName(".sman")
	viper.AddConfigPath("$HOME")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("sman")
	viper.ReadInConfig()
}
