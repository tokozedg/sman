package sman

import (
	"strconv"
	"strings"
)

// Placeholder struct
type Placeholder struct {
	Name, Desc, Input string
	Options           []string
	Patterns          []string
}

// DisplayName returns formatted name
func (p *Placeholder) DisplayName() string {
	return cyan("[" + p.Name + "]")
}

// AddPattern adds new pattern to Patterns slice if not exists
func (p *Placeholder) AddPattern(pattern string) {
	if !sliceContains(p.Patterns, pattern) {
		p.Patterns = append(p.Patterns, pattern)
	}
}

// ParseOptions splits string by comma and sets Options value
func (p *Placeholder) ParseOptions(in string) {
	split := strings.Split(in, ",")
	// Loop through options, and check for escaped comma
	var toAppend string
	var op []string
	for _, o := range split {
		if o[len(o)-1:] == `\` {
			toAppend += o[:len(o)-1]
			toAppend += ","
			continue
		} else {
			toAppend += o
		}
		op = append(op, toAppend)
		toAppend = ""
	}
	p.Options = op
}

// SetInput sets input variable of Placeholder.
// If input is string and Placeholder has options, we try to set option value.
// The same way, if input is empty we try to set first option as input, assuming default value.
func (p *Placeholder) SetInput(input string) {
	if len(p.Options) != 0 {
		if len(input) == 0 {
			input = p.Options[0]
		} else if i, err := strconv.Atoi(input); err == nil {
			if (i > 0) && (i <= len(p.Options)) {
				i--
				input = p.Options[i]
			}
		}
	}
	p.Input = input
}

// searchPlaceholder searches placeholder by name returning position and bool.
func searchPlaceholder(in []Placeholder, n string) (i int, ok bool) {
	for i, p := range in {
		if p.Name == n {
			return i, true
		}
	}
	return i, false
}
