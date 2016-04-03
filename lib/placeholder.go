package sman

import (
	"fmt"
	"strings"
)

// Placeholder struct
type Placeholder struct {
	Name, Desc, Input string
	Options           []string
	Patterns          []string
}

func ParseOptions(in string) (options []string) {
	split := strings.Split(in, ",")
	// Loop Through Options, and check for escaped comma
	var toAppend string
	for _, o := range split {
		if o[len(o)-1:] == `\` {
			toAppend += o[:len(o)-1]
			toAppend += ","
			continue
		} else {
			toAppend += o
		}
		options = append(options, toAppend)
		toAppend = ""
	}
	return options
}

func SearchPlaceholder(in []Placeholder, n string) (i int, ok bool) {
	for i, p := range in {
		if p.Name == n {
			return i, true
		}
	}
	return i, false
}

func (p *Placeholder) DisplayName() string {
	return fmt.Sprintf("%s[ %s ]%s", `\033[95m`, p.Name, "-")
}

func (p *Placeholder) AddPattern(pattern string) {
	if !SliceContains(p.Patterns, pattern) {
		p.Patterns = append(p.Patterns, pattern)
	}
}
