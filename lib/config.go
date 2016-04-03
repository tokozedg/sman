package sman

import (
	"github.com/spf13/viper"
)

type Config struct {
	SnippetDir                 string
	ExecConfirm, AppendHistory bool
}

func GetConfig() (c Config) {
	c.SnippetDir = ExpandPath(viper.GetString("snippet_dir"))
	c.AppendHistory = viper.GetBool("append_history")
	c.ExecConfirm = viper.GetBool("exec_confirm")
	return c
}
