package types

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
)

type Localizable[T any] map[enums.LanguageKey]T

func (localizable Localizable[T]) MarshalJSON(ctx context.Context) ([]byte, error) {
	var lang = ExtractLanguage(ctx)
	response := Localize(localizable, lang)

	return json.Marshal(response)
}

func ExtractLanguage(ctx context.Context) enums.LanguageKey {
	lang, _ := ctx.Value(enums.CtxLangKey).(string)
	extractLanguage := enums.LanguageKey(lang)

	if ctx == nil || extractLanguage == "" {
		return enums.EN
	}

	return extractLanguage
}

func (localizable Localizable[T]) GetValue(ctx context.Context) T {
	lang := ExtractLanguage(ctx)
	return LocalizeValue(localizable, lang)
}

func LocalizeValue[T any](localizable Localizable[T], languageKey enums.LanguageKey) T {
	if value, exists := localizable[languageKey]; exists {
		return value
	}
	// Fallback to English if no translation found
	if value, exists := localizable["en"]; exists {
		return value
	}
	return *new(T)
}

func Localize[T any](localizable Localizable[T], languageKey enums.LanguageKey) map[string]any {
	localizedValue := LocalizeValue(localizable, languageKey)
	response := make(map[string]any)

	// Spread all translations into the response
	for key, value := range localizable {
		response[string(key)] = value
	}

	response["defaultTranslation"] = localizedValue

	return response
}
