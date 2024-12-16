package hw03frequencyanalysis

import (
	"regexp"
	"slices"
	"sort"
	"strings"
)

var (
	reSpace = regexp.MustCompile(`\s+`)
	reClean = regexp.MustCompile(`^[^\p{L}]+|[^\p{L}]+$`)
)

type word struct {
	str string
	hit int
}

func (w *word) incHit() {
	w.hit++
}

func findWordByStr(str string, words []*word) *word {
	for _, val := range words {
		if val.str == str {
			return val
		}
	}
	return nil
}

func cleanStr(str string) string {
	if len(str) == 1 && str == "-" {
		return ""
	}
	str = strings.ToLower(str)
	return reClean.ReplaceAllString(str, "")
}

func Top10(str string) []string {
	if len(str) == 0 {
		return []string{}
	}

	var words []*word
	var result []string
	limit := 10
	strTrimmed := reSpace.ReplaceAllString(str, " ")
	strSplited := strings.Split(strTrimmed, " ")

	for _, val := range strSplited {
		val = cleanStr(val)

		if val == "" {
			continue
		}

		existWord := findWordByStr(val, words)

		if existWord != nil {
			existWord.incHit()
		} else {
			words = append(words, &word{str: val, hit: 1})
		}
	}

	slices.SortFunc(words, func(a, b *word) int {
		if a.hit > b.hit {
			return -1
		}

		if a.hit == b.hit {
			return 0
		}

		return 1
	})

	sort.Slice(words, func(i, j int) bool {
		if words[i].hit == words[j].hit {
			return words[i].str < words[j].str
		}
		return words[i].hit > words[j].hit
	})

	if len(words) < limit {
		limit = len(words)
	}

	for i := 0; i < limit; i++ {
		result = append(result, words[i].str)
	}

	return result
}
