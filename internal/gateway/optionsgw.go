package gateway

import (
	"context"
	"errors"
	"larsa-e-commerce-microservice-v2/internal/common/util"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"
)

type OptionsGw struct{}

func NewOptionsGw(i *do.Injector) (*OptionsGw, error) {
	return &OptionsGw{}, nil
}

func (o *OptionsGw) GetOption(ctx context.Context, option string) (any, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return "", err
	}
	return common.GetOptionValue(option, cfg.Hp)
}

func (o *OptionsGw) GetHRLink(ctx context.Context) (string, error) {

	HR_LINK_val, err := o.GetOption(ctx, "HR_LINK")

	if err != nil {
		return "", errors.New("service options not available")
	}
	hrLink, ok := HR_LINK_val.(string)
	if !ok {
		return "", errors.New("error option value")
	}

	return hrLink, nil
}
