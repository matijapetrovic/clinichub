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
	r.Post("/clinics/<id>/prices", res.addPrice)
	r.Put("/clinics/<id>/prices", res.updatePrice)
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
	var request UpdateClinicRequest
	if err := c.Read(&request); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	clinic, err := r.service.Update(c.Request.Context(), c.Param("id"), request)
	if err != nil {
		return err
	}

	return c.Write(clinic)
}

func (r resource) addPrice(c *routing.Context) error {
	var request AddAppointmentTypePriceRequest
	if err := c.Read(&request); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	price, err := r.service.AddAppointmentTypePrice((c.Request.Context()), c.Param("id"), request)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(price, http.StatusCreated)
}

func (r resource) updatePrice(c *routing.Context) error {
	var request UpdateAppointmentTypePriceRequest
	if err := c.Read(&request); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	clinic, err := r.service.UpdateAppointmentTypePrice(c.Request.Context(), c.Param("id"), request)
	if err != nil {
		return err
	}

	return c.Write(clinic)
}
