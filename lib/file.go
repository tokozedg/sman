package sman

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// baseFileName receives full path and returns file name without extension
func baseFileName(file string) string {
	return strings.TrimSuffix(filepath.Base(file), ".yml")
}

// expandPath receives path string and returns absolute path
// such as expanding ~ to user home path.
func expandPath(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	if len(p) > 2 && p[:2] == "~/" {
		return strings.Replace(p, "~", os.Getenv("HOME"), 1)
	}
	p, err := filepath.Abs(p)
	checkError(err, "Can find absolute path for: "+p)
	return p
}

// ymlFiles returns slice of yml files names in a dir without extension
func ymlFiles(dir string) (files []string) {
	fio, _ := ioutil.ReadDir(dir)
	for _, f := range fio {
		if filepath.Ext(f.Name()) == ".yml" {
			n := baseFileName(f.Name())
			files = append(files, n)
		}
	}
	sort.Strings(files)
	return files
}

// unmarshalFile reads yml file and returns snippet objects map
// where snippet name is the map key and Snippet instance is the value
func unmarshalFile(file string) (snippetsMap map[string]Snippet) {
	yamlFile, err := ioutil.ReadFile(file)
	checkError(err, "Can't read file: "+file)
	err = yaml.Unmarshal(yamlFile, &snippetsMap)
	checkError(err, "Can't unmarshal file: "+file)
	return snippetsMap
}
