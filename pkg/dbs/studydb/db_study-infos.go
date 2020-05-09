package studydb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/influenzanet/study-service/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (dbService *StudyDBService) GetStudyByStudyKey(instanceID string, studyKey string) (study models.Study, err error) {
	if err = dbService.collectionRefStudyInfos(instanceID).FindOne(
		context.Background(),
		bson.D{
			primitive.E{Key: "key", Value: studyKey}, //{"studyKey", studyKey},
		},
		options.FindOne(),
	).Decode(&study); err != nil {
		return study, err
	}
	return
}

func (dbService *StudyDBService) GetStudiesByStatus(instanceID string, status string, onlyKeys bool) (studies []models.Study, err error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	filter := bson.M{"status": status}

	var opts *options.FindOptions
	if onlyKeys {
		projection := bson.D{
			primitive.E{Key: "key", Value: 1},       // {"secretKey", 1},
			primitive.E{Key: "secretKey", Value: 1}, // {"secretKey", 1},
		}
		opts = options.Find().SetProjection(projection)
	}

	cur, err := dbService.collectionRefStudyInfos(instanceID).Find(
		ctx,
		filter,
		opts,
	)

	if err != nil {
		return studies, err
	}
	defer cur.Close(ctx)

	studies = []models.Study{}
	for cur.Next(ctx) {
		var result models.Study
		err := cur.Decode(&result)
		if err != nil {
			return studies, err
		}

		studies = append(studies, result)
	}
	if err := cur.Err(); err != nil {
		return studies, err
	}

	return studies, nil
}

func (dbService *StudyDBService) GetStudySecretKey(instanceID string, studyKey string) (secretKey string, err error) {
	projection := bson.D{
		primitive.E{Key: "secretKey", Value: 1}, // {"secretKey", 1},
	}

	var study models.Study
	if err = dbService.collectionRefStudyInfos(instanceID).FindOne(
		context.Background(),
		bson.D{
			primitive.E{Key: "key", Value: studyKey}, //{"studyKey", studyKey},
		},
		options.FindOne().SetProjection(projection),
	).Decode(&study); err != nil {
		return "", err
	}
	return study.SecretKey, nil
}

func (dbService *StudyDBService) GetStudyMembers(instanceID string, studyKey string) (members []models.StudyMember, err error) {
	projection := bson.D{
		primitive.E{Key: "members", Value: 1}, // {"members", 1},
	}

	var study models.Study
	if err = dbService.collectionRefStudyInfos(instanceID).FindOne(
		context.Background(),
		bson.D{
			primitive.E{Key: "key", Value: studyKey}, //{"studyKey", studyKey},
		},
		options.FindOne().SetProjection(projection),
	).Decode(&study); err != nil {
		return []models.StudyMember{}, err
	}
	return study.Members, nil
}

func (dbService *StudyDBService) GetStudyRules(instanceID string, studyKey string) (rules []models.Expression, err error) {
	projection := bson.D{
		primitive.E{Key: "rules", Value: 1}, // {"members", 1},
	}

	var study models.Study
	if err = dbService.collectionRefStudyInfos(instanceID).FindOne(
		context.Background(),
		bson.D{
			primitive.E{Key: "key", Value: studyKey}, //{"studyKey", studyKey},
		},
		options.FindOne().SetProjection(projection),
	).Decode(&study); err != nil {
		return []models.Expression{}, err
	}
	return study.Rules, nil
}

// saveParticipantStateDB creates or replaces the participant states in the DB
func (dbService *StudyDBService) CreateStudy(instanceID string, study models.Study) (models.Study, error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	filter := bson.M{"key": study.Key}
	if res := dbService.collectionRefStudyInfos(instanceID).FindOne(ctx, filter); res.Err() == nil {
		return study, fmt.Errorf("studyKey already used: %s", study.Key)
	}

	res, err := dbService.collectionRefStudyInfos(instanceID).InsertOne(ctx, study)
	id, ok := res.InsertedID.(primitive.ObjectID)
	if ok {
		study.ID = id
	}
	return study, err
}

func (dbService *StudyDBService) UpdateStudyKey(instanceID string, oldKey string, newKey string) error {
	ctx, cancel := dbService.getContext()
	defer cancel()

	_, err := dbService.GetStudyByStudyKey(instanceID, newKey)
	if err == nil {
		return errors.New("newKey already exists")
	}

	filter := bson.M{
		"key": oldKey,
	}
	update := bson.M{"$set": bson.M{"key": newKey}}
	_, err = dbService.collectionRefStudyInfos(instanceID).UpdateOne(ctx, filter, update)
	return err
}

func (dbService *StudyDBService) UpdateStudyStatus(instanceID string, studyKey string, status string) error {
	ctx, cancel := dbService.getContext()
	defer cancel()

	filter := bson.M{
		"key": studyKey,
	}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := dbService.collectionRefStudyInfos(instanceID).UpdateOne(ctx, filter, update)
	return err
}

func (dbService *StudyDBService) UpdateStudyInfo(instanceID string, study models.Study) (models.Study, error) {
	ctx, cancel := dbService.getContext()
	defer cancel()

	elem := models.Study{}
	filter := bson.M{"key": study.Key}
	rd := options.After
	fro := options.FindOneAndReplaceOptions{
		ReturnDocument: &rd,
	}
	err := dbService.collectionRefStudyInfos(instanceID).FindOneAndReplace(ctx, filter, study, &fro).Decode(&elem)
	return elem, err
}

func (dbService *StudyDBService) ShouldPerformTimerEvent(instanceID string, studyKey string, timerEventFrequency int64) error {
	ctx, cancel := dbService.getContext()
	defer cancel()

	filter := bson.M{
		"key":                 studyKey,
		"nextTimerEventAfter": bson.M{"$lt": time.Now().Unix()},
	}
	update := bson.M{"$set": bson.M{"nextTimerEventAfter": time.Now().Unix() + timerEventFrequency}}
	res, err := dbService.collectionRefStudyInfos(instanceID).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount < 1 {
		return errors.New("not modified")
	}
	return nil
}