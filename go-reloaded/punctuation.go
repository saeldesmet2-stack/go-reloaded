package main

import "strings"

const punctuation = ".,!?:;"

func isPunctuation(r rune) bool {
	return strings.ContainsRune(punctuation, r)
}

// fixPunctuation colle la ponctuation initiale d'un morceau au morceau précédent.
// Comme les morceaux sont ensuite joints par une espace, la ponctuation se retrouve
// collée à gauche et séparée du mot suivant. Les groupes (... !?) restent entiers.
func fixPunctuation(words [][]rune) [][]rune {
	result := make([][]rune, 0, len(words))
	for _, word := range words {
		n := 0
		for n < len(word) && isPunctuation(word[n]) {
			n++
		}
		lead, rest := word[:n], word[n:]

		if n > 0 && len(result) > 0 {
			last := len(result) - 1
			result[last] = append(result[last], lead...)
		} else {
			rest = word // rien à qui coller : on garde le morceau tel quel
		}
		if len(rest) > 0 {
			result = append(result, rest)
		}
	}
	return result
}

// parseQuote indique si le morceau est une apostrophe isolée. Elle peut être
// suivie de ponctuation : fixPunctuation a pu coller un "." ou une "," derrière
// elle (cas de ' hello ' .). tail contient cette ponctuation éventuelle.
func parseQuote(word []rune) (ok bool, tail []rune) {
	if len(word) == 0 || word[0] != '\'' {
		return false, nil
	}
	for _, r := range word[1:] {
		if !isPunctuation(r) {
			return false, nil
		}
	}
	return true, word[1:]
}

// fixQuotes gère les ' isolées, par paires : la 1re se colle au mot suivant,
// la 2e au mot précédent. Les apostrophes dans un mot (don't) ne sont pas touchées.
func fixQuotes(words [][]rune) [][]rune {
	result := make([][]rune, 0, len(words))
	var pending []rune // apostrophe ouvrante, en attente de son mot
	open := false

	for _, word := range words {
		isQuote, tail := parseQuote(word)
		switch {
		case isQuote && open && len(pending) == 0:
			// apostrophe fermante : collée au mot précédent, avec la ponctuation qui la suit
			last := len(result) - 1
			result[last] = append(result[last], word...)
			open = false
		case isQuote && !open && len(tail) == 0:
			pending = []rune{'\''}
			open = true
		default:
			result = append(result, append(pending, word...))
			pending = nil
			if isQuote {
				open = false
			}
		}
	}
	if len(pending) > 0 { // apostrophe ouvrante sans mot après
		result = append(result, pending)
	}
	return result
}
