package sman

import (
	"github.com/renstrom/fuzzysearch/fuzzy"
	"sort"
	"strings"
)

// topsFromRanks iterates through fuzzy.Ranks and returns results
// whith the best distance
func topsFromRanks(ranks fuzzy.Ranks) (matched []string) {
	if len(ranks) == 0 {
		return matched
	}
	sort.Sort(ranks)
	topDistance := ranks[0].Distance
	for _, r := range ranks {
		if r.Distance == topDistance {
			matched = append(matched, r.Target)
		} else {
			break
		}
	}
	return matched
}

//fSearchFileName fuzzy searches pattern within available files in a dir
func fSearchFileName(pattern string, dir string) (matched []string) {
	files := ymlFiles(dir)
	ranks := fuzzy.RankFind(pattern, files)
	return topsFromRanks(ranks)
}

// fSearchSnippet matches pattern to snippet name in SnippetSlice
// returns SnippetSlice of best matched snippets.
func fSearchSnippet(snippets SnippetSlice, pattern string) (matched SnippetSlice) {
	// special case handling if pattern == snippet name
	wholeNameMatchTest := func(s Snippet) bool { return s.Name == pattern }
	wholeNameMatch := snippets.FilterView(wholeNameMatchTest)
	if wholeNameMatch.Len() == 1 {
		return wholeNameMatch
	}

	topRank := -1
	for _, s := range snippets {
		for _, part := range nameCombinations(s.Name) {
			r := fuzzy.RankMatch(pattern, part)
			switch {
			case r == -1:
				continue
			case topRank == -1 || r < topRank:
				matched = SnippetSlice{s}
				topRank = r
			case r == topRank:
				matched = append(matched, s)
			}
		}
	}
	return matched
}

// construct name combinations using ':' as separator of name parts
// ex: 1:2:3 -> [1, 2, 3, 1:2, 2:3, 1:2:3]
func nameCombinations(pattern string) (combinations []string) {
	if len(pattern) > 0 {
		if strings.Contains(pattern, ":") {
			singleNames := strings.Split(pattern, ":")

			// single names
			combinations = append(combinations, singleNames...)

			// combinations with length 2 to n-1
			for combLength := 2; combLength < len(singleNames); combLength++ {
				for startAt := 0; startAt < len(singleNames)-combLength+1; startAt++ {
					var combination []string
					for idx := startAt; idx < startAt+combLength; idx++ {
						combination = append(combination, singleNames[idx])
					}
					combinations = append(combinations, strings.Join(combination, ":"))
				}
			}
		}

		// complete name
		combinations = append(combinations, pattern)
	}
	return
}
