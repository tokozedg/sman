package sman

import (
	//"fmt"
	"github.com/fatih/color"
	"regexp"
	"strings"
)

type Snippet struct {
	Desc, Command, Name, Do, File string
	Tags                          []string `,flow`
	Placeholders                  []Placeholder
}

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

func (s *Snippet) DisplayTags() string {
	if len(s.Tags) > 0 {
		return strings.Join(s.Tags, " | ")
	} else {
		return "----"
	}
}

func (s *Snippet) DisplayDesc() string {
	return strings.Title(s.Desc)
}

func (s *Snippet) DisplayDo() string {
	if len(s.Do) > 0 {
		return strings.Title(s.Do)
	} else {
		return "----"
	}
}

func (s *Snippet) DisplayFile() string {
	magenta := color.New(color.FgMagenta).SprintFunc()
	return magenta(s.File)
}

func (snippet *Snippet) SetInputs(inputs []string) {
	for i, v := range inputs {
		if i > len(snippet.Placeholders) {
			return
		}
		snippet.Placeholders[i].SetInput(v)
	}
}

func (s *Snippet) ReplacePlaceholders() {
	for _, p := range s.Placeholders {
		for _, pattern := range p.Patterns {
			s.Command = strings.Replace(s.Command, pattern, p.Input, -1)
		}
	}
}

func (s *Snippet) parseCommand() {
	r, err := regexp.Compile(`<<(\w+)(?:\((.*?)\))?(#.*?)?>>`)
	CheckError(err, "Invalid regexp")
	m := r.FindAllStringSubmatch(s.Command, -1)
	for _, v := range m {
		pattern := v[0]
		name := v[1]
		options := v[2]
		desc := v[3]
		// If placeholder already exists
		if i, ok := SearchPlaceholder(s.Placeholders, name); ok {
			s.Placeholders[i].AddPattern(pattern)
		} else {
			// Create new placeholder
			var p Placeholder
			p.Name = name
			p.Desc = desc
			if len(options) > 0 {
				p.Options = ParseOptions(options)
			}
			p.AddPattern(pattern)
			s.Placeholders = append(s.Placeholders, p)
		}
	}
}

func initSnippets(snippetMap map[string]Snippet, file string, dir string) (snippets SnippetSlice) {
	for n, s := range snippetMap {
		s.Name = n
		s.File = file
		if len(s.Desc) == 0 {
			s.Desc = "----"
		}
		if len(s.Command) == 0 {
			c := ReadFile(dir + s.File + "/" + s.Name)
			if len(c) > 0 {
				s.Command = strings.TrimSpace(c)
			} else {
				continue
			}
		}
		s.parseCommand()
		snippets = append(snippets, s)
	}
	return snippets
}

func filterByTag(snippets SnippetSlice, tag string) (matched SnippetSlice) {
	for _, s := range snippets {
		if SliceContains(s.Tags, tag) {
			matched = append(matched, s)
		}
	}
	return matched
}

func snippetsInFile(file, dir string) (snippets SnippetSlice) {
	matchedFile := CheckFileFlag(file, dir)
	fullPath := FullSnippetPath(matchedFile, dir)
	snippets = initSnippets(UnmarshalFile(fullPath), file, dir)
	return snippets
}

func snippetsInDir(dir string) (snippets SnippetSlice) {
	for _, f := range YmlFiles(dir) {
		snippets = append(snippets, snippetsInFile(f, dir)...)
	}
	return snippets
}

func GetSnippets(name, file, dir, tag string) SnippetSlice {
	var snippets SnippetSlice
	if len(file) > 0 {
		snippets = snippetsInFile(file, dir)
	} else {
		snippets = snippetsInDir(dir)
	}
	if len(tag) > 0 {
		return filterByTag(snippets, tag)
	}
	return snippets
}

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

func SnippetNames(slice SnippetSlice) (names []string, snippetMap map[string]Snippet) {
	snippetMap = make(map[string]Snippet)
	for _, s := range slice {
		names = append(names, s.Name)
		snippetMap[s.Name] = s
	}
	return names, snippetMap
}
