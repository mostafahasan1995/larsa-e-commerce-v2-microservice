package gateway

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"git.larsa.io/mahdawi/microservices-commons.git/nats"
	"github.com/samber/do"
)

type EmailGw struct {
}

func NewEmailGw(i *do.Injector) (*EmailGw, error) {
	return &EmailGw{}, nil
}

type Email struct {
	To      string                 `bson:"to" json:"to"`
	Subject string                 `bson:"subject" json:"subject"`
	Message string                 `bson:"message" json:"message"`
	Headers *common.HeaderParams   `bson:"headers" json:"headers"`
	Others  map[string]interface{} `bson:"others" json:"others"`
}

func (egw *EmailGw) SendEmail(ctx context.Context, msg *types.Message) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	email := &Email{
		To:      msg.Email,
		Subject: msg.Subject,
		Message: msg.MessageHtml,
		Headers: cfg.Hp,
		Others:  msg.Others,
	}

	nats.Publish("email:send", email)
	nats.PublishJet("email:send", email)
	return nil
}
