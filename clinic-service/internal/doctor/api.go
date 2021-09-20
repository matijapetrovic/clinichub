package doctor

import (
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/matijapetrovic/clinichub/clinic-service/internal/errors"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/log"
)

func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Use(authHandler)

	r.Get("/doctors", res.getAll)
	r.Get("/doctors/<id>", res.getById)
	r.Get("/clinics/<clinicId>/doctors", res.getByClinicId)

	r.Post("/doctors", res.create)
	r.Put("/doctors/<id>", res.update)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) getByClinicId(c *routing.Context) error {
	appointmentTypeId := c.Request.URL.Query().Get("appointmentTypeId")
	date := c.Request.URL.Query().Get("date")
	doctors, err := r.service.GetByClinicId(c.Request, c.Param("clinicId"), GetByClinicIdRequest{
		AppointmentTypeId: appointmentTypeId,
		Date:              date,
	})
	if err != nil {
		return err
	}

	return c.Write(doctors)
}

func (r resource) getById(c *routing.Context) error {
	doctor, err := r.service.GetById(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(doctor)
}

func (r resource) getAll(c *routing.Context) error {
	doctors, err := r.service.GetAll(c.Request.Context())
	if err != nil {
		return err
	}

	return c.Write(doctors)
}

func (r resource) create(c *routing.Context) error {
	var request CreateDoctorRequest
	if err := c.Read(&request); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	doctor, err := r.service.Create((c.Request.Context()), request)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(doctor, http.StatusCreated)
}

func (r resource) update(c *routing.Context) error {
	var input UpdateDoctorRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	doctor, err := r.service.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		return err
	}

	return c.Write(doctor)
}
