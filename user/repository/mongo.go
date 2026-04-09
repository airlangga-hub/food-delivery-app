package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/airlangga-hub/food-delivery-app/user/model"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type EmailStatus string

const (
	StatusPending EmailStatus = "PENDING"
	StatusSending EmailStatus = "SENDING"
	StatusSent    EmailStatus = "SENT"
)

type mongoRepository struct {
	paymentsCol     *mongo.Collection
	validate        *validator.Validate
	mailjetSender   string
	mailjetURL      string
	mailjetUsername string
	mailjetPassword string
}

func NewMongoRepository(db *mongo.Database, validate *validator.Validate, mailjetSender, mailjetURL, mailjetUsername, mailjetPassword string) *mongoRepository {
	paymentsCol := db.Collection("payments")
	return &mongoRepository{
		paymentsCol:     paymentsCol,
		validate:        validate,
		mailjetSender:   mailjetSender,
		mailjetURL:      mailjetURL,
		mailjetUsername: mailjetUsername,
		mailjetPassword: mailjetPassword,
	}
}

func (r *mongoRepository) CreatePaymentRecord(ctx context.Context, paymentRecord model.PaymentRecord) error {
	_, err := r.paymentsCol.InsertOne(ctx, paymentRecord)
	if err != nil {
		return fmt.Errorf("user.mongoRepository.CreatePaymentRecord: %w", err)
	}
	return nil
}

func (r *mongoRepository) GetPendingEmailRecords(ctx context.Context) ([]model.PaymentRecord, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"email_status": StatusPending},
			{
				"email_status": StatusSending,
				"updated_at":        bson.M{"$lt": time.Now().Add(-20 * time.Minute)}, // stuck for 20 mins
			},
		},
	}

	cursor, err := r.paymentsCol.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("payment.mongoRepository.GetPendingEmailRecords: %w", err)
	}
	defer cursor.Close(ctx)

	var records []model.PaymentRecord
	if err := cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("payment.mongoRepository.GetPendingEmailRecords: %w", err)
	}

	return records, nil
}

func (r *mongoRepository) SendEmail(ctx context.Context, rec model.PaymentRecord) error {
	var itemsList strings.Builder
	for _, p := range rec.PaymentGatewayResponse.Items {
		fmt.Fprintf(
			&itemsList,
			"- %s (Qty: %d) - Rp %d\n",
			p.Name, p.Quantity, p.NetUnitAmount)
	}

	var messageBuilder strings.Builder
	fmt.Fprintf(
		&messageBuilder,
		"Thank you for your order, %s!\n\nItems:\n%s\nTotal Amount: Rp %d\n\n",
		rec.Email,
		itemsList.String(),
		rec.PaymentGatewayResponse.Amount,
	)

	payload := model.MailjetRequest{
		Messages: []model.MessageRequest{
			{
				From: model.Person{
					Email: r.mailjetSender,
					Name:  "Transaction FTGO 14",
				},
				To: []model.Person{
					{
						Email: rec.Email,
						Name:  "Library Renter",
					},
				},
				Subject:  fmt.Sprintf("Order Confirmation #%s", rec.PaymentGatewayResponse.ReferenceID),
				TextPart: messageBuilder.String(),
			},
		},
	}

	ppayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("payment.repo.SendEmail: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", r.mailjetURL, bytes.NewReader(ppayload))
	if err != nil {
		return fmt.Errorf("payment.repo.SendEmail: %w", err)
	}

	req.SetBasicAuth(r.mailjetUsername, r.mailjetPassword)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("payment.repo.SendEmail: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("payment.repo.SendEmail: %v", string(body))
	}

	var mailjetResp model.MailjetResponse
	if err := json.NewDecoder(res.Body).Decode(&mailjetResp); err != nil {
		return fmt.Errorf("payment.repo.SendEmail: %w", err)
	}

	if mailjetResp.StatusCode >= 400 {
		return fmt.Errorf("payment.repo.SendEmail: %s", mailjetResp.ErrorMessage)
	}

	return nil
}

func (r *mongoRepository) MarkEmailAsSending(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("payment.mongoRepository.MarkEmailAsSent: %w", err)
	}

	filter := bson.M{
		"_id": objID,
	}

	update := bson.M{
		"$set": bson.M{
			"email_status": StatusSending,
			"updated_at":        time.Now(),
		},
	}

	result, err := r.paymentsCol.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("payment.mongoRepository.MarkEmailAsSent: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("payment.mongoRepository.MarkEmailAsSent: %w", fmt.Errorf("no pending record found with ID %s", id))
	}

	return nil
}

func (r *mongoRepository) IfErrorMarkEmailAsPending(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("payment.mongoRepository.MarkEmailAsSent: %w", err)
	}

	filter := bson.M{
		"_id": objID,
	}

	update := bson.M{
		"$set": bson.M{
			"email_status": StatusPending,
			"updated_at":        time.Now(),
		},
	}

	result, err := r.paymentsCol.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("payment.mongoRepository.IfErrorMarkEmailAsPending: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("payment.mongoRepository.IfErrorMarkEmailAsPending: %w", fmt.Errorf("no sending record found with ID %s", id))
	}

	return nil
}

func (r *mongoRepository) MarkEmailAsSent(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("payment.mongoRepository.MarkEmailAsSent: %w", err)
	}

	filter := bson.M{
		"_id": objID,
	}

	update := bson.M{
		"$set": bson.M{
			"email_status": StatusSent,
			"updated_at":        time.Now(),
		},
	}

	result, err := r.paymentsCol.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("payment.mongoRepository.MarkEmailAsSent: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("payment.mongoRepository.MarkEmailAsSent: %w", fmt.Errorf("no sending record found with ID %s", id))
	}

	return nil
}
