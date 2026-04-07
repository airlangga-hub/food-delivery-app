package repository

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type xenditRepository struct {
	topupCol                *mongo.Collection
	xenditPaymentSessionURL string
	xenditAPIkey            string
	validate                *validator.Validate
}

func NewXenditRepository(db *mongo.Database, xenditPaymentSessionURL, xenditAPIkey string, validate *validator.Validate) *xenditRepository {
	topupCol := db.Collection("topup")

	return &xenditRepository{
		topupCol:                topupCol,
		xenditPaymentSessionURL: xenditPaymentSessionURL,
		xenditAPIkey:            xenditAPIkey,
		validate:                validate,
	}
}
