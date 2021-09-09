package entity

import (
	"fmt"
	"strconv"
	"strings"
)

type Doctor struct {
	Id                   string `json:"id"`
	ClinicId             string `json:"clinicId"`
	FirstName            string `json:"firstName"`
	LastName             string `json:"lastName"`
	WorkStart            string `json:"workStart"`
	WorkEnd              string `json:"workEnd"`
	SpecializationId     string `json:"-" db:"specialization_id"`
	AppointmentType      `json:"specialization" db:"-"`
	AppointmentTypePrice uint     `json:"specializationPrice" db:"-"`
	AvailableHours       []string `json:"availableHours" db:"-"`
}

type Time struct {
	Hour   uint `json:"hour"`
	Minute uint `json:"minute"`
}

func ParseTime(s string) (Time, error) {
	split := strings.Split(s, ":")
	hour, err := strconv.Atoi(split[0])
	if err != nil {
		return Time{}, err
	}
	minute, err := strconv.Atoi(split[1])
	if err != nil {
		return Time{}, err
	}
	return Time{Hour: uint(hour), Minute: uint(minute)}, nil
}

func (t Time) ToString() string {
	hour := fmt.Sprintf("%d", t.Hour)
	if t.Hour < 10 {
		hour = "0" + hour
	}
	minute := fmt.Sprintf("%d", t.Minute)
	if t.Minute < 10 {
		minute = "0" + minute
	}
	return fmt.Sprintf("%s:%s", hour, minute)
}

func (t Time) Before(o Time) bool {
	if t.Hour < o.Hour {
		return true
	}

	if t.Hour > o.Hour {
		return false
	}

	if t.Minute < o.Minute {
		return true
	}
	return false
}

func GetHours(from Time, to Time) map[uint]struct{} {
	hours := make(map[uint]struct{})

	if from.Before(to) {
		for hour := from.Hour; hour < to.Hour; hour++ {
			hours[hour] = struct{}{}
		}
	} else {
		for hour := from.Hour; hour < 24; hour++ {
			hours[hour] = struct{}{}
		}

		for hour := uint(0); hour < to.Hour; hour++ {
			hours[hour] = struct{}{}
		}
	}
	return hours

}
