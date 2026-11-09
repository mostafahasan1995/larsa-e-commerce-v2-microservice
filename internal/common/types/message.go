package types

type Message struct {
	Email       string
	Subject     string
	MessageHtml string
	Others      map[string]interface{}
}
