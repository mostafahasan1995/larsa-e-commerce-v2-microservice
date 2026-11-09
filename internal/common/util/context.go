package util

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"net/http"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"
)

type CtxKey int

const (
	_ CtxKey = iota
	ReqAppCfg
	ReqUser
)

// add AppCfg sruct to ctx for use in repos- usually in handlers
func AddCtxAppCfg(r *http.Request) (context.Context, *types.AppCfg) {
	hp := common.ExtractHeaderParams(r)
	db := hp.Client
	lang := hp.AcceptLanguage

	user, _ := r.Context().Value(ReqUser).(*types.User)

	cfg := &types.AppCfg{
		Db:   db,
		Lang: lang,
		Hp:   hp,
		User: user,
	}

	return context.WithValue(r.Context(), ReqAppCfg, cfg), cfg
}

// retrive AppConfig form ctx - usually in repos
func GetReqAppCfg(ctx context.Context) (*types.AppCfg, error) {
	cfg, _ := ctx.Value(ReqAppCfg).(*types.AppCfg)

	return cfg, nil
}

// add user to ctx- usually in middleware
func SetReqUser(ctx context.Context, user *types.User) context.Context {
	return context.WithValue(ctx, ReqUser, user)
}

func GetInjector(ctx context.Context) *do.Injector {
	return ctx.Value(enums.CtxInjectorKey).(*do.Injector)
}
