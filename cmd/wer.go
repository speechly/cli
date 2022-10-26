package cmd

import (
	"math"
	"strings"

	"github.com/agnivade/levenshtein"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type EditDistance struct {
	dist int
	base int
}

func (e EditDistance) Add(b EditDistance) EditDistance {
	return EditDistance{
		dist: e.dist + b.dist,
		base: e.base + b.base,
	}
}

func (e EditDistance) AsER() float64 {
	if e.base == 0 {
		return math.NaN()
	}
	return float64(e.dist) / float64(e.base)
}

func wordDistance(expected string, actual string) (EditDistance, error) {
	w2r := make(map[string]rune)
	exp, base := wordsToString(expected, w2r)
	act, _ := wordsToString(actual, w2r)

	distance := levenshtein.ComputeDistance(exp, act)
	return EditDistance{dist: distance, base: base}, nil
}

func wordsToString(s string, w2r map[string]rune) (string, int) {
	c := cases.Upper(language.English)
	words := strings.Split(c.String(s), " ")
	wordString := ""
	for _, w := range words {
		r := runifyWord(w, w2r)
		wordString = wordString + string(r)
	}
	return wordString, len(words)
}

func runifyWord(word string, w2r map[string]rune) rune {
	val, ok := w2r[word]
	if !ok {
		val = rune(len(w2r))
		w2r[word] = val
	}
	return val
}
