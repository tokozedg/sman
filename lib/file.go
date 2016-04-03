package sman

import (
	//"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ExpandPath(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	if len(p) > 2 && p[:2] == "~/" {
		p = strings.Replace(p, "~", os.Getenv("HOME"), 1)
	}
	return p
}

func SearchCommandFile(name string, dir string) (content string) {
	p := FullCommandPath(name, dir)
	d, e := ioutil.ReadFile(p)
	CheckError(e, "Can't open command file: "+p)
	content = string(d)
	return content
}

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

func UnmarshalFile(file string) (snippetsMap map[string]Snippet) {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		PrintlnError("Can't read file: " + file + " for unmarshal")
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, &snippetsMap)
	if err != nil {
		PrintlnError("Can't unmarshal file: " + file)
		panic(err)
	}
	return snippetsMap
}
