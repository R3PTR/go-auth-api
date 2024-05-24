package absences

import (
	"time"
)

type DateRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type Absence struct {
	Id                 string    `json:"id"`
	TypeOfAbsence      string    `json:"typeOfAbsence"`
	UserId             string    `json:"userId"`
	DateRange          DateRange `json:"dateRange"`
	TotalDays          int       `json:"totalDays"`
	Reason             string    `json:"reason"`
	Status             string    `json:"status"`
	ReasonForRejection string    `json:"reasonForRejection,omitempty"`
}

type newAbsence struct {
	TypeOfAbsence string    `json:"typeOfAbsence"`
	DateRange     DateRange `json:"dateRange"`
	TotalDays     int       `json:"totalDays"`
	Reason        string    `json:"reason,omitempty"`
}

type UpdateOwnAbsence struct {
	Id        string    `json:"id"`
	DateRange DateRange `json:"dateRange"`
	TotalDays int       `json:"totalDays"`
	Reason    string    `json:"reason,omitempty"`
}

type UpdateAbsenceAsAdmin struct {
	Id                 string `json:"id"`
	Status             string `json:"status"`
	ReasonForRejection string `json:"reasonForRejection,omitempty"`
}
