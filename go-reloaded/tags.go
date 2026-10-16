package main

import (
	"strconv"
	"unicode"
)

// tagNames liste les noms de balises reconnus : (hex) (bin) (up) (low) (cap).
var tagNames = []string{"hex", "bin", "up", "low", "cap"}

// maxCount borne le nombre de mots visés par (up, N) : un N démesuré revient de toute façon à "tous les mots".
const maxCount = 1 << 30

// readTag lit une balise au début de r : (up), (hex)... ou (up, 2). Elle renvoie le nom, le nombre de mots visés (1 par défaut) et le nombre de runes lues. size vaut 0 si r ne commence pas par une balise.
func readTag(r []rune) (name string, count, size int) {
	if len(r) == 0 || r[0] != '(' {
		return "", 0, 0
	}
	i := 1
	for _, candidate := range tagNames {
		if hasPrefix(r[i:], candidate) {
			name = candidate
			break
		}
	}
	if name == "" {
		return "", 0, 0
	}
	i += len(name) // les noms sont en ASCII : 1 octet = 1 rune

	count = 1
	if i < len(r) && r[i] == ',' {
		i++
		for i < len(r) && unicode.IsSpace(r[i]) {
			i++
		}
		digits := i
		count = 0
		for i < len(r) && r[i] >= '0' && r[i] <= '9' {
			count = min(count*10+int(r[i]-'0'), maxCount)
			i++
		}
		if i == digits {
			return "", 0, 0 // virgule sans nombre : pas une balise
		}
	}
	if i >= len(r) || r[i] != ')' {
		return "", 0, 0
	}
	return name, count, i + 1
}

// hasPrefix indique si r commence par les caractères (ASCII) de s.
func hasPrefix(r []rune, s string) bool {
	if len(r) < len(s) {
		return false
	}
	for i, c := range s {
		if r[i] != c {
			return false
		}
	}
	return true
}

// applyTags parcourt les morceaux : un mot normal est recopié,
// une balise est appliquée aux mots déjà recopiés puis disparaît.
func applyTags(words [][]rune) [][]rune {
	result := make([][]rune, 0, len(words))
	for _, word := range words {
		name, count, size := readTag(word)
		if size == 0 || size != len(word) {
			result = append(result, word)
			continue
		}
		applyToLast(result, name, count)
	}
	return result
}

// applyToLast transforme les `count` derniers mots (la ponctuation seule ne compte pas).
func applyToLast(words [][]rune, name string, count int) {
	for i := len(words) - 1; i >= 0 && count > 0; i-- {
		if !isWord(words[i]) {
			continue
		}
		words[i] = transform(name, words[i])
		count--
	}
}

// isWord indique si le morceau contient au moins une lettre ou un chiffre.
func isWord(word []rune) bool {
	for _, r := range word {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func transform(name string, word []rune) []rune {
	switch name {
	case "hex":
		return convertBase(word, 16)
	case "bin":
		return convertBase(word, 2)
	case "up":
		return mapRunes(word, unicode.ToUpper)
	case "low":
		return mapRunes(word, unicode.ToLower)
	case "cap":
		return capitalize(word)
	}
	return word
}

// convertBase convertit un nombre écrit en base 16 ou 2 vers la base 10.
// strconv ne travaille que sur des string : la conversion se fait ici, à la frontière.
func convertBase(word []rune, base int) []rune {
	n, err := strconv.ParseInt(string(word), base, 64)
	if err != nil {
		return word // pas un nombre valide : on ne touche à rien
	}
	return []rune(strconv.FormatInt(n, 10))
}

// mapRunes renvoie une copie de word dont chaque rune est passée dans f.
func mapRunes(word []rune, f func(rune) rune) []rune {
	out := make([]rune, len(word))
	for i, r := range word {
		out[i] = f(r)
	}
	return out
}

// capitalize met la première rune en majuscule et toutes les autres en minuscules.
func capitalize(word []rune) []rune {
	out := mapRunes(word, unicode.ToLower)
	if len(out) > 0 {
		out[0] = unicode.ToUpper(out[0])
	}
	return out
}
