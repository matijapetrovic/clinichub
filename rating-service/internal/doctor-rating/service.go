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
	GetAvaialableRatings(request *http.Request) ([]Doctor, error)
	RateDoctor(ctx context.Context, doctorId string, request RateDoctorRequest) (entity.DoctorRating, error)
	GetDoctorRating(ctx context.Context, doctorID string) (entity.AverageRating, error)
}

type RateDoctorRequest struct {
	Rating float32 `json:"rating"`
}

func (m RateDoctorRequest) Validate() error {
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

type Doctor struct {
	Id        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Name      string `json:"name"`
}

func (s service) GetAvaialableRatings(request *http.Request) ([]Doctor, error) {
	ctx := request.Context()
	appointments, err := getPatientAppointments(request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}
	doctorsToRate := make(map[string]bool)

	for _, appointment := range appointments {
		_, err := s.repo.GetRating(ctx, appointment.PatientId, appointment.DoctorId)
		if err == sql.ErrNoRows {
			doctorsToRate[appointment.DoctorId] = true
		}
	}

	result := make([]Doctor, 0)

	for doctorId := range doctorsToRate {
		doctor, err := getDoctor(request.Header.Get("Authorization"), doctorId)
		if err != nil {
			return nil, err
		}
		doctor.Name = fmt.Sprintf("%s %s", doctor.FirstName, doctor.LastName)
		result = append(result, doctor)
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

func getDoctor(token string, doctorId string) (Doctor, error) {
	url, err := url.Parse("http://localhost:8081/v1/doctors/" + doctorId)
	if err != nil {
		return Doctor{}, err
	}

	client := httpclient.NewJsonClient(
		"GET",
		url,
		func(ctx context.Context, r *http.Response) (interface{}, error) {
			var doctor Doctor
			err := json.NewDecoder(r.Body).Decode(&doctor)
			if err != nil {
				return nil, err
			}
			return doctor, nil
		},
		token,
		nil,
	)

	res, err := client.Endpoint()(context.Background(), struct{}{})
	if err != nil {
		return Doctor{}, err
	}
	doctor, ok := res.(Doctor)
	if !ok {
		return Doctor{}, errors.New("unexpected error")
	}

	return doctor, nil
}

func (s service) GetDoctorRating(ctx context.Context, doctorID string) (entity.AverageRating, error) {
	rating, err := s.repo.GetDoctorRating(ctx, doctorID)
	if err != nil {
		return entity.AverageRating{}, err
	}
	return rating, nil
}

func (s service) RateDoctor(ctx context.Context, doctorId string, req RateDoctorRequest) (entity.DoctorRating, error) {
	if err := req.Validate(); err != nil {
		return entity.DoctorRating{}, err
	}
	user := auth.CurrentUser(ctx)
	id := entity.GenerateID()
	err := s.repo.RateDoctor(ctx, entity.DoctorRating{
		ID:        id,
		DoctorId:  doctorId,
		PatientId: user.GetID(),
		Rating:    req.Rating,
	})
	if err != nil {
		return entity.DoctorRating{}, err
	}
	return s.repo.GetById(ctx, id)
}
