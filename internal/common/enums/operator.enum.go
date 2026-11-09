package enums

type Operator string

const (
	EQ     Operator = "eq"
	NE     Operator = "ne"
	GT     Operator = "gt"
	GTE    Operator = "gte"
	LT     Operator = "lt"
	LTE    Operator = "lte"
	IN     Operator = "in"
	NOT_IN Operator = "not_in"
)
