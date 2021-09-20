package doctor_rating

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/matijapetrovic/clinichub/rating-service/pkg/httpclient"
	"github.com/matijapetrovic/clinichub/rating-service/pkg/log"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/matijapetrovic/clinichub/rating-service/internal/auth"
	"github.com/matijapetrovic/clinichub/rating-service/internal/entity"
)

type Service interface {
	GetAvaialableRatings(request *http.Request) ([]Clinic, error)
	RateClinic(ctx context.Context, clinicId string, request RateClinicRequest) (entity.ClinicRating, error)
	GetClinicRating(ctx context.Context, clinicId string) (entity.AverageRating, error)
}

type RateClinicRequest struct {
	Rating float32 `json:"rating"`
}

func (m RateClinicRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Rating, validation.Min(0.0), validation.Max(5.0)),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

type Appointment struct {
	Id                string    `json:"id"`
	ClinicId          string    `json:"clinicId"`
	DoctorId          string    `json:"doctorId"`
	PatientId         string    `json:"patientId"`
	AppointmentTypeId string    `json:"appointmentTypeId"`
	Price             uint      `json:"price"`
	Time              time.Time `json:"time"`
}

type Clinic struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (s service) GetAvaialableRatings(request *http.Request) ([]Clinic, error) {
	ctx := request.Context()
	appointments, err := getPatientAppointments(request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}
	clinicsToRate := make(map[string]bool)

	for _, appointment := range appointments {
		_, err := s.repo.GetRating(ctx, appointment.PatientId, appointment.ClinicId)
		if err == sql.ErrNoRows {
			clinicsToRate[appointment.ClinicId] = true
		}
	}

	result := make([]Clinic, 0)

	for clinicId := range clinicsToRate {
		clinic, err := getClinic(request.Header.Get("Authorization"), clinicId)
		if err != nil {
			return nil, err
		}
		result = append(result, clinic)
	}

	return result, nil
}

func getPatientAppointments(token string) ([]Appointment, error) {
	url, err := url.Parse("http://localhost:8083/v1/appointments")
	if err != nil {
		return nil, err
	}
	endDate := strings.Split(time.Now().String(), " ")[0]
	fmt.Printf(endDate)
	queryParamMap := map[string]string{
		"endDate": endDate,
	}

	client := httpclient.NewJsonClient(
		"GET",
		url,
		func(ctx context.Context, r *http.Response) (interface{}, error) {
			var appointment []Appointment
			err := json.NewDecoder(r.Body).Decode(&appointment)
			if err != nil {
				return nil, err
			}
			return appointment, nil
		},
		token,
		httpclient.QueryParamBeforeFunc(queryParamMap),
	)

	res, err := client.Endpoint()(context.Background(), struct{}{})
	if err != nil {
		return nil, err
	}
	appointments, ok := res.([]Appointment)
	if !ok {
		return nil, errors.New("unexpected error")
	}

	return appointments, nil
}

func getClinic(token string, clinicId string) (Clinic, error) {
	url, err := url.Parse("http://localhost:8081/v1/clinics/" + clinicId)
	if err != nil {
		return Clinic{}, err
	}

	client := httpclient.NewJsonClient(
		"GET",
		url,
		func(ctx context.Context, r *http.Response) (interface{}, error) {
			var clinic Clinic
			err := json.NewDecoder(r.Body).Decode(&clinic)
			if err != nil {
				return nil, err
			}
			return clinic, nil
		},
		token,
		nil,
	)

	res, err := client.Endpoint()(context.Background(), struct{}{})
	if err != nil {
		return Clinic{}, err
	}
	clincic, ok := res.(Clinic)
	if !ok {
		return Clinic{}, errors.New("unexpected error")
	}

	return clincic, nil
}

func (s service) GetClinicRating(ctx context.Context, clinicId string) (entity.AverageRating, error) {
	rating, err := s.repo.GetClinicRating(ctx, clinicId)
	if err != nil {
		return entity.AverageRating{}, err
	}
	return rating, nil
}

func (s service) RateClinic(ctx context.Context, clinicId string, req RateClinicRequest) (entity.ClinicRating, error) {
	if err := req.Validate(); err != nil {
		return entity.ClinicRating{}, err
	}
	user := auth.CurrentUser(ctx)
	id := entity.GenerateID()
	err := s.repo.RateClinic(ctx, entity.ClinicRating{
		ID:        id,
		ClinicId:  clinicId,
		PatientId: user.GetID(),
		Rating:    req.Rating,
	})
	if err != nil {
		return entity.ClinicRating{}, err
	}
	return s.repo.GetById(ctx, id)
}
