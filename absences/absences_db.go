package absences

import (
	"context"

	"github.com/R3PTR/go-auth-api/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AbsencesDbService struct {
	mongoClient *database.MongoDBClient
}

func NewAbsencesDbService(mongoClient *database.MongoDBClient) *AbsencesDbService {
	return &AbsencesDbService{mongoClient: mongoClient}
}

// getAbsenceCollection returns the absence collection.
func (a *AbsencesDbService) getAbsenceCollection() *mongo.Collection {
	return a.mongoClient.GetCollection(a.mongoClient.Config.AbsencesDatabase, a.mongoClient.Config.AbsenceCollection)
}

// GetAbsences returns all absences.
func (a *AbsencesDbService) GetAbsences() ([]Absence, error) {
	var absences []Absence

	cursor, err := a.getAbsenceCollection().Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var absence Absence
		err := cursor.Decode(&absence)
		if err != nil {
			return nil, err
		}
		absences = append(absences, absence)
	}

	return absences, nil
}

// GetAbsencesByUserId returns all absences for a user.
func (a *AbsencesDbService) GetAbsencesByUserId(userId string) ([]Absence, error) {
	var absences []Absence

	cursor, err := a.getAbsenceCollection().Find(context.Background(), bson.M{"userId": userId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var absence Absence
		err := cursor.Decode(&absence)
		if err != nil {
			return nil, err
		}
		absences = append(absences, absence)
	}

	return absences, nil
}

// GetAbsenceById returns an absence by ID.
func (a *AbsencesDbService) GetAbsenceById(id string) (Absence, error) {
	var absence Absence
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return absence, err
	}
	err = a.getAbsenceCollection().FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&absence)
	return absence, err
}

// CreateAbsence creates a new absence.
func (a *AbsencesDbService) CreateAbsence(newAbsence newAbsence, userId string) error {
	absence := Absence{
		TypeOfAbsence: newAbsence.TypeOfAbsence,
		UserId:        userId,
		StartDate:     newAbsence.StartDate,
		EndDate:       newAbsence.EndDate,
		Reason:        newAbsence.Reason,
		Status:        "pending",
	}
	_, err := a.getAbsenceCollection().InsertOne(context.Background(), absence)
	if err != nil {
		return err
	}
	return nil
}

// UpdateOwnAbsence updates an absence.
func (a *AbsencesDbService) UpdateOwnAbsence(updateOwnAbsence UpdateOwnAbsence) error {
	objectId, err := primitive.ObjectIDFromHex(updateOwnAbsence.Id)
	if err != nil {
		return err
	}
	_, err = a.getAbsenceCollection().UpdateOne(context.Background(), bson.M{"_id": objectId}, bson.M{"$set": bson.M{"startDate": updateOwnAbsence.StartDate, "endDate": updateOwnAbsence.EndDate, "reason": updateOwnAbsence.Reason}})
	if err != nil {
		return err
	}
	return nil
}

// UpdateAbsenceAsAdmin updates an absence as an admin.
func (a *AbsencesDbService) UpdateAbsenceAsAdmin(updateAbsenceAsAdmin UpdateAbsenceAsAdmin) error {
	objectId, err := primitive.ObjectIDFromHex(updateAbsenceAsAdmin.Id)
	if err != nil {
		return err
	}
	_, err = a.getAbsenceCollection().UpdateOne(context.Background(), bson.M{"_id": objectId}, bson.M{"$set": bson.M{"status": updateAbsenceAsAdmin.Status, "reasonForRejection": updateAbsenceAsAdmin.ReasonForRejection}})
	if err != nil {
		return err
	}
	return nil
}

// DeleteAbsence deletes an absence.
func (a *AbsencesDbService) DeleteAbsence(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = a.getAbsenceCollection().DeleteOne(context.Background(), bson.M{"_id": objectId})
	return err
}
