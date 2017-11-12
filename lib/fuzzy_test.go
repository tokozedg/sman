package sman

import (
	"github.com/renstrom/fuzzysearch/fuzzy"
	"reflect"
	"testing"
)

func TestTopsFromRanks(t *testing.T) {
	tests := []struct {
		name        string
		ranks       fuzzy.Ranks
		wantMatched []string
	}{
		{"empty ranks",
			fuzzy.Ranks{},
			[]string(nil),
		},
		{"single matched",
			fuzzy.Ranks{{"fir", "first", 3}, {"fir", "second", 8}},
			[]string{"first"},
		},
		{"multiple matched",
			fuzzy.Ranks{{"fir", "first", 3}, {"fir", "second", 3}, {"fir", "invalid", 8}},
			[]string{"first", "second"},
		},
	}
	for _, tt := range tests {
		if gotMatched := topsFromRanks(tt.ranks); !reflect.DeepEqual(gotMatched, tt.wantMatched) {
			t.Errorf("%q. topsFromRanks() = %v, want %v", tt.name, gotMatched, tt.wantMatched)
		}
	}
}

func TestFSearchSnippet(t *testing.T) {
	tests := []struct {
		name        string
		snippets    SnippetSlice
		pattern     string
		wantMatched SnippetSlice
	}{
		{"no match",
			SnippetSlice{
				Snippet{Name: "first"},
				Snippet{Name: "nonfirst"},
			}, "bird",
			SnippetSlice(nil),
		},
		{"single matched",
			SnippetSlice{
				Snippet{Name: "first"},
				Snippet{Name: "nonfirst"},
			}, "first",
			SnippetSlice{Snippet{Name: "first"}},
		},
		{"multiple matched",
			SnippetSlice{
				Snippet{Name: "first"},
				Snippet{Name: "firbe"},
				Snippet{Name: "non:match"},
			}, "fir",
			SnippetSlice{
				Snippet{Name: "first"},
				Snippet{Name: "firbe"},
			},
		},
		{"multiple matched subtasks",
			SnippetSlice{
				Snippet{Name: "user:add"},
				Snippet{Name: "alias:add"},
				Snippet{Name: "non:match"},
			}, "add",
			SnippetSlice{
				Snippet{Name: "user:add"},
				Snippet{Name: "alias:add"},
			},
		},
		{"single matched fully qualified",
			SnippetSlice{
				Snippet{Name: "user:add"},
				Snippet{Name: "alias:add"},
				Snippet{Name: "non:match"},
			}, "alias:add",
			SnippetSlice{
				Snippet{Name: "alias:add"},
			},
		},
		{"single matched fuzzy fully qualified",
			SnippetSlice{
				Snippet{Name: "user:add"},
				Snippet{Name: "alias:add"},
				Snippet{Name: "non:match"},
			}, "als:dd",
			SnippetSlice{
				Snippet{Name: "alias:add"},
			},
		},
	}
	for _, tt := range tests {
		if gotMatched := fSearchSnippet(tt.snippets, tt.pattern); !reflect.DeepEqual(gotMatched, tt.wantMatched) {
			t.Errorf("%q. fSearchSnippet() = %#v, want %#v", tt.name, gotMatched, tt.wantMatched)
		}
	}
}
