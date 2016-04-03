package sman

import (
	//"fmt"
	"github.com/renstrom/fuzzysearch/fuzzy"
	"sort"
)

func snippetNames(slice SnippetSlice) (names []string, snippetMap map[string]Snippet) {
	snippetMap = make(map[string]Snippet)
	for _, s := range slice {
		names = append(names, s.Name)
		snippetMap[s.Name] = s
	}
	return names, snippetMap
}

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

func FSearchFileName(name string, dir string) (matched []string) {
	files := YmlFiles(dir)
	ranks := fuzzy.RankFind(name, files)
	return topsFromRanks(ranks)
}

func FSearchSnippet(snippets SnippetSlice, name string) (matched SnippetSlice) {
	names, snippetMap := snippetNames(snippets)
	ranks := fuzzy.RankFind(name, names)
	for _, n := range topsFromRanks(ranks) {
		matched = append(matched, snippetMap[n])
	}
	return matched
}
