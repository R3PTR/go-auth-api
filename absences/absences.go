package absences

type AbsencesService struct {
	absencesDbService *AbsencesDbService
}

func NewAbsencesService(absencesDbService *AbsencesDbService) *AbsencesService {
	return &AbsencesService{absencesDbService: absencesDbService}
}

// GetAbsences returns all absences.
func (a *AbsencesService) GetAbsences() ([]Absence, error) {
	return a.absencesDbService.GetAbsences()
}

// GetAbsencesByUserId returns all absences for a user.
func (a *AbsencesService) GetAbsencesByUserId(userId string) ([]Absence, error) {
	return a.absencesDbService.GetAbsencesByUserId(userId)
}

// GetAbsenceById returns an absence by ID.
func (a *AbsencesService) GetAbsenceById(id string) (Absence, error) {
	return a.absencesDbService.GetAbsenceById(id)
}

// CreateAbsence creates a new absence.
func (a *AbsencesService) CreateAbsence(newAbsence newAbsence, userId string) error {
	return a.absencesDbService.CreateAbsence(newAbsence, userId)
}

// UpdateOwnAbsence updates an absence.
func (a *AbsencesService) UpdateOwnAbsence(updateOwnAbsence UpdateOwnAbsence) error {
	return a.absencesDbService.UpdateOwnAbsence(updateOwnAbsence)
}

// UpdateAbsenceAsAdmin updates an absence as an admin.
func (a *AbsencesService) UpdateAbsenceAsAdmin(updateAbsenceAsAdmin UpdateAbsenceAsAdmin) error {
	return a.absencesDbService.UpdateAbsenceAsAdmin(updateAbsenceAsAdmin)
}

// DeleteAbsence deletes an absence.
func (a *AbsencesService) DeleteAbsence(id string) error {
	return a.absencesDbService.DeleteAbsence(id)
}
