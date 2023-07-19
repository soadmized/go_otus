package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	splitted := strings.Fields(text)
	freq := make(map[string]int, 0)

	for _, s := range splitted {
		if s != "" {
			freq[s]++
		}
	}

	var top []string

	sorted := sortByWordCount(freq)
	if len(sorted) == 0 {
		return top
	}

	top = append(top, sorted...)

	return top[:10]
}

func sortByWordCount(wordFrequencies map[string]int) []string {
	words := collectWords(wordFrequencies)
	sort.Slice(words, func(i, j int) bool {
		curr := wordFrequencies[words[i]]
		next := wordFrequencies[words[j]]

		if curr == next {
			return words[i] < words[j]
		}

		return curr > next
	})

	return words
}

func collectWords(wordFreq map[string]int) []string {
	words := []string{}

	for w := range wordFreq {
		words = append(words, w)
	}

	return words
}
