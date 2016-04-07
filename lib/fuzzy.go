package sman

import (
	"github.com/renstrom/fuzzysearch/fuzzy"
	"sort"
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

//FSearchFileName fuzzy searches pattern within available files in a dir
func FSearchFileName(pattern string, dir string) (matched []string) {
	files := YmlFiles(dir)
	ranks := fuzzy.RankFind(pattern, files)
	return topsFromRanks(ranks)
}

// FSearchSnippet fuzzy searches pattern in a SnippetSlice
// which returns SnippetSlice containing only top ranked snippets
func FSearchSnippet(snippets SnippetSlice, pattern string) (matched SnippetSlice) {
	names, snippetMap := SnippetNames(snippets)
	ranks := fuzzy.RankFind(pattern, names)
	for _, n := range topsFromRanks(ranks) {
		matched = append(matched, snippetMap[n])
	}
	return matched
}
