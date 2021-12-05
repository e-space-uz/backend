package mongodb

import (
	"context"
	"encoding/json"
	"time"

	"github.com/e-space-uz/backend/config"
	"github.com/e-space-uz/backend/models"
	"github.com/e-space-uz/backend/pkg/utils"
	"github.com/e-space-uz/backend/storage/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type applicantRepo struct {
	collection *mongo.Collection
}

func NewApplicantRepo(db *mongo.Database) repo.ApplicantI {
	return &applicantRepo{
		collection: db.Collection(config.ApplicantCollection),
	}
}

func (ar *applicantRepo) Create(applicant *models.Applicant) (*models.Applicant, error) {

	passportIssueDate, err := time.Parse("2006-02-01", applicant.PassportIssueDate)
	if err != nil {
		return nil, err
	}
	passportExpiryDate, err := time.Parse("2006-02-01", applicant.PassportExpiryDate)
	if err != nil {
		return nil, err
	}
	createApplicant := &models.CreateUpdateApplicant{
		ID:                 objectID,
		Login:              applicant.Login,
		FirstName:          applicant.FirstName,
		LastName:           applicant.LastName,
		Gender:             applicant.Gender,
		PhoneNumber:        applicant.PhoneNumber,
		UserType:           applicant.UserType,
		MiddleName:         applicant.MiddleName,
		FullName:           applicant.FullName,
		Nationality:        applicant.Nationality,
		PermanentAddress:   applicant.PermanentAddress,
		PassportNumber:     applicant.PassportNumber,
		PassportIssueDate:  passportIssueDate,
		PassportExpiryDate: passportExpiryDate,
		PassportIssuePlace: applicant.PassportIssuePlace,
		Pin:                applicant.Pin,
		Email:              applicant.Email,
		Inn:                applicant.Inn,
		BirthDate:          applicant.BirthDate,
		BirthPlace:         applicant.BirthPlace,
		Citizenship:        applicant.Citizenship,
		ApplicantType:      applicant.ApplicantType,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	_, err = ar.collection.InsertOne(
		context.Background(),
		createApplicant,
	)
	if err != nil {
		return nil, err
	}

	return applicant, nil
}

func (ar *applicantRepo) Get(id string) (*models.Applicant, error) {
	var (
		response  *models.Applicant
		applicant *models.Applicant
	)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if err = ar.collection.FindOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		}).Decode(&applicant); err != nil {
		return nil, err
	}
	byte, err := json.Marshal(&applicant)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(byte, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (ar *applicantRepo) GetAll(req *models.GetAllApplicantsRequest) ([]*models.Applicant, uint32, error) {
	var (
		response   []*models.Applicant
		applicants []*models.Applicant
		filter     = bson.D{}
	)

	if req.FullName != "" {
		filter = append(filter, primitive.E{Key: "full_name", Value: bson.D{primitive.E{Key: "$regex", Value: req.FullName},
			primitive.E{Key: "$options", Value: "im"}}})
	}
	if req.PhoneNumber != "" {
		filter = append(filter, primitive.E{Key: "phone_number", Value: bson.D{primitive.E{Key: "$regex", Value: req.PhoneNumber},
			primitive.E{Key: "$options", Value: "im"}}})
	}

	if req.PassportNumber != "" {
		filter = append(filter, primitive.E{Key: "passport_number", Value: bson.D{primitive.E{Key: "$regex", Value: req.PassportNumber},
			primitive.E{Key: "$options", Value: "im"}}})
	}
	if req.UserType != "" {
		filter = append(filter, bson.E{Key: "user_type", Value: req.UserType})
	}

	if req.Pinfl != "" {
		filter = append(filter, bson.E{Key: "pin", Value: req.Pinfl})
	}

	opts := options.Find()

	skip := (req.Page - 1) * req.Limit
	opts.SetLimit(int64(req.Limit))
	opts.SetSkip(int64(skip))
	opts.SetSort(bson.M{
		"created_at": -1,
	})

	count, err := ar.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	rows, err := ar.collection.Find(
		context.Background(),
		filter,
		opts,
	)
	if err != nil {
		return nil, 0, err
	}

	if err = rows.All(context.Background(), &applicants); err != nil {
		return nil, 0, err
	}

	if err := utils.MarshalUnmarshal(applicants, &response); err != nil {
		return nil, 0, err
	}

	return response, uint32(count), nil
}

func (ar *applicantRepo) Update(applicant *models.Applicant) error {
	objectID, err := primitive.ObjectIDFromHex(applicant.Id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"id":           applicant.Id,
			"first_name":   applicant.FirstName,
			"last_name":    applicant.LastName,
			"gender":       applicant.Gender,
			"phone_number": applicant.PhoneNumber,
			"user_type":    applicant.UserType,
		}}

	filter := bson.M{"_id": bson.M{"$eq": objectID}}
	_, err = ar.collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return err
	}

	return nil
}

func (ar *applicantRepo) Exists(id string) (bool, error) {
	//objectID, err := primitive.ObjectIDFromHex(id)
	//if err != nil {
	//	return false, err
	//}
	filter := bson.M{"login": id}
	count, err := ar.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (ar *applicantRepo) GetByUserID(id string) (*models.Applicant, error) {
	var (
		response  *models.Applicant
		applicant *models.Applicant
	)

	if err := ar.collection.FindOne(
		context.Background(),
		bson.M{
			"login": id,
		}).Decode(&applicant); err != nil {
		return nil, err
	}
	byte, err := json.Marshal(&applicant)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(byte, &response); err != nil {
		return nil, err
	}

	return response, nil
}
