package sman

import (
	"github.com/spf13/viper"
	"github.com/fatih/color"
)

// Config file struct
type Config struct {
	SnippetDir                 string
	ExecConfirm, AppendHistory bool
	LsFilesColor               *color.Color
}

func init() {
	viper.SetDefault("snippet_dir", "~/snippets")
	viper.SetDefault("append_history", "true")
	viper.SetDefault("exec_confirm", "true")
	viper.SetDefault("ls_color_files", "34")
}

//getConfig reads config and returns struct
func getConfig() (c Config) {
	c.SnippetDir = expandPath(viper.GetString("snippet_dir"))
	c.AppendHistory = viper.GetBool("append_history")
	c.ExecConfirm = viper.GetBool("exec_confirm")
	c.LsFilesColor = parseColor(viper.GetString("ls_color_files"))
	return c
}
