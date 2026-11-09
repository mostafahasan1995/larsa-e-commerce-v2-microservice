package enums

type LanguageKey string

const (
	EN LanguageKey = "en"
	AR LanguageKey = "ar"
	TR LanguageKey = "tr"
)

var LanguageList = []LanguageKey{EN, AR, TR}

func (l LanguageKey) GetLanguageKey() string {
	return string(l)
}
