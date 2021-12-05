package mongodb

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/e-space-uz/backend/config"
	"github.com/e-space-uz/backend/models"
	"github.com/e-space-uz/backend/storage/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type entityFilesRepo struct {
	collection *mongo.Collection
}

func NewEntityFilesRepo(db *mongo.Database) repo.EntityFilesI {
	return &entityFilesRepo{
		collection: db.Collection(config.EntityFilesCollection),
	}
}
func (sr *entityFilesRepo) Create(ctx context.Context, entityFiles *models.CreateEntityFiles) (string, error) {
	_, err := sr.collection.InsertOne(
		context.Background(),
		entityFiles,
	)
	if err != nil {
		return "", err
	}

	return entityFiles.ID.Hex(), nil
}

func (sr *entityFilesRepo) Get(ctx context.Context, id string) (*models.EntityFiles, error) {
	var entityFilesDecode models.EntityFiles
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if err := sr.collection.FindOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		}).Decode(&entityFilesDecode); err != nil {
		return nil, err
	}

	return &models.EntityFiles{
		ID:       entityFilesDecode.ID,
		FileName: entityFilesDecode.FileName,
		Url:      entityFilesDecode.Url,
		Comment:  entityFilesDecode.Comment,
		User:     entityFilesDecode.User,
	}, nil
}

func (sr *entityFilesRepo) GetAll(ctx context.Context, page, limit uint32, search string) ([]*models.EntityFiles, uint32, error) {
	var (
		entityFileses []*models.EntityFiles
		filter        = bson.D{}
	)
	start := time.Now()
	if search != "" {
		filter = append(filter, bson.E{Key: "name", Value: primitive.E{Key: "$regex", Value: "*." + search + ".*"}})
	}

	opts := options.Find()

	skip := (page - 1) * limit
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(skip))
	opts.SetSort(bson.M{
		"created_at": -1,
	})
	count, err := sr.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return nil, 0, err
	}

	rows, err := sr.collection.Find(
		context.Background(),
		filter,
		opts,
	)
	defer func() {
		_ = rows.Close(context.Background())
	}()
	if err != nil {
		return nil, 0, err
	}
	for rows.Next(context.Background()) {
		var entityFiles *models.EntityFiles
		var entityFilesDecode models.EntityFiles

		if err := rows.Decode(&entityFilesDecode); err != nil {
			return nil, 0, err
		}
		entityFiles = &models.EntityFiles{
			ID:       entityFilesDecode.ID,
			FileName: entityFilesDecode.FileName,
			Url:      entityFilesDecode.Url,
			Comment:  entityFilesDecode.Comment,
			User:     entityFilesDecode.User,
		}
		entityFileses = append(entityFileses, entityFiles)
	}
	r := new(big.Int)
	fmt.Println(r.Binomial(1000, 10))

	elapsed := time.Since(start)
	fmt.Printf("Binomial took %s", elapsed)
	return entityFileses, uint32(count), nil
}
func (sr *entityFilesRepo) Update(ctx context.Context, entityFiles *models.EntityFiles) error {
	update := bson.M{
		"$set": bson.M{
			"file_name":  entityFiles.FileName,
			"url":        entityFiles.Url,
			"comment":    entityFiles.Comment,
			"user":       entityFiles.User,
			"updated_at": time.Now(),
		}}

	filter := bson.M{"_id": bson.M{"$eq": entityFiles.ID}}
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
func (sr *entityFilesRepo) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": bson.M{"$eq": id}}
	_, err := sr.collection.DeleteOne(
		context.Background(),
		filter)

	if err != nil {
		return err
	}

	return nil
}

func (sr *entityFilesRepo) EntityFileExists(ctx context.Context, id string) (bool, error) {
	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": ObjectID}
	res, err := sr.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, err
	}

	return res == 1, nil
}
