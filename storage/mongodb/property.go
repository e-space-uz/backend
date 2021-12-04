package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/e-space-uz/backend/config"
	"github.com/e-space-uz/backend/models"
	"github.com/e-space-uz/backend/storage/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type propertyRepo struct {
	collection *mongo.Collection
}

func NewPropertyRepo(db *mongo.Database) repo.PropertyI {
	return &propertyRepo{
		collection: db.Collection(config.PropertyCollection),
	}
}
func (pr *propertyRepo) Create(property *models.CreateUpdateProperty) (string, error) {
	createUpdateProperty := &models.CreateUpdateProperty{
		ID:          property.ID,
		Name:        property.Name,
		Type:        property.Type,
		Label:       property.Label,
		Placeholder: property.Placeholder,
		IsRequired:  property.IsRequired,
		Description: property.Description,
		Validation:  property.Validation,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	for _, option := range property.PropertyOptions {
		createUpdateProperty.PropertyOptions = append(createUpdateProperty.PropertyOptions, &models.PropertyOption{
			Name:  option.Name,
			Value: option.Value,
		})
	}

	resp, err := pr.collection.InsertOne(
		context.Background(),
		createUpdateProperty,
	)
	if err != nil {
		return "", err
	}
	return resp.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (pr *propertyRepo) Get(id string) (*models.Property, error) {
	var (
		response *models.Property
		property []*models.Property
	)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	matchPropertyID := bson.D{
		primitive.E{Key: "$match", Value: bson.D{
			primitive.E{Key: "_id", Value: objectID}}}}

	row, err := pr.collection.Aggregate(
		context.Background(),
		mongo.Pipeline{
			matchPropertyID,
		})
	defer func() {
		_ = row.Close(context.Background())
	}()
	if err != nil {
		return nil, err
	}

	if err := row.All(context.Background(), &property); err != nil {
		return nil, err
	}

	if len(property) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	byte, err := json.Marshal(property[0])
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(byte, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (pr *propertyRepo) GetAll(page, limit uint32, search string) ([]*models.Property, uint32, error) {
	var (
		response    []*models.Property
		properties  []*models.Property
		pipeline    = mongo.Pipeline{}
		filter      = bson.D{}
		filterCount = bson.D{}
		skip        = (page - 1) * limit
	)
	start := time.Now()
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

	count, err := pr.collection.CountDocuments(context.Background(), filterCount)

	if err != nil {
		return nil, 0, err
	}

	pipeline = append(pipeline,
		bson.D{primitive.E{Key: "$sort", Value: bson.D{primitive.E{Key: "name", Value: -1}}}},
		bson.D{primitive.E{Key: "$skip", Value: skip}},
		bson.D{primitive.E{Key: "$limit", Value: limit}})

	rows, err := pr.collection.Aggregate(
		context.Background(),
		pipeline)
	defer func() {
		_ = rows.Close(context.Background())
	}()

	if err != nil {
		return nil, 0, err
	}

	if err := rows.All(context.Background(), &properties); err != nil {
		return nil, 0, err
	}

	byte, err := json.Marshal(properties)
	if err != nil {
		return nil, 0, err
	}
	if err := json.Unmarshal(byte, &response); err != nil {
		return nil, 0, err
	}
	r := new(big.Int)
	fmt.Println(r.Binomial(1000, 10))

	elapsed := time.Since(start)
	fmt.Printf("Binomial took %s", elapsed)
	return response, uint32(count), nil
}

func (pr *propertyRepo) Delete(id string) error {
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		}}

	filter := bson.M{"_id": bson.M{"$eq": id}}
	_, err := pr.collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)

	if err != nil {
		return err
	}

	return nil
}
func (pr *propertyRepo) Update(property *models.CreateUpdateProperty) error {
	updateProperty := &models.CreateUpdateProperty{
		ID:          property.ID,
		Name:        property.Name,
		Type:        property.Type,
		Label:       property.Label,
		Placeholder: property.Placeholder,
		IsRequired:  property.IsRequired,
		Description: property.Description,
		Validation:  property.Validation,
	}
	for _, option := range property.PropertyOptions {
		updateProperty.PropertyOptions = append(updateProperty.PropertyOptions, &models.PropertyOption{
			Name:  option.Name,
			Value: option.Value,
		})
	}

	update := bson.M{
		"$set": bson.M{
			"name":             updateProperty.Name,
			"type":             updateProperty.Type,
			"label":            updateProperty.Label,
			"placeholder":      updateProperty.Placeholder,
			"is_required":      updateProperty.IsRequired,
			"validation":       updateProperty.Validation,
			"description":      updateProperty.Description,
			"property_options": updateProperty.PropertyOptions,
			"updated_at":       time.Now(),
		}}

	filter := bson.M{"_id": bson.M{"$eq": property.ID}}
	_, err := pr.collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)

	if err != nil {
		return err
	}

	return nil
}

func (cr *propertyRepo) PropertyExists(id string) (bool, error) {
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
