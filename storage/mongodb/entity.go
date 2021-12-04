package mongodb

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/e-space-uz/backend/config"
	"github.com/e-space-uz/backend/models"
	"github.com/e-space-uz/backend/pkg/utils"
	"github.com/e-space-uz/backend/storage/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type entityRepo struct {
	collection *mongo.Collection
}

func NewEntityRepo(db *mongo.Database) repo.EntityI {
	return &entityRepo{
		collection: db.Collection(config.EntityCollection),
	}
}

func (er *entityRepo) Create(ctx context.Context, entity *models.CreateUpdateEntity) (string, error) {
	var (
		entitySoato = strconv.Itoa(int(entity.District.Soato))
		filter      = bson.M{"entity_soato": entitySoato}
	)

	createEntity := &models.CreateUpdateEntity{
		ID:                 entity.ID,
		Status:             "statusID",
		EntitySoato:        entitySoato,
		EntityTypeCode:     entity.EntityTypeCode,
		Version:            uint64(2),
		Address:            entity.Address,
		EntityNumber:       entity.EntityNumber,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		EntityStatusUpdate: time.Now(),
		City: &models.City{
			ID:     entity.City.ID,
			Name:   entity.City.Name,
			RuName: entity.City.RuName,
			Code:   entity.City.Code,
			Soato:  entity.City.Soato,
		},
		Region: &models.Region{
			ID:     entity.Region.ID,
			Name:   entity.Region.Name,
			RuName: entity.Region.RuName,
			Code:   entity.Region.Code,
			Soato:  entity.Region.Soato,
		},
		District: &models.District{
			ID:     entity.District.ID,
			Name:   entity.District.Name,
			RuName: entity.District.RuName,
			Code:   entity.District.Code,
			Soato:  entity.District.Soato,
		},
	}
	count, err := er.collection.CountDocuments(ctx, filter)
	if err != nil {
		return "", err
	}
	createEntity.EntityNumber = "B" + entitySoato + "-" + strconv.Itoa(int(count+1))
	if entity.EntityProperties != nil {
		for _, property := range entity.EntityProperties {
			createEntity.EntityProperties = append(createEntity.EntityProperties, &models.CreateEntityProperty{
				PropertyID: property.PropertyID,
				Value:      property.Value,
			})
		}
	} else {
		createEntity.EntityProperties = []*models.CreateEntityProperty{}
	}

	if entity.EntityGallery != nil {
		for _, galleryID := range entity.EntityGallery {
			if err != nil {
				return "", err
			}
			createEntity.EntityGallery = append(createEntity.EntityGallery, galleryID)
		}
	} else {
		createEntity.EntityGallery = []string{}
	}

	createEntity.EntityDrafts = []primitive.ObjectID{}
	resp, err := er.collection.InsertOne(
		ctx,
		createEntity,
	)
	if err != nil {
		return "", err
	}
	return resp.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (er *entityRepo) Get(ctx context.Context, id string) (*models.Entity, error) {
	var (
		response              *models.Entity
		appDecode             []*models.Entity
		pipeline              = mongo.Pipeline{}
		propertyCollection    = config.PropertyCollection
		entityFilesCollection = config.EntityFilesCollection
		entityDraftCollection = config.EntityDraftCollection
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
			primitive.E{Key: "$unwind", Value: bson.D{
				primitive.E{Key: "path", Value: "$status"}, {
					Key: "preserveNullAndEmptyArrays", Value: false}}}},
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: entityFilesCollection},
				primitive.E{Key: "localField", Value: "entity_files"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "entity_files"}}}},
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: entityDraftCollection},
				primitive.E{Key: "localField", Value: "entity_drafts"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "entity_drafts"}}}},

		bson.D{
			primitive.E{Key: "$unwind", Value: "$entity_properties"}},
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: propertyCollection},
				primitive.E{Key: "localField", Value: "entity_properties.property_id"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "entity_properties.property"}}}},
		bson.D{primitive.E{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$entity_properties.property"}}}},

		bson.D{
			primitive.E{Key: "$group", Value: bson.D{
				primitive.E{Key: "_id", Value: "$_id"},
				primitive.E{Key: "external_id", Value: bson.D{
					primitive.E{Key: "$first", Value: "$external_id"}}},
				primitive.E{Key: "address", Value: bson.D{
					primitive.E{Key: "$first", Value: "$address"}}},
				primitive.E{Key: "entity_number", Value: bson.D{
					primitive.E{Key: "$first", Value: "$entity_number"}}},
				primitive.E{Key: "revert_comment", Value: bson.D{
					primitive.E{Key: "$first", Value: "$revert_comment"}}},
				primitive.E{Key: "entity_soato", Value: bson.D{
					primitive.E{Key: "$first", Value: "$entity_soato"}}},
				primitive.E{Key: "status", Value: bson.D{
					primitive.E{Key: "$first", Value: "$status"}}},
				primitive.E{Key: "version", Value: bson.D{
					primitive.E{Key: "$first", Value: "$version"}}},
				primitive.E{Key: "city", Value: bson.D{
					primitive.E{Key: "$first", Value: "$city"}}},
				primitive.E{Key: "region", Value: bson.D{
					primitive.E{Key: "$first", Value: "$region"}}},
				primitive.E{Key: "district", Value: bson.D{
					primitive.E{Key: "$first", Value: "$district"}}},
				primitive.E{Key: "entity_files", Value: bson.D{
					primitive.E{Key: "$first", Value: "$entity_files"}}},
				primitive.E{Key: "entity_drafts", Value: bson.D{
					primitive.E{Key: "$first", Value: "$entity_drafts"}}},
				primitive.E{Key: "created_at", Value: bson.D{
					primitive.E{Key: "$first", Value: "$created_at"}}},
				primitive.E{Key: "updated_at", Value: bson.D{
					primitive.E{Key: "$first", Value: "$updated_at"}}},
				primitive.E{Key: "entity_gallery", Value: bson.D{
					primitive.E{Key: "$first", Value: "$entity_gallery"}}},
				primitive.E{Key: "organizations", Value: bson.D{
					primitive.E{Key: "$first", Value: "$organizations"}}},
				primitive.E{Key: "deadline", Value: bson.D{
					primitive.E{Key: "$first", Value: "$deadline"}}},
				primitive.E{Key: "entity_properties", Value: bson.D{
					primitive.E{Key: "$push", Value: "$entity_properties"}}}}},
		})

	row, err := er.collection.Aggregate(
		ctx,
		pipeline)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = row.Close(ctx)
	}()

	if err := row.All(ctx, &appDecode); err != nil {
		return nil, err
	}
	if len(appDecode) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	byteObject, err := json.Marshal(appDecode[0])
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(byteObject, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (er *entityRepo) GetAll(ctx context.Context, req *models.GetAllEntitiesRequest) ([]*models.GetAllEntities, uint64, error) {
	var (
		res      []*models.GetAllEntities
		entities []*models.GetAllEntities
		skip     = (req.Page - 1) * req.Limit
		count    = make(chan int64)
		response = make(chan []*models.GetAllEntities)
		errChan  = make(chan error, 2)
	)

	filter, pipeline, err := getAllFilter(req)
	if err != nil {
		return nil, 0, err
	}
	go func(filter bson.D) {
		c, err := er.collection.CountDocuments(ctx, filter)
		errChan <- err
		count <- c
	}(filter)

	go func(pipeline mongo.Pipeline) {
		pipeline = append(pipeline,
			bson.D{primitive.E{Key: "$project", Value: bson.D{
				primitive.E{Key: "_id", Value: 1},
				primitive.E{Key: "entity_number", Value: 1},
				primitive.E{Key: "status_id", Value: 1},
				primitive.E{Key: "entity_soato", Value: 1},
				primitive.E{Key: "version", Value: 1},
				primitive.E{Key: "address", Value: 1},
				primitive.E{Key: "city", Value: 1},
				primitive.E{Key: "region", Value: 1},
				primitive.E{Key: "district", Value: 1},
				primitive.E{Key: "entity_properties", Value: 1},
				primitive.E{Key: "created_at", Value: 1},
			}}},
			bson.D{primitive.E{Key: "$sort", Value: bson.D{primitive.E{Key: "created_at", Value: -1}}}},
			bson.D{primitive.E{Key: "$skip", Value: skip}},
			bson.D{primitive.E{Key: "$limit", Value: req.Limit}},

			bson.D{
				primitive.E{Key: "$unwind", Value: bson.D{
					primitive.E{Key: "path", Value: "$status"}, {
						Key: "preserveNullAndEmptyArrays", Value: false}}}})
		rows, err := er.collection.Aggregate(
			ctx,
			pipeline)
		if err != nil {
			return
		}
		defer func() {
			errChan <- err
			response <- res
			rows.Close(ctx)
		}()
		if err = rows.All(ctx, &entities); err != nil {
			return
		}
		err = utils.MarshalUnmarshal(entities, &res)

		if err != nil {
			return
		}
	}(pipeline)
	if err := <-errChan; err != nil {
		return nil, 0, <-errChan
	}
	return <-response, uint64(<-count), <-errChan
}

func (er *entityRepo) GetAllWithProperties(ctx context.Context, req *models.GetAllEntitiesRequest) ([]*models.GetAllEntities, error) {
	var (
		res      []*models.GetAllEntities
		entities []*models.GetAllEntities
		skip     = (req.Page - 1) * req.Limit
	)

	_, pipeline, err := getAllFilter(req)
	if err != nil {
		return nil, err
	}
	pipeline = append(pipeline,
		bson.D{primitive.E{Key: "$project", Value: bson.D{
			primitive.E{Key: "_id", Value: 1},
			primitive.E{Key: "entity_number", Value: 1},
			primitive.E{Key: "status_id", Value: 1},
			primitive.E{Key: "entity_soato", Value: 1},
			primitive.E{Key: "version", Value: 1},
			primitive.E{Key: "address", Value: 1},
			primitive.E{Key: "city", Value: 1},
			primitive.E{Key: "region", Value: 1},
			primitive.E{Key: "district", Value: 1},
			primitive.E{Key: "entity_properties", Value: 1},
			primitive.E{Key: "created_at", Value: 1},
		}}},
		bson.D{primitive.E{Key: "$sort", Value: bson.D{primitive.E{Key: "created_at", Value: -1}}}},
		bson.D{primitive.E{Key: "$skip", Value: skip}},
		bson.D{primitive.E{Key: "$limit", Value: req.Limit}},
		bson.D{
			primitive.E{Key: "$unwind", Value: bson.D{
				primitive.E{Key: "path", Value: "$status"}, {
					Key: "preserveNullAndEmptyArrays", Value: false}}}})

	rows, err := er.collection.Aggregate(
		ctx,
		pipeline)
	if err != nil {
		return nil, err
	}
	defer func() {
		rows.Close(ctx)
	}()
	if err = rows.All(ctx, &entities); err != nil {
		return nil, err
	}
	err = utils.MarshalUnmarshal(entities, &res)
	return res, err
}

func (er *entityRepo) Delete(ctx context.Context, id string) error {
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		}}

	filter := bson.M{"_id": bson.M{"$eq": id}}
	_, err := er.collection.UpdateOne(
		ctx,
		filter,
		update,
	)

	if err != nil {
		return err
	}

	return nil
}

func (er *entityRepo) UpdateEntityDrafts(ctx context.Context, entityID, entityDraftID string) error {
	entityObjectID, err := primitive.ObjectIDFromHex(entityID)
	if err != nil {
		return err
	}
	entityDraftObjectID, err := primitive.ObjectIDFromHex(entityDraftID)
	if err != nil {
		return err
	}
	update := bson.M{
		"$push": bson.M{
			"entity_drafts": entityDraftObjectID,
		}}
	filter := bson.M{"_id": bson.M{"$eq": entityObjectID}}
	_, err = er.collection.UpdateOne(
		ctx,
		filter,
		update,
	)
	return err
}

// filter in one function
func getAllFilter(req *models.GetAllEntitiesRequest) (bson.D, mongo.Pipeline, error) {
	var (
		filter   = bson.D{}
		pipeline = mongo.Pipeline{}
	)

	if req.RegionID != "" {
		filter = append(filter, primitive.E{Key: "region.id", Value: req.RegionID})
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$match", Value: bson.D{primitive.E{Key: "region.id", Value: req.RegionID}}}})
	}
	if req.CityID != "" {
		filter = append(filter, primitive.E{Key: "city.id", Value: req.CityID})
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$match", Value: bson.D{primitive.E{Key: "city.id", Value: req.CityID}}}})
	}

	return filter, pipeline, nil

}
