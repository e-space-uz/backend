package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Applicant struct {
	ID                 string             `json:"id" bson:"_id"`
	FirstName          string             `json:"first_name" bson:"first_name"`
	LastName           string             `json:"last_name" bson:"last_name"`
	Gender             string             `json:"gender" bson:"gender"`
	PhoneNumber        string             `json:"phone_number" bson:"phone_number"`
	UserType           string             `json:"user_type" bson:"user_type" example:"applicant"`
	MiddleName         string             `json:"middle_name" bson:"middle_name"`
	FullName           string             `json:"full_name" bson:"full_name"`
	Login              string             `json:"login" bson:"login"`
	Nationality        string             `json:"nationality" bson:"nationality"`
	PermanentAddress   string             `json:"permanent_address" bson:"permanent_address"`
	PassportNumber     string             `json:"passport_number" bson:"passport_number"`
	PassportIssuePlace string             `json:"passport_issue_place" bson:"passport_issue_place"`
	Pin                string             `json:"pin" bson:"pin"`
	Email              string             `json:"email" bson:"email"`
	Inn                string             `json:"inn" bson:"inn"`
	BirthDate          string             `json:"birth_date" bson:"birth_date"`
	BirthPlace         string             `json:"birth_place" bson:"birth_place"`
	Citizenship        string             `json:"citizenship" bson:"citizenship"`
	ApplicantType      string             `json:"applicant_type" bson:"applicant_type"`
	PassportIssueDate  primitive.DateTime `json:"passport_issue_date" bson:"passport_issue_date"`
	PassportExpiryDate primitive.DateTime `json:"passport_expiry_date" bson:"passport_expiry_date"`
	CreatedAt          primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt          primitive.DateTime `json:"updated_at" bson:"updated_at"`
}

type CreateUpdateApplicant struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id"`
	FirstName          string             `json:"first_name" bson:"first_name"`
	LastName           string             `json:"last_name" bson:"last_name"`
	Gender             string             `json:"gender" bson:"gender"`
	PhoneNumber        string             `json:"phone_number" bson:"phone_number"`
	UserType           string             `json:"user_type" bson:"user_type" example:"applicant"`
	MiddleName         string             `json:"middle_name" bson:"middle_name"`
	FullName           string             `json:"full_name" bson:"full_name"`
	Login              string             `json:"login" bson:"login"`
	Nationality        string             `json:"nationality" bson:"nationality"`
	PermanentAddress   string             `json:"permanent_address" bson:"permanent_address"`
	PassportNumber     string             `json:"passport_number" bson:"passport_number"`
	PassportIssueDate  time.Time          `json:"passport_issue_date" bson:"passport_issue_date"`
	PassportExpiryDate time.Time          `json:"passport_expiry_date" bson:"passport_expiry_date"`
	PassportIssuePlace string             `json:"passport_issue_place" bson:"passport_issue_place"`
	Pin                string             `json:"pin" bson:"pin"`
	Email              string             `json:"email" bson:"email"`
	Inn                string             `json:"inn" bson:"inn"`
	BirthDate          string             `json:"birth_date" bson:"birth_date"`
	BirthPlace         string             `json:"birth_place" bson:"birth_place"`
	Citizenship        string             `json:"citizenship" bson:"citizenship"`
	ApplicantType      string             `json:"applicant_type" bson:"applicant_type"`
	CreatedAt          time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at" bson:"updated_at"`
}

type CreateUpdateApplicantSwag struct {
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	Gender             string `json:"gender"`
	PhoneNumber        string `json:"phone_number"`
	UserType           string `json:"user_type" example:"applicant"`
	MiddleName         string `json:"middle_name"`
	FullName           string `json:"full_name"`
	Login              string `json:"login" bson:"login"`
	Nationality        string `json:"nationality"`
	PermanentAddress   string `json:"permanent_address"`
	PassportNumber     string `json:"passport_number"`
	PassportIssueDate  string `json:"passport_issue_date"`
	PassportExpiryDate string `json:"passport_expiry_date"`
	PassportIssuePlace string `json:"passport_issue_place"`
	Pin                string `json:"pin"`
	Email              string `json:"email"`
	Inn                string `json:"inn"`
	BirthDate          string `json:"birth_date"`
	BirthPlace         string `json:"birth_place"`
	Citizenship        string `json:"citizenship"`
	ApplicantType      string `json:"applicant_type"`
}

type GetAllApplicantsRequestSwag struct {
	FullName       string `json:"full_name"`
	UserType       string `json:"user_type"`
	PhoneNumber    string `json:"phone_number"`
	PassportNumber string `json:"passport_number"`
	Pinfl          string `json:"pinfl"`
	Page           uint32 `json:"page"`
	Limit          uint32 `json:"limit"`
}

type GetAllApplicantsResponse struct {
	Applicants []*Applicant `json:"applicants"`
	Count      int64        `json:"count"`
}
