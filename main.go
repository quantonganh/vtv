package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"strings"
	"sync"
	"text/tabwriter"

	"golang.org/x/sync/errgroup"
)

const wordsPerLine = 10

var (
	consonants         = []rune{'b', 'c', 'd', 'đ', 'g', 'h', 'k', 'l', 'm', 'n', 'p', 'q', 'r', 's', 't', 'v', 'x'}
	doubleConsonants   = []string{"ch", "gh", "gi", "kh", "ng", "ngh", "nh", "ph", "qu", "th", "tr"}
	finalConsonants    = []string{"m", "p", "n", "t", "nh", "ch", "ng", "c"}
	acuteDotConsonants = []string{"p", "t", "ch", "c"}

	frontVowels                 = []string{"a", "ê", "i"}
	frontVowelConsonants        = []string{"nh", "ch"}
	vowels                      = []rune{'a', 'ă', 'â', 'e', 'ê', 'i', 'o', 'ô', 'ơ', 'u', 'ư', 'y'}
	doubleVowels                = []string{"ai", "ao", "au", "ay", "âu", "ây", "êu", "eo", "ia", "iê", "yê", "iu", "oa", "oe", "oă", "oi", "oo", "ôô", "ơi", "ua", "uâ", "uă", "uâ", "uê", "ua", "ui", "ưi", "uo", "ươ", "ưu", "uơ", "uy"}
	vowelWithoutFinalConsonants = []string{"uyu", "uya", "ươu", "ươi", "uôi", "uây", "uai", "ưu", "ưi", "ui", "ưa", "oeo", "oay", "oao", "oai", "ơi", "ôi", "oi", "iu", "iêu", "yêu", "ia", "êu", "eo", "ây", "ay", "âu", "au", "ao", "ai"}
)

type vowel struct {
	name                  string
	single                bool
	double                bool
	grave                 bool
	acute                 bool
	hook                  bool
	tilde                 bool
	dot                   bool
	withoutFinalConsonant bool
}

func main() {
	input := os.Args[1]
	cl, vl := splitIntoConsonantsAndVowels(input)

	ws := makeWords(cl, vl)
	// printWords(ws)
	cws := make([]string, 0)
	for _, w1 := range ws {
		for _, w2 := range ws {
			cws = append(cws, fmt.Sprintf("%s %s", w1, w2))
		}
	}
	findInWordlist(cws)
}

func splitIntoConsonantsAndVowels(s string) ([]rune, []rune) {
	cl := make([]rune, 0)
	vl := make([]rune, 0)
	for _, char := range s {
		if isConsonant(char) {
			cl = append(cl, char)
		} else if isVowel(char) {
			vl = append(vl, char)
		}
	}
	return cl, vl
}

func isConsonant(r rune) bool {
	for _, c := range consonants {
		if c == r {
			return true
		}
	}
	return false
}

func isDoubleConsonant(s string) bool {
	for _, dc := range doubleConsonants {
		if dc == s {
			return true
		}
	}
	return false
}

func isVowel(r rune) bool {
	for _, v := range vowels {
		if v == r {
			return true
		}
	}
	return false
}

func isDoubleVowel(s string) bool {
	for _, dv := range doubleVowels {
		if dv == s {
			return true
		}
	}
	return false
}

func isVowelWithoutFinalConsonant(s string) bool {
	for _, vwfc := range vowelWithoutFinalConsonants {
		if vwfc == s {
			return true
		}
	}
	return false
}

func makeWords(cl []rune, vl []rune) []string {
	words := make([]string, 0)

	cs := makeConsonants(cl)
	// fmt.Println(cs)

	vs := makeVowels(vl)
	// for _, v := range vs {
	// 	fmt.Printf("%+v\n", v)
	// }
	for _, v := range vs {
		words = append(words, v.name)
	}
	// fmt.Println(words)

	vcs := makeVowelConsonants(vs, cs)
	// fmt.Printf("vowel consonants: %+v\n", vcs)
	for _, vc := range vcs {
		words = append(words, vc)
	}

	for _, c := range cs {
		for _, v := range vs {
			switch v.name {
			case "â":
			case "uâ":
			case "iê":
			case "u":
				if c != "q" {
					words = append(words, fmt.Sprintf("%s%s", c, v.name))
				}
			case "âu":
				if c != "q" {
					words = append(words, fmt.Sprintf("%s%s", c, v.name))
				}
			case "iu":
				if c != "m" {
					words = append(words, fmt.Sprintf("%s%s", c, v.name))
				}
			case "ê":
				if c != "c" && c != "q" {
					words = append(words, fmt.Sprintf("%s%s", c, v.name))
				}
			case "êu":
				if c != "c" && c != "ch" && c != "q" && c != "nh" {
					words = append(words, fmt.Sprintf("%s%s", c, v.name))
				}
			default:
				words = append(words, fmt.Sprintf("%s%s", c, v.name))
			}
		}

		for _, vc := range vcs {
			switch vc {
			case "uân":
				if c != "c" && c != "n" {
					words = append(words, fmt.Sprintf("%s%s", c, vc))
				}
			case "ân":
				if c != "n" && c != "q" {
					words = append(words, fmt.Sprintf("%s%s", c, vc))
				}
			case "uc":
				if c != "q" {
					words = append(words, fmt.Sprintf("%s%s", c, vc))
				}
			case "êm":
				if c != "tr" && c != "r" {
					words = append(words, fmt.Sprintf("%s%s", c, vc))
				}
			case "iêm":
				if c != "m" && c != "r" && c != "tr" {
					words = append(words, fmt.Sprintf("%s%s", c, vc))
				}
			case "ên":
				if c != "c" && c != "q" && c != "ch" && c != "nh" {
					words = append(words, fmt.Sprintf("%s%s", c, vc))
				}
			default:
				words = append(words, fmt.Sprintf("%s%s", c, vc))
			}
		}
	}
	return words
}

//go:embed Viet74K.txt
var fs embed.FS

func findInWordlist(dws []string) {
	b, err := fs.ReadFile("Viet74K.txt")
	if err != nil {
		panic(err)
	}
	words := strings.Split(string(b), "\n")

	finalWords := make([]string, 0)
	g, _ := errgroup.WithContext(context.Background())
	mu := &sync.Mutex{}
	for _, source := range dws {
		source := source
		g.Go(func() error {
			if isFound(source, words) {
				mu.Lock()
				finalWords = append(finalWords, source)
				mu.Unlock()
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		panic(err)
	}

	printWords(finalWords)
}

func isFound(source string, words []string) bool {
	for _, w := range words {
		if w == source {
			return true
		}
	}
	return false
}

func makeConsonants(cl []rune) []string {
	cs := make([]string, 0)
	for _, r1 := range cl {
		cs = append(cs, fmt.Sprintf("%c", r1))
		for _, r2 := range cl {
			dc := fmt.Sprintf("%c%c", r1, r2)
			if isDoubleConsonant(dc) {
				cs = append(cs, dc)
			}
		}
	}
	return cs
}

func makeVowels(vl []rune) []*vowel {
	vs := make([]*vowel, 0)
	for _, r1 := range vl {
		r1InStr := toStr(r1)
		vs = append(vs, &vowel{
			name:   r1InStr,
			single: true,
		})
		vs = append(vs, &vowel{
			name:   addGraveToVowel(r1InStr),
			single: true,
			grave:  true,
		})
		vs = append(vs, &vowel{
			name:   addAcuteToVowel(r1InStr),
			single: true,
			acute:  true,
		})
		vs = append(vs, &vowel{
			name:   addHookToVowel(r1InStr),
			single: true,
			hook:   true,
		})
		vs = append(vs, &vowel{
			name:   addTildeToVowel(r1InStr),
			single: true,
			tilde:  true,
		})
		vs = append(vs, &vowel{
			name:   addDotToVowel(r1InStr),
			single: true,
			dot:    true,
		})
		for _, r2 := range vl {
			dv := fmt.Sprintf("%c%c", r1, r2)
			if isDoubleVowel(dv) {
				lv := &vowel{
					name:   dv,
					double: true,
				}
				if isVowelWithoutFinalConsonant(dv) {
					lv.withoutFinalConsonant = true
				}
				vs = append(vs, lv)

				if dv != "uâ" && dv != "uê" && dv != "oo" {
					gv := &vowel{
						name:   fmt.Sprintf("%s%c", addGraveToVowel(r1InStr), r2),
						double: true,
						grave:  true,
					}

					av := &vowel{
						name:   fmt.Sprintf("%s%c", addAcuteToVowel(r1InStr), r2),
						double: true,
						acute:  true,
					}

					hv := &vowel{
						name:   fmt.Sprintf("%s%c", addHookToVowel(r1InStr), r2),
						double: true,
						hook:   true,
					}

					tv := &vowel{
						name:   fmt.Sprintf("%s%c", addTildeToVowel(r1InStr), r2),
						double: true,
						tilde:  true,
					}

					dbv := &vowel{
						name:   fmt.Sprintf("%s%c", addDotToVowel(r1InStr), r2),
						double: true,
						dot:    true,
					}

					if isVowelWithoutFinalConsonant(dv) {
						gv.withoutFinalConsonant = true
						av.withoutFinalConsonant = true
						hv.withoutFinalConsonant = true
						tv.withoutFinalConsonant = true
						dbv.withoutFinalConsonant = true
					}
					vs = append(vs, gv, av, hv, tv, dbv)
				}
			}
		}
	}
	return vs
}

func toStr(r rune) string {
	return fmt.Sprintf("%c", r)
}

func makeVowelConsonants(vs []*vowel, cs []string) []string {
	vcs := make([]string, 0)
	for _, v := range vs {
		if isFrontVowel(v) {
			for _, c := range cs {
				if isFrontVowelConsonant(c) {
					vcs = append(vcs, fmt.Sprintf("%s%s", v.name, c))
				}
			}
		} else if isAcuteDotVowel(v) && !v.withoutFinalConsonant {
			for _, c := range cs {
				if isAcuteDotConsonant(c) {
					if isFrontVowelConsonant(c) {
						if isFrontVowel(v) {
							vcs = append(vcs, fmt.Sprintf("%s%s", v.name, c))
						}
					} else {
						vcs = append(vcs, fmt.Sprintf("%s%s", v.name, c))
					}
				}
			}
		} else if !v.withoutFinalConsonant {
			for _, c := range cs {
				if isNotFrontBackAcuteDotVowelConsonant(c) {
					vcs = append(vcs, fmt.Sprintf("%s%s", v.name, c))
				}
			}
		}
	}
	return vcs
}

func isFrontVowel(v *vowel) bool {
	for _, fv := range frontVowels {
		if fv == v.name {
			return true
		}
	}
	return false
}

func isAcuteDotVowel(v *vowel) bool {
	return v.acute || v.dot
}

func isAcuteDotConsonant(s string) bool {
	for _, adc := range acuteDotConsonants {
		if adc == s {
			return true
		}
	}
	return false
}

func isFrontVowelConsonant(s string) bool {
	for _, fvc := range frontVowelConsonants {
		if fvc == s {
			return true
		}
	}
	return false
}

func isNotFrontBackAcuteDotVowelConsonant(s string) bool {
	for _, fc := range finalConsonants {
		if fc == s && !isFrontVowelConsonant(s) && !isAcuteDotConsonant(s) {
			return true
		}
	}
	return false
}

func printWords(dws []string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	defer w.Flush()

	i := 0
	for i <= len(dws)-wordsPerLine {
		fmt.Fprintf(w, fmt.Sprintf("%s\n", format(wordsPerLine)), toAny(dws[i:i+wordsPerLine])...)
		i += wordsPerLine
	}

	fmt.Fprintf(w, fmt.Sprintf("%s\n", format(len(dws)-i)), toAny(dws[i:])...)
}

func format(count int) string {
	return strings.Repeat("%s\t", count)
}

func toAny(words []string) []any {
	a := make([]any, 0)
	for _, w := range words {
		a = append(a, w)
	}
	return a
}
