package absences

type Absence struct {
	Id                 string `json:"id"`
	TypeOfAbsence      string `json:"typeOfAbsence"`
	UserId             string `json:"userId"`
	StartDate          string `json:"startDate"`
	EndDate            string `json:"endDate"`
	Reason             string `json:"reason"`
	Status             string `json:"status"`
	ReasonForRejection string `json:"reasonForRejection,omitempty"`
}

type newAbsence struct {
	TypeOfAbsence string `json:"typeOfAbsence"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	Reason        string `json:"reason"`
}

type UpdateOwnAbsence struct {
	Id        string `json:"id"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Reason    string `json:"reason"`
}

type UpdateAbsenceAsAdmin struct {
	Id                 string `json:"id"`
	Status             string `json:"status"`
	ReasonForRejection string `json:"reasonForRejection,omitempty"`
}
