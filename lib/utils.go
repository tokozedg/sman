package sman

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadFromCli() string {
	reader := bufio.NewReader(os.Stdin)
	i, _ := reader.ReadString('\n')
	i = strings.TrimSpace(i)
	return i
}

func PrintError(line ...interface{}) {
	for _, l := range line {
		s := fmt.Sprint(l)
		os.Stderr.WriteString(s)
	}
}

func PrintlnError(line ...interface{}) {
	PrintError(line...)
	PrintError("\n")
}

func DashLineError() {
	PrintlnError("----")
}

func DashLine() {
	fmt.Println("----")
}

func CheckFileFlag(file string, dir string) (f string) {
	results := FSearchFileName(file, dir)
	switch len(results) {
	case 0:
		PrintlnError("Unable to find file: " + fileFlag)
		os.Exit(1)
	case 1:
		f = results[0]
	default:
		PrintError("Multiple file found. Choose one")
		f = results[MakeChoise(results)]
	}
	return f
}

func MakeChoise(from []string) (result int) {
	return result
}

func FullSnippetPath(file string, dir string) string {
	return dir + "/" + file + ".yml"
}

func FullCommandPath(name string, dir string) string {
	return dir + "/commands/" + name
}
func BaseFileName(file string) string {
	return strings.TrimSuffix(filepath.Base(file), ".yml")
}

func CheckError(e error, msg string) {
	if e != nil {
		PrintlnError(msg)
		panic(e)
	}
}

func SliceContains(in []string, dst string) bool {
	for _, i := range in {
		if i == dst {
			return true
		}
	}
	return false
}
