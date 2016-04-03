package sman

import (
	"regexp"
	"strings"
)

type Snippet struct {
	Desc, Command, Name, Do, File string
	Tags                          []string `,flow`
	Placeholders                  []Placeholder
}

func (snippet *Snippet) SetInputs(inputs []string) {
	for i, v := range inputs {
		if i > len(snippet.Placeholders) {
			return
		}
		snippet.Placeholders[i].Input = v
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

func initSnippets(snippetMap map[string]Snippet, file string, snippetDir string) (snippets SnippetSlice) {
	for n, s := range snippetMap {
		s.Name = n
		s.File = file
		if len(s.Desc) == 0 {
			s.Desc = "----"
		}
		if len(s.Command) == 0 {
			c := SearchCommandFile(s.Name, snippetDir)
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

func (snippet *Snippet) DisplayCommand() (out string) {
	out = snippet.Command
	for _, p := range snippet.Placeholders {
		for _, t := range p.Patterns {
			out = strings.Replace(out, t,
				p.DisplayName(), -1)
		}
	}
	out = strings.TrimSpace(out)
	return out
}

func SearchSnippetInFile(name string, file string, snippetDir string) SnippetSlice {
	var snippets SnippetSlice
	fullPath := FullSnippetPath(file, snippetDir)
	snippets = initSnippets(UnmarshalFile(fullPath), file, snippetDir)
	return FSearchSnippet(snippets, name)
}

func SearchSnippetInDir(name string, dir string) (matched SnippetSlice) {
	for _, f := range YmlFiles(dir) {
		matched = append(matched, SearchSnippetInFile(name, f, dir)...)
	}
	return matched
}

type SnippetSlice []Snippet

func (s SnippetSlice) Len() int {
	return len(s)
}

func (s SnippetSlice) Less(a, b int) bool {
	return s[a].Name < s[b].Name
}

func (s SnippetSlice) Swap(a, b int) {
	s[a], s[b] = s[b], s[a]
}
