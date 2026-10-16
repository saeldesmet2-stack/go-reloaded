package main

import (
	"strings"
	"unicode"
)

// fixArticles transforme "a" en "an" (et "A" en "An") devant une voyelle ou un h.
// Elle modifie le slice en place.
func fixArticles(words [][]rune) {
	for i := 0; i+1 < len(words); i++ {
		word := words[i]
		isArticle := len(word) == 1 && (word[0] == 'a' || word[0] == 'A')
		if isArticle && startsWithVowelOrH(words[i+1:]) {
			words[i] = append(word, 'n')
		}
	}
}

// startsWithVowelOrH regarde la première rune du mot qui suit.
// Les apostrophes sont ignorées : dans "a ' honest man '", le mot suivant est bien honest.
func startsWithVowelOrH(next [][]rune) bool {
	for _, word := range next {
		for _, r := range word {
			if r == '\'' {
				continue
			}
			return strings.ContainsRune("aeiouh", unicode.ToLower(r))
		}
	}
	return false
}
