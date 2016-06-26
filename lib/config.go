package sman

import (
	"github.com/spf13/viper"
)

// Config file struct
type Config struct {
	SnippetDir                 string
	ExecConfirm, AppendHistory bool
}

func init() {
	viper.SetDefault("snippet_dir", "~/snippets")
	viper.SetDefault("append_history", "true")
	viper.SetDefault("exec_confirm", "true")
}

//getConfig reads config and returns struct
func getConfig() (c Config) {
	c.SnippetDir = expandPath(viper.GetString("snippet_dir"))
	c.AppendHistory = viper.GetBool("append_history")
	c.ExecConfirm = viper.GetBool("exec_confirm")
	return c
}
