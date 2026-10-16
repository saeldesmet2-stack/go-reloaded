package main

import (
	"slices"
	"strings"
	"unicode"
)

// processText applique toutes les règles ligne par ligne,
// ce qui permet de conserver les retours à la ligne du fichier d'origine.
func processText(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = processLine(line)
	}
	return strings.Join(lines, "\n")
}

// processLine est le "pipeline" : l'ordre des étapes compte.
// La ligne est convertie en runes une seule fois, au début. Tout le travail
// se fait ensuite sur des []rune, et on ne repasse en string qu'à la fin.
func processLine(line string) string {
	words := splitWords([]rune(line))
	words = applyTags(words)
	fixArticles(words)
	words = fixPunctuation(words)
	words = fixQuotes(words)
	return string(joinWords(words))
}

// splitWords découpe une ligne (en runes) en morceaux : soit une balise complète
// comme (up) ou (cap, 6), soit n'importe quelle suite de runes non blanches.
// Chaque morceau est une copie indépendante : on peut le modifier ou lui
// ajouter des runes sans toucher à la ligne ni aux autres morceaux.
func splitWords(line []rune) [][]rune {
	var words [][]rune
	for i := 0; i < len(line); {
		if unicode.IsSpace(line[i]) {
			i++
			continue
		}
		if _, _, size := readTag(line[i:]); size > 0 {
			words = append(words, slices.Clone(line[i:i+size]))
			i += size
			continue
		}
		start := i
		for i < len(line) && !unicode.IsSpace(line[i]) {
			i++
		}
		words = append(words, slices.Clone(line[start:i]))
	}
	return words
}

// joinWords rejoint les morceaux avec une seule espace entre chacun.
func joinWords(words [][]rune) []rune {
	var out []rune
	for i, word := range words {
		if i > 0 {
			out = append(out, ' ')
		}
		out = append(out, word...)
	}
	return out
}
