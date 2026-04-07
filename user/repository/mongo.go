package repository

import "go.mongodb.org/mongo-driver/v2/mongo"

type mongoRepository struct {
	topUpCol *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *mongoRepository {
	topUpCol := db.Collection("top_up")
	return &mongoRepository{topUpCol: topUpCol}
}