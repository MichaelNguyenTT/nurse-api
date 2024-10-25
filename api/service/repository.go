package service

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) List() (Services, error) {
	services := make([]*Service, 0)
	if err := r.db.Find(&services).Error; err != nil {
		return nil, err
	}

	return services, nil
}

func (r *Repository) Create(s *Service) (*Service, error) {
	if err := r.db.Create(s).Error; err != nil {
		return nil, err
	}

	return s, nil
}

func (r *Repository) Read(id uuid.UUID) (*Service, error) {
	service := &Service{}
	if err := r.db.Where("id = ?", id).First(&service).Error; err != nil {
		return nil, err
	}

	return service, nil
}

func (r *Repository) Update(s *Service) (int64, error) {
	response := r.db.Model(&Service{}).
		Select("Name", "Category", "Notes", "UpdatedAt").
		Where("id = ?", s.ID).
		Updates(s)

	return response.RowsAffected, response.Error
}

func (r *Repository) Delete(id uuid.UUID) (int64, error) {
	response := r.db.Where("id = ?", id).Delete(&Service{})

	return response.RowsAffected, response.Error

}
