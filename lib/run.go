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

func appendHistory(snippet Snippet) {
	histLine := "s run " + snippet.Name
	for _, p := range snippet.Placeholders {
		if len(p.Input) == 0 {
			return
		}
		arg := p.Input
		arg = strings.Replace(arg, `"`, `\"`, -1)
		arg = strings.Replace(p.Input, `$`, `\$`, -1)
		// If Input Conaints ' - Choose Escaped Version
		if strings.Contains(p.Input, `'`) {
			arg = strings.Replace(arg, `'`, `\\\\\'`, -1)
			histLine += fmt.Sprintf(` $'$\'%v\'' `, arg)
		} else {
			//Normal Version
			histLine += fmt.Sprintf(` "'%v'" `, arg)
		}
	}
	// TODO bash and zsh versions differ
	//sh = "history -s " + histLine + ";"

	fmt.Println("print -s " + histLine + ";")
}

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
	for i := range snippet.Placeholders {
		p := &snippet.Placeholders[i]
		if len(p.Input) > 0 {
			continue
		}
		if len(p.Desc) > 0 {
			PrintlnError(p.Desc)
		}
		PrintError(p.DisplayName())
		if len(p.Options) > 0 {
			PrintError(" ", ChoicePrompt(p.Options))
		}
		PrintError(": ")
		r := ReadFromCli()
		p.SetInput(r)
	}
}

//func chooseSnippet(from SnippetSlice) Snippet {
//var choices []string
//for _, s := range from {
//k := s.File + ":" + s.Name
//choices = append(choices, k)
//}
//PrintlnError(ChoicePrompt(choices))
//c := ReadFromCli()
//if c
//return from[c]
//}

func run(name string, inputs ...string) {
	c := GetConfig()
	snippets := GetSnippets(name, fileFlag, c.SnippetDir, tagFlag)
	matchedSnippets := FSearchSnippet(snippets, name)
	var snippet Snippet
	switch len(matchedSnippets) {
	case 0:
		PrintlnError("No snippets matched...")
		os.Exit(1)
	case 1:
		snippet = matchedSnippets[0]
	default:
		PrintlnError("Multiple snippets matched...")
		os.Exit(1)
		//snippet = chooseSnippet(matchedSnippets)
	}
	snippet.SetInputs(inputs)
	DashLineError()
	if len(inputs) < len(snippet.Placeholders) {
		PrintlnError(snippet.DisplayCommand())
		DashLineError()
		requestInput(&snippet)
		DashLineError()
	}
	snippet.ReplacePlaceholders()
	PrintlnError(snippet.Command)
	if (!printFlag && copyFlag) || (!execFlag && snippet.Do == "copy") {
		DashLineError()
		err := clipboard.WriteAll(snippet.Command)
		CheckError(err, "Error while copying")
		PrintlnError("Snippet Copied...")
	}
	if (!printFlag && execFlag) || (!copyFlag && snippet.Do == "exec") {
		execute(snippet.Command, c.ExecConfirm)
	}
	if c.AppendHistory {
		appendHistory(snippet)
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
			PrintlnError("Run what?")
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
