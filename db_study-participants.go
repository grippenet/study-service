package main

import (
	"github.com/influenzanet/study-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// findParticipantsByStudyStatusDB retrieve all participant states from a study by status (e.g. active)
func findParticipantsByStudyStatusDB(instanceID string, studyKey string, studyStatus string, useProjection bool) (pStates []models.ParticipantState, err error) {
	ctx, cancel := getContext()
	defer cancel()

	filter := bson.M{"studyStatus": studyStatus}

	var opts *options.FindOptions
	if useProjection {
		projection := bson.D{
			primitive.E{Key: "studyStatus", Value: 1},   // {"secretKey", 1},
			primitive.E{Key: "participantID", Value: 1}, // {"secretKey", 1},
		}
		opts = options.Find().SetProjection(projection)
	}

	cur, err := collectionRefStudyParticipant(instanceID, studyKey).Find(
		ctx,
		filter,
		opts,
	)

	if err != nil {
		return pStates, err
	}
	defer cur.Close(ctx)

	pStates = []models.ParticipantState{}
	for cur.Next(ctx) {
		var result models.ParticipantState
		err := cur.Decode(&result)
		if err != nil {
			return pStates, err
		}

		pStates = append(pStates, result)
	}
	if err := cur.Err(); err != nil {
		return pStates, err
	}

	return pStates, nil
}

// findParticipantStateDB retrieves the participant state for a given participant from a study
func findParticipantStateDB(instanceID string, studyKey string, participantID string) (models.ParticipantState, error) {
	ctx, cancel := getContext()
	defer cancel()

	filter := bson.M{"participantID": participantID}

	elem := models.ParticipantState{}
	err := collectionRefStudyParticipant(instanceID, studyKey).FindOne(ctx, filter).Decode(&elem)
	return elem, err
}

// saveParticipantStateDB creates or replaces the participant states in the DB
func saveParticipantStateDB(instanceID string, studyKey string, pState models.ParticipantState) (models.ParticipantState, error) {
	ctx, cancel := getContext()
	defer cancel()

	filter := bson.M{"participantID": pState.ParticipantID}

	upsert := true
	rd := options.After
	options := options.FindOneAndReplaceOptions{
		Upsert:         &upsert,
		ReturnDocument: &rd,
	}
	elem := models.ParticipantState{}
	err := collectionRefStudyParticipant(instanceID, studyKey).FindOneAndReplace(
		ctx, filter, pState, &options,
	).Decode(&elem)
	return elem, err
}