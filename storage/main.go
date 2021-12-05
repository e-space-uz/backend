package storage

import (
	"github.com/e-space-uz/backend/storage/mongodb"
	"github.com/e-space-uz/backend/storage/repo"
	db "go.mongodb.org/mongo-driver/mongo"
)

type StorageI interface {
	Applicant() repo.ApplicantI
	Staff() repo.StaffI
	City() repo.CityI
	Region() repo.RegionI
	District() repo.DistrictI
	Property() repo.PropertyI
	Entity() repo.EntityI
	EntityDraft() repo.EntityDraftI
	EntityFiles() repo.EntityFilesI
	GroupProperty() repo.GroupPropertyI
}

type storageMongo struct {
	applicantRepo     repo.ApplicantI
	staffRepo         repo.StaffI
	propertyRepo      repo.PropertyI
	cityRepo          repo.CityI
	regionRepo        repo.RegionI
	districtRepo      repo.DistrictI
	entityRepo        repo.EntityI
	entityDraftRepo   repo.EntityDraftI
	groupPropertyRepo repo.GroupPropertyI
	entityFilesRepo   repo.EntityFilesI
}

func NewStorageMongo(db *db.Database) StorageI {
	return &storageMongo{
		// applicantRepo:     mongodb.NewApplicantRepo(db),
		// staffRepo:         mongodb.NewStaffRepo(db),
		cityRepo:          mongodb.NewCityRepo(db),
		regionRepo:        mongodb.NewRegionRepo(db),
		districtRepo:      mongodb.NewDistrictRepo(db),
		propertyRepo:      mongodb.NewPropertyRepo(db),
		entityRepo:        mongodb.NewEntityRepo(db),
		groupPropertyRepo: mongodb.NewGroupPropertyRepo(db),
		entityFilesRepo:   mongodb.NewEntityFilesRepo(db),
		entityDraftRepo:   mongodb.NewEntityDraftRepo(db),
	}
}

func (s *storageMongo) City() repo.CityI {
	return s.cityRepo
}
func (s *storageMongo) Applicant() repo.ApplicantI {
	return s.applicantRepo
}
func (s *storageMongo) Region() repo.RegionI {
	return s.regionRepo
}
func (s *storageMongo) District() repo.DistrictI {
	return s.districtRepo
}
func (s *storageMongo) Property() repo.PropertyI {
	return s.propertyRepo
}
func (s *storageMongo) GroupProperty() repo.GroupPropertyI {
	return s.groupPropertyRepo
}
func (s *storageMongo) Entity() repo.EntityI {
	return s.entityRepo
}
func (s *storageMongo) EntityFiles() repo.EntityFilesI {
	return s.entityFilesRepo
}
func (s *storageMongo) EntityDraft() repo.EntityDraftI {
	return s.entityDraftRepo
}

func (s *storageMongo) Staff() repo.StaffI {
	return s.staffRepo
}
