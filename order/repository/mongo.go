package repository

import (
	"context"
	"fmt"

	"github.com/airlangga-hub/food-delivery-app/order/model"
	userpb "github.com/airlangga-hub/food-delivery-app/order/user_pb"
)

type mongoRepository struct {
	userClient userpb.UserServiceClient
}

func NewMongoRepository(userClient userpb.UserServiceClient) *mongoRepository {
	return &mongoRepository{userClient: userClient}
}

func (r *mongoRepository) CreatePaymentRecord(ctx context.Context, paymentRecord model.PaymentRecord) error {
	items := make([]*userpb.PaymentGatewayItem, len(paymentRecord.PaymentGatewayResponse.Items))

	for i, itm := range paymentRecord.PaymentGatewayResponse.Items {
		items[i] = &userpb.PaymentGatewayItem{
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

	pgResponse := &userpb.PaymentGatewayResponse{
		PaymentSessionId: paymentRecord.PaymentGatewayResponse.PaymentSessionID,
		Created:          paymentRecord.PaymentGatewayResponse.Created,
		Updated:          paymentRecord.PaymentGatewayResponse.Updated,
		Status:           paymentRecord.PaymentGatewayResponse.Status,
		ReferenceId:      paymentRecord.PaymentGatewayResponse.ReferenceID,
		Currency:         paymentRecord.PaymentGatewayResponse.Currency,
		Amount:           paymentRecord.PaymentGatewayResponse.Amount,
		Country:          paymentRecord.PaymentGatewayResponse.Country,
		ExpiresAt:        paymentRecord.PaymentGatewayResponse.ExpiresAt,
		SessionType:      paymentRecord.PaymentGatewayResponse.SessionType,
		Mode:             paymentRecord.PaymentGatewayResponse.Mode,
		Locale:           paymentRecord.PaymentGatewayResponse.Locale,
		BusinessId:       paymentRecord.PaymentGatewayResponse.BusinessID,
		CustomerId:       paymentRecord.PaymentGatewayResponse.CustomerID,
		CaptureMethod:    paymentRecord.PaymentGatewayResponse.CaptureMethod,
		Description:      paymentRecord.PaymentGatewayResponse.Description,
		Items:            items,
		SuccessReturnUrl: paymentRecord.PaymentGatewayResponse.SuccessReturnURL,
		CancelReturnUrl:  paymentRecord.PaymentGatewayResponse.CancelReturnURL,
		PaymentLinkUrl:   paymentRecord.PaymentGatewayResponse.PaymentLinkURL,
	}

	_, err := r.userClient.CreatePaymentRecord(ctx, &userpb.CreatePaymentRecordRequest{
		Id:                     paymentRecord.ID,
		Email:                  paymentRecord.Email,
		EmailStatus:            string(paymentRecord.EmailStatus),
		PaymentType:            string(paymentRecord.PaymentType),
		PaymentGatewayResponse: pgResponse,
	})

	if err != nil {
		return fmt.Errorf("order.repository.CreatePaymentRecord: %w", err)
	}
	
	return nil
}
