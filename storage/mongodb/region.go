package mongodb

import (
	"context"
	"encoding/json"

	"github.com/e-space-uz/backend/config"
	"github.com/e-space-uz/backend/models"
	"github.com/e-space-uz/backend/pkg/utils"
	"github.com/e-space-uz/backend/storage/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type regionRepo struct {
	collection *mongo.Collection
}

func NewRegionRepo(db *mongo.Database) repo.RegionI {
	return &regionRepo{
		collection: db.Collection(config.RegionCollection)}
}
func (rr *regionRepo) Get(ctx context.Context, id string) (*models.Region, error) {
	var (
		region         []*models.Region
		response       *models.Region
		pipeline       = mongo.Pipeline{}
		cityCollection = config.CityCollection
	)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline,
		bson.D{
			primitive.E{Key: "$match", Value: bson.D{
				primitive.E{Key: "_id", Value: objectID}}}},
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: cityCollection},
				primitive.E{Key: "localField", Value: "city_id"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "city"}}}},
		bson.D{primitive.E{Key: "$unwind", Value: bson.D{
			primitive.E{Key: "path", Value: "$city"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}})

	row, err := rr.collection.Aggregate(
		ctx,
		pipeline)
	defer func() {
		_ = row.Close(ctx)
	}()
	if err != nil {
		return nil, err
	}

	if err := row.All(ctx, &region); err != nil {
		return nil, err
	}

	if len(region) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	byte, err := json.Marshal(region[0])
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(byte, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (rr *regionRepo) GetAll(ctx context.Context, page, limit uint32) ([]*models.Region, uint32, error) {
	var (
		regions        []*models.Region
		response       []*models.Region
		filter         = bson.D{}
		cityCollection = config.CityCollection
		pipeline       = mongo.Pipeline{}
		skip           = (page - 1) * limit
	)

	pipeline = append(pipeline,
		bson.D{primitive.E{Key: "$skip", Value: skip}},
		bson.D{primitive.E{Key: "$limit", Value: limit}},
		bson.D{primitive.E{Key: "$sort", Value: bson.M{"name": -1}}},
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: cityCollection},
				primitive.E{Key: "localField", Value: "city_id"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "city"}}}},
		bson.D{primitive.E{Key: "$unwind", Value: bson.D{
			primitive.E{Key: "path", Value: "$city"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}},
	)

	count, err := rr.collection.CountDocuments(ctx, filter)

	if err != nil {
		return nil, 0, err
	}

	rows, err := rr.collection.Aggregate(
		ctx, pipeline)

	if err != nil {
		return nil, 0, err
	}

	if err := rows.All(ctx, &regions); err != nil {
		return nil, 0, err
	}
	if err := utils.MarshalUnmarshal(regions, &response); err != nil {
		return nil, 0, err
	}
	return response, uint32(count), nil
}
func (rr *regionRepo) GetAllByCity(ctx context.Context, cityID, name string) ([]*models.Region, uint32, error) {
	var (
		regions        []*models.Region
		response       []*models.Region
		filter         = bson.D{}
		cityCollection = config.CityCollection
		pipeline       = mongo.Pipeline{}
	)

	if name != "" {
		filter = append(filter, bson.E{Key: "name", Value: bson.D{
			primitive.E{Key: "$regex", Value: name},
			primitive.E{Key: "$options", Value: "im"},
		}})
		pipeline = append(pipeline, bson.D{
			primitive.E{Key: "$match", Value: bson.D{
				primitive.E{Key: "name", Value: bson.D{
					primitive.E{Key: "$regex", Value: name},
					primitive.E{Key: "$options", Value: "im"},
				}}}}})
	}

	ID, err := primitive.ObjectIDFromHex(cityID)
	if err != nil {
		return nil, 0, err
	}
	filter = append(filter, bson.E{Key: "city_id", Value: ID})

	pipeline = append(pipeline,
		bson.D{
			primitive.E{Key: "$match", Value: bson.D{
				primitive.E{Key: "city_id", Value: ID}}}},
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: cityCollection},
				primitive.E{Key: "localField", Value: "city_id"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "city"}}}},
		bson.D{primitive.E{Key: "$unwind", Value: bson.D{
			primitive.E{Key: "path", Value: "$city"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}})

	count, err := rr.collection.CountDocuments(ctx, filter)

	if err != nil {
		return nil, 0, err
	}

	rows, err := rr.collection.Aggregate(
		ctx, pipeline)

	if err != nil {
		return nil, 0, err
	}

	if err := rows.All(ctx, &regions); err != nil {
		return nil, 0, err
	}

	if err := utils.MarshalUnmarshal(regions, &response); err != nil {
		return nil, 0, err
	}
	return response, uint32(count), nil
}
