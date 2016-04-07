package sman

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ExpandPath receives path string and extracts to absolute path
func ExpandPath(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	if len(p) > 2 && p[:2] == "~/" {
		p = strings.Replace(p, "~", os.Getenv("HOME"), 1)
	}
	return p
}

// ReadFile reads file and returns string content
func ReadFile(file string) string {
	c, _ := ioutil.ReadFile(file)
	return string(c)
}

// YmlFiles return slice of yml files in a dir
func YmlFiles(dir string) (files []string) {
	fio, _ := ioutil.ReadDir(dir)
	for _, f := range fio {
		if filepath.Ext(f.Name()) == ".yml" {
			n := strings.TrimSuffix(f.Name(), ".yml")
			files = append(files, n)
		}
	}
	sort.Strings(files)
	return files
}

// UnmarshalFile reads file and returns snippet objects map
// where name is the map key and Snippet is the value
func UnmarshalFile(file string) (snippetsMap map[string]Snippet) {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		PrintlnError("Can't read file: " + file + " for unmarshal")
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, &snippetsMap)
	if err != nil {
		return
	}
	return snippetsMap
}
