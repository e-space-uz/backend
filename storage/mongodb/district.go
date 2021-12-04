package mongodb

import (
	"context"
	"fmt"

	"github.com/e-space-uz/backend/config"
	"github.com/e-space-uz/backend/models"
	"github.com/e-space-uz/backend/pkg/utils"
	"github.com/e-space-uz/backend/storage/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type districtRepo struct {
	collection *mongo.Collection
}

func NewDistrictRepo(db *mongo.Database) repo.DistrictI {
	return &districtRepo{
		collection: db.Collection(config.DistrictCollection)}
}

func (cr *districtRepo) Get(ctx context.Context, id string) (*models.District, error) {
	var districtDecode models.District
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if err := cr.collection.FindOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		}).Decode(&districtDecode); err != nil {
		return nil, err
	}

	return &districtDecode, nil
}

// This method to get all districts.
func (cr *districtRepo) GetAll(ctx context.Context, page, limit, code uint32, name string) ([]*models.District, uint32, error) {
	var (
		response         []*models.District
		districts        []*models.District
		filter           = bson.D{}
		cityCollection   = config.CityCollection
		regionCollection = config.RegionCollection
		skip             = (page - 1) * limit
		pipeline         = mongo.Pipeline{}
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
	if code != 0 {
		filter = append(filter, bson.E{Key: "code", Value: code})

		pipeline = append(pipeline, bson.D{
			primitive.E{Key: "$match", Value: bson.D{
				primitive.E{Key: "code", Value: code}}}})
	}

	count, err := cr.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return nil, 0, err
	}
	// projectID := bson.D{primitive.E{Key: "$project", Value: bson.D{primitive.E{Key: "id", Value: "$_id"}}}}
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
		bson.D{
			primitive.E{Key: "$unwind", Value: bson.D{
				primitive.E{Key: "path", Value: "$city"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}},
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: regionCollection},
				primitive.E{Key: "localField", Value: "region_id"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "region"}}}},
		bson.D{
			primitive.E{Key: "$unwind", Value: bson.D{
				primitive.E{Key: "path", Value: "$region"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}},
	)

	rows, err := cr.collection.Aggregate(
		context.Background(), pipeline)

	if err != nil {
		return nil, 0, err
	}
	if err := rows.All(context.Background(), &districts); err != nil {
		return nil, 0, err
	}
	if err := utils.MarshalUnmarshal(districts, &response); err != nil {
		return nil, 0, err
	}
	fmt.Println(districts)

	return response, uint32(count), nil
}
func (cr *districtRepo) GetAllByCityRegion(ctx context.Context, regionID, cityID, name string) ([]*models.District, uint32, error) {
	var (
		response         []*models.District
		districts        []*models.District
		filter           = bson.D{}
		cityCollection   = config.CityCollection
		regionCollection = config.RegionCollection
		pipeline         = mongo.Pipeline{}
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

	regionObjectID, err := primitive.ObjectIDFromHex(regionID)
	if err != nil {
		return nil, 0, err
	}
	cityObjectID, err := primitive.ObjectIDFromHex(cityID)
	if err != nil {
		return nil, 0, err
	}
	filter = append(filter, bson.E{Key: "region_id", Value: regionObjectID}, bson.E{Key: "city_id", Value: cityObjectID})

	pipeline = append(pipeline, bson.D{
		primitive.E{Key: "$match", Value: bson.D{
			primitive.E{Key: "city_id", Value: cityObjectID},
			primitive.E{Key: "region_id", Value: regionObjectID}}}})

	count, err := cr.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return nil, 0, err
	}
	// projectID := bson.D{primitive.E{Key: "$project", Value: bson.D{primitive.E{Key: "id", Value: "$_id"}}}}
	pipeline = append(pipeline,
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: cityCollection},
				primitive.E{Key: "localField", Value: "city_id"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "city"}}}},
		bson.D{
			primitive.E{Key: "$unwind", Value: bson.D{
				primitive.E{Key: "path", Value: "$city"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}},
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: regionCollection},
				primitive.E{Key: "localField", Value: "region_id"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "region"}}}},
		bson.D{
			primitive.E{Key: "$unwind", Value: bson.D{
				primitive.E{Key: "path", Value: "$region"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}},
		bson.D{
			primitive.E{Key: "$sort", Value: bson.M{"name": -1}}},
	)

	rows, err := cr.collection.Aggregate(
		context.Background(), pipeline)

	if err != nil {
		return nil, 0, err
	}
	if err := rows.All(context.Background(), &districts); err != nil {
		return nil, 0, err
	}
	if err := utils.MarshalUnmarshal(districts, &response); err != nil {
		return nil, 0, err
	}
	return response, uint32(count), nil
}
