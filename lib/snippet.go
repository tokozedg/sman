package sman

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// Snippet struct
type Snippet struct {
	Desc, Command, Name, Do, File string
	Tags                          []string
	Placeholders                  []Placeholder
}

// DisplayCommand returns formatted command
func (s *Snippet) DisplayCommand() (out string) {
	out = s.Command
	for _, p := range s.Placeholders {
		for _, t := range p.Patterns {
			out = strings.Replace(out, t,
				p.DisplayName(), -1)
		}
	}
	out = strings.TrimSpace(out)
	return out
}

// SetInputs receives slice of inputs and sets Placeholders inputs in the same order
func (s *Snippet) SetInputs(inputs []string) {
	for i, v := range inputs {
		if i > len(s.Placeholders) {
			return
		}
		s.Placeholders[i].SetInput(v)
	}
}

// ReplacePlaceholders replaces placeholders patterns by placeholders input
func (s *Snippet) ReplacePlaceholders() {
	for _, p := range s.Placeholders {
		for _, pattern := range p.Patterns {
			s.Command = strings.Replace(s.Command, pattern, p.Input, -1)
		}
	}
}

// ParseCommand reads snippet command and creates Placeholders instance
func (s *Snippet) ParseCommand() {
	r, err := regexp.Compile(`<<(\w+)(?:\((.*?)\))?(#.*?)?>>`)
	checkError(err, "Invalid regexp")
	m := r.FindAllStringSubmatch(s.Command, -1)
	for _, v := range m {
		pattern := v[0]
		name := v[1]
		options := v[2]
		desc := v[3]
		// If placeholder already exists
		if i, ok := searchPlaceholder(s.Placeholders, name); ok {
			s.Placeholders[i].AddPattern(pattern)
		} else {
			// Create new placeholder
			var p Placeholder
			p.Name = name
			p.Desc = desc
			if len(options) > 0 {
				p.ParseOptions(options)
			}
			p.AddPattern(pattern)
			s.Placeholders = append(s.Placeholders, p)
		}
	}
}

// initSnippets initializes snippet after unmarshal
func initSnippets(snippetMap map[string]Snippet, file string, dir string) (snippets SnippetSlice) {
	for n, s := range snippetMap {
		s.Name = n
		s.File = file
		if len(s.Command) == 0 {
			// Search command file
			c, _ := ioutil.ReadFile(dir + "/" + s.File + "/" + s.Name)
			if len(c) > 0 {
				s.Command = strings.TrimSpace(string(c))
			} else {
				continue
			}
		}
		s.ParseCommand()
		snippets = append(snippets, s)
	}
	return snippets
}

// filterByTag filters snippet slice by tag
func filterByTag(snippets SnippetSlice, tag string) (matched SnippetSlice) {
	for _, s := range snippets {
		if sliceContains(s.Tags, tag) {
			matched = append(matched, s)
		}
	}
	return matched
}

// snippetsInFile returns snippet slice in file
func snippetsInFile(file, dir string) (snippets SnippetSlice) {
	fullPath := dir + "/" + file + ".yml"
	snippets = initSnippets(unmarshalFile(fullPath), file, dir)
	return snippets
}

// snippetsInDir returns snippet in dir
func snippetsInDir(dir string) (snippets SnippetSlice) {
	for _, f := range ymlFiles(dir) {
		snippets = append(snippets, snippetsInFile(f, dir)...)
	}
	return snippets
}

func getSnippets(name, file, dir, tag string) SnippetSlice {
	var snippets SnippetSlice
	// file flag is defined. Use fuzzy search
	if len(file) > 0 {
		fileMatched := fSearchFileName(file, dir)
		switch len(fileMatched) {
		case 0:
			printlnError("Unable to find any file with pattern: " + file)
			os.Exit(1)
		case 1:
			snippets = snippetsInFile(fileMatched[0], dir)
		default:
			printError("Multiple files matched.")
			os.Exit(1)
		}
	} else {
		snippets = snippetsInDir(dir)
	}
	// filter snippets by tag
	if len(tag) > 0 {
		return filterByTag(snippets, tag)
	}
	return snippets
}

// SnippetSlice for sorting
type SnippetSlice []Snippet

func (s SnippetSlice) Len() int {
	return len(s)
}

func (s SnippetSlice) Less(a, b int) bool {
	return s[a].File+s[a].Name < s[b].File+s[b].Name
}

func (s SnippetSlice) Swap(a, b int) {
	s[a], s[b] = s[b], s[a]
}

// filter input by test function and return new slice with matched snippets
func (ss SnippetSlice) FilterView(test func(Snippet) bool) (ret SnippetSlice) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
