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
)

func shellString(s string) string {
	return `$'` + s + `'`
}
func appendHistory(snippet Snippet) {
	histLine := "s run -f " + snippet.File + " " + snippet.Name
	for _, p := range snippet.Placeholders {
		if len(p.Input) == 0 {
			return
		}
		histLine += " " + shellString(p.Input)
	}
	fmt.Println("_append_history ", histLine)
}

func executeConfirmed() bool {
	dashLineError()
	msg := "Execute Snippet? [Y/n]: "
	for {
		printError(msg)
		in := readFromCli()
		if sliceContains([]string{
			"N", "n", "no", "NO"}, in) {
			return false
		} else if sliceContains([]string{
			"Y", "y", "yes", "Yes", "YES", ""}, in) {
			return true
		}
	}
}

func execute(cmd string, confirm bool) {
	if confirm && !executeConfirmed() {
		return
	}
	printlnError("Executing...")
	fmt.Println(strings.TrimSpace(cmd))
}

func requestInput(snippet *Snippet) {
	for i := range snippet.Placeholders {
		p := &snippet.Placeholders[i]
		if len(p.Input) > 0 {
			continue
		}
		if len(p.Desc) > 0 {
			printlnError(p.Desc)
		}
		printError(p.DisplayName())
		if len(p.Options) > 0 {
			printError(" ", choicePrompt(p.Options))
		}
		printError(": ")
		r := readFromCli()
		p.SetInput(r)
	}
}

func run(name string, inputs ...string) {
	c := getConfig()
	snippets := getSnippets(name, fileFlag, c.SnippetDir, tagFlag)
	matchedSnippets := fSearchSnippet(snippets, name)
	var snippet Snippet
	switch len(matchedSnippets) {
	case 0:
		printlnError("No snippets matched...")
		os.Exit(1)
	case 1:
		snippet = matchedSnippets[0]
	default:
		printlnError("Multiple snippets matched...")
		os.Exit(1)
	}

	if len(inputs) > len(snippet.Placeholders) {
		printlnError("You gave ", len(inputs), " argument(s), limit is ", len(snippet.Placeholders))
		os.Exit(1)
	}

	snippet.SetInputs(inputs)
	dashLineError()
	if len(inputs) < len(snippet.Placeholders) {
		printlnError(snippet.DisplayCommand())
		dashLineError()
		requestInput(&snippet)
		dashLineError()
	}
	snippet.ReplacePlaceholders()
	printlnError(snippet.Command)
	dashLineError()
	if c.AppendHistory {
		appendHistory(snippet)
	}
	if printFlag {
		return
	}
	if copyFlag || (!execFlag && snippet.Do == "copy") {
		err := clipboard.WriteAll(snippet.Command)
		checkError(err, "Error while copying")
		printlnError("Snippet Copied...")
	}
	if execFlag || (!copyFlag && snippet.Do == "exec") {
		execute(snippet.Command, c.ExecConfirm)
	}
}

var runCmd = &cobra.Command{
	Use:     "run [-f FILE] [-f TAG] SNIPPET [PLACEHOLDER VALUES...] [-cxp]",
	Aliases: []string{"r"},
	Short:   "Run snippet",
	Long: `
Runs snippet and execute action specified with flags or snippet "do".
Flags action overrides "do".

The next arguments after snippet name will be used to fill placeholder values.
The first argument will be considered as a first placeholder input.
Use "s show <snippet>" to get placeholder order numbers.

Examples:
s run -f mysql db:dump -x
	- run 'db:dump' snippet from file 'mysql' and execute when done
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			printlnError("need snippet name...")
			os.Exit(1)
		}
		run(args[0], args[1:]...)
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
	RootCmd.SetOutput(os.Stderr)
	runCmd.Flags().BoolVarP(&copyFlag, "copy", "c", false, "copy snippet")
	runCmd.Flags().BoolVarP(&execFlag, "exec", "x", false, "execute snippet")
	runCmd.Flags().BoolVarP(&printFlag, "print", "p", false, "print snippet")
}
