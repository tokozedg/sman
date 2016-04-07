package sman

import (
	"github.com/spf13/viper"
)

// Config object
type Config struct {
	SnippetDir                 string
	ExecConfirm, AppendHistory bool
}

//GetConfig reads and return config struct
func GetConfig() (c Config) {
	c.SnippetDir = ExpandPath(viper.GetString("snippet_dir"))
	c.AppendHistory = viper.GetBool("append_history")
	c.ExecConfirm = viper.GetBool("exec_confirm")
	return c
}
