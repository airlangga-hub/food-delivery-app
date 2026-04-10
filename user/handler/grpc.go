package handler

import (
	"context"
	"errors"

	"github.com/airlangga-hub/food-delivery-app/user/model"
	pb "github.com/airlangga-hub/food-delivery-app/user/pb"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService interface {
	RegisterCustomer(ctx context.Context, input model.UserRegister) (model.UserInfo, error)
	Login(ctx context.Context, email string, password string) (string, error)
	GetUserInfo(ctx context.Context, email string) (model.UserInfo, error)
	PaymentGatewayWebhook(ctx context.Context, userID uuid.UUID, paymentType model.PaymentType, amount int) error
	TopUpBalance(ctx context.Context, userID uuid.UUID, userEmail string, amount int) (model.PaymentLink, error)
	CreatePaymentRecord(ctx context.Context, paymentRecord model.PaymentRecord) error
}

type Handler struct {
	pb.UnimplementedUserServiceServer
	Svc UserService
}

func NewHandler(svc UserService) *Handler {
	return &Handler{Svc: svc}
}

func (h *Handler) RegisterCustomer(ctx context.Context, req *pb.RegisterCustomerRequest) (*pb.RegisterCustomerResponse, error) {
	userinfo, err := h.Svc.RegisterCustomer(ctx, model.UserRegister{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    req.Password,
		Address:     req.Address,
		PhoneNumber: req.PhoneNumber,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "user.handler.RegisterCustomer: %v", err)
	}

	return &pb.RegisterCustomerResponse{
		FirstName:   userinfo.FirstName,
		LastName:    userinfo.LastName,
		Email:       userinfo.Email,
		Address:     userinfo.Address,
		PhoneNumber: userinfo.PhoneNumber,
		Balance:     int64(userinfo.Balance),
	}, nil
}

func (h *Handler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := h.Svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "user.handler.Login (user not found): %v", err)
		}
		return nil, status.Errorf(codes.Internal, "user.handler.Login: %v", err)
	}
	return &pb.LoginResponse{Token: token}, nil
}

func (h *Handler) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	userinfo, err := h.Svc.GetUserInfo(ctx, req.Email)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "user.handler.GetUserInfo (user not found): %v", err)
		}
		return nil, status.Errorf(codes.Internal, "user.handler.GetUserInfo: %v", err)
	}
	return &pb.GetUserInfoResponse{
		FirstName:   userinfo.FirstName,
		LastName:    userinfo.LastName,
		Email:       userinfo.Email,
		Address:     userinfo.Address,
		PhoneNumber: userinfo.PhoneNumber,
		Balance:     int64(userinfo.Balance),
	}, nil
}

func (h *Handler) PaymentGatewayWebhook(ctx context.Context, req *pb.PaymentGatewayWebhookRequest) (*pb.PaymentGatewayWebhookResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "user.handler.PaymentGatewayWebhook (Parse): %v", err)
	}

	if err := h.Svc.PaymentGatewayWebhook(ctx, userID, model.PaymentType(req.PaymentType), int(req.Amount)); err != nil {
		return nil, status.Errorf(codes.Internal, "user.handler.PaymentGatewayWebhook: %v", err)
	}

	return nil, nil
}

func (h *Handler) TopUpBalance(ctx context.Context, req *pb.TopUpBalanceRequest) (*pb.TopUpBalanceResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "user.handler.TopUpBalance (Parse): %v", err)
	}

	paymentLink, err := h.Svc.TopUpBalance(ctx, userID, req.UserEmail, int(req.Amount))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "user.handler.TopUpBalance: %v", err)
	}

	return &pb.TopUpBalanceResponse{
		PaymentLink: paymentLink.PaymentLink,
	}, nil
}

func (h *Handler) CreatePaymentRecord(ctx context.Context, req *pb.CreateCreatePaymentRecordRequest) (*pb.CreateCreatePaymentRecordResponse, error) {
	items := make([]model.PaymentGatewayItem, len(req.PaymentGatewayResponse.Items))

	for i, itm := range req.PaymentGatewayResponse.Items {
		items[i] = model.PaymentGatewayItem{
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

	pgResponse := model.PaymentGatewayResponse{
		PaymentSessionID: req.PaymentGatewayResponse.PaymentSessionId,
		Created:          req.PaymentGatewayResponse.Created,
		Updated:          req.PaymentGatewayResponse.Updated,
		Status:           req.PaymentGatewayResponse.Status,
		ReferenceID:      req.PaymentGatewayResponse.ReferenceId,
		Currency:         req.PaymentGatewayResponse.Currency,
		Amount:           req.PaymentGatewayResponse.Amount,
		Country:          req.PaymentGatewayResponse.Country,
		ExpiresAt:        req.PaymentGatewayResponse.ExpiresAt,
		SessionType:      req.PaymentGatewayResponse.SessionType,
		Mode:             req.PaymentGatewayResponse.Mode,
		Locale:           req.PaymentGatewayResponse.Locale,
		BusinessID:       req.PaymentGatewayResponse.BusinessId,
		CustomerID:       req.PaymentGatewayResponse.CustomerId,
		CaptureMethod:    req.PaymentGatewayResponse.CaptureMethod,
		Description:      req.PaymentGatewayResponse.Description,
		Items:            items,
		SuccessReturnURL: req.PaymentGatewayResponse.SuccessReturnUrl,
		CancelReturnURL:  req.PaymentGatewayResponse.CancelReturnUrl,
		PaymentLinkURL:   req.PaymentGatewayResponse.PaymentLinkUrl,
	}

	objID, err := bson.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "user.handler.CreatePaymentRecord (ObjectIDFromHex): %v", err)
	}

	if err := h.Svc.CreatePaymentRecord(ctx, model.PaymentRecord{
		ID:                     objID,
		Email:                  req.Email,
		EmailStatus:            model.EmailStatus(req.EmailStatus),
		PaymentType:            model.PaymentType(req.PaymentType),
		PaymentGatewayResponse: pgResponse,
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "user.handler.CreatePaymentRecord: %v", err)
	}

	return nil, nil
}
