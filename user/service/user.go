package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/airlangga-hub/food-delivery-app/user/helper"
	"github.com/airlangga-hub/food-delivery-app/user/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserPaymentGatewayRepository interface {
	CreatePaymentSession(ctx context.Context, paymentType model.PaymentType, userID uuid.UUID, userEmail string, amount int, items []model.PaymentGatewayItem) (model.PaymentGatewayResponse, error)
}

type UserMongoRepository interface {
	CreatePaymentRecord(ctx context.Context, paymentRecord model.PaymentRecord) error
	GetPendingEmailRecords(ctx context.Context) ([]model.PaymentRecord, error)
	SendEmail(ctx context.Context, rec model.PaymentRecord) error
	MarkEmailAsSending(ctx context.Context, id string) error
	IfErrorMarkEmailAsPending(ctx context.Context, id string) error
	MarkEmailAsSent(ctx context.Context, id string) error
}

type UserSQLRepository interface {
	UpdateLedger(ctx context.Context, userID uuid.UUID, reason model.LedgerReason, amount int) error
	RegisterCustomer(ctx context.Context, user model.UserRegister) (model.UserInfo, error)
	Login(ctx context.Context, email string) (model.UserLogin, error)
	GetUserInfo(ctx context.Context, email string, role model.RoleUser) (model.UserInfo, error)
}

type userService struct {
	userPaymentGatewayRepository UserPaymentGatewayRepository
	userMongoRepository          UserMongoRepository
	userSqlRepository            UserSQLRepository
	logger                       *slog.Logger
	jwtKey                       string
}

func NewUserService(userpaymentGatewayRepo UserPaymentGatewayRepository, userMongoRepo UserMongoRepository, userSqlRepo UserSQLRepository, logger *slog.Logger, jwtKey string) *userService {
	return &userService{
		userPaymentGatewayRepository: userpaymentGatewayRepo,
		userMongoRepository:          userMongoRepo,
		userSqlRepository:            userSqlRepo,
		logger:                       logger,
		jwtKey:                       jwtKey,
	}
}

func (s *userService) RegisterCustomer(ctx context.Context, input model.UserRegister) (model.UserInfo, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("user.service.Register (GenerateFromPassword): %w", err)
	}

	input.Password = string(hashed)

	userInfo, err := s.userSqlRepository.RegisterCustomer(ctx, input)
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("user.service.Register (Register): %w", err)
	}

	return userInfo, nil
}

func (s *userService) Login(ctx context.Context, email string, password string) (string, error) {
	userLogin, err := s.userSqlRepository.Login(ctx, email)
	if err != nil {
		return "", fmt.Errorf("user.service.Login (Login): %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(userLogin.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("user.service.Login (CompareHashAndPassword): %w", err)
	}

	token, err := helper.MakeJWT(userLogin.UserID.String(), userLogin.Email, model.RoleUser(userLogin.Role), []byte(s.jwtKey))
	if err != nil {
		return "", fmt.Errorf("user.service.Login (MakeJWT): %w", err)
	}

	return token, nil
}

func (s *userService) GetUserInfo(ctx context.Context, email string, role model.RoleUser) (model.UserInfo, error) {
	userInfo, err := s.userSqlRepository.GetUserInfo(ctx, email, role)
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("user.service.GetUserInfo: %w", err)
	}
	return userInfo, nil
}

func (s *userService) TopUpBalance(ctx context.Context, userID uuid.UUID, userEmail string, amount int) (model.PaymentLink, error) {
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
		return model.PaymentLink{}, fmt.Errorf("user.service.TopUpBalance: %w", err)
	}

	if err := s.userMongoRepository.CreatePaymentRecord(ctx, model.PaymentRecord{
		Email:                  userEmail,
		EmailStatus:            model.EmailStatusPending,
		PaymentType:            model.PaymentTypeTopUp,
		PaymentGatewayResponse: paymentGatewayResp,
	}); err != nil {
		return model.PaymentLink{}, fmt.Errorf("user.service.TopUpBalance: %w", err)
	}

	return model.PaymentLink{PaymentLink: paymentGatewayResp.PaymentLinkURL}, nil
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

func (s *userService) ProcessUnsentEmail(ctx context.Context) error {
	records, err := s.userMongoRepository.GetPendingEmailRecords(ctx)
	if err != nil {
		return fmt.Errorf("service.ProcessEmailQueue: %w", err)
	}

	for _, rec := range records {

		recordID := rec.ID.Hex()

		if err := s.userMongoRepository.MarkEmailAsSending(ctx, recordID); err != nil {
			s.logger.Error(fmt.Sprintf("Failed to update status to sending for record %s", recordID), slog.Any("error", err))
			continue
		}

		if err := s.userMongoRepository.SendEmail(ctx, rec); err != nil {
			s.logger.Error(fmt.Sprintf("Failed to send email to %s", rec.Email), slog.Any("error", err))

			if err := s.userMongoRepository.IfErrorMarkEmailAsPending(ctx, recordID); err != nil {
				s.logger.Error(fmt.Sprintf("Failed to update status to pending for record %s", recordID), slog.Any("error", err))
				continue
			}

			continue
		}

		if err := s.userMongoRepository.MarkEmailAsSent(ctx, recordID); err != nil {
			s.logger.Error(fmt.Sprintf("Failed to update status to sent for record %s", recordID), slog.Any("error", err))
		}
	}
	return nil
}

func (s *userService) CreatePaymentRecord(ctx context.Context, paymentRecord model.PaymentRecord) error {
	if err := s.userMongoRepository.CreatePaymentRecord(ctx, paymentRecord); err != nil {
		return fmt.Errorf("user.service.CreatePaymentRecord: %w", err)
	}
	return nil
}
