package service

import (
	"context"
	"fmt"

	"github.com/airlangga-hub/food-delivery-app/user/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserPaymentGatewayRepository interface {
	CreatePaymentSession(ctx context.Context, paymentType model.PaymentType, userID uuid.UUID, userEmail string, amount int, items []model.PaymentGatewayItem) (model.PaymentGatewayResponse, error)
}

type UserMongoRepository interface {
	CreatePaymentRecord(ctx context.Context, paymentRecord model.PaymentRecord) error
}

type UserSQLRepository interface {
	UpdateLedger(ctx context.Context, userID uuid.UUID, reason model.LedgerReason, amount int) error
	Register(ctx context.Context, user model.UserRegister) (model.UserInfo, error)
	Login(ctx context.Context, email string) (model.UserInfo, string, error)
	GetUserInfo(ctx context.Context, email string) (model.UserInfo, error)
}

type userService struct {
	userPaymentGatewayRepository UserPaymentGatewayRepository
	userMongoRepository          UserMongoRepository
	userSqlRepository            UserSQLRepository
}

func NewUserService(userpaymentGatewayRepo UserPaymentGatewayRepository, userMongoRepo UserMongoRepository, userSqlRepo UserSQLRepository) *userService {
	return &userService{
		userPaymentGatewayRepository: userpaymentGatewayRepo,
		userMongoRepository:          userMongoRepo,
		userSqlRepository:            userSqlRepo,
	}
}

func (s *userService) Register(ctx context.Context, input model.UserRegister) (model.UserInfo, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("user.service.Register (GenerateFromPassword): %w", err)
	}

	input.Password = string(hashed)

	userInfo, err := s.userSqlRepository.Register(ctx, input)
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("user.service.Register (Register): %w", err)
	}

	return userInfo, nil
}

func (s *userService) Login(ctx context.Context, email string, password string) (string, error) {
	_, hashedPassword, err := s.userSqlRepository.Login(ctx, email)
	if err != nil {
		return "", err
	}

	// compare password input with hash from DB
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", err
	}

	token := ""
	return token, nil
}

func (s *userService) GetUserInfo(ctx context.Context, email string) (model.UserInfo, error) {
	return s.userSqlRepository.GetUserInfo(ctx, email)
}

func (s *userService) TopUpBalance(ctx context.Context, userID uuid.UUID, userEmail string, amount int) (string, error) {
	items := []model.PaymentGatewayItem{
		{
			ReferenceID:   uuid.NewString(),
			Name:          "Top Up Customer Balance",
			Type:          "DIGITAL_SERVICE",
			Category:      "Top Up",
			NetUnitAmount: amount,
			Quantity:      1,
		},
	}

	paymentGatewayResp, err := s.userPaymentGatewayRepository.CreatePaymentSession(ctx, model.PaymentTypeTopUp, userID, userEmail, amount, items)
	if err != nil {
		return "", fmt.Errorf("user.service.TopUpBalance: %w", err)
	}

	if err := s.userMongoRepository.CreatePaymentRecord(ctx, model.PaymentRecord{
		Email:                  userEmail,
		EmailStatus:            model.EmailStatusPending,
		PaymentType:            model.PaymentTypeTopUp,
		PaymentGatewayResponse: paymentGatewayResp,
	}); err != nil {
		return "", fmt.Errorf("user.service.TopUpBalance: %w", err)
	}

	return paymentGatewayResp.PaymentLinkURL, nil
}

func (s *userService) PaymentGatewayWebhook(ctx context.Context, userID uuid.UUID, paymentType model.PaymentType, amount int) error {
	var reason model.LedgerReason

	switch paymentType {
	case model.PaymentTypeOrder:
		reason = model.LedgerReasonCustomerOrder
	case model.PaymentTypeTopUp:
		reason = model.LedgerReasonTopUp
	default:
		return fmt.Errorf("order.service.PaymentGatewayWebhook (invalid payment type: %s): %w", string(paymentType), model.ErrNotFound)
	}

	if err := s.userSqlRepository.UpdateLedger(ctx, userID, reason, amount); err != nil {
		return fmt.Errorf("order.service.PaymentGatewayWebhook (UpdateLedger): %w", err)
	}

	return nil
}
