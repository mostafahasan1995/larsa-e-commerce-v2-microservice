package gateway

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/util"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"
)

type NotificationGw struct {
}

func NewNotificationGw(i *do.Injector) (*NotificationGw, error) {
	return &NotificationGw{}, nil
}

type NotificationRequest struct {
	Content   string
	Topic     string
	TopicType string
	Data      NotificationPayload
}

type NotificationPayload struct {
	Id    string
	Model string
}

func (ngw *NotificationGw) SendNotification(ctx context.Context, request *NotificationRequest) error {
	return ngw.sendNotificationViaHttp(ctx, request)
}

func (ngw *NotificationGw) sendNotificationViaHttp(ctx context.Context, request *NotificationRequest) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	clientId := cfg.Hp.Client

	serviceToken, err := common.GetServiceToken("ecommerce", &common.HeaderParams{})
	if err != nil {
		return err
	}

	hp := common.HeaderParams{
		Client:       clientId,
		ServiceToken: serviceToken,
	}

	reqOp := &common.RequestParams{
		Service: "notifications",
		Path:    "notification/topic",
		Method:  "POST",
		Data: map[string]interface{}{
			"content": request.Content,
			"topic":   request.Topic,
			"type":    request.TopicType,
			"data":    request.Data,
		},
		Header: &hp,
	}

	_, err = common.CallService(reqOp)
	if err != nil {
		return err
	}

	return nil

}
