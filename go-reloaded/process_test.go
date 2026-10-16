package main

import (
	"reflect"
	"testing"
)

func TestProcessText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"hex", "1E (hex) files were added", "30 files were added"},
		{"bin", "It has been 10 (bin) years", "It has been 2 years"},
		{"up", "Ready, set, go (up) !", "Ready, set, GO!"},
		{"low", "I should stop SHOUTING (low)", "I should stop shouting"},
		{"cap", "Welcome to the Brooklyn bridge (cap)", "Welcome to the Brooklyn Bridge"},
		{"up avec nombre", "This is so exciting (up, 2)", "This is SO EXCITING"},
		{"ponctuation", "I was sitting over there ,and then BAMM !!", "I was sitting over there, and then BAMM!!"},
		{"groupes", "I was thinking ... You were right", "I was thinking... You were right"},
		{"apostrophes", "I am exactly how they describe me: ' awesome '", "I am exactly how they describe me: 'awesome'"},
		{"apostrophes multi-mots",
			"As Elton John said: ' I am the most well-known homosexual in the world '",
			"As Elton John said: 'I am the most well-known homosexual in the world'"},
		{"a -> an", "There it was. A amazing rock!", "There it was. An amazing rock!"},
		{"mix hex+bin", "Simply add 42 (hex) and 10 (bin) and you will see the result is 68.",
			"Simply add 66 and 2 and you will see the result is 68."},
		{"ponctuation finale", "Punctuation tests are ... kinda boring ,what do you think ?",
			"Punctuation tests are... kinda boring, what do you think?"},
		{"balise puis virgule", "it was the worst of times (up) , it was", "it was the worst of TIMES, it was"},
		{"nombre trop grand", "so exciting (cap, 6)", "So Exciting"},
		{"balise sans mot avant", "(up) hello", "hello"},
		{"hex invalide", "zz (hex) test", "zz test"},
		{"cap sur mot en majuscules", "BRIDGE (cap)", "Bridge"},
		{"a majuscule devant h", "A hour", "An hour"},
		{"a devant consonne", "a dog", "a dog"},
		{"apostrophe dans un mot", "don't stop", "don't stop"},
		{"retours a la ligne conserves", "hello ,world\nbye (up)\n", "hello, world\nBYE\n"},
		{"espaces multiples", "a   b    (up)", "a B"},
		{"decimaux intacts", "pi is 3.14 !", "pi is 3.14!"},
		{"fichier vide", "", ""},
		{"citation avec ponctuation", "harold wilson (cap, 2) : ' I am a optimist ,but a optimist who carries a raincoat . '",
			"Harold Wilson: 'I am an optimist, but an optimist who carries a raincoat.'"},
		{"ponctuation seule avant balise", "Ready , set (up, 2)", "READY, SET"},
		{"ponctuation apres apostrophe fermante", "he said ' hello ' .", "he said 'hello'."},
		{"virgule apres apostrophe fermante", "he said ' hello ' , and left", "he said 'hello', and left"},
		{"deux citations", "' one ' , ' two '", "'one', 'two'"},
		{"a devant citation", "a ' honest ' man", "an 'honest' man"},
		{"a dans une citation", "' a amazing rock '", "'an amazing rock'"},
		{"deux apostrophes sans mot", "I am ' ' ok", "I am '' ok"},
		{"up accents", "l'été (up)", "L'ÉTÉ"},
		{"cap accents", "ÉCOLE (cap)", "École"},
		{"low accents", "ÇA VA (low, 2)", "ça va"},
		{"up japonais inchange", "日本語 (up) ok", "日本語 ok"},
		{"ponctuation apres mot accentue", "il est là ,vraiment ?", "il est là, vraiment?"},
		{"nombre N enorme", "so exciting (up, 99999999999999999999)", "SO EXCITING"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processText(tt.input); got != tt.want {
				t.Errorf("\n entrée : %q\n obtenu : %q\n attendu: %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSplitWords(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"mots simples", "a  b\tc", []string{"a", "b", "c"}},
		{"balise avec nombre gardee entiere", "so (up, 2) now", []string{"so", "(up, 2)", "now"}},
		{"balise collee a la ponctuation", "go (up),then", []string{"go", "(up)", ",then"}},
		{"balise inconnue = simple mot", "(foo) (up, x)", []string{"(foo)", "(up,", "x)"}},
		{"balise collee a un mot", "go(up)", []string{"go(up)"}},
		{"multi-octets", "été 日本語", []string{"été", "日本語"}},
		{"ligne vide", "   ", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			for _, word := range splitWords([]rune(tt.input)) {
				got = append(got, string(word))
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\n entrée : %q\n obtenu : %q\n attendu: %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestReadTag(t *testing.T) {
	tests := []struct {
		input     string
		wantName  string
		wantCount int
		wantSize  int
	}{
		{"(up)", "up", 1, 4},
		{"(hex)", "hex", 1, 5},
		{"(cap, 6) suite", "cap", 6, 8},
		{"(low,3)", "low", 3, 7},
		{"(up, )", "", 0, 0},
		{"(up, x)", "", 0, 0},
		{"(up", "", 0, 0},
		{"(foo)", "", 0, 0},
		{"up", "", 0, 0},
		{"", "", 0, 0},
	}
	for _, tt := range tests {
		name, count, size := readTag([]rune(tt.input))
		if name != tt.wantName || count != tt.wantCount || size != tt.wantSize {
			t.Errorf("readTag(%q) = (%q, %d, %d), attendu (%q, %d, %d)",
				tt.input, name, count, size, tt.wantName, tt.wantCount, tt.wantSize)
		}
	}
}
