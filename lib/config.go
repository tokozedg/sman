package sman

import (
	"github.com/spf13/viper"
)

// Config object
type Config struct {
	SnippetDir, Shell          string
	ExecConfirm, AppendHistory bool
}

func init() {
	viper.SetDefault("snippet_dir", "~/snippets")
	viper.SetDefault("append_history", "true")
	viper.SetDefault("exec_confirm", "true")
}

//GetConfig reads and return config struct
func GetConfig() (c Config) {
	c.SnippetDir = ExpandPath(viper.GetString("snippet_dir"))
	c.Shell = viper.GetString("shell")
	c.AppendHistory = viper.GetBool("append_history")
	c.ExecConfirm = viper.GetBool("exec_confirm")
	return c
}
