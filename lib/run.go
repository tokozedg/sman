package sman

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	copyFlag, execFlag, printFlag bool
	fileFlag                      string
)

func executeConfirmed() bool {
	DashLineError()
	msg := "Execute Snippet? [Y/n]: "
	for {
		PrintError(msg)
		in := ReadFromCli()
		if SliceContains([]string{
			"N", "n", "no", "NO"}, in) {
			return false
		} else if SliceContains([]string{
			"Y", "y", "yes", "Yes", "YES", ""}, in) {
			return true
		}
	}
}

func execute(cmd string, confirm bool) {
	if confirm && !executeConfirmed() {
		return
	}
	DashLineError()
	PrintlnError("Executing...")
	fmt.Println(strings.TrimSpace(cmd))
}

func requestInput(snippet *Snippet) {
	for i, _ := range snippet.Placeholders {
		p := &snippet.Placeholders[i]
		if len(p.Input) > 0 {
			continue
		}
		if len(p.Desc) > 0 {
			PrintlnError(p.Desc)
		}
		PrintError(p.DisplayName(), ": ")
		r := ReadFromCli()
		p.Input = r
	}
}

func run(name string, inputs ...string) {
	c := GetConfig()
	var fileName string
	var matchedSnippets SnippetSlice
	if len(fileFlag) > 0 {
		fileName = CheckFileFlag(fileFlag, c.SnippetDir)
		matchedSnippets = SearchSnippetInFile(name, fileName, c.SnippetDir)
	} else {
		matchedSnippets = SearchSnippetInDir(name, c.SnippetDir)
	}
	var snippet Snippet
	switch len(matchedSnippets) {
	case 0:
		PrintlnError("No snippets matched...")
		os.Exit(1)
	case 1:
		snippet = matchedSnippets[0]
	default:
		//snippet = chooseSnippet(matchedSnippets)
	}
	snippet.SetInputs(inputs)
	DashLineError()
	if len(inputs) < len(snippet.Placeholders) {
		PrintlnError(snippet.DisplayCommand())
		requestInput(&snippet)
		DashLineError()
	}
	snippet.ReplacePlaceholders()
	PrintlnError(snippet.Command)
	if printFlag {
		return
	}
	if (copyFlag) || (!execFlag && snippet.Do == "copy") {
		clipboard.WriteAll(snippet.Command)
		PrintlnError("Snippet Copied...")
	}
	if (execFlag) || (!copyFlag && snippet.Do == "exec") {
		execute(snippet.Command, c.ExecConfirm)
	}
}

var runCmd = &cobra.Command{
	Use:     "run <snippet> [<placeholder1>] [<placeholder2>]...",
	Aliases: []string{"r"},
	Short:   "Run snippet",
	Long: `
Run snippet and execute action specified with flags or "do" variable
in yml file. Flags action overrides "do" action defined in yml.

The next arguments after snippet name will be used to fill placeholder values.
Use "s show <snippet>" to get placeholder order numbers.

Examples:
s run -f mysql db/dump -x
	- run db/dump snippet from file mysql and execute when done
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Run what?")
			os.Exit(1)
		}
		run(args[0], args[1:]...)
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&copyFlag, "copy", "c", false, "copy snippet")
	runCmd.Flags().BoolVarP(&execFlag, "exec", "x", false, "execute snippet")
	runCmd.Flags().BoolVarP(&printFlag, "print", "p", false, "print snippet")
	runCmd.Flags().StringVarP(&fileFlag, "file", "f", "", "snippet file")
}
