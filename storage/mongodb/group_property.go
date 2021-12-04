package mongodb

import (
	"context"
	"encoding/json"
	"time"

	"github.com/e-space-uz/backend/config"
	"github.com/e-space-uz/backend/models"
	"github.com/e-space-uz/backend/storage/repo"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type groupProperty struct {
	collection *mongo.Collection
}

func NewGroupPropertyRepo(db *mongo.Database) repo.GroupPropertyI {
	return &groupProperty{
		collection: db.Collection(config.GroupPropertyCollection),
	}
}
func (sr *groupProperty) Create(ctx context.Context, groupProperty *models.CreateGroupProperty) (string, error) {
	var groupPropertiesArray = make([]*models.CreateProperties, len(groupProperty.Properties))

	for _, property := range groupProperty.Properties {
		groupPropertiesArray[property.Order] = &models.CreateProperties{
			PropertyID: property.PropertyID,
			Order:      property.Order}
	}

	resp, err := sr.collection.InsertOne(
		context.Background(),
		groupPropertiesArray,
	)
	if err != nil {
		return "", err
	}
	return resp.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (sr *groupProperty) Get(ctx context.Context, id string) (*models.GroupProperty, error) {
	var (
		groupPropertyDecode []*models.GroupProperty
		response            *models.GroupProperty
		propertyCollection  = config.PropertyCollection
		pipeline            = mongo.Pipeline{}
	)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "error")
	}

	pipeline = append(pipeline,
		bson.D{
			bson.E{Key: "$match", Value: bson.D{
				bson.E{Key: "_id", Value: objectID}}}},
		bson.D{
			bson.E{Key: "$lookup", Value: bson.D{
				bson.E{Key: "from", Value: propertyCollection},
				bson.E{Key: "localField", Value: "properties.property_id"},
				bson.E{Key: "foreignField", Value: "_id"},
				bson.E{Key: "as", Value: "properties"}}}},
	)

	rows, err := sr.collection.Aggregate(context.Background(), pipeline)
	defer func() {
		_ = rows.Close(context.Background())
	}()

	if err != nil {
		return nil, errors.Wrap(err, "error")
	}
	if err := rows.All(context.Background(), &groupPropertyDecode); err != nil {
		return nil, errors.Wrap(err, "error")
	}
	if len(groupPropertyDecode) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	byte, err := json.Marshal(groupPropertyDecode[0])
	if err != nil {
		return nil, errors.Wrap(err, "error")
	}
	if err := json.Unmarshal(byte, &response); err != nil {
		return nil, errors.Wrap(err, "error")
	}

	return response, nil
}

func (sr *groupProperty) GetAll(ctx context.Context, page, limit uint32, search string) ([]*models.GetAllGroupProperty, uint32, error) {
	var (
		groupPropertiesDecode   []*models.GetAllGroupProperty
		groupPropertiesResponse []*models.GetAllGroupProperty
		filter                  = bson.D{}
		filterCount             = bson.D{}
		skip                    = (page - 1) * limit
		pipeline                = mongo.Pipeline{}
	)

	if search != "" {
		filter = bson.D{primitive.E{Key: "$match", Value: bson.D{
			primitive.E{Key: "name", Value: bson.D{
				primitive.E{Key: "$regex", Value: search},
				primitive.E{Key: "$options", Value: "im"}}}}}}
		pipeline = append(pipeline, filter)
		filterCount = bson.D{primitive.E{Key: "name", Value: bson.D{
			primitive.E{Key: "$regex", Value: search},
			primitive.E{Key: "$options", Value: "im"}}}}
	}

	count, err := sr.collection.CountDocuments(context.Background(), filterCount)

	if err != nil {
		return nil, 0, err
	}
	pipeline = append(pipeline,
		bson.D{primitive.E{Key: "$skip", Value: skip}},
		bson.D{primitive.E{Key: "$limit", Value: limit}})

	rows, err := sr.collection.Aggregate(context.Background(), pipeline)
	defer func() {
		_ = rows.Close(context.Background())
	}()

	if err != nil {
		return nil, 0, err
	}
	if err := rows.All(context.Background(), &groupPropertiesDecode); err != nil {
		return nil, 0, err
	}
	byte, err := json.Marshal(groupPropertiesDecode)
	if err != nil {
		return nil, 0, err
	}
	if err := json.Unmarshal(byte, &groupPropertiesResponse); err != nil {
		return nil, 0, err
	}

	return groupPropertiesResponse, uint32(count), nil
}

func (sr *groupProperty) Delete(ctx context.Context, id string) error {
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		}}

	filter := bson.M{"_id": bson.M{"$eq": id}}
	_, err := sr.collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)

	if err != nil {
		return err
	}

	return nil
}
func (sr *groupProperty) Update(ctx context.Context, groupProperty *models.CreateGroupProperty) error {
	org := models.OrganizationCreate{}

	updateGroupProperty := models.CreateGroupProperty{
		Name:         groupProperty.Name,
		Step:         groupProperty.Step,
		Type:         groupProperty.Type,
		Status:       groupProperty.Status,
		Description:  groupProperty.Description,
		Organization: org,
	}
	for _, property := range groupProperty.Properties {
		updateGroupProperty.Properties = append(updateGroupProperty.Properties, &models.CreateProperties{
			PropertyID: property.PropertyID,
			Order:      property.Order,
		})
	}

	update := bson.M{
		"$set": bson.M{
			"name":           updateGroupProperty.Name,
			"step":           updateGroupProperty.Step,
			"type":           updateGroupProperty.Type,
			"status":         updateGroupProperty.Status,
			"description":    updateGroupProperty.Description,
			"properties":     updateGroupProperty.Properties,
			"write_statuses": updateGroupProperty.WriteStatuses,
			"read_statuses":  updateGroupProperty.ReadStatuses,
			"organization":   updateGroupProperty.Organization,
			"updated_at":     time.Now(),
		}}
	_ = update
	return nil
}

func (cr *groupProperty) GroupPropertyExists(ctx context.Context, id string) (bool, error) {
	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": ObjectID}
	res, err := cr.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, err
	}

	return res == 1, nil
}

func (sr *groupProperty) GetAllByType(ctx context.Context, typeOf, step uint32) ([]*models.GroupProperty, uint32, error) {
	var (
		groupPropertiesDecode []*models.GroupProperty
		groupProperties       []*models.GroupProperty
		filter                = bson.D{}
		propertyCollection    = config.PropertyCollection
		pipeline              = mongo.Pipeline{}
	)

	if step != 0 {
		filter = append(filter, primitive.E{Key: "step", Value: step})
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$match", Value: bson.D{primitive.E{Key: "step", Value: step}}}})
	}
	if typeOf != 0 {
		filter = append(filter, primitive.E{Key: "type", Value: typeOf})
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$match", Value: bson.D{primitive.E{Key: "type", Value: typeOf}}}})
	}

	count, err := sr.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return nil, 0, err
	}
	pipeline = append(pipeline, bson.D{
		bson.E{Key: "$lookup", Value: bson.D{
			bson.E{Key: "from", Value: propertyCollection},
			bson.E{Key: "localField", Value: "properties.property_id"},
			bson.E{Key: "foreignField", Value: "_id"},
			bson.E{Key: "as", Value: "properties"}}}},
	)

	rows, err := sr.collection.Aggregate(context.Background(), pipeline)
	defer func() {
		_ = rows.Close(context.Background())
	}()

	if err != nil {
		return nil, 0, err
	}
	if err := rows.All(context.Background(), &groupPropertiesDecode); err != nil {
		return nil, 0, err
	}
	byte, err := json.Marshal(groupPropertiesDecode)
	if err != nil {
		return nil, 0, err
	}
	if err := json.Unmarshal(byte, &groupProperties); err != nil {
		return nil, 0, err
	}

	return groupProperties, uint32(count), nil
}
