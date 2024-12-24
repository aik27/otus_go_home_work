package hw03frequencyanalysis

import (
	"regexp"
	"slices"
	"sort"
	"strings"
)

var (
	reSpace  = regexp.MustCompile(`\s+`)
	reClean  = regexp.MustCompile(`^[^\p{L}]+|[^\p{L}]+$`)
	reExWord = regexp.MustCompile(`^-{2,}$`)
)

type word struct {
	str string
	hit int
}

func cleanStr(str string) string {
	if len(str) == 1 && str == "-" {
		return ""
	}

	if reExWord.Match([]byte(str)) {
		return str
	}

	return reClean.ReplaceAllString(strings.ToLower(str), "")
}

func Top10(str string) []string {
	if len(str) == 0 {
		return []string{}
	}

	limit := 10
	strTrm := reSpace.ReplaceAllString(str, " ")
	strSpt := strings.Split(strTrm, " ")

	words := make(map[string]int, len(strSpt))
	sortedWords := make([]word, 0, len(strSpt))
	result := make([]string, 0, len(strSpt))

	for _, v := range strSpt {
		v = cleanStr(v)
		if v == "" {
			continue
		}

		_, ok := words[v]
		if ok {
			words[v]++
		} else {
			words[v] = 1
		}
	}

	for k, v := range words {
		sortedWords = append(sortedWords, word{str: k, hit: v})
	}

	slices.SortFunc(sortedWords, func(a, b word) int {
		if a.hit > b.hit {
			return -1
		}

		if a.hit == b.hit {
			return 0
		}

		return 1
	})

	sort.Slice(sortedWords, func(i, j int) bool {
		if sortedWords[i].hit == sortedWords[j].hit {
			return sortedWords[i].str < sortedWords[j].str
		}
		return sortedWords[i].hit > sortedWords[j].hit
	})

	if len(sortedWords) < limit {
		limit = len(sortedWords)
	}

	for i := 0; i < limit; i++ {
		result = append(result, sortedWords[i].str)
	}

	return result
}
