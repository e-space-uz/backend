package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/e-space-uz/backend/config"
	"github.com/e-space-uz/backend/models"
	"github.com/e-space-uz/backend/storage/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type staffRepo struct {
	collection *mongo.Collection
}

func NewStaffRepo(db *mongo.Database) repo.StaffI {
	return &staffRepo{
		collection: db.Collection(config.StaffCollection),
	}
}

func (sr *staffRepo) Create(staff *models.CreateUpdateStaff) (string, error) {

	objectID, err := primitive.ObjectIDFromHex(staff.Id)
	if err != nil {
		return "", err
	}

	organizationID, err := primitive.ObjectIDFromHex(staff.OrganizationId)
	if err != nil {
		return "", err
	}

	passportIssueDate, err := time.Parse(config.TimeLayout, staff.PassportIssueDate)
	if err != nil {
		return "", err
	}

	roleID, err := primitive.ObjectIDFromHex(staff.RoleId)
	if err != nil {
		return "", err
	}

	createStaff := &models.CreateUpdateStaff{
		ID:                 objectID,
		RoleID:             roleID,
		OrganizationID:     organizationID,
		ExternalID:         staff.ExternalId,
		FirstName:          staff.FirstName,
		LastName:           staff.LastName,
		MiddleName:         staff.MiddleName,
		UniqueName:         staff.UniqueName,
		PhoneNumber:        staff.PhoneNumber,
		UserType:           staff.UserType,
		Pinfl:              staff.Pinfl,
		Address:            staff.Address,
		Inn:                staff.Inn,
		Login:              staff.Login,
		Password:           staff.Password,
		LastLogin:          staff.LastLogin,
		ExtraInfo:          staff.ExtraInfo,
		PassportIssueDate:  passportIssueDate,
		Policy:             staff.Policy,
		PassportNumber:     staff.PassportNumber,
		PassportIssuePlace: staff.PassportIssuePlace,
		Email:              staff.Email,
		Soato:              staff.Soato,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		Verified:           true,
	}

	if staff.City != nil {
		createStaff.City = &models.City{
			ID:     staff.City.Id,
			Name:   staff.City.Name,
			RuName: staff.City.RuName,
			Soato:  staff.City.Soato,
			Code:   staff.City.Code,
		}
	}

	if staff.Region != nil {
		createStaff.Region = &models.Region{
			ID:     staff.Region.Id,
			Name:   staff.Region.Name,
			RuName: staff.Region.RuName,
			Soato:  staff.City.Soato,
			Code:   staff.Region.Code,
		}
	}

	_, err = sr.collection.InsertOne(
		context.Background(),
		createStaff,
	)

	return createStaff.ID.Hex(), err
}
func (sr *staffRepo) GetAll(req *models.GetAllStaffsRequest) ([]*models.Staff, uint32, error) {
	var (
		filter   = bson.D{}
		skip     = (req.Page - 1) * req.Limit
		pipeline = mongo.Pipeline{}
	)

	fmt.Println(req)
	if req.PhoneNumber != "" {
		filter = append(filter, bson.E{Key: "phone_number", Value: bson.D{primitive.E{Key: "$regex", Value: req.PhoneNumber}}})
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$match", Value: bson.D{primitive.E{Key: "phone_number", Value: bson.D{primitive.E{Key: "$regex", Value: req.PhoneNumber}}}}}})
	}
	if req.Soato != "" {
		filter = append(filter, bson.E{Key: "soato", Value: bson.D{primitive.E{Key: "$regex", Value: req.Soato}}})
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$match", Value: bson.D{primitive.E{Key: "soato", Value: bson.D{primitive.E{Key: "$regex", Value: req.Soato}}}}}})
	}
	if req.RoleId != "" {
		roleObjectID, err := primitive.ObjectIDFromHex(req.RoleId)
		if err != nil {
			return nil, 0, err
		}
		filter = append(filter, primitive.E{Key: "role_id", Value: roleObjectID})
		pipeline = append(pipeline, bson.D{
			primitive.E{Key: "$match", Value: bson.D{
				primitive.E{Key: "role_id", Value: roleObjectID},
			}},
		})
	}
	if req.Status {
		filter = append(filter, bson.E{Key: "status", Value: req.Status})
		pipeline = append(pipeline, bson.D{
			primitive.E{Key: "$match", Value: bson.D{
				primitive.E{Key: "status", Value: req.Status},
			}},
		})
	}
	if req.OrganizationId != "" {
		organizationObjectID, err := primitive.ObjectIDFromHex(req.OrganizationId)
		if err != nil {
			return nil, 0, err
		}
		filter = append(filter, primitive.E{Key: "organization_id", Value: organizationObjectID})
		pipeline = append(pipeline, bson.D{
			primitive.E{Key: "$match", Value: bson.D{
				primitive.E{Key: "organization_id", Value: organizationObjectID},
			}},
		})
	}
	if req.SearchString != "" {
		pipeline = append(pipeline, bson.D{
			primitive.E{Key: "$match", Value: bson.D{
				primitive.E{Key: "$or", Value: []bson.M{
					{"first_name": bson.M{"$regex": req.SearchString}},
					{"unique_name": bson.M{"$regex": req.SearchString}},
					{"phone_number": bson.M{"$regex": req.SearchString}},
				}},
			},
			},
		})
		filter = append(filter,
			primitive.E{Key: "$or", Value: []bson.M{
				{"first_name": bson.M{"$regex": req.SearchString}},
				{"unique_name": bson.M{"$regex": req.SearchString}},
				{"phone_number": bson.M{"$regex": req.SearchString}},
			}})
	}

	pipeline = append(pipeline,
		bson.D{
			primitive.E{Key: "$unwind", Value: bson.D{
				primitive.E{Key: "path", Value: "$role"}, {
					Key: "preserveNullAndEmptyArrays", Value: false}}}},

		bson.D{
			primitive.E{Key: "$unwind", Value: bson.D{
				primitive.E{Key: "path", Value: "$organization"}, {
					Key: "preserveNullAndEmptyArrays", Value: false}}}},

		bson.D{primitive.E{Key: "$sort", Value: bson.D{primitive.E{Key: "created_at", Value: -1}}}},
		bson.D{primitive.E{Key: "$skip", Value: skip}},
		bson.D{primitive.E{Key: "$limit", Value: req.Limit}},
	)

	count, err := sr.collection.CountDocuments(
		context.Background(),
		filter,
	)
	if err != nil {
		return nil, 0, err
	}

	rows, err := sr.collection.Aggregate(
		context.Background(),
		pipeline,
	)
	if err != nil {
		return nil, 0, err
	}

	return rows, uint32(count), nil
}

func (sr *staffRepo) GetCount(soato string, organizationID string) (int32, error) {
	var (
		filter = bson.D{}
	)

	organizationObjectID, err := primitive.ObjectIDFromHex(organizationID)
	if err != nil {
		return 0, err
	}
	filter = append(filter, primitive.E{Key: "organization_id", Value: organizationObjectID}, primitive.E{Key: "soato", Value: soato})

	count, err := sr.collection.CountDocuments(
		context.Background(),
		filter,
	)
	return int32(count), err
}

func (sr *staffRepo) UpdatePassword(ctx context.Context, oldPassword, newPassword string, userID string) error {
	var (
		update = primitive.M{}
	)
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	update = bson.M{
		"$set": bson.M{
			"password": newPassword,
			"verified": true,
		},
	}

	filter := bson.M{"_id": bson.M{"$eq": objectID}}
	_, err = sr.collection.UpdateOne(ctx, filter, update)

	return err
}

func (sr *staffRepo) SetRoleID(ctx context.Context, roleID, staffID string) error {
	var (
		update = primitive.M{}
	)
	objectID, err := primitive.ObjectIDFromHex(staffID)
	if err != nil {
		return err
	}
	roleObjectID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		return err
	}
	update = bson.M{
		"$set": bson.M{
			"role_id":  roleObjectID,
			"verified": false,
		},
	}

	filter := bson.M{"_id": bson.M{"$eq": objectID}}
	_, err = sr.collection.UpdateOne(ctx, filter, update)

	return err
}
func (sr *staffRepo) SetStaffSoato(ctx context.Context, staffID, soato string) error {
	var (
		update = primitive.M{}
	)
	objectID, err := primitive.ObjectIDFromHex(staffID)
	if err != nil {
		return err
	}
	update = bson.M{
		"$set": bson.M{
			"soato": soato,
		},
	}

	filter := bson.M{"_id": bson.M{"$eq": objectID}}
	_, err = sr.collection.UpdateOne(ctx, filter, update)

	return err
}

// func (sr *staffRepo) ImportStaff(staff *models.CreateUpdateStaff) (string, error) {

// 	objectID, err := primitive.ObjectIDFromHex(staff.Id)
// 	if err != nil {
// 		return "", err
// 	}

// 	roleID, err := primitive.ObjectIDFromHex(staff.RoleId)
// 	if err != nil {
// 		return "", err
// 	}
// 	organizationID, err := primitive.ObjectIDFromHex(staff.OrganizationID)
// 	if err != nil {
// 		return "", err
// 	}
// 	createStaff := &models.CreateUpdateStaffImport{
// 		ID:                 objectID,
// 		RoleID:             roleID,
// 		BranchID:           staff.BranchId,
// 		OrganizationID:     organizationID,
// 		ExternalID:         staff.ExternalId,
// 		FirstName:          staff.FirstName,
// 		LastName:           staff.LastName,
// 		MiddleName:         staff.MiddleName,
// 		UniqueName:         staff.UniqueName,
// 		PhoneNumber:        staff.PhoneNumber,
// 		UserType:           staff.UserType,
// 		Pinfl:              staff.Pinfl,
// 		Address:            staff.Address,
// 		Inn:                staff.Inn,
// 		Login:              staff.Login,
// 		Password:           staff.Password,
// 		LastLogin:          staff.LastLogin,
// 		ExtraInfo:          staff.ExtraInfo,
// 		PassportIssueDate:  staff.PassportIssueDate,
// 		Policy:             staff.Policy,
// 		PassportNumber:     staff.PassportNumber,
// 		PassportIssuePlace: staff.PassportIssuePlace,
// 		Email:              staff.Email,
// 		CreatedAt:          time.Now(),
// 		UpdatedAt:          time.Now(),
// 	}

// 	resp, err := sr.collection.InsertOne(
// 		context.Background(),
// 		createStaff,
// 	)

// 	if err != nil {
// 		return "", err
// 	}

// 	return resp.InsertedID.(primitive.ObjectID).Hex(), nil
// }
