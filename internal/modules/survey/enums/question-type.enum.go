package survey_enums

type QuestionType string

const (
	RateType        QuestionType = "rate"
	FreeTextType    QuestionType = "freeText"
	MultiChoiceType QuestionType = "multiChoice"
)

func (qt QuestionType) String() string {
	return string(qt)
}
