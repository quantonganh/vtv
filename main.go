package main

import (
	"bufio"
	"context"
	"embed"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/template"
	"time"
	"unicode/utf8"

	"github.com/gorilla/mux"
	"github.com/quantonganh/go-cache"
	"github.com/quantonganh/vtv/ui"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/runes"
	"golang.org/x/text/secure/precis"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

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
	var port string
	flag.StringVar(&port, "port", "8043", "which port to listen to")
	flag.Parse()

	zlog := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()

	// Here is your final handler
	r := mux.NewRouter()
	r.Use(hlog.NewHandler(zlog))
	r.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))
	r.Use(hlog.UserAgentHandler("user_agent"))
	r.Use(hlog.RefererHandler("referer"))
	r.Use(hlog.RequestIDHandler("req_id", "Request-Id"))

	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"mod": func(i, j, r int) bool {
			return i%j == r
		},
	}).ParseFS(content, "ui/html/*.html")
	if err != nil {
		zlog.Fatal().Err(err).Msg("error creating new template")
	}

	c := cache.New()

	r.PathPrefix("/static/").Handler(http.FileServer(http.FS(ui.StaticFS)))
	r.Handle("/", errorHandler(indexHandler(tmpl)))
	r.Handle("/search", errorHandler(searchHandler(tmpl, c)))

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		zlog.Fatal().Err(err).Msg("Startup failed")
	}
}

func normalize(s string) string {
	return removeAccentsFromVowel(strings.ToLower(strings.ReplaceAll(s, " ", "")))
}

func splitIntoConsonantsAndVowels(s string) ([]rune, []rune) {
	cl := make([]rune, 0)
	vl := make([]rune, 0)
	for _, char := range s {
		if isConsonant(char) {
			cl = appendIfMissing(cl, char)
		} else if isVowel(char) {
			vl = appendIfMissing(vl, char)
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

func appendIfMissing[T comparable](s []T, r T) []T {
	for _, l := range s {
		if l == r {
			return s
		}
	}
	return append(s, r)
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
	words = append(words, vcs...)

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

//go:embed Viet39K.txt
var fs embed.FS

func findInWordlist(dws []string) []string {
	f, err := fs.Open("Viet39K.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	m := make(map[string]struct{})
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		m[scanner.Text()] = struct{}{}
	}

	finalWords := make([]string, 0)
	g, _ := errgroup.WithContext(context.Background())
	mu := &sync.Mutex{}
	for _, word := range dws {
		word := word
		g.Go(func() error {
			if _, ok := m[word]; ok {
				mu.Lock()
				finalWords = append(finalWords, word)
				mu.Unlock()
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		panic(err)
	}

	return finalWords
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

func addGraveToVowel(v string) string {
	t := transform.Chain(
		norm.NFD,
		precis.UsernameCaseMapped.NewTransformer(),
		runes.Map(func(r rune) rune {
			switch r {
			case 'a':
				return 'à'
			case 'ă':
				return 'ằ'
			case 'â':
				return 'ầ'
			case 'e':
				return 'è'
			case 'ê':
				return 'ề'
			case 'i':
				return 'ì'
			case 'o':
				return 'ò'
			case 'ô':
				return 'ồ'
			case 'ơ':
				return 'ờ'
			case 'u':
				return 'ù'
			case 'ư':
				return 'ừ'
			case 'y':
				return 'ỳ'
			}
			return r
		}),
		norm.NFC,
	)
	result, _, _ := transform.String(t, v)
	return result
}

func removeGraveFromVowel(v string) string {
	t := transform.Chain(
		norm.NFD,
		precis.UsernameCaseMapped.NewTransformer(),
		runes.Map(func(r rune) rune {
			switch r {
			case 'à':
				return 'a'
			case 'ằ':
				return 'ă'
			case 'ầ':
				return 'â'
			case 'è':
				return 'e'
			case 'ề':
				return 'ê'
			case 'ì':
				return 'i'
			case 'ò':
				return 'o'
			case 'ồ':
				return 'ô'
			case 'ờ':
				return 'ơ'
			case 'ù':
				return 'u'
			case 'ừ':
				return 'ư'
			case 'ỳ':
				return 'y'
			}
			return r
		}),
		norm.NFC,
	)
	result, _, _ := transform.String(t, v)
	return result
}

func addAcuteToVowel(v string) string {
	t := transform.Chain(
		norm.NFD,
		precis.UsernameCaseMapped.NewTransformer(),
		runes.Map(func(r rune) rune {
			switch r {
			case 'a':
				return 'á'
			case 'ă':
				return 'ắ'
			case 'â':
				return 'ấ'
			case 'e':
				return 'é'
			case 'ê':
				return 'ế'
			case 'i':
				return 'í'
			case 'o':
				return 'ó'
			case 'ô':
				return 'ố'
			case 'ơ':
				return 'ớ'
			case 'u':
				return 'ú'
			case 'ư':
				return 'ứ'
			case 'y':
				return 'ý'
			}
			return r
		}),
		norm.NFC,
	)
	result, _, _ := transform.String(t, v)
	return result
}

func removeAcuteFromVowel(v string) string {
	t := transform.Chain(
		norm.NFD,
		precis.UsernameCaseMapped.NewTransformer(),
		runes.Map(func(r rune) rune {
			switch r {
			case 'á':
				return 'a'
			case 'ắ':
				return 'ă'
			case 'ấ':
				return 'â'
			case 'é':
				return 'e'
			case 'ế':
				return 'ê'
			case 'í':
				return 'i'
			case 'ó':
				return 'o'
			case 'ố':
				return 'ô'
			case 'ớ':
				return 'ơ'
			case 'ú':
				return 'u'
			case 'ứ':
				return 'ư'
			case 'ý':
				return 'y'
			}
			return r
		}),
		norm.NFC,
	)
	result, _, _ := transform.String(t, v)
	return result
}

func addHookToVowel(v string) string {
	t := transform.Chain(
		norm.NFD,
		precis.UsernameCaseMapped.NewTransformer(),
		runes.Map(func(r rune) rune {
			switch r {
			case 'a':
				return 'ả'
			case 'ă':
				return 'ẳ'
			case 'â':
				return 'ẩ'
			case 'e':
				return 'ẻ'
			case 'ê':
				return 'ể'
			case 'i':
				return 'ỉ'
			case 'o':
				return 'ỏ'
			case 'ô':
				return 'ổ'
			case 'ơ':
				return 'ở'
			case 'u':
				return 'ủ'
			case 'ư':
				return 'ử'
			case 'y':
				return 'ỷ'
			}
			return r
		}),
		norm.NFC,
	)
	result, _, _ := transform.String(t, v)
	return result
}

func removeHookFromVowel(v string) string {
	t := transform.Chain(
		norm.NFD,
		precis.UsernameCaseMapped.NewTransformer(),
		runes.Map(func(r rune) rune {
			switch r {
			case 'ả':
				return 'a'
			case 'ẳ':
				return 'ă'
			case 'ẩ':
				return 'â'
			case 'ẻ':
				return 'e'
			case 'ể':
				return 'ê'
			case 'ỉ':
				return 'i'
			case 'ỏ':
				return 'o'
			case 'ổ':
				return 'ô'
			case 'ở':
				return 'ơ'
			case 'ủ':
				return 'u'
			case 'ử':
				return 'ư'
			case 'ỷ':
				return 'y'
			}
			return r
		}),
		norm.NFC,
	)
	result, _, _ := transform.String(t, v)
	return result
}

func addTildeToVowel(v string) string {
	t := transform.Chain(
		norm.NFD,
		precis.UsernameCaseMapped.NewTransformer(),
		runes.Map(func(r rune) rune {
			switch r {
			case 'a':
				return 'ã'
			case 'ă':
				return 'ẵ'
			case 'â':
				return 'ẫ'
			case 'e':
				return 'ẽ'
			case 'ê':
				return 'ễ'
			case 'i':
				return 'ĩ'
			case 'o':
				return 'õ'
			case 'ô':
				return 'ỗ'
			case 'ơ':
				return 'ỡ'
			case 'u':
				return 'ũ'
			case 'ư':
				return 'ữ'
			case 'y':
				return 'ỹ'
			}
			return r
		}),
		norm.NFC,
	)
	result, _, _ := transform.String(t, v)
	return result
}

func removeTildeFromVowel(v string) string {
	t := transform.Chain(
		norm.NFD,
		precis.UsernameCaseMapped.NewTransformer(),
		runes.Map(func(r rune) rune {
			switch r {
			case 'ã':
				return 'a'
			case 'ẵ':
				return 'ă'
			case 'ẫ':
				return 'â'
			case 'ẽ':
				return 'e'
			case 'ễ':
				return 'ê'
			case 'ĩ':
				return 'i'
			case 'õ':
				return 'o'
			case 'ỗ':
				return 'ô'
			case 'ỡ':
				return 'ơ'
			case 'ũ':
				return 'u'
			case 'ữ':
				return 'ư'
			case 'ỹ':
				return 'y'
			}
			return r
		}),
		norm.NFC,
	)
	result, _, _ := transform.String(t, v)
	return result
}

func addDotToVowel(v string) string {
	t := transform.Chain(
		norm.NFD,
		precis.UsernameCaseMapped.NewTransformer(),
		runes.Map(func(r rune) rune {
			switch r {
			case 'a':
				return 'ạ'
			case 'ă':
				return 'ặ'
			case 'â':
				return 'ậ'
			case 'e':
				return 'ẹ'
			case 'ê':
				return 'ệ'
			case 'i':
				return 'ị'
			case 'o':
				return 'ọ'
			case 'ô':
				return 'ộ'
			case 'ơ':
				return 'ợ'
			case 'u':
				return 'ụ'
			case 'ư':
				return 'ự'
			case 'y':
				return 'ỵ'
			}
			return r
		}),
		norm.NFC,
	)
	result, _, _ := transform.String(t, v)
	return result
}

func removeDotFromVowel(v string) string {
	t := transform.Chain(
		norm.NFD,
		precis.UsernameCaseMapped.NewTransformer(),
		runes.Map(func(r rune) rune {
			switch r {
			case 'ạ':
				return 'a'
			case 'ặ':
				return 'ă'
			case 'ậ':
				return 'â'
			case 'ẹ':
				return 'e'
			case 'ệ':
				return 'ê'
			case 'ị':
				return 'i'
			case 'ọ':
				return 'o'
			case 'ộ':
				return 'ô'
			case 'ợ':
				return 'ơ'
			case 'ụ':
				return 'u'
			case 'ự':
				return 'ư'
			case 'ỵ':
				return 'y'
			}
			return r
		}),
		norm.NFC,
	)
	result, _, _ := transform.String(t, v)
	return result
}

func removeAccentsFromVowel(s string) string {
	return removeGraveFromVowel(removeAcuteFromVowel(removeHookFromVowel(removeTildeFromVowel(removeDotFromVowel(s)))))
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

const (
	minLength    = 6
	maxLength    = 15
	wordsPerLine = 5
)

//go:embed ui/html/*.html
var content embed.FS

type appHandler func(w http.ResponseWriter, r *http.Request) error

type PageData struct {
	Query     string
	Results   [][]string
	Message   string
	Classes   []string
	Remainder int
}

func indexHandler(tmpl *template.Template) appHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		return tmpl.ExecuteTemplate(w, "base", PageData{})
	}
}

func searchHandler(tmpl *template.Template, c *cache.Cache) appHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		query := r.FormValue("q")
		length := utf8.RuneCountInString(query)
		if length < minLength || length > maxLength {
			data := PageData{
				Query:   query,
				Message: fmt.Sprintf("Độ dài truy vấn tìm kiếm phải từ %d đến %d ký tự.", minLength, maxLength),
			}
			return tmpl.ExecuteTemplate(w, "base", data)
		}

		value, found := c.Get(query)
		if found {
			hlog.FromRequest(r).Info().Msgf("found search results for '%s' in the memory cache", query)
		}

		words, ok := value.([]string)
		if !ok {
			words = search(query)
			c.Set(query, words, 7*24*time.Hour)
		}

		total := len(words)
		var message string
		if total == 0 {
			message = "Không tìm thấy từ nào."
		} else {
			message = fmt.Sprintf("Kết quả: %d từ.", total)
		}

		results := make([][]string, 0)
		for i := 0; i < total; i += wordsPerLine {
			end := i + wordsPerLine
			if end > total {
				end = total
			}
			results = append(results, words[i:end])
		}

		// https://getbootstrap.com/docs/5.0/components/list-group/#contextual-classes
		classes := []string{
			"",
			" list-group-item-primary",
			" list-group-item-secondary",
			" list-group-item-success",
			" list-group-item-danger",
			" list-group-item-warning",
			" list-group-item-info",
			" list-group-item-light",
			" list-group-item-dark",
		}
		data := PageData{
			Query:   query,
			Results: results,
			Message: message,
			Classes: classes,
		}

		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			return fmt.Errorf("error applying template: %w", err)
		}

		return nil
	}
}

func search(query string) []string {
	input := normalize(query)
	cl, vl := splitIntoConsonantsAndVowels(input)

	ws := makeWords(cl, vl)
	cws := make([]string, 0)
	for _, w1 := range ws {
		for _, w2 := range ws {
			cws = append(cws, fmt.Sprintf("%s %s", w1, w2))
		}
	}

	return findInWordlist(cws)
}

func errorHandler(handler appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			// Handle the error and send an appropriate response
			fmt.Println("Error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
