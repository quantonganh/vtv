package main

import (
	"golang.org/x/text/runes"
	"golang.org/x/text/secure/precis"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

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
