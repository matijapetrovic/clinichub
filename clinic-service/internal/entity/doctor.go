package entity

import (
	"fmt"
	"strconv"
	"strings"
)

type Doctor struct {
	Id               string `json:"id"`
	ClinicId         string `json:"clinicId"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	WorkStart        string `json:"workStart"`
	WorkEnd          string `json:"workEnd"`
	SpecializationId string `json:"-" db:"specialization_id"`
	AppointmentType  `json:"specialization" db:"-"`
}

type Time struct {
	Hour   uint `json:"hour"`
	Minute uint `json:"minute"`
}

func (t Time) Parse(s string) (Time, error) {
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
