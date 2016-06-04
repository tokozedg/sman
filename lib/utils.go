package sman

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
)

// readFromCli reads CLI input and return string
func readFromCli() string {
	reader := bufio.NewReader(os.Stdin)
	i, _ := reader.ReadString('\n')
	i = strings.TrimSpace(i)
	return i
}

func printError(line ...interface{}) {
	for _, l := range line {
		s := fmt.Sprint(l)
		_, _ = os.Stderr.WriteString(s)
	}
}

func printlnError(line ...interface{}) {
	printError(line...)
	printError("\n")
}

func dashLineError() {
	printlnError("----")
}

func dashLine() {
	fmt.Println("----")
}

// displayString titles string or returns dashes if empty
func displayString(s string) string {
	if len(s) > 0 {
		return strings.Title(s)
	}
	return "----"
}

// displaySlice joins slice and returns titled string or dashes if empty
func displaySlice(s []string) string {
	return displayString(strings.Join(s, " | "))
}

func magenta(s string) string {
	m := color.New(color.FgMagenta).SprintFunc()
	return m(s)
}

func cyan(s string) string {
	c := color.New(color.FgCyan).SprintFunc()
	return c(s)
}

func choicePrompt(from []string) (result string) {
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

func checkError(e error, msg string) {
	if e != nil {
		printlnError(msg)
		panic(e)
	}
}

// sliceContains checks if string exists in slice
func sliceContains(in []string, dst string) bool {
	for _, i := range in {
		if i == dst {
			return true
		}
	}
	return false
}
