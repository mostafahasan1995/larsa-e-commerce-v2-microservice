package types

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
)

type LocalizableArray[T any] []Localizable[T]

func (localizableArray LocalizableArray[T]) MarshalJSON(ctx context.Context) ([]byte, error) {
	var lang = ExtractLanguage(ctx)

	localizedArray := make([]map[string]any, len(localizableArray))
	for i, localizable := range localizableArray {
		localizedArray[i] = Localize(localizable, enums.LanguageKey(lang))
	}

	return json.Marshal(localizedArray)
}
