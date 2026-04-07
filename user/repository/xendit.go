package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/airlangga-hub/food-delivery-app/user/model"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type xenditRepository struct {
	xenditPaymentSessionURL string
	xenditAPIkey            string
	validate                *validator.Validate
}

func NewXenditRepository(xenditPaymentSessionURL, xenditAPIkey string, validate *validator.Validate) *xenditRepository {

	return &xenditRepository{
		xenditPaymentSessionURL: xenditPaymentSessionURL,
		xenditAPIkey:            xenditAPIkey,
		validate:                validate,
	}
}

func (r *xenditRepository) CreatePaymentSession(ctx context.Context, userID uuid.UUID, userEmail string, amount int, items []model.PaymentGatewayItem) (model.PaymentGatewayResponse, error) {
	iitems := make([]Item, len(items))
	for i, item := range items {
		iitems[i] = Item{
			ReferenceID:   item.ReferenceID,
			Name:          item.Name,
			Description:   item.Description,
			Type:          item.Type,
			Category:      item.Category,
			NetUnitAmount: item.NetUnitAmount,
			Quantity:      item.Quantity,
			URL:           item.URL,
		}
	}

	payload := XenditPaymentSessionRequest{
		ReferenceID: userID.String(),
		SessionType: "PAY",
		Mode:        "PAYMENT_LINK",
		Amount:      amount,
		Currency:    "IDR",
		Country:     "ID",
		Customer: Customer{
			ReferenceID: uuid.NewString(),
			Type:        "INDIVIDUAL",
			IndividualDetail: IndividualDetail{
				GivenNames: userEmail,
				Surname:    "email",
			},
			Email: userEmail,
		},
		Items: iitems,
	}

	if err := r.validate.Struct(payload); err != nil {
		return model.PaymentGatewayResponse{}, fmt.Errorf("user.repository.CreatePaymentSession (validate.Struct): %w", err)
	}

	var buf *bytes.Buffer
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		return model.PaymentGatewayResponse{}, fmt.Errorf("user.repository.CreatePaymentSession (JSON encoding): %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.xenditPaymentSessionURL, buf)
	if err != nil {
		return model.PaymentGatewayResponse{}, fmt.Errorf("user.repository.CreatePaymentSession (NewRequest): %w", err)
	}

	req.SetBasicAuth(r.xenditAPIkey, "")
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.PaymentGatewayResponse{}, fmt.Errorf("user.repository.CreatePaymentSession (DefaultClient): %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(res.Body)
		return model.PaymentGatewayResponse{}, fmt.Errorf("user.repository.CreatePaymentSession (StatusCode): %s", string(body))
	}

	var xenditResp XenditPaymentSessionResponse
	if err := json.NewDecoder(res.Body).Decode(&xenditResp); err != nil {
		return model.PaymentGatewayResponse{}, fmt.Errorf("user.repository.CreatePaymentSession (JSON decoding): %w", err)
	}

	return model.PaymentGatewayResponse{
		PaymentSessionID: xenditResp.PaymentSessionID,
		Created:          xenditResp.Created,
		Updated:          xenditResp.Updated,
		Status:           xenditResp.Status,
		ReferenceID:      xenditResp.ReferenceID,
		Currency:         xenditResp.Currency,
		Amount:           xenditResp.Amount,
		Country:          xenditResp.Country,
		ExpiresAt:        xenditResp.ExpiresAt,
		SessionType:      xenditResp.SessionType,
		Mode:             xenditResp.Mode,
		Locale:           xenditResp.Locale,
		BusinessID:       xenditResp.BusinessID,
		CustomerID:       xenditResp.CustomerID,
		CaptureMethod:    xenditResp.CaptureMethod,
		Description:      xenditResp.Description,
		Items:            items,
		SuccessReturnURL: xenditResp.SuccessReturnURL,
		CancelReturnURL:  xenditResp.CancelReturnURL,
		PaymentLinkURL:   xenditResp.PaymentLinkURL,
	}, nil
}
