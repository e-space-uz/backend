package mongodb

import (
	"context"

	"github.com/e-space-uz/backend/config"
	"github.com/e-space-uz/backend/models"
	"github.com/e-space-uz/backend/pkg/utils"
	"github.com/e-space-uz/backend/storage/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type cityRepo struct {
	collection *mongo.Collection
}

func NewCityRepo(db *mongo.Database) repo.CityI {
	return &cityRepo{
		collection: db.Collection(config.CityCollection)}
}

func (cr *cityRepo) Get(ctx context.Context, id string) (*models.City, error) {
	var cityDecode models.City
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if err := cr.collection.FindOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		}).Decode(&cityDecode); err != nil {
		return nil, err
	}

	return &cityDecode, nil
}

func (cr *cityRepo) GetAll(ctx context.Context, page, limit, code uint32, name string) ([]*models.City, uint32, error) {
	var (
		response []*models.City
		cities   []*models.City
		filter   = bson.D{}
	)
	if name != "" {
		filter = append(filter, primitive.E{Key: "name", Value: bson.D{
			primitive.E{Key: "$regex", Value: name},
			primitive.E{Key: "$options", Value: "im"},
		}})
	}
	if code != 0 {
		filter = append(filter, primitive.E{Key: "code", Value: code})
	}
	opts := options.Find()
	skip := (page - 1) * limit
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(skip))
	opts.SetSort(bson.M{
		"name": 1,
	})
	count, err := cr.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return nil, 0, err
	}

	rows, err := cr.collection.Find(
		context.Background(),
		filter,
		opts,
	)
	if err != nil {
		return nil, 0, err
	}

	if err := rows.All(context.Background(), &cities); err != nil {
		return nil, 0, err
	}
	if err := utils.MarshalUnmarshal(cities, &response); err != nil {
		return nil, 0, err
	}
	return response, uint32(count), nil
}
