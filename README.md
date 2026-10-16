# Go-reloaded

Outil en ligne de commande, écrit en Go qui complète, édite et autocorrige un texte.

Il reçoit en arguments un fichier d'entrée et un fichier de sortie, applique une série de règles de modification, puis écrit le texte corrigé dans le fichier de sortie.

## Prérequis

- Go 1.21 ou plus (`go version`)

## Utilisation

```bash
go run . <fichier_entrée> <fichier_sortie>
```

Exemple :

```bash
$ cat examples/sample.txt
Simply add 42 (hex) and 10 (bin) and you will see the result is 68.

$ go run . examples/sample.txt result.txt

$ cat result.txt
Simply add 66 and 2 and you will see the result is 68.
```

Le dossier `examples/` contient un fichier d'entrée (`sample.txt`) et le résultat attendu (`expected.txt`) :

```bash
go run . examples/sample.txt result.txt
diff result.txt examples/expected.txt && echo "OK"
```

En cas d'erreur (mauvais nombre d'arguments, fichier d'entrée introuvable), le programme affiche un message sur `stderr` et se termine avec le code 1.

## Règles appliquées

| Règle | Effet | Exemple |
|---|---|---|
| `(hex)` | le mot précédent (hexadécimal) devient décimal | `1E (hex) files` → `30 files` |
| `(bin)` | le mot précédent (binaire) devient décimal | `10 (bin) years` → `2 years` |
| `(up)` | mot précédent en majuscules | `go (up)` → `GO` |
| `(low)` | mot précédent en minuscules | `SHOUTING (low)` → `shouting` |
| `(cap)` | mot précédent avec majuscule initiale | `bridge (cap)` → `Bridge` |
| `(up, N)` `(low, N)` `(cap, N)` | idem sur les N mots précédents | `so exciting (up, 2)` → `SO EXCITING` |
| `. , ! ? : ;` | collés au mot précédent, espace après | `there ,and` → `there, and` |
| `...` `!?` | groupes de ponctuation gardés entiers | `thinking ... You` → `thinking... You` |
| `'` … `'` | collées aux mots qu'elles entourent | `' awesome '` → `'awesome'` |
| `a` | devient `an` devant une voyelle ou un `h` | `A amazing rock` → `An amazing rock` |

Les balises disparaissent du texte final, sans laisser d'espace en trop.

## Architecture

```
go-reloaded/
├── main.go           point d'entrée : arguments, lecture et écriture des fichiers
├── process.go        découpage de la ligne en runes et pipeline de transformations
├── tags.go           lecture et application des balises (hex) (bin) (up) (low) (cap) et (x, N)
├── articles.go       règle a → an
├── punctuation.go    ponctuation et apostrophes
├── main_test.go      tests de run() : fichiers et arguments invalides
├── process_test.go   tests des règles (exemples du sujet, cas limites, accents) et du découpage en runes
├── examples/         exemple d'entrée et sortie attendue
├── go.mod
└── README.md
```

### Fonctionnement

Le texte est traité ligne par ligne, ce qui conserve les retours à la ligne. Chaque ligne passe dans ce pipeline :

```
ligne → []rune → mots → balises → a/an → ponctuation → apostrophes → jointure
```

Tout le traitement se fait sur des `[]rune` (et non sur des `string` indexées en octets) : les caractères accentués ou non latins (`été`, `日本語`) sont donc toujours traités comme un seul caractère. On ne repasse en `string` qu'aux frontières : lecture/écriture des fichiers, et `strconv` pour `(hex)` / `(bin)`.

1. **Runes → mots** : la ligne est convertie en `[]rune`, puis découpée à la main (sans regexp) en `[][]rune`. Une balise complète comme `(up, 2)` reste un seul mot.
2. **Balises** : chaque balise est appliquée aux mots déjà vus, puis supprimée.
3. **a / an** : correction en place des articles.
4. **Ponctuation** : la ponctuation en début de mot est collée au mot précédent.
5. **Apostrophes** : les `'` isolées sont appariées, la première collée au mot suivant, la seconde au mot précédent.
6. **Jointure** : les mots sont rejoints par une seule espace, puis reconvertis en `string`.

### Choix de conception

- L'ordre du pipeline compte : les balises passent avant la ponctuation (sinon `(up),` ne serait plus reconnu), et la ponctuation avant les apostrophes (pour obtenir `raincoat.'`).
- Dans `(up, N)`, seuls les vrais mots sont comptés : la ponctuation isolée est ignorée.
- Un nombre N supérieur au nombre de mots disponibles ne provoque pas d'erreur : tous les mots précédents sont transformés (même si N est gigantesque).
- Chaque mot est une copie indépendante de la ligne d'origine : ajouter des runes à un mot (`a` → `an`, ponctuation collée) ne peut jamais écraser le mot voisin.
- Une apostrophe fermante peut être suivie de ponctuation (`' hello ' .` → `'hello'.`) et `a` regarde le mot situé derrière l'apostrophe ouvrante (`a ' honest ' man` → `an 'honest' man`).
- Un mot qui n'est pas un nombre valide devant `(hex)` ou `(bin)` est laissé tel quel.
- La ponctuation n'est traitée qu'en début de mot, pour ne pas casser les nombres décimaux comme `3.14`.

## Tests

```bash
go test ./...        # lance tous les tests
go test -v ./...     # détail de chaque cas
gofmt -l .           # liste les fichiers mal formatés (aucune sortie = OK)
go vet ./...         # analyse statique
```

## Limites connues

- Les lignes sont traitées indépendamment : une balise en début de ligne ne s'applique pas aux mots de la ligne précédente.
- Une balise collée à un mot (`go(up)`) n'est pas reconnue : le sujet les écrit toujours séparées par un espace.
- Comme l'impose le sujet, `a` devient `an` devant un `h` (`a house` devient `an house`).
- Le fichier d'entrée doit être en UTF-8 : un octet invalide est remplacé par `�` lors de la conversion en runes.
