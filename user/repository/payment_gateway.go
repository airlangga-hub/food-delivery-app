package repository

import (
	"context"
	"fmt"

	"github.com/airlangga-hub/food-delivery-app/user/model"
	orderpb "github.com/airlangga-hub/food-delivery-app/user/order_pb"
	"github.com/google/uuid"
)

type paymentGatewayRepository struct {
	orderClient orderpb.OrderServiceClient
}

func NewPaymentGatewayRepository(orderClient orderpb.OrderServiceClient) *paymentGatewayRepository {
	return &paymentGatewayRepository{orderClient: orderClient}
}

func (r *paymentGatewayRepository) CreatePaymentSession(ctx context.Context, paymentType model.PaymentType, userID uuid.UUID, userEmail string, amount int, items []model.PaymentGatewayItem) (model.PaymentGatewayResponse, error) {
	iitems := make([]*orderpb.PaymentGatewayItem, len(items))

	for i, itm := range items {
		iitems[i] = &orderpb.PaymentGatewayItem{
			ReferenceId:   itm.ReferenceID,
			Name:          itm.Name,
			Description:   itm.Description,
			Type:          itm.Type,
			Category:      itm.Category,
			NetUnitAmount: int64(itm.NetUnitAmount),
			Quantity:      int64(itm.Quantity),
			Url:           itm.URL,
		}
	}

	resp, err := r.orderClient.CreatePaymentSession(ctx, &orderpb.CreatePaymentSessionRequest{
		PaymentType: string(paymentType),
		UserId:      userID.String(),
		UserEmail:   userEmail,
		Amount:      int64(amount),
		Items:       iitems,
	})

	if err != nil {
		return model.PaymentGatewayResponse{}, fmt.Errorf("user.repository.CreatePaymentSession: %w", err)
	}

	pgItem := make([]model.PaymentGatewayItem, len(resp.Items))

	for i, itm := range resp.Items {
		pgItem[i] = model.PaymentGatewayItem{
			ReferenceID:   itm.ReferenceId,
			Name:          itm.Name,
			Description:   itm.Description,
			Type:          itm.Type,
			Category:      itm.Category,
			NetUnitAmount: int(itm.NetUnitAmount),
			Quantity:      int(itm.Quantity),
			URL:           itm.Url,
		}
	}

	return model.PaymentGatewayResponse{
		PaymentSessionID: resp.PaymentSessionId,
		Created:          resp.Created,
		Updated:          resp.Updated,
		Status:           resp.Status,
		ReferenceID:      resp.ReferenceId,
		Currency:         resp.Currency,
		Amount:           resp.Amount,
		Country:          resp.Country,
		ExpiresAt:        resp.ExpiresAt,
		SessionType:      resp.SessionType,
		Mode:             resp.Mode,
		Locale:           resp.Locale,
		BusinessID:       resp.BusinessId,
		CustomerID:       resp.CustomerId,
		CaptureMethod:    resp.CaptureMethod,
		Description:      resp.Description,
		Items:            pgItem,
		SuccessReturnURL: resp.SuccessReturnUrl,
		CancelReturnURL:  resp.CancelReturnUrl,
		PaymentLinkURL:   resp.PaymentLinkUrl,
	}, nil
}
