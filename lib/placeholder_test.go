package sman

import (
	"reflect"
	"testing"
)

func TestParseOptions(t *testing.T) {
	tests := []struct {
		name        string
		in          string
		wantOptions []string
	}{
		{"1", "word", []string{"word"}},
		{"2", "word,bored", []string{"word", "bored"}},
		{"3", `word,\,bored`, []string{"word", ",bored"}},
		{"4", `word,\\,bored`, []string{"word", `\,bored`}},
	}
	var p Placeholder
	for _, tt := range tests {
		p.ParseOptions(tt.in)
		if !reflect.DeepEqual(p.Options, tt.wantOptions) {
			t.Errorf("%q. ParseOptions() = %v, want %v", tt.name, p.Options, tt.wantOptions)
		}
	}
}

func TestPlaceholderSetInput(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		placeholder Placeholder
		wantInput   string
	}{
		{"empty options", "test", Placeholder{}, "test"},
		{"string while has options", "test",
			Placeholder{
				Options: []string{"one", "two"},
			}, "test"},
		{"valid option", "1",
			Placeholder{
				Options: []string{"one", "two"},
			}, "one"},
		{"number", "3",
			Placeholder{
				Options: []string{"one", "two"},
			}, "3"},
		{"empty input", "",
			Placeholder{
				Options: []string{"one", "two"},
			}, "one"},
	}
	for _, tt := range tests {
		tt.placeholder.SetInput(tt.input)
		if tt.placeholder.Input != tt.wantInput {
			t.Errorf("%q. SetIntput() = %v, want %v", tt.name, tt.placeholder.Input, tt.wantInput)
		}
	}
}
