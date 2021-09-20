package clinic

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	appointment_type "github.com/matijapetrovic/clinichub/clinic-service/internal/appointment-type"
	"github.com/matijapetrovic/clinichub/clinic-service/internal/entity"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/httpclient"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/log"
)

type Service interface {
	GetById(request *http.Request, id string) (entity.Clinic, error)
	Count(ctx context.Context) (int, error)
	Query(request *http.Request, req QueryClinicsRequest) ([]entity.Clinic, error)
	Create(ctx context.Context, req CreateClinicRequest) (entity.Clinic, error)
	Update(ctx context.Context, clinicId string, req UpdateClinicRequest) (entity.Clinic, error)
	AddAppointmentTypePrice(ctx context.Context, clinicId string, req AddAppointmentTypePriceRequest) (entity.AppointmentTypePrice, error)
	GetAppointmentTypePrices(ctx context.Context, clinicId string) ([]entity.AppointmentTypePrice, error)
	UpdateAppointmentTypePrice(ctx context.Context, clinicId string, req UpdateAppointmentTypePriceRequest) (entity.AppointmentTypePrice, error)
}

type QueryClinicsRequest struct {
	AppointmentTypeId string `json:"appointmentTypeId"`
	Date              string `json:"date"`
	Limit             int    `json:"limit"`
	Offset            int    `json:"offset"`
}

func (m QueryClinicsRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.AppointmentTypeId, validation.Length(36, 36)),
		validation.Field(&m.Date, validation.Date("1970-12-30"), validation.Min(time.Now())),
		validation.Field(&m.Limit, validation.Min(1)),
		validation.Field(&m.Offset, validation.Min(0)),
	)
}

type CreateClinicRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Address     entity.Address `json:"address"`
}

func (m CreateClinicRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(1, 50)),
		validation.Field(&m.Description, validation.Required, validation.Length(1, 256)),
	)
}

type UpdateClinicRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Address     entity.Address `json:"address"`
}

func (m UpdateClinicRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(1, 50)),
		validation.Field(&m.Description, validation.Required, validation.Length(1, 256)),
	)
}

type AddAppointmentTypePriceRequest struct {
	AppointmentTypeId string `json:"appointmentTypeId"`
	Price             uint   `json:"price"`
}

func (m AddAppointmentTypePriceRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.AppointmentTypeId, validation.Required, validation.Length(36, 36)),
		validation.Field(&m.Price, validation.Required, validation.Min(1)),
	)
}

type UpdateAppointmentTypePriceRequest struct {
	AppointmentTypeId string `json:"appointmentTypeId"`
	Price             uint   `json:"price"`
}

func (m UpdateAppointmentTypePriceRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.AppointmentTypeId, validation.Required, validation.Length(36, 36)),
		validation.Field(&m.Price, validation.Required, validation.Min(1)),
	)
}

type service struct {
	repo                Repository
	appointmentTypeRepo appointment_type.Repository
	logger              log.Logger
}

func NewService(repo Repository, appointmentTypeRepo appointment_type.Repository, logger log.Logger) Service {
	return service{repo, appointmentTypeRepo, logger}
}

func (s service) GetById(request *http.Request, id string) (entity.Clinic, error) {
	ctx := request.Context()
	clinic, err := s.repo.GetById(ctx, id)
	if err != nil {
		return entity.Clinic{}, err
	}

	rating, err := getClinicRating(clinic.Id, request.Header.Get("Authorization"))
	if err != nil {
		return entity.Clinic{}, err
	}
	clinic.Rating = rating

	return clinic, nil
}

func (s service) Create(ctx context.Context, req CreateClinicRequest) (entity.Clinic, error) {
	if err := req.Validate(); err != nil {
		return entity.Clinic{}, err
	}

	id := entity.GenerateID()
	err := s.repo.Create(ctx, entity.Clinic{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
	})

	if err != nil {
		return entity.Clinic{}, err
	}

	return s.repo.GetById(ctx, id)
}

func (s service) Update(ctx context.Context, clinicId string, req UpdateClinicRequest) (entity.Clinic, error) {
	if err := req.Validate(); err != nil {
		return entity.Clinic{}, err
	}

	clinic, err := s.repo.GetById(ctx, clinicId)
	if err != nil {
		return entity.Clinic{}, err
	}

	clinic.Name = req.Name
	clinic.Description = req.Description
	clinic.Address = req.Address

	err = s.repo.Update(ctx, clinic)
	if err != nil {
		return entity.Clinic{}, err
	}

	return clinic, nil
}

func (s service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

func (s service) Query(request *http.Request, req QueryClinicsRequest) ([]entity.Clinic, error) {
	ctx := request.Context()
	if req.AppointmentTypeId != "" && req.Date == "" || req.AppointmentTypeId == "" && req.Date != "" {
		return nil, errors.New("bad request")
	}
	if req.AppointmentTypeId == "" && req.Date == "" {
		clinics, err := s.repo.GetPaged(ctx, req.Offset, req.Limit)
		if err != nil {
			return nil, err
		}

		return clinics, nil
	} else {
		clinicIds, err := s.repo.GetIdsByHasPrice(ctx, req.AppointmentTypeId)
		if err != nil {
			return nil, err
		}

		clinics, err := s.repo.GetByIdList(ctx, clinicIds)
		if err != nil {
			return nil, err
		}

		for idx, clinic := range clinics {
			rating, err := getClinicRating(clinic.Id, request.Header.Get("Authorization"))
			if err != nil {
				return nil, err
			}
			clinic.Rating = rating

			price, err := s.repo.GetAppointmentTypePrice(ctx, clinic.Id, req.AppointmentTypeId)
			if err != nil {
				return nil, err
			}
			clinic.Price = price.Price
			clinics[idx] = clinic
		}

		return clinics, err
	}
}

func getClinicRating(clinicId string, token string) (entity.Rating, error) {
	url, err := url.Parse("http://localhost:8082/v1/clinics/" + clinicId + "/average-rating")
	if err != nil {
		return entity.Rating{}, err
	}

	client := httpclient.NewJsonClient(
		"GET",
		url,
		func(ctx context.Context, r *http.Response) (interface{}, error) {
			var rating entity.Rating
			err := json.NewDecoder(r.Body).Decode(&rating)
			if err != nil {
				return nil, err
			}
			return rating, nil
		},
		token,
		nil,
	)

	res, err := client.Endpoint()(context.Background(), struct{}{})
	if err != nil {
		return entity.Rating{}, err
	}
	rating, ok := res.(entity.Rating)
	if !ok {
		return entity.Rating{}, errors.New("unexpected error")
	}

	return rating, nil
}

func (s service) AddAppointmentTypePrice(ctx context.Context, clinicId string, req AddAppointmentTypePriceRequest) (entity.AppointmentTypePrice, error) {
	clinic, err := s.repo.GetById(ctx, clinicId)
	if err != nil {
		return entity.AppointmentTypePrice{}, err
	}

	appointmentType, err := s.appointmentTypeRepo.GetById(ctx, req.AppointmentTypeId)
	if err != nil {
		return entity.AppointmentTypePrice{}, err
	}

	price := entity.AppointmentTypePrice{
		ClinicId:          clinic.Id,
		AppointmentTypeId: appointmentType.Id,
		Price:             req.Price,
	}

	err = s.repo.AddAppointmentTypePrice(ctx, price)

	return price, err
}

func (s service) UpdateAppointmentTypePrice(ctx context.Context, clinicId string, req UpdateAppointmentTypePriceRequest) (entity.AppointmentTypePrice, error) {
	clinic, err := s.repo.GetById(ctx, clinicId)
	if err != nil {
		return entity.AppointmentTypePrice{}, err
	}

	appointmentType, err := s.appointmentTypeRepo.GetById(ctx, req.AppointmentTypeId)
	if err != nil {
		return entity.AppointmentTypePrice{}, err
	}

	price := entity.AppointmentTypePrice{
		ClinicId:          clinic.Id,
		AppointmentTypeId: appointmentType.Id,
		Price:             req.Price,
	}

	err = s.repo.UpdateAppointmentTypePrice(ctx, price)

	return price, err
}

func (s service) GetAppointmentTypePrices(ctx context.Context, clinicId string) ([]entity.AppointmentTypePrice, error) {
	appointmentTypePrices, err := s.repo.GetAppointmentTypePrices(ctx, clinicId)
	if err != nil {
		return nil, err
	}

	for i, price := range appointmentTypePrices {
		appointmentType, err := s.appointmentTypeRepo.GetById(ctx, price.AppointmentTypeId)
		if err != nil {
			return nil, err
		}
		price.AppointmentType = appointmentType
		appointmentTypePrices[i] = price
	}

	return appointmentTypePrices, nil
}
