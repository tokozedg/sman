package sman

import (
	"reflect"
	"testing"
)

func TestSnippetReplacePlaceholders(t *testing.T) {
	tests := []struct {
		name         string
		command      string
		placeholders []Placeholder
		wantCommand  string
	}{
		{"no placeholder", "hello world", []Placeholder(nil), "hello world"},
		{"single patterns", "hello <<name>>",
			[]Placeholder{
				Placeholder{
					Name:     "name",
					Patterns: []string{"<<name>>"},
					Input:    "test",
				},
			},
			"hello test",
		},
		{"multiple patterns", "hello <<name#desc>> sup <<name>>",
			[]Placeholder{
				Placeholder{
					Name:     "name",
					Patterns: []string{"<<name#desc>>", "<<name>>"},
					Input:    "test",
				},
			},
			"hello test sup test",
		},
	}
	for _, tt := range tests {
		s := &Snippet{
			Command:      tt.command,
			Placeholders: tt.placeholders,
		}
		s.ReplacePlaceholders()
		if s.Command != tt.wantCommand {
			t.Errorf("%q. ReplacePlaceholders() = %#v, want %#v", tt.name, s.Command, tt.wantCommand)
		}
	}
}

func TestSnippetParseCommand(t *testing.T) {
	tests := []struct {
		name             string
		command          string
		wantPlaceholders []Placeholder
	}{
		{"no placeholder", "hello world", []Placeholder(nil)},
		{"full placeholder", "hello <<name(one,two)#desc>>",
			[]Placeholder{
				Placeholder{
					Name:     "name",
					Desc:     "#desc",
					Options:  []string{"one", "two"},
					Patterns: []string{"<<name(one,two)#desc>>"},
				},
			},
		},
		{"multiple patterns", `hello <<name#desc>> <<name>>`,
			[]Placeholder{
				Placeholder{
					Name:     "name",
					Desc:     "#desc",
					Patterns: []string{"<<name#desc>>", "<<name>>"},
				},
			},
		},
		{"multiple placeholders", `hello <<name>> <<last>>`,
			[]Placeholder{
				Placeholder{
					Name:     "name",
					Patterns: []string{"<<name>>"},
				},
				Placeholder{
					Name:     "last",
					Patterns: []string{"<<last>>"},
				},
			},
		},
	}
	for _, tt := range tests {
		s := &Snippet{
			Command: tt.command,
		}
		s.ParseCommand()
		if !reflect.DeepEqual(s.Placeholders, tt.wantPlaceholders) {
			t.Errorf("%q. ParseCommand() = %#v, want %#v", tt.name, s.Placeholders, tt.wantPlaceholders)
		}
	}
}

func TestInitSnippets(t *testing.T) {
	tests := []struct {
		name         string
		snippetMap   map[string]Snippet
		file         string
		dir          string
		wantSnippets SnippetSlice
	}{
		{"t",
			map[string]Snippet{"echo": Snippet{Command: "hello world"}},
			"file", "",
			SnippetSlice{
				Snippet{
					Name:    "echo",
					Command: "hello world",
					File:    "file",
				},
			},
		},
		{"ext_command",
			map[string]Snippet{"ext_command": Snippet{}},
			"examples", testPath,
			SnippetSlice{
				Snippet{
					Name:    "ext_command",
					Command: "test command",
					File:    "examples",
				},
			},
		},
		{"invalid snippet",
			map[string]Snippet{"invalid_snippet": Snippet{}},
			"examples", testPath,
			SnippetSlice(nil),
		},
	}
	makeTestFiles(t)
	for _, tt := range tests {
		if gotSnippets := initSnippets(tt.snippetMap, tt.file, tt.dir); !reflect.DeepEqual(gotSnippets, tt.wantSnippets) {
			t.Errorf("%q. initSnippets() = %#v, want %#v", tt.name, gotSnippets, tt.wantSnippets)
		}
	}
	defer cleanTestFiles(t)
}

func TestFilterByTag(t *testing.T) {
	tests := []struct {
		name        string
		snippets    SnippetSlice
		tag         string
		wantMatched SnippetSlice
	}{
		{"1",
			SnippetSlice{
				Snippet{
					Name: "skipped",
				},
				Snippet{
					Name: "matched",
					Tags: []string{"tag"},
				},
			},
			"tag",
			SnippetSlice{
				Snippet{
					Name: "matched",
					Tags: []string{"tag"},
				},
			},
		},
	}
	for _, tt := range tests {
		if gotMatched := filterByTag(tt.snippets, tt.tag); !reflect.DeepEqual(gotMatched, tt.wantMatched) {
			t.Errorf("%q. filterByTag() = %v, want %v", tt.name, gotMatched, tt.wantMatched)
		}
	}
}
