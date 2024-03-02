package absences

import (
	"net/http"

	"github.com/R3PTR/go-auth-api/auth"
	"github.com/gin-gonic/gin"
)

// AbsencesController
type AbsencesController struct {
	absencesService *AbsencesService
}

// NewAbsencesController
func NewAbsencesController(absencesService *AbsencesService) *AbsencesController {
	return &AbsencesController{absencesService: absencesService}
}

// GetAbsences
func (a *AbsencesController) GetAllAbsences(c *gin.Context) {
	absences, err := a.absencesService.GetAbsences()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, absences)
}

// GetAbsencesByUserId
func (a *AbsencesController) GetAbsences(c *gin.Context) {
	user_unasserted, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
	}
	user, ok := user_unasserted.(*auth.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
	}
	absences, err := a.absencesService.GetAbsencesByUserId(user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, absences)
}

// CreateAbsence
func (a *AbsencesController) CreateAbsence(c *gin.Context) {
	user_unasserted, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
	}
	user, ok := user_unasserted.(*auth.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
	}
	var newAbsence newAbsence
	if err := c.ShouldBindJSON(&newAbsence); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err := a.absencesService.CreateAbsence(newAbsence, user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Absence created"})
}

// UpdateOwnAbsence
func (a *AbsencesController) UpdateOwnAbsence(c *gin.Context) {
	user_unasserted, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
	}
	user, ok := user_unasserted.(*auth.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
	}
	var updateOwnAbsence UpdateOwnAbsence
	if err := c.ShouldBindJSON(&updateOwnAbsence); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	absence, err := a.absencesService.GetAbsenceById(updateOwnAbsence.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	if absence.UserId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this absence"})
	}
	err = a.absencesService.UpdateOwnAbsence(updateOwnAbsence)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Absence updated"})
}

// UpdateAbsenceAsAdmin
func (a *AbsencesController) UpdateAbsenceAsAdmin(c *gin.Context) {
	var updateAbsenceAsAdmin UpdateAbsenceAsAdmin
	if err := c.ShouldBindJSON(&updateAbsenceAsAdmin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err := a.absencesService.UpdateAbsenceAsAdmin(updateAbsenceAsAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Absence updated"})
}

// DeleteAbsence
func (a *AbsencesController) DeleteAbsence(c *gin.Context) {
	id := c.Param("id")
	err := a.absencesService.DeleteAbsence(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Absence deleted"})
}
