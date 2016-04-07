package sman

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	cfgFile, fileFlag, tagFlag string
)

//RootCmd for cobra
var RootCmd = &cobra.Command{
	Use:   "sman",
	Short: "CLI Snippet Manager",
	Long:  ``,
}

//Execute for cobra
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().BoolVarP(&color.NoColor, "no-color", "", false, "disable colors")
	RootCmd.PersistentFlags().StringVarP(&fileFlag, "file", "f", "", "snippet file")
	RootCmd.PersistentFlags().StringVarP(&tagFlag, "tags", "t", "", "tags filter")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}
	viper.SetConfigName(".sman")
	viper.AddConfigPath("$HOME")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("sman")
	_ = viper.ReadInConfig()
}
