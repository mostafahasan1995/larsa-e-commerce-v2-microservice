package base_utils

type AdditionalOptions struct {
	ReturnError bool
}

type NumericFilter struct {
	Gte *int `json:"gte,omitempty"`
	Lte *int `json:"lte,omitempty"`
}
