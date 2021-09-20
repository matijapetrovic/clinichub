package doctor_rating

import (
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/matijapetrovic/clinichub/rating-service/internal/errors"
	"github.com/matijapetrovic/clinichub/rating-service/pkg/log"
	"net/http"
)

func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	// r.Use(authHandler)

	r.Get("/doctors/<id>/average-rating", res.getRating)
	r.Get("/doctors/to-rate", res.getAvailableRatings)
	r.Post("/doctors/<id>/ratings", res.rateDoctor)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) getRating(c *routing.Context) error {
	rating, err := r.service.GetDoctorRating(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(rating)
}

func (r resource) getAvailableRatings(c *routing.Context) error {
	doctors, err := r.service.GetAvaialableRatings(c.Request)
	if err != nil {
		return err
	}

	return c.Write(doctors)
}

func (r resource) rateDoctor(c *routing.Context) error {
	var request RateDoctorRequest
	if err := c.Read(&request); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	rating, err := r.service.RateDoctor(c.Request.Context(), c.Param("id"), request)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(rating, http.StatusCreated)
}
