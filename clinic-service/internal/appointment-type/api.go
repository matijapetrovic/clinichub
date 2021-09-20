package appointment_type

import (
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/matijapetrovic/clinichub/clinic-service/internal/errors"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/log"
)

func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Use(authHandler)

	r.Get("/appointment-types", res.getAll)

	r.Post("/appointment-types", res.create)
	r.Put("/appointment-types/<id>", res.update)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) getAll(c *routing.Context) error {
	appointmentTypes, err := r.service.GetAll(c.Request.Context())
	if err != nil {
		return err
	}

	return c.Write(appointmentTypes)
}

func (r resource) create(c *routing.Context) error {
	var request CreateAppointmentTypeRequest
	if err := c.Read(&request); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	appointmentType, err := r.service.Create((c.Request.Context()), request)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(appointmentType, http.StatusCreated)
}

func (r resource) update(c *routing.Context) error {
	var input UpdateAppointmentTypeRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	appointmentType, err := r.service.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		return err
	}

	return c.Write(appointmentType)
}
