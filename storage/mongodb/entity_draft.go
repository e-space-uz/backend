package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
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

type entityDraftRepo struct {
	collection *mongo.Collection
}

func NewEntityDraftRepo(db *mongo.Database) repo.EntityDraftI {
	return &entityDraftRepo{
		collection: db.Collection(config.EntityDraftCollection),
	}
}
func (cr entityDraftRepo) Create(ctx context.Context, req *models.CreateEntityDraft) (string, error) {
	soatoCount, err := cr.collection.CountDocuments(
		ctx,
		bson.D{bson.E{
			Key: "entity_draft_soato", Value: req.EntityDraftSoato,
		}},
	)
	if err != nil {
		return "", err
	}

	draftNumber := "T" + strconv.Itoa(int(req.Region.Soato)) + "-" + strconv.Itoa(int(soatoCount)+1)

	createEntity := &models.CreateEntityDraft{
		ID:                req.ID,
		Status:            "new",
		Comment:           req.Comment,
		EntityDraftNumber: draftNumber,
		EntityDraftSoato:  req.EntityDraftSoato,

		City: models.City{
			ID:     req.City.ID,
			Name:   req.City.Name,
			RuName: req.City.RuName,
			Code:   req.City.Code,
			Soato:  req.City.Soato,
		},
		Region: models.Region{
			ID:     req.Region.ID,
			Name:   req.Region.Name,
			RuName: req.Region.RuName,
			Code:   req.Region.Code,
			Soato:  req.Region.Soato,
		},
		District: models.District{
			ID:     req.District.ID,
			Name:   req.District.Name,
			RuName: req.District.RuName,
			Code:   req.District.Code,
			Soato:  req.District.Soato,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	for _, property := range req.EntityProperties {
		createEntity.EntityProperties = append(createEntity.EntityProperties, &models.CreateEntityProperty{
			PropertyID: property.PropertyID,
			Value:      property.Value,
		})
	}
	if len(req.EntityGallery) != 0 {
		for _, galleryID := range req.EntityGallery {
			if err != nil {
				return "", err
			}
			createEntity.EntityGallery = append(createEntity.EntityGallery, galleryID)
		}

	} else {
		createEntity.EntityGallery = []string{}
	}
	resp, err := cr.collection.InsertOne(
		ctx,
		createEntity)

	return resp.InsertedID.(primitive.ObjectID).Hex(), err
}

func (cr entityDraftRepo) Get(ctx context.Context, id string) (*models.EntityDraft, error) {
	var (
		response           models.EntityDraft
		appDecode          []*models.EntityDraft
		pipeline           = mongo.Pipeline{}
		entityCollection   = config.EntityCollection
		propertyCollection = config.PropertyCollection
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
			primitive.E{Key: "$unwind", Value: "$entity_properties"}},
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: propertyCollection},
				primitive.E{Key: "localField", Value: "entity_properties.property_id"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "entity_properties.property"}}}},
		bson.D{primitive.E{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$entity_properties.property"}}}},
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: entityCollection},
				primitive.E{Key: "localField", Value: "entity_id"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "entity"}}}},
		bson.D{
			primitive.E{Key: "$unwind", Value: bson.D{
				primitive.E{Key: "path", Value: "$entity"}, {
					Key: "preserveNullAndEmptyArrays", Value: true}}}},
		bson.D{
			primitive.E{Key: "$group", Value: bson.D{
				primitive.E{Key: "_id", Value: "$_id"},
				primitive.E{Key: "entity_draft_number", Value: bson.D{
					primitive.E{Key: "$first", Value: "$entity_draft_number"}}},
				primitive.E{Key: "entity_draft_soato", Value: bson.D{
					primitive.E{Key: "$first", Value: "$entity_draft_soato"}}},
				primitive.E{Key: "status", Value: bson.D{
					primitive.E{Key: "$first", Value: "$status"}}},
				primitive.E{Key: "comment", Value: bson.D{
					primitive.E{Key: "$first", Value: "$comment"}}},
				primitive.E{Key: "city", Value: bson.D{
					primitive.E{Key: "$first", Value: "$city"}}},
				primitive.E{Key: "region", Value: bson.D{
					primitive.E{Key: "$first", Value: "$region"}}},
				primitive.E{Key: "district", Value: bson.D{
					primitive.E{Key: "$first", Value: "$district"}}},
				primitive.E{Key: "applicant", Value: bson.D{
					primitive.E{Key: "$first", Value: "$applicant"}}},
				primitive.E{Key: "created_at", Value: bson.D{
					primitive.E{Key: "$first", Value: "$created_at"}}},
				primitive.E{Key: "updated_at", Value: bson.D{
					primitive.E{Key: "$first", Value: "$updated_at"}}},
				primitive.E{Key: "entity_gallery", Value: bson.D{
					primitive.E{Key: "$first", Value: "$entity_gallery"}}},
				primitive.E{Key: "entity", Value: bson.D{
					primitive.E{Key: "$first", Value: "$entity"}}},
				primitive.E{Key: "entity_properties", Value: bson.D{
					primitive.E{Key: "$push", Value: "$entity_properties"}}}}},
		})

	row, err := cr.collection.Aggregate(
		ctx,
		pipeline,
	)
	defer func() {
		_ = row.Close(ctx)
	}()
	if err != nil {
		return nil, err
	}

	if err := row.All(ctx, &appDecode); err != nil {
		return nil, err
	}

	if len(appDecode) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	if err := utils.MarshalUnmarshal(appDecode[0], &response); err != nil {
		return nil, err
	}

	//for i, _ := range response.EntityGallery {
	//	response.EntityGallery[i] = config.MinioUrl + response.EntityGallery[i]
	//}

	return &response, nil
}

func (cr entityDraftRepo) GetAll(ctx context.Context, req *models.GetAllEntityDraftsRequest) ([]*models.GetAllEntityDrafts, uint64, error) {
	var (
		response         []*models.GetAllEntityDrafts
		entitiesDrafts   []*models.GetAllEntityDrafts
		entityCollection = config.EntityCollection
		//propertyCollection = config.PropertyCollection
		skip                  = (req.Page - 1) * req.Limit
		pipeline, filter, err = filterDraft(req)
	)

	count, err := cr.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	fmt.Println(pipeline)
	fmt.Println(filter)
	pipeline = append(pipeline,
		bson.D{primitive.E{Key: "$sort", Value: bson.D{primitive.E{Key: "created_at", Value: -1}}}},
		bson.D{primitive.E{Key: "$skip", Value: skip}},
		bson.D{primitive.E{Key: "$limit", Value: req.Limit}},
		bson.D{
			primitive.E{Key: "$unwind", Value: bson.D{
				primitive.E{Key: "path", Value: "$status"}, {
					Key: "preserveNullAndEmptyArrays", Value: false}}}},
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: entityCollection},
				primitive.E{Key: "localField", Value: "entity_id"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "entity"}}}},
		bson.D{
			primitive.E{Key: "$unwind", Value: bson.D{
				primitive.E{Key: "path", Value: "$entity"}, {
					Key: "preserveNullAndEmptyArrays", Value: true}}}})
	rows, err := cr.collection.Aggregate(
		ctx,
		pipeline)
	defer func() {
		rows.Close(ctx)
	}()

	if err != nil {
		return nil, 0, err
	}
	if err = rows.All(ctx, &entitiesDrafts); err != nil {
		return nil, 0, err
	}
	fmt.Println(&entitiesDrafts)
	byte, err := json.Marshal(&entitiesDrafts)
	if err != nil {
		return nil, 0, err
	}
	if err := json.Unmarshal(byte, &response); err != nil {
		return nil, 0, err
	}
	return response, uint64(count), nil
}

func (cr entityDraftRepo) Delete(ctx context.Context, id string) error {
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		}}

	filter := bson.M{"_id": bson.M{"$eq": id}}
	_, err := cr.collection.UpdateOne(
		ctx,
		filter,
		update,
	)

	if err != nil {
		return err
	}

	return nil
}

func (cr entityDraftRepo) DeleteFromDB(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	result, err := cr.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil || result.DeletedCount == 0 {
		return err
	}
	return nil
}

func (cr entityDraftRepo) UpdateEntityDraftStatus(ctx context.Context, entityDraftID, status string) error {
	entityDraftObjectID, err := primitive.ObjectIDFromHex(entityDraftID)
	if err != nil {
		return err
	}

	statusObjectID, err := primitive.ObjectIDFromHex(status)
	if err != nil {
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"status_id": statusObjectID,
		},
	}
	filter := bson.M{"_id": bson.M{"$eq": entityDraftObjectID}}
	_, err = cr.collection.UpdateOne(
		ctx,
		filter,
		update)

	return err
}

func filterDraft(req *models.GetAllEntityDraftsRequest) (pipeline mongo.Pipeline, filter bson.D, err error) {
	if req.EntityDraftNumber != "" {
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$match", Value: bson.D{
			primitive.E{Key: "entity_draft_number", Value: bson.D{
				primitive.E{Key: "$regex", Value: req.EntityDraftNumber},
				primitive.E{Key: "$options", Value: "im"}}}}}})
		filter = append(filter, primitive.E{Key: "entity_draft_number", Value: bson.D{
			primitive.E{Key: "$regex", Value: req.EntityDraftNumber},
			primitive.E{Key: "$options", Value: "im"}}})

	}
	if req.CityID != "" {
		filter = append(filter, primitive.E{Key: "city.id", Value: req.CityID})
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$match", Value: bson.D{primitive.E{Key: "city.id", Value: req.CityID}}}})

	}
	if req.RegionID != "" {
		filter = append(filter, primitive.E{Key: "region.id", Value: req.RegionID})
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$match", Value: bson.D{primitive.E{Key: "region.id", Value: req.RegionID}}}})

	}
	if req.Status != "" {
		statusObjectID, err := primitive.ObjectIDFromHex(req.Status)
		if err != nil {
			return nil, nil, err
		}
		filter = append(filter, bson.E{Key: "status_id", Value: statusObjectID})
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$match", Value: bson.D{primitive.E{Key: "status_id", Value: statusObjectID}}}})
	}
	// if req. != "" {
	// 	userObjectID, err := primitive.ObjectIDFromHex(req.)
	// 	if err != nil {
	// 		return nil, nil, err
	// 	}
	// 	filter = append(filter, bson.E{Key: "applicant.user_id", Value: userObjectID})
	// 	pipeline = append(pipeline, bson.D{primitive.E{Key: "$match", Value: bson.D{primitive.E{Key: "applicant.user_id", Value: userObjectID}}}})
	// }

	if req.EntityDraftNumber != "" {
		filter = append(filter, bson.E{Key: "entity_draft_soato", Value: bson.D{primitive.E{Key: "$regex", Value: req.EntityDraftNumber}}})
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$match", Value: bson.D{primitive.E{Key: "entity_draft_soato", Value: bson.D{primitive.E{Key: "$regex", Value: req.EntityDraftNumber}}}}}})
	}

	// if req.FromDate != "" {
	// 	if req.ToDate != "" {
	// 		fromDate, err := time.Parse(config.TimeLayout, req.FromDate)
	// 		if err != nil {
	// 			return nil, nil, err
	// 		}
	// 		toDate, err := time.Parse(config.TimeLayout, req.ToDate)
	// 		if err != nil {
	// 			return nil, nil, err
	// 		}
	// 		filter = append(filter, primitive.E{Key: "created_at", Value: bson.D{
	// 			primitive.E{Key: "$gte", Value: fromDate},
	// 			primitive.E{Key: "$lte", Value: toDate},
	// 		}})
	// 		pipeline = append(pipeline, bson.D{primitive.E{Key: "$match", Value: bson.D{
	// 			primitive.E{Key: "created_at", Value: bson.M{"$gte": fromDate}},
	// 			primitive.E{Key: "created_at", Value: bson.M{"$lte": toDate}},
	// 		}}})
	// 	} else {
	// 		fromDate, err := time.Parse(config.TimeLayout, req.FromDate)
	// 		if err != nil {
	// 			return nil, nil, err
	// 		}
	// 		toDate := fromDate.AddDate(0, 0, 1)
	// 		filter = append(filter, primitive.E{Key: "created_at", Value: bson.D{
	// 			primitive.E{Key: "$gte", Value: fromDate},
	// 			primitive.E{Key: "$lte", Value: toDate},
	// 		}})
	// 		pipeline = append(pipeline, bson.D{primitive.E{Key: "$match", Value: bson.D{
	// 			primitive.E{Key: "created_at", Value: bson.M{"$gte": fromDate}},
	// 			primitive.E{Key: "created_at", Value: bson.M{"$lte": toDate}},
	// 		}}})
	// 	}
	// }
	return
}
