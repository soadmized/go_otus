package hw03frequencyanalysis

import (
	"regexp"
	"sort"
)

func Top10(text string) []string {
	space := regexp.MustCompile(`\s`)
	splitted := space.Split(text, -1)

	freq := make(map[string]int, 0)

	for _, s := range splitted {
		if s != "" {
			freq[s] = freq[s] + 1 //nolint:gocritic
		}
	}

	sorted := sortByWordCount(freq)
	if len(sorted) == 0 {
		return []string{}
	}

	keys := []string{}
	for _, s := range sorted {
		keys = append(keys, s.Key)
	}

	return keys[:10]
}

func sortByWordCount(wordFrequencies map[string]int) PairList {
	list := make(PairList, len(wordFrequencies))
	i := 0

	for k, v := range wordFrequencies {
		list[i] = Pair{k, v}
		i++
	}

	sort.Stable(sort.Reverse(list))

	return list
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int {
	return len(p)
}

func (p PairList) Less(i, j int) bool {
	if p[i].Value == p[j].Value {
		return p[i].Key > p[j].Key
	}

	return p[i].Value < p[j].Value
}

func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
