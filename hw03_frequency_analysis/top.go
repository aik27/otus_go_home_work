package main

import (
	"regexp"
	"slices"
	"sort"
	"strings"
)

var reg = regexp.MustCompile(`\s+`)

type Word struct {
	str string
	hit int
}

func (w *Word) increase() {
	w.hit++
}

func Top10(str string) []string {
	if len(str) == 0 {
		return []string{}
	}

	var words []*Word

	str = reg.ReplaceAllString(str, " ")
	sptStr := strings.Split(str, " ")

	for _, val := range sptStr {
		word := findWordByStr(val, words)

		if word != nil {
			word.increase()
		} else {
			words = append(words, &Word{str: val, hit: 1})
		}
	}

	slices.SortFunc(words, func(a, b *Word) int {
		if a.hit > b.hit {
			return -1
		} else if a.hit == b.hit {
			return 0
		} else {
			return 1
		}
	})

	sort.Slice(words, func(i, j int) bool {
		if words[i].hit == words[j].hit {
			return words[i].str < words[j].str
		}
		return words[i].hit > words[j].hit
	})

	var res []string

	for i := 0; i < 10; i++ {
		res = append(res, words[i].str)
	}

	/*
		for _, person := range res {
			fmt.Printf("%s\n", person)
		}
		for _, person := range words {
			fmt.Printf("%+v\n", *person)
		}
	*/

	return res
}

func findWordByStr(str string, words []*Word) *Word {
	for _, word := range words {
		if word.str == str {
			return word
		}
	}
	return nil
}
