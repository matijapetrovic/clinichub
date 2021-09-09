package entity

import "time"

type Appointment struct {
	Id                string    `json:"id"`
	ClinicId          string    `json:"clinicId"`
	DoctorId          string    `json:"doctorId"`
	PatientId         string    `json:"patientId"`
	AppointmentTypeId string    `json:"appointmentTypeId"`
	Price             int       `json:"price"`
	Time              time.Time `json:"time"`
}
