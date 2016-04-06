package sman

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"sort"
	"text/tabwriter"
)

func filterSnippets(p string, slice SnippetSlice) (matched SnippetSlice) {
	r, err := regexp.Compile(p)
	CheckError(err, "Invalid search pattern")
	for _, s := range slice {
		if r.MatchString(s.Name) ||
			r.MatchString(s.Command) ||
			r.MatchString(s.Desc) {
			matched = append(matched, s)
		}
	}
	return matched
}

func doLs(pattern string) {
	c := GetConfig()
	snippets := GetSnippets(pattern, fileFlag, c.SnippetDir, tagFlag)
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 25, 2, 0, ' ', 0)
	snippets = filterSnippets(pattern, snippets)
	sort.Sort(snippets)
	var prevFile string
	blue := color.New(color.FgBlue).SprintFunc()
	for _, s := range snippets {
		if s.File != prevFile {
			fmt.Fprintln(w, blue(s.File+":"))
			prevFile = s.File
		}
		line := fmt.Sprintf("   %v\t[%v]\t%v", s.Name, s.DisplayTags(),
			s.DisplayDesc())
		fmt.Fprintln(w, line)
	}
	w.Flush()
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:     "ls <pattern>",
	Aliases: []string{"l"},
	Short:   "List or search in all available snippets",
	Long: `
List or search in all available snippets.

<pattern>
Uses regexp. Pattern is matched against snippet name, description and command.

Examples:
s ls docker
	- Show all snippet matching "docker"
s ls -f docker
	- Show all snippets in file "docker.yml"
s ls -t docker,cli
	- Show all snippets tagged with "docker" OR "cli"

	`,
	Run: func(cmd *cobra.Command, args []string) {
		var p string
		if len(args) > 0 {
			p = args[0]
		}
		doLs(p)
	},
}

func init() {
	RootCmd.AddCommand(lsCmd)
}
