package service

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// response to users
type ServiceResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Notes    string `json:"notes"`
}

// used to convert to satisfy db
type ServiceRequest struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Notes    string `json:"notes"`
}

// gorm usage
type Service struct {
	ID        uuid.UUID `gorm:"primarykey"`
	Name      string
	Category  string
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type Services []*Service

func (s *Service) ToClient() *ServiceResponse {
	return &ServiceResponse{
		ID:       s.ID.String(),
		Name:     s.Name,
		Category: s.Category,
		Notes:    s.Notes,
	}
}

func (s Services) GetAll() []*ServiceResponse {
	services := make([]*ServiceResponse, 0, len(s))

	for _, data := range s {
		services = append(services, data.ToClient())
	}

	return services
}

func (v *ServiceRequest) ToDB() *Service {
	return &Service{
		Name:     v.Name,
		Category: v.Category,
		Notes:    v.Notes,
	}
}
