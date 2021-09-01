package appointment_type

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/matijapetrovic/clinichub/clinic-service/internal/entity"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/log"
)

type Service interface {
	GetById(ctx context.Context, id string) (entity.AppointmentType, error)
	GetAll(ctx context.Context) ([]entity.AppointmentType, error)
	Create(ctx context.Context, req CreateAppointmentTypeRequest) (entity.AppointmentType, error)
	Update(ctx context.Context, id string, req UpdateAppointmentTypeRequest) (entity.AppointmentType, error)
}

type CreateAppointmentTypeRequest struct {
	Name string `json:"name"`
}

func (m CreateAppointmentTypeRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(1, 50)),
	)
}

type UpdateAppointmentTypeRequest struct {
	Name string `json:"name"`
}

func (m UpdateAppointmentTypeRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(1, 50)),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

func (s service) GetById(ctx context.Context, id string) (entity.AppointmentType, error) {
	appointmentType, err := s.repo.GetById(ctx, id)
	if err != nil {
		return entity.AppointmentType{}, err
	}

	return appointmentType, nil
}

func (s service) Create(ctx context.Context, req CreateAppointmentTypeRequest) (entity.AppointmentType, error) {
	if err := req.Validate(); err != nil {
		return entity.AppointmentType{}, err
	}

	id := entity.GenerateID()
	err := s.repo.Create(ctx, entity.AppointmentType{
		Id:   id,
		Name: req.Name,
	})

	if err != nil {
		return entity.AppointmentType{}, err
	}

	return s.repo.GetById(ctx, id)
}

func (s service) Update(ctx context.Context, appointmentTypeId string, req UpdateAppointmentTypeRequest) (entity.AppointmentType, error) {
	if err := req.Validate(); err != nil {
		return entity.AppointmentType{}, err
	}

	appointmentType, err := s.repo.GetById(ctx, appointmentTypeId)
	if err != nil {
		return entity.AppointmentType{}, err
	}

	appointmentType.Name = req.Name

	err = s.repo.Update(ctx, appointmentType)
	if err != nil {
		return entity.AppointmentType{}, err
	}

	return appointmentType, nil
}

func (s service) GetAll(ctx context.Context) ([]entity.AppointmentType, error) {
	appointmentTypes, err := s.repo.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	return appointmentTypes, nil
}
