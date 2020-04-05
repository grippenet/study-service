package main

import (
	"testing"

	"github.com/influenzanet/study-service/models"
)

func TestDbCreateStudyInfos(t *testing.T) {
	testStudy := models.Study{
		Key:       "testStudyKey1",
		SecretKey: "supersecret",
		Members: []models.StudyMember{
			models.StudyMember{
				UserID: "testuser",
				Role:   "maintainer",
			},
		},
	}

	t.Run("Create study with not existing key", func(t *testing.T) {
		study, err := createStudyInDB(testInstanceID, testStudy)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		if study.ID.IsZero() {
			t.Errorf("unexpected id: %s", study.ID.Hex())
		}
	})

	t.Run("Try to create study with existing key", func(t *testing.T) {
		_, err := createStudyInDB(testInstanceID, testStudy)
		if err == nil {
			t.Error("should return error")
		}
	})
}

func TestDbUpdateStudyInfos(t *testing.T) {
	testStudies := []models.Study{
		models.Study{Key: "test1", Status: "active", Members: []models.StudyMember{
			models.StudyMember{
				UserID: "testuser",
				Role:   "maintainer",
			},
		}},
		models.Study{Key: "test2", Status: "active", Members: []models.StudyMember{
			models.StudyMember{
				UserID: "testuser",
				Role:   "maintainer",
			},
		}},
	}

	for _, s := range testStudies {
		_, err := createStudyInDB(testInstanceID, s)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}

	t.Run("Try to update study key with existing key", func(t *testing.T) {
		if err := updateStudyKeyInDB(testInstanceID, "test1", "test2"); err == nil {
			t.Error("should fail with error when key exists")
		}
	})

	t.Run("Try to update study key with not existing key", func(t *testing.T) {
		if err := updateStudyKeyInDB(testInstanceID, "test1", "test3"); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	})

	t.Run("Update study status", func(t *testing.T) {
		if err := updateStudyStatusInDB(testInstanceID, "test1", "inactive"); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	})

	t.Run("Try to update other properties with ok key", func(t *testing.T) {
		testStudies[1].SecretKey = "343434"
		upd, err := updateStudyInfoInDB(testInstanceID, testStudies[1])
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if upd.SecretKey != "343434" {
			t.Errorf("unexpected value: %s, %s (have, want)", upd.SecretKey, "343434")
		}
	})

	t.Run("Try to update other properties with wrong key", func(t *testing.T) {
		testStudies[1].Key = "wrong"
		testStudies[1].SecretKey = "34343w4"
		_, err := updateStudyInfoInDB(testInstanceID, testStudies[1])
		if err == nil {
			t.Error("should return error since key not known")
			return
		}
	})
}

func TestDbGetStudyInfos(t *testing.T) {
	testStudies := []models.Study{
		models.Study{Key: "test1", SecretKey: "testsecret", Status: "active", Members: []models.StudyMember{
			models.StudyMember{
				UserID: "testuser",
				Role:   "maintainer",
			},
		}},
	}

	for _, s := range testStudies {
		_, err := createStudyInDB(testInstanceID, s)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}

	t.Run("Get secret key", func(t *testing.T) {
		secret, err := getStudySecretKey(testInstanceID, testStudies[0].Key)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if secret != testStudies[0].SecretKey {
			t.Errorf("unexpected value: %s, %s (have, want)", secret, testStudies[0].SecretKey)
		}
	})

	t.Run("Get members", func(t *testing.T) {
		members, err := getStudyMembers(testInstanceID, testStudies[0].Key)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if len(members) != 1 {
			t.Errorf("unexpected number of members: %d", len(members))
		}
	})
}
