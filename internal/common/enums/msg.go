package enums

import "git.larsa.io/mahdawi/microservices-commons.git/common"

type MsgTyps int

const (
	None MsgTyps = iota
)

type authKey string

var CtxAuthKey authKey = "authKey"

var CtxLangKey common.CommonContextKey = "lang"

var CtxInjectorKey common.CommonContextKey = "injector"

var CtxQueryKeysKey common.CommonContextKey = "queryKeys"
