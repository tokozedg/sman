package sman

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadFromCli reads CLI input and return as string
func ReadFromCli() string {
	reader := bufio.NewReader(os.Stdin)
	i, _ := reader.ReadString('\n')
	i = strings.TrimSpace(i)
	return i
}

// PrintError prints interface to stderr
func PrintError(line ...interface{}) {
	for _, l := range line {
		s := fmt.Sprint(l)
		_, _ = os.Stderr.WriteString(s)
	}
}

// PrintlnError prints interface to stderr
func PrintlnError(line ...interface{}) {
	PrintError(line...)
	PrintError("\n")
}

// DashLineError print dashes to stderr
func DashLineError() {
	PrintlnError("----")
}

// DashLine print dashes
func DashLine() {
	fmt.Println("----")
}

// CheckFileFlag returns best matched file
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
		//f = results[MakeChoice(results)]
	}
	return f
}

// ChoicePrompt return choice prompt from slice
func ChoicePrompt(from []string) (result string) {
	for i, s := range from {
		i++
		var bs, be, sp string
		if i == 1 {
			bs = "["
			be = "]"
		} else {
			bs = "("
			be = ")"
		}
		if i != len(from) {
			sp = " "
		}
		result += fmt.Sprintf("%s%v%s %s%s", bs, i, be, s, sp)
	}
	return result
}

// FullYmlPath returns yml file absolute path
func FullYmlPath(file string, dir string) string {
	return dir + "/" + file + ".yml"
}

// BaseFileName returns file name from absolute path
func BaseFileName(file string) string {
	return strings.TrimSuffix(filepath.Base(file), ".yml")
}

// CheckError checks error
func CheckError(e error, msg string) {
	if e != nil {
		PrintlnError(msg)
		panic(e)
	}
}

// SliceContains checks if var exists in a slice
func SliceContains(in []string, dst string) bool {
	for _, i := range in {
		if i == dst {
			return true
		}
	}
	return false
}
