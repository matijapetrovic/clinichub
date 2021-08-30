package clinic

import (
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/matijapetrovic/clinichub/clinic-service/internal/errors"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/log"
)

func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	// r.Use(authHandler)

	r.Get("/clinics", res.getAll)
	r.Get("/clinics/<id>/prices", res.getPrices)
	r.Post("/clinics", res.create)
	r.Put("/clinics/<id>", res.update)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) getAll(c *routing.Context) error {
	clinics, err := r.service.GetAll(c.Request.Context())
	if err != nil {
		return err
	}

	return c.Write(clinics)
}

func (r resource) getPrices(c *routing.Context) error {
	prices, err := r.service.GetAppointmentTypePrices(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(prices)
}

func (r resource) create(c *routing.Context) error {
	var request CreateClinicRequest
	if err := c.Read(&request); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	clinic, err := r.service.Create((c.Request.Context()), request)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(clinic, http.StatusCreated)
}

func (r resource) update(c *routing.Context) error {
	var input UpdateClinicRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	clinic, err := r.service.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		return err
	}

	return c.Write(clinic)
}
