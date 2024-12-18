package studyengine

import (
	"testing"
	"time"
	"fmt"
	"github.com/influenzanet/study-service/pkg/dbs/studydb"
	"github.com/influenzanet/study-service/pkg/types"
)

// Reference/Lookup methods
func TestEvalCheckEventType(t *testing.T) {
	exp := types.Expression{Name: "checkEventType", Data: []types.ExpressionArg{
		{DType: "str", Str: "ENTER"},
	}}

	t.Run("for matching", func(t *testing.T) {
		EvalContext := EvalContext{
			Event: types.StudyEvent{Type: "ENTER"},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected type or value: %s", ret)
		}
	})

	t.Run("for not matching", func(t *testing.T) {
		EvalContext := EvalContext{
			Event: types.StudyEvent{Type: "enter"},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected type or value: %s", ret)
		}
	})
}

func TestEvalCheckSurveyResponseKey(t *testing.T) {
	exp := types.Expression{Name: "checkSurveyResponseKey", Data: []types.ExpressionArg{
		{DType: "str", Str: "weekly"},
	}}

	t.Run("for no survey responses at all", func(t *testing.T) {
		EvalContext := EvalContext{
			Event: types.StudyEvent{Type: "SUBMIT"},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected type or value: %s", ret)
		}
	})

	t.Run("not matching key", func(t *testing.T) {
		EvalContext := EvalContext{
			Event: types.StudyEvent{
				Type: "SUBMIT",
				Response: types.SurveyResponse{
					Key:       "intake",
					Responses: []types.SurveyItemResponse{},
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected type or value: %s", ret)
		}
	})

	t.Run("for matching key", func(t *testing.T) {
		EvalContext := EvalContext{
			Event: types.StudyEvent{
				Type: "SUBMIT",
				Response: types.SurveyResponse{
					Key:       "weekly",
					Responses: []types.SurveyItemResponse{},
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected type or value: %s", ret)
		}
	})
}

func TestEvalHasStudyStatus(t *testing.T) {
	t.Run("with not matching state", func(t *testing.T) {
		exp := types.Expression{Name: "hasStudyStatus", Data: []types.ExpressionArg{
			{DType: "str", Str: types.PARTICIPANT_STUDY_STATUS_EXITED},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("with matching state", func(t *testing.T) {
		exp := types.Expression{Name: "hasStudyStatus", Data: []types.ExpressionArg{
			{DType: "str", Str: types.PARTICIPANT_STUDY_STATUS_ACTIVE},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
}

type MockStudyDBService struct {
	Responses []types.SurveyResponse
}

func (db MockStudyDBService) FindSurveyResponses(instanceID string, studyKey string, query studydb.ResponseQuery) (responses []types.SurveyResponse, err error) {
	for _, r := range db.Responses {
		if query.ParticipantID != r.ParticipantID {
			continue
		}
		if len(query.SurveyKey) > 0 && query.SurveyKey != r.Key {
			continue
		}
		if query.Since > 0 && r.SubmittedAt < query.Since {
			continue
		}
		if query.Until > 0 && r.SubmittedAt > query.Until {
			continue
		}
		responses = append(responses, r)
	}

	return responses, nil
}

func (db MockStudyDBService) DeleteConfidentialResponses(instanceID string, studyKey string, participantID string, key string) (count int64, err error) {
	return
}

func (db MockStudyDBService) SaveResearcherMessage(instanceID string, studyKey string, message types.StudyMessage) error {
	return nil
}

func TestEvalCheckConditionForOldResponses(t *testing.T) {
	testResponses := []types.SurveyResponse{
		{
			Key: "S1", ParticipantID: "P1", SubmittedAt: 10, Responses: []types.SurveyItemResponse{
				{Key: "S1.Q1", Response: &types.ResponseItem{
					Key: "rg", Items: []*types.ResponseItem{
						{Key: "scg", Items: []*types.ResponseItem{{Key: "1"}}},
					},
				}}},
		},
		{
			Key: "S1", ParticipantID: "P1", SubmittedAt: 13, Responses: []types.SurveyItemResponse{
				{Key: "S1.Q1", Response: &types.ResponseItem{
					Key: "rg", Items: []*types.ResponseItem{
						{Key: "scg", Items: []*types.ResponseItem{{Key: "1"}}},
					},
				}}},
		},
		{
			Key: "S1", ParticipantID: "P2", SubmittedAt: 13, Responses: []types.SurveyItemResponse{
				{Key: "S1.Q1", Response: &types.ResponseItem{
					Key: "rg", Items: []*types.ResponseItem{
						{Key: "scg", Items: []*types.ResponseItem{{Key: "1"}}},
					},
				}}},
		},
		{
			Key: "S2", ParticipantID: "P1", SubmittedAt: 15, Responses: []types.SurveyItemResponse{
				{Key: "S2.Q1", Response: &types.ResponseItem{
					Key: "rg", Items: []*types.ResponseItem{
						{Key: "scg", Items: []*types.ResponseItem{{Key: "1"}}},
					},
				}}},
		},
		{
			Key: "S1", ParticipantID: "P1", SubmittedAt: 17, Responses: []types.SurveyItemResponse{
				{Key: "S1.Q1", Response: &types.ResponseItem{
					Key: "rg", Items: []*types.ResponseItem{
						{Key: "scg", Items: []*types.ResponseItem{{Key: "1"}}},
					},
				}}},
		},
		{
			Key: "S1", ParticipantID: "P1", SubmittedAt: 22, Responses: []types.SurveyItemResponse{
				{Key: "S1.Q1", Response: &types.ResponseItem{
					Key: "rg", Items: []*types.ResponseItem{
						{Key: "scg", Items: []*types.ResponseItem{{Key: "2"}}},
					},
				}}},
		},
	}

	t.Run("missing DB config", func(t *testing.T) {
		exp := types.Expression{Name: "checkConditionForOldResponses"}

		EvalContext := EvalContext{
			Configs: ActionConfigs{
				DBService: nil,
			},
		}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should return error")
			return
		}
	})

	t.Run("missing instanceID", func(t *testing.T) {
		exp := types.Expression{Name: "checkConditionForOldResponses"}

		EvalContext := EvalContext{
			Configs: ActionConfigs{
				DBService: MockStudyDBService{},
			},
			Event: types.StudyEvent{
				StudyKey: "testStudy",
			},
		}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should return error")
			return
		}
	})

	t.Run("missing studyKey", func(t *testing.T) {
		exp := types.Expression{Name: "checkConditionForOldResponses"}

		EvalContext := EvalContext{
			Configs: ActionConfigs{
				DBService: MockStudyDBService{
					Responses: testResponses,
				},
			},
			Event: types.StudyEvent{
				StudyKey:   "testStudy",
				InstanceID: "testInstance",
			},
		}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should return error")
			return
		}
	})

	t.Run("missing condition", func(t *testing.T) {
		exp := types.Expression{Name: "checkConditionForOldResponses"}

		EvalContext := EvalContext{
			Configs: ActionConfigs{
				DBService: MockStudyDBService{
					Responses: testResponses,
				},
			},
			Event: types.StudyEvent{
				StudyKey:   "testStudy",
				InstanceID: "testInstance",
			},
			ParticipantState: types.ParticipantState{
				ParticipantID: "P1",
			},
		}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should return error")
			return
		}
	})

	t.Run("checkType all", func(t *testing.T) {
		exp := types.Expression{Name: "checkConditionForOldResponses", Data: []types.ExpressionArg{
			{Exp: &types.Expression{
				Name: "responseHasKeysAny",
				Data: []types.ExpressionArg{
					{Str: "S1.Q1", DType: "str"},
					{Str: "rg.scg", DType: "str"},
					{Str: "1", DType: "str"},
				},
			}, DType: "exp"},
		}}

		EvalContext := EvalContext{
			Configs: ActionConfigs{
				DBService: MockStudyDBService{
					Responses: testResponses,
				},
			},
			Event: types.StudyEvent{
				StudyKey:   "testStudy",
				InstanceID: "testInstance",
			},
			ParticipantState: types.ParticipantState{
				ParticipantID: "P1",
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value retrieved: %d", ret)
		}
	})

	t.Run("checkType any", func(t *testing.T) {
		exp := types.Expression{Name: "checkConditionForOldResponses", Data: []types.ExpressionArg{
			{Exp: &types.Expression{
				Name: "responseHasKeysAny",
				Data: []types.ExpressionArg{
					{Str: "S1.Q1", DType: "str"},
					{Str: "rg.scg", DType: "str"},
					{Str: "1", DType: "str"},
				},
			}, DType: "exp"},
			{Str: "any", DType: "str"},
		}}

		EvalContext := EvalContext{
			Configs: ActionConfigs{
				DBService: MockStudyDBService{
					Responses: testResponses,
				},
			},
			Event: types.StudyEvent{
				StudyKey:   "testStudy",
				InstanceID: "testInstance",
			},
			ParticipantState: types.ParticipantState{
				ParticipantID: "P1",
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value retrieved: %d", ret)
		}
	})

	t.Run("checkType count - with enough", func(t *testing.T) {
		exp := types.Expression{Name: "checkConditionForOldResponses", Data: []types.ExpressionArg{
			{Exp: &types.Expression{
				Name: "responseHasKeysAny",
				Data: []types.ExpressionArg{
					{Str: "S1.Q1", DType: "str"},
					{Str: "rg.scg", DType: "str"},
					{Str: "1", DType: "str"},
				},
			}, DType: "exp"},
			{Num: 3, DType: "num"},
		}}

		EvalContext := EvalContext{
			Configs: ActionConfigs{
				DBService: MockStudyDBService{
					Responses: testResponses,
				},
			},
			Event: types.StudyEvent{
				StudyKey:   "testStudy",
				InstanceID: "testInstance",
			},
			ParticipantState: types.ParticipantState{
				ParticipantID: "P1",
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value retrieved: %d", ret)
		}
	})

	t.Run("checkType count - with not enough", func(t *testing.T) {
		exp := types.Expression{Name: "checkConditionForOldResponses", Data: []types.ExpressionArg{
			{Exp: &types.Expression{
				Name: "responseHasKeysAny",
				Data: []types.ExpressionArg{
					{Str: "S1.Q1", DType: "str"},
					{Str: "rg.scg", DType: "str"},
					{Str: "1", DType: "str"},
				},
			}, DType: "exp"},
			{Num: 4, DType: "num"},
		}}

		EvalContext := EvalContext{
			Configs: ActionConfigs{
				DBService: MockStudyDBService{
					Responses: testResponses,
				},
			},
			Event: types.StudyEvent{
				StudyKey:   "testStudy",
				InstanceID: "testInstance",
			},
			ParticipantState: types.ParticipantState{
				ParticipantID: "P1",
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value retrieved: %d", ret)
		}
	})

	t.Run("filter for survey type", func(t *testing.T) {
		exp := types.Expression{Name: "checkConditionForOldResponses", Data: []types.ExpressionArg{
			{Exp: &types.Expression{
				Name: "responseHasKeysAny",
				Data: []types.ExpressionArg{
					{Str: "S1.Q1", DType: "str"},
					{Str: "rg.scg", DType: "str"},
					{Str: "1", DType: "str"},
				},
			}, DType: "exp"},
			{Num: 4, DType: "num"},
			{Str: "S2", DType: "str"},
		}}

		EvalContext := EvalContext{
			Configs: ActionConfigs{
				DBService: MockStudyDBService{
					Responses: testResponses,
				},
			},
			Event: types.StudyEvent{
				StudyKey:   "testStudy",
				InstanceID: "testInstance",
			},
			ParticipantState: types.ParticipantState{
				ParticipantID: "P1",
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value retrieved: %d", ret)
		}
	})

	t.Run("filter for interval type", func(t *testing.T) {
		exp := types.Expression{Name: "checkConditionForOldResponses", Data: []types.ExpressionArg{
			{Exp: &types.Expression{
				Name: "responseHasKeysAny",
				Data: []types.ExpressionArg{
					{Str: "S1.Q1", DType: "str"},
					{Str: "rg.scg", DType: "str"},
					{Str: "1", DType: "str"},
				},
			}, DType: "exp"},
			{Num: 2, DType: "num"},
			{Str: "", DType: "str"},
			{Num: 16, DType: "num"},
			{Num: 18, DType: "num"},
		}}

		EvalContext := EvalContext{
			Configs: ActionConfigs{
				DBService: MockStudyDBService{
					Responses: testResponses,
				},
			},
			Event: types.StudyEvent{
				StudyKey:   "testStudy",
				InstanceID: "testInstance",
			},
			ParticipantState: types.ParticipantState{
				ParticipantID: "P1",
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value retrieved: %d", ret)
		}
	})
}

func TestEvalGetStudyEntryTime(t *testing.T) {
	t.Run("try retrieve entered at time", func(t *testing.T) {
		exp := types.Expression{Name: "getStudyEntryTime"}
		tStart := time.Now().Unix()
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				EnteredAt:   tStart,
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if ret.(float64) != float64(tStart) {
			t.Errorf("unexpected value retrieved: %d", ret)
		}
	})
}

func TestEvalHasSurveyKeyAssigned(t *testing.T) {
	t.Run("has survey assigned", func(t *testing.T) {
		exp := types.Expression{Name: "hasSurveyKeyAssigned", Data: []types.ExpressionArg{
			{DType: "str", Str: "test1"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				AssignedSurveys: []types.AssignedSurvey{
					{SurveyKey: "test1"},
					{SurveyKey: "test2"},
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value retrieved: %d", ret)
		}
	})

	t.Run("doesn't have the survey assigned", func(t *testing.T) {
		exp := types.Expression{Name: "hasSurveyKeyAssigned", Data: []types.ExpressionArg{
			{DType: "str", Str: "test1"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				AssignedSurveys: []types.AssignedSurvey{
					{SurveyKey: "test2"},
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value retrieved: %d", ret)
		}
	})

	t.Run("missing argument", func(t *testing.T) {
		exp := types.Expression{Name: "hasSurveyKeyAssigned", Data: []types.ExpressionArg{}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				AssignedSurveys: []types.AssignedSurvey{
					{SurveyKey: "test2"},
				},
			},
		}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should throw an error about missing arg")
			return
		}
	})

	t.Run("wrong argument", func(t *testing.T) {
		exp := types.Expression{Name: "hasSurveyKeyAssigned", Data: []types.ExpressionArg{
			{DType: "exp", Exp: &types.Expression{}},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				AssignedSurveys: []types.AssignedSurvey{
					{SurveyKey: "test2"},
				},
			},
		}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should throw an error about arg type")
			return
		}
	})
}

func TestEvalGetSurveyKeyAssignedFrom(t *testing.T) {
	t.Run("has survey assigned", func(t *testing.T) {
		exp := types.Expression{Name: "getSurveyKeyAssignedFrom", Data: []types.ExpressionArg{
			{DType: "str", Str: "test1"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				AssignedSurveys: []types.AssignedSurvey{
					{SurveyKey: "test1", ValidFrom: 10, ValidUntil: 100},
					{SurveyKey: "test2", ValidFrom: 10, ValidUntil: 100},
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if ret.(float64) != 10 {
			t.Errorf("unexpected value retrieved: %d", ret)
		}
	})

	t.Run("doesn't have the survey assigned", func(t *testing.T) {
		exp := types.Expression{Name: "getSurveyKeyAssignedFrom", Data: []types.ExpressionArg{
			{DType: "str", Str: "test1"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				AssignedSurveys: []types.AssignedSurvey{
					{SurveyKey: "test2", ValidFrom: 10, ValidUntil: 100},
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if ret.(float64) != -1 {
			t.Errorf("unexpected value retrieved: %d", ret)
		}
	})

	t.Run("missing argument", func(t *testing.T) {
		exp := types.Expression{Name: "getSurveyKeyAssignedFrom", Data: []types.ExpressionArg{}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				AssignedSurveys: []types.AssignedSurvey{
					{SurveyKey: "test1", ValidFrom: 10, ValidUntil: 100},
					{SurveyKey: "test2", ValidFrom: 10, ValidUntil: 100},
				},
			},
		}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should throw an error about missing arg")
			return
		}
	})

	t.Run("wrong argument", func(t *testing.T) {
		exp := types.Expression{Name: "getSurveyKeyAssignedFrom", Data: []types.ExpressionArg{
			{DType: "exp", Exp: &types.Expression{}},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				AssignedSurveys: []types.AssignedSurvey{
					{SurveyKey: "test1", ValidFrom: 10, ValidUntil: 100},
					{SurveyKey: "test2", ValidFrom: 10, ValidUntil: 100},
				},
			},
		}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should throw an error about arg type")
			return
		}
	})
}

func TestEvalGetSurveyKeyAssignedUntil(t *testing.T) {
	t.Run("has survey assigned", func(t *testing.T) {
		exp := types.Expression{Name: "getSurveyKeyAssignedUntil", Data: []types.ExpressionArg{
			{DType: "str", Str: "test1"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				AssignedSurveys: []types.AssignedSurvey{
					{SurveyKey: "test1", ValidFrom: 10, ValidUntil: 100},
					{SurveyKey: "test2", ValidFrom: 10, ValidUntil: 100},
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if ret.(float64) != 100 {
			t.Errorf("unexpected value retrieved: %d", ret)
		}
	})

	t.Run("doesn't have the survey assigned", func(t *testing.T) {
		exp := types.Expression{Name: "getSurveyKeyAssignedUntil", Data: []types.ExpressionArg{
			{DType: "str", Str: "test1"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				AssignedSurveys: []types.AssignedSurvey{
					{SurveyKey: "test2", ValidFrom: 10, ValidUntil: 100},
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if ret.(float64) != -1 {
			t.Errorf("unexpected value retrieved: %d", ret)
		}
	})

	t.Run("missing argument", func(t *testing.T) {
		exp := types.Expression{Name: "getSurveyKeyAssignedUntil", Data: []types.ExpressionArg{}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				AssignedSurveys: []types.AssignedSurvey{
					{SurveyKey: "test1", ValidFrom: 10, ValidUntil: 100},
					{SurveyKey: "test2", ValidFrom: 10, ValidUntil: 100},
				},
			},
		}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should throw an error about missing arg")
			return
		}
	})

	t.Run("wrong argument", func(t *testing.T) {
		exp := types.Expression{Name: "getSurveyKeyAssignedUntil", Data: []types.ExpressionArg{
			{DType: "exp", Exp: &types.Expression{}},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				AssignedSurveys: []types.AssignedSurvey{
					{SurveyKey: "test1", ValidFrom: 10, ValidUntil: 100},
					{SurveyKey: "test2", ValidFrom: 10, ValidUntil: 100},
				},
			},
		}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should throw an error about arg type")
			return
		}
	})
}

func TestEvalHasParticipantFlag(t *testing.T) {
	t.Run("participant hasn't got any participant flags (empty / nil)", func(t *testing.T) {
		exp := types.Expression{Name: "hasParticipantFlag", Data: []types.ExpressionArg{
			{DType: "str", Str: "key1"},
			{DType: "str", Str: "value1"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if ret.(bool) {
			t.Error("should be false")
		}
	})

	t.Run("participant has other participant flags, but this key is missing", func(t *testing.T) {
		exp := types.Expression{Name: "hasParticipantFlag", Data: []types.ExpressionArg{
			{DType: "str", Str: "key1"},
			{DType: "str", Str: "value1"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				Flags: map[string]string{
					"key2": "value1",
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if ret.(bool) {
			t.Error("should be false")
		}
	})

	t.Run("participant has correct participant flag's key, but value is different", func(t *testing.T) {
		exp := types.Expression{Name: "hasParticipantFlag", Data: []types.ExpressionArg{
			{DType: "str", Str: "key1"},
			{DType: "str", Str: "value1"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				Flags: map[string]string{
					"key1": "value2",
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if ret.(bool) {
			t.Error("should be false")
		}
	})

	t.Run("participant has correct participant flag's key and value is same", func(t *testing.T) {
		exp := types.Expression{Name: "hasParticipantFlag", Data: []types.ExpressionArg{
			{DType: "str", Str: "key1"},
			{DType: "str", Str: "value1"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				Flags: map[string]string{
					"key1": "value1",
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if !ret.(bool) {
			t.Error("should be true")
		}
	})

	t.Run("missing arguments", func(t *testing.T) {
		exp := types.Expression{Name: "hasParticipantFlag", Data: []types.ExpressionArg{}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				Flags: map[string]string{
					"key1": "value1",
				},
			},
		}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should throw error")
			return
		}
	})

	t.Run("using num at 1st argument (expressions allowed, should return string)", func(t *testing.T) {
		exp := types.Expression{Name: "hasParticipantFlag", Data: []types.ExpressionArg{
			{DType: "num", Num: 22},
			{DType: "str", Str: "value1"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				Flags: map[string]string{
					"key1": "value1",
				},
			},
		}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should throw error")
			return
		}
	})

	t.Run("missing arguments", func(t *testing.T) {
		exp := types.Expression{Name: "hasParticipantFlag", Data: []types.ExpressionArg{
			{DType: "str", Str: "key1"},
			{DType: "num", Num: 22},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				Flags: map[string]string{
					"key1": "value1",
				},
			},
		}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should throw error")
			return
		}
	})
}

func TestEvalHasParticipantFlagKey(t *testing.T) {
	t.Run("participant hasn't got any participant flags (empty / nil)", func(t *testing.T) {
		exp := types.Expression{Name: "hasParticipantFlagKey", Data: []types.ExpressionArg{
			{DType: "str", Str: "key1"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if ret.(bool) {
			t.Error("should be false")
		}
	})

	t.Run("participant has other key", func(t *testing.T) {
		exp := types.Expression{Name: "hasParticipantFlagKey", Data: []types.ExpressionArg{
			{DType: "str", Str: "key1"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				Flags: map[string]string{
					"key2": "1",
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if ret.(bool) {
			t.Error("should be false")
		}
	})

	t.Run("participant has correct key", func(t *testing.T) {
		exp := types.Expression{Name: "hasParticipantFlagKey", Data: []types.ExpressionArg{
			{DType: "str", Str: "key1"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				Flags: map[string]string{
					"key2": "1",
					"key1": "1",
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if !ret.(bool) {
			t.Error("should be true")
		}
	})
}

func TestEvalHasResponseKey(t *testing.T) {
	testEvalContext := EvalContext{
		Event: types.StudyEvent{
			Type: "SUBMIT",
			Response: types.SurveyResponse{
				Key: "weekly",
				Responses: []types.SurveyItemResponse{
					{
						Key: "weekly.Q1", Response: &types.ResponseItem{
							Key: "rg", Items: []*types.ResponseItem{
								{Key: "1", Value: "something"},
								{Key: "2"},
							}},
					},
					{
						Key: "weekly.Q2", Response: &types.ResponseItem{
							Key: "rg", Items: []*types.ResponseItem{
								{Key: "1", Value: "123.23", Dtype: "date"},
							}},
					},
				},
			},
		},
	}

	//
	t.Run("no survey item response found", func(t *testing.T) {
		exp := types.Expression{Name: "hasResponseKey", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q3"},
			{DType: "str", Str: "rg.1"},
		}}
		v, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if v.(bool) {
			t.Errorf("unexpected value: %b", v)
		}
	})

	t.Run("repsonse item in question missing", func(t *testing.T) {
		exp := types.Expression{Name: "hasResponseKey", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q1"},
			{DType: "str", Str: "rg.wrong"},
		}}
		v, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if v.(bool) {
			t.Errorf("unexpected value: %b", v)
		}
	})

	t.Run("repsonse item in partly there missing", func(t *testing.T) {
		exp := types.Expression{Name: "hasResponseKey", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q1"},
			{DType: "str", Str: "rg.1.1"},
		}}
		v, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if v.(bool) {
			t.Errorf("unexpected value: %b", v)
		}
	})

	t.Run("has key", func(t *testing.T) {
		exp := types.Expression{Name: "hasResponseKey", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q1"},
			{DType: "str", Str: "rg.2"},
		}}
		v, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !v.(bool) {
			t.Errorf("unexpected value: %b", v)
		}
	})
}
func TestEvalHasResponseKeyWithValue(t *testing.T) {
	testEvalContext := EvalContext{
		Event: types.StudyEvent{
			Type: "SUBMIT",
			Response: types.SurveyResponse{
				Key: "weekly",
				Responses: []types.SurveyItemResponse{
					{
						Key: "weekly.Q1", Response: &types.ResponseItem{
							Key: "rg", Items: []*types.ResponseItem{
								{Key: "1", Value: "something"},
								{Key: "2"},
							}},
					},
					{
						Key: "weekly.Q2", Response: &types.ResponseItem{
							Key: "rg", Items: []*types.ResponseItem{
								{Key: "1", Value: "123.23", Dtype: "date"},
							}},
					},
				},
			},
		},
	}

	//
	t.Run("no survey item response found", func(t *testing.T) {
		exp := types.Expression{Name: "hasResponseKeyWithValue", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q3"},
			{DType: "str", Str: "rg.1"},
			{DType: "str", Str: "something"},
		}}
		v, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if v.(bool) {
			t.Errorf("unexpected value: %b", v)
		}
	})

	t.Run("repsonse item in question missing", func(t *testing.T) {
		exp := types.Expression{Name: "hasResponseKeyWithValue", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q1"},
			{DType: "str", Str: "rg.wrong"},
			{DType: "str", Str: "something"},
		}}
		v, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if v.(bool) {
			t.Errorf("unexpected value: %b", v)
		}
	})

	t.Run("has empty value", func(t *testing.T) {
		exp := types.Expression{Name: "hasResponseKeyWithValue", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q1"},
			{DType: "str", Str: "rg.2"},
			{DType: "str", Str: "something"},
		}}
		v, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if v.(bool) {
			t.Errorf("unexpected value: %b", v)
		}
	})

	t.Run("normal", func(t *testing.T) {
		exp := types.Expression{Name: "hasResponseKeyWithValue", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q1"},
			{DType: "str", Str: "rg.1"},
			{DType: "str", Str: "something"},
		}}
		v, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !v.(bool) {
			t.Errorf("unexpected value: %b", v)
		}
	})
}
func TestEvalGetResponseValueAsNum(t *testing.T) {
	testEvalContext := EvalContext{
		Event: types.StudyEvent{
			Type: "SUBMIT",
			Response: types.SurveyResponse{
				Key: "weekly",
				Responses: []types.SurveyItemResponse{
					{
						Key: "weekly.Q1", Response: &types.ResponseItem{
							Key: "rg", Items: []*types.ResponseItem{
								{Key: "1", Value: "something"},
								{Key: "2"},
							}},
					},
					{
						Key: "weekly.Q2", Response: &types.ResponseItem{
							Key: "rg", Items: []*types.ResponseItem{
								{Key: "1", Value: "123.23", Dtype: "date"},
							}},
					},
				},
			},
		},
	}

	//
	t.Run("no survey item response found", func(t *testing.T) {
		exp := types.Expression{Name: "getResponseValueAsNum", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q3"},
			{DType: "str", Str: "rg.1"},
		}}
		_, err := ExpressionEval(exp, testEvalContext)
		if err == nil {
			t.Error("should return an error")
			return
		}
	})

	t.Run("repsonse item in question missing", func(t *testing.T) {
		exp := types.Expression{Name: "getResponseValueAsNum", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q1"},
			{DType: "str", Str: "rg.wrong"},
		}}
		_, err := ExpressionEval(exp, testEvalContext)
		if err == nil {
			t.Error("should return an error")
			return
		}
	})

	t.Run("has empty value", func(t *testing.T) {
		exp := types.Expression{Name: "getResponseValueAsNum", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q1"},
			{DType: "str", Str: "rg.2"},
		}}
		_, err := ExpressionEval(exp, testEvalContext)
		if err == nil {
			t.Error("should return an error")
			return
		}
	})

	t.Run("repsonse item's value is not a number", func(t *testing.T) {
		exp := types.Expression{Name: "getResponseValueAsNum", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q1"},
			{DType: "str", Str: "rg.1"},
		}}
		_, err := ExpressionEval(exp, testEvalContext)
		if err == nil {
			t.Error("should return an error")
			return
		}
	})

	t.Run("is number", func(t *testing.T) {
		exp := types.Expression{Name: "getResponseValueAsNum", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q2"},
			{DType: "str", Str: "rg.1"},
		}}
		v, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if v.(float64) != 123.23 {
			t.Errorf("unexpected value: %b", v)
		}
	})
}

func TestEvalCountResponseItems(t *testing.T) {
	testEvalContext := EvalContext{
		Event: types.StudyEvent{
			Type: "SUBMIT",
			Response: types.SurveyResponse{
				Key: "weekly",
				Responses: []types.SurveyItemResponse{
					{
						Key: "weekly.Q1", Response: &types.ResponseItem{
							Key: "rg", Items: []*types.ResponseItem{
								{Key: "mcg", Items: []*types.ResponseItem{
									{Key: "1"},
									{Key: "2"},
									{Key: "3"},
								}},
							}},
					},
					{
						Key: "weekly.Q2", Response: &types.ResponseItem{
							Key: "rg", Items: []*types.ResponseItem{
								{Key: "mcg", Items: []*types.ResponseItem{}},
							}},
					},
				},
			},
		},
	}

	//
	t.Run("no survey item response found", func(t *testing.T) {
		exp := types.Expression{Name: "countResponseItems", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q3"},
			{DType: "str", Str: "rg.mcg"},
		}}
		_, err := ExpressionEval(exp, testEvalContext)
		if err == nil {
			t.Error("should return an error")
			return
		}
	})

	t.Run("repsonse item in question missing", func(t *testing.T) {
		exp := types.Expression{Name: "countResponseItems", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q1"},
			{DType: "str", Str: "rg.wrong"},
		}}
		_, err := ExpressionEval(exp, testEvalContext)
		if err == nil {
			t.Error("should return an error")
			return
		}
	})

	t.Run("has empty value", func(t *testing.T) {
		exp := types.Expression{Name: "countResponseItems", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q2"},
			{DType: "str", Str: "rg.mcg"},
		}}
		v, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if v.(float64) != 0.0 {
			t.Errorf("unexpected value: %b", v)
		}
	})

	t.Run("has 3 values", func(t *testing.T) {
		exp := types.Expression{Name: "countResponseItems", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q1"},
			{DType: "str", Str: "rg.mcg"},
		}}
		v, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if v.(float64) != 3.0 {
			t.Errorf("unexpected value: %b", v)
		}
	})
}

func TestEvalGetResponseValueAsStr(t *testing.T) {
	testEvalContext := EvalContext{
		Event: types.StudyEvent{
			Type: "SUBMIT",
			Response: types.SurveyResponse{
				Key: "weekly",
				Responses: []types.SurveyItemResponse{
					{
						Key: "weekly.Q1", Response: &types.ResponseItem{
							Key: "rg", Items: []*types.ResponseItem{
								{Key: "1", Value: "something"},
								{Key: "2"},
							}},
					},
				},
			},
		},
	}

	//
	t.Run("no survey item response found", func(t *testing.T) {
		exp := types.Expression{Name: "getResponseValueAsStr", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q3"},
			{DType: "str", Str: "rg.1"},
		}}
		_, err := ExpressionEval(exp, testEvalContext)
		if err == nil {
			t.Error("should return an error")
			return
		}
	})

	t.Run("repsonse item in question missing", func(t *testing.T) {
		exp := types.Expression{Name: "getResponseValueAsStr", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q1"},
			{DType: "str", Str: "rg.wrong"},
		}}
		_, err := ExpressionEval(exp, testEvalContext)
		if err == nil {
			t.Error("should return an error")
			return
		}
	})

	t.Run("has empty value", func(t *testing.T) {
		exp := types.Expression{Name: "getResponseValueAsStr", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q1"},
			{DType: "str", Str: "rg.2"},
		}}
		v, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if v != "" {
			t.Errorf("unexpected value: %s instead of %s", v, "blank")
		}
	})

	t.Run("has value", func(t *testing.T) {
		exp := types.Expression{Name: "getResponseValueAsStr", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.Q1"},
			{DType: "str", Str: "rg.1"},
		}}
		v, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if v != "something" {
			t.Errorf("unexpected value: %s instead of %s", v, "something")
		}
	})
}

func TestMustGetStrValue(t *testing.T) {
	testEvalContext := EvalContext{}

	t.Run("not string value", func(t *testing.T) {
		_, err := testEvalContext.mustGetStrValue(types.ExpressionArg{
			Num:   0,
			DType: "num",
		})
		if err == nil {
			t.Error("should produce error")
		}
	})

	t.Run("string value", func(t *testing.T) {
		v, err := testEvalContext.mustGetStrValue(types.ExpressionArg{
			Str:   "hello",
			DType: "str",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if v != "hello" {
			t.Errorf("unexpected value: %s", v)
		}
	})
}

func TestEvalResponseHasOnlyKeysOtherThan(t *testing.T) {
	testEvalContext := EvalContext{
		Event: types.StudyEvent{
			Type: "SUBMIT",
			Response: types.SurveyResponse{
				Key:       "wwekly",
				Responses: []types.SurveyItemResponse{},
			},
		},
	}

	t.Run("no survey item response found", func(t *testing.T) {
		exp := types.Expression{Name: "responseHasOnlyKeysOtherThan", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.G1.Q1"},
			{DType: "str", Str: "rg.mcg"},
			{DType: "str", Str: "1"},
			{DType: "str", Str: "2"},
		}}
		testEvalContext.Event.Response.Responses = []types.SurveyItemResponse{
			{Key: "weekly.G1.Q2", Response: &types.ResponseItem{Key: "rg", Items: []*types.ResponseItem{{Key: "mcg", Items: []*types.ResponseItem{
				{Key: "0"},
			}}}}},
		}
		ret, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}

	})

	t.Run("with response item found, but no response parent group", func(t *testing.T) {
		exp := types.Expression{Name: "responseHasOnlyKeysOtherThan", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.G1.Q1"},
			{DType: "str", Str: "rg.mcg"},
			{DType: "str", Str: "1"},
			{DType: "str", Str: "2"},
		}}
		testEvalContext.Event.Response.Responses = []types.SurveyItemResponse{
			{Key: "weekly.G1.Q1", Response: &types.ResponseItem{Key: "rg", Items: []*types.ResponseItem{{Key: "scg", Items: []*types.ResponseItem{
				{Key: "0"},
			}}}}},
		}
		ret, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}

	})

	t.Run("response group does include at least one", func(t *testing.T) {
		exp := types.Expression{Name: "responseHasOnlyKeysOtherThan", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.G1.Q1"},
			{DType: "str", Str: "rg.mcg"},
			{DType: "str", Str: "1"},
			{DType: "str", Str: "2"},
		}}
		testEvalContext.Event.Response.Responses = []types.SurveyItemResponse{
			{Key: "weekly.G1.Q1", Response: &types.ResponseItem{Key: "rg", Items: []*types.ResponseItem{{Key: "mcg", Items: []*types.ResponseItem{
				{Key: "0"},
				{Key: "1"},
				{Key: "3"},
			}}}}},
		}
		ret, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}

	})

	t.Run("response group is empty", func(t *testing.T) {
		exp := types.Expression{Name: "responseHasOnlyKeysOtherThan", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.G1.Q1"},
			{DType: "str", Str: "rg.mcg"},
			{DType: "str", Str: "1"},
			{DType: "str", Str: "2"},
		}}
		testEvalContext.Event.Response.Responses = []types.SurveyItemResponse{
			{Key: "weekly.G1.Q1", Response: &types.ResponseItem{Key: "rg", Items: []*types.ResponseItem{{Key: "mcg", Items: []*types.ResponseItem{}}}}},
		}
		ret, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}

	})

	t.Run("response group includes all and other responses", func(t *testing.T) {
		exp := types.Expression{Name: "responseHasOnlyKeysOtherThan", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.G1.Q1"},
			{DType: "str", Str: "rg.mcg"},
			{DType: "str", Str: "1"},
			{DType: "str", Str: "2"},
		}}
		testEvalContext.Event.Response.Responses = []types.SurveyItemResponse{
			{Key: "weekly.G1.Q1", Response: &types.ResponseItem{Key: "rg", Items: []*types.ResponseItem{{Key: "mcg", Items: []*types.ResponseItem{
				{Key: "0"},
				{Key: "1"},
				{Key: "2"},
			}}}}},
		}
		ret, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}

	})

	t.Run("response group includes none of the options", func(t *testing.T) {
		exp := types.Expression{Name: "responseHasOnlyKeysOtherThan", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.G1.Q1"},
			{DType: "str", Str: "rg.mcg"},
			{DType: "str", Str: "1"},
			{DType: "str", Str: "2"},
		}}
		testEvalContext.Event.Response.Responses = []types.SurveyItemResponse{
			{Key: "weekly.G1.Q1", Response: &types.ResponseItem{Key: "rg", Items: []*types.ResponseItem{{Key: "mcg", Items: []*types.ResponseItem{
				{Key: "0"},
				{Key: "3"},
			}}}}},
		}
		ret, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
}

func TestEvalResponseHasKeysAny(t *testing.T) {
	testEvalContext := EvalContext{
		Event: types.StudyEvent{
			Type: "SUBMIT",
			Response: types.SurveyResponse{
				Key:       "wwekly",
				Responses: []types.SurveyItemResponse{},
			},
		},
	}
	t.Run("no survey item response found", func(t *testing.T) {
		exp := types.Expression{Name: "responseHasKeysAny", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.G1.Q1"},
			{DType: "str", Str: "rg.mcg"},
			{DType: "str", Str: "1"},
			{DType: "str", Str: "2"},
		}}
		testEvalContext.Event.Response.Responses = []types.SurveyItemResponse{
			{Key: "weekly.G1.Q2", Response: &types.ResponseItem{Key: "rg", Items: []*types.ResponseItem{{Key: "mcg", Items: []*types.ResponseItem{
				{Key: "0"},
			}}}}},
		}
		ret, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}

	})
	t.Run("with response item found, but no response parent group", func(t *testing.T) {
		exp := types.Expression{Name: "responseHasKeysAny", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.G1.Q1"},
			{DType: "str", Str: "rg.mcg"},
			{DType: "str", Str: "1"},
			{DType: "str", Str: "2"},
		}}
		testEvalContext.Event.Response.Responses = []types.SurveyItemResponse{
			{Key: "weekly.G1.Q1", Response: &types.ResponseItem{Key: "rg", Items: []*types.ResponseItem{{Key: "scg", Items: []*types.ResponseItem{
				{Key: "0"},
			}}}}},
		}
		ret, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}

	})

	t.Run("response group does not include any", func(t *testing.T) {
		exp := types.Expression{Name: "responseHasKeysAny", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.G1.Q1"},
			{DType: "str", Str: "rg.mcg"},
			{DType: "str", Str: "1"},
			{DType: "str", Str: "2"},
		}}
		testEvalContext.Event.Response.Responses = []types.SurveyItemResponse{
			{Key: "weekly.G1.Q1", Response: &types.ResponseItem{Key: "rg", Items: []*types.ResponseItem{{Key: "mcg", Items: []*types.ResponseItem{
				{Key: "0"},
				{Key: "3"},
			}}}}},
		}
		ret, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}

	})

	t.Run("response group includes all and other responses", func(t *testing.T) {
		exp := types.Expression{Name: "responseHasKeysAny", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.G1.Q1"},
			{DType: "str", Str: "rg.mcg"},
			{DType: "str", Str: "1"},
			{DType: "str", Str: "2"},
		}}
		testEvalContext.Event.Response.Responses = []types.SurveyItemResponse{
			{Key: "weekly.G1.Q1", Response: &types.ResponseItem{Key: "rg", Items: []*types.ResponseItem{{Key: "mcg", Items: []*types.ResponseItem{
				{Key: "0"},
				{Key: "1"},
				{Key: "2"},
			}}}}},
		}
		ret, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}

	})
	t.Run("response group includes only of the multiple options", func(t *testing.T) {
		exp := types.Expression{Name: "responseHasKeysAny", Data: []types.ExpressionArg{
			{DType: "str", Str: "weekly.G1.Q1"},
			{DType: "str", Str: "rg.mcg"},
			{DType: "str", Str: "1"},
			{DType: "str", Str: "2"},
		}}
		testEvalContext.Event.Response.Responses = []types.SurveyItemResponse{
			{Key: "weekly.G1.Q1", Response: &types.ResponseItem{Key: "rg", Items: []*types.ResponseItem{{Key: "mcg", Items: []*types.ResponseItem{
				{Key: "0"},
				{Key: "1"},
			}}}}},
		}
		ret, err := ExpressionEval(exp, testEvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}

	})

}

func TestEvalLastSubmissionDateOlderThan(t *testing.T) {
	t.Run("with not older", func(t *testing.T) {
		exp := types.Expression{Name: "lastSubmissionDateOlderThan", Data: []types.ExpressionArg{
			{DType: "exp", Exp: &types.Expression{Name: "timestampWithOffset", Data: []types.ExpressionArg{
				{DType: "num", Num: -10},
			}}},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				LastSubmissions: map[string]int64{
					"s1": time.Now().Unix() - 2,
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("with specific survey is older", func(t *testing.T) {
		exp := types.Expression{Name: "lastSubmissionDateOlderThan", Data: []types.ExpressionArg{
			{DType: "exp", Exp: &types.Expression{Name: "timestampWithOffset", Data: []types.ExpressionArg{
				{DType: "num", Num: -10},
			}}},
			{DType: "str", Str: "s2"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				LastSubmissions: map[string]int64{
					"s1": time.Now().Unix() - 2,
					"s2": time.Now().Unix() - 20,
				}},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("with only one type of survey is older", func(t *testing.T) {
		exp := types.Expression{Name: "lastSubmissionDateOlderThan", Data: []types.ExpressionArg{
			{DType: "num", Num: 10},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				LastSubmissions: map[string]int64{
					"s1": time.Now().Unix() - 2,
					"s2": time.Now().Unix() - 20,
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("with all types are older", func(t *testing.T) {
		exp := types.Expression{Name: "lastSubmissionDateOlderThan", Data: []types.ExpressionArg{
			{DType: "exp", Exp: &types.Expression{Name: "timestampWithOffset", Data: []types.ExpressionArg{
				{DType: "num", Num: -10},
			}}},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
				LastSubmissions: map[string]int64{
					"s1": time.Now().Unix() - 25,
					"s2": time.Now().Unix() - 20,
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
}

// Comparisons
func TestEvalEq(t *testing.T) {
	t.Run("for eq numbers", func(t *testing.T) {
		exp := types.Expression{Name: "eq", Data: []types.ExpressionArg{
			{DType: "num", Num: 23},
			{DType: "num", Num: 23},
		}}
		EvalContext := EvalContext{
			Event: types.StudyEvent{Type: "TIMER"},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("for not equal numbers", func(t *testing.T) {
		exp := types.Expression{Name: "eq", Data: []types.ExpressionArg{
			{DType: "num", Num: 13},
			{DType: "num", Num: 23},
		}}
		EvalContext := EvalContext{
			Event: types.StudyEvent{Type: "enter"},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("for equal strings", func(t *testing.T) {
		exp := types.Expression{Name: "eq", Data: []types.ExpressionArg{
			{DType: "str", Str: "enter"},
			{DType: "str", Str: "enter"},
		}}
		EvalContext := EvalContext{
			Event: types.StudyEvent{Type: "enter"},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("for not equal strings", func(t *testing.T) {
		exp := types.Expression{Name: "eq", Data: []types.ExpressionArg{
			{DType: "str", Str: "enter"},
			{DType: "str", Str: "time..."},
		}}
		EvalContext := EvalContext{
			Event: types.StudyEvent{Type: "enter"},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
}

func TestEvalLT(t *testing.T) {
	t.Run("2 < 2", func(t *testing.T) {
		exp := types.Expression{Name: "lt", Data: []types.ExpressionArg{
			{DType: "num", Num: 2},
			{DType: "num", Num: 2},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("2 < 1", func(t *testing.T) {
		exp := types.Expression{Name: "lt", Data: []types.ExpressionArg{
			{DType: "num", Num: 2},
			{DType: "num", Num: 1},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("1 < 2", func(t *testing.T) {
		exp := types.Expression{Name: "lt", Data: []types.ExpressionArg{
			{DType: "num", Num: 1},
			{DType: "num", Num: 2},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("a < b", func(t *testing.T) {
		exp := types.Expression{Name: "lt", Data: []types.ExpressionArg{
			{DType: "str", Str: "a"},
			{DType: "str", Str: "b"},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("b < b", func(t *testing.T) {
		exp := types.Expression{Name: "lt", Data: []types.ExpressionArg{
			{DType: "str", Str: "b"},
			{DType: "str", Str: "b"},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("b < a", func(t *testing.T) {
		exp := types.Expression{Name: "lt", Data: []types.ExpressionArg{
			{DType: "str", Str: "b"},
			{DType: "str", Str: "a"},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
}

func TestEvalLTE(t *testing.T) {
	t.Run("2 <= 2", func(t *testing.T) {
		exp := types.Expression{Name: "lte", Data: []types.ExpressionArg{
			{DType: "num", Num: 2},
			{DType: "num", Num: 2},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("2 <= 1", func(t *testing.T) {
		exp := types.Expression{Name: "lte", Data: []types.ExpressionArg{
			{DType: "num", Num: 2},
			{DType: "num", Num: 1},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("1 <= 2", func(t *testing.T) {
		exp := types.Expression{Name: "lte", Data: []types.ExpressionArg{
			{DType: "num", Num: 1},
			{DType: "num", Num: 2},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("a <= b", func(t *testing.T) {
		exp := types.Expression{Name: "lte", Data: []types.ExpressionArg{
			{DType: "str", Str: "a"},
			{DType: "str", Str: "b"},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("b <= b", func(t *testing.T) {
		exp := types.Expression{Name: "lte", Data: []types.ExpressionArg{
			{DType: "str", Str: "b"},
			{DType: "str", Str: "b"},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("b <= a", func(t *testing.T) {
		exp := types.Expression{Name: "lte", Data: []types.ExpressionArg{
			{DType: "str", Str: "b"},
			{DType: "str", Str: "a"},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
}

func TestEvalGT(t *testing.T) {
	t.Run("2 > 2", func(t *testing.T) {
		exp := types.Expression{Name: "gt", Data: []types.ExpressionArg{
			{DType: "num", Num: 2},
			{DType: "num", Num: 2},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("2 > 1", func(t *testing.T) {
		exp := types.Expression{Name: "gt", Data: []types.ExpressionArg{
			{DType: "num", Num: 2},
			{DType: "num", Num: 1},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("1 > 2", func(t *testing.T) {
		exp := types.Expression{Name: "gt", Data: []types.ExpressionArg{
			{DType: "num", Num: 1},
			{DType: "num", Num: 2},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("a > b", func(t *testing.T) {
		exp := types.Expression{Name: "gt", Data: []types.ExpressionArg{
			{DType: "str", Str: "a"},
			{DType: "str", Str: "b"},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("b > b", func(t *testing.T) {
		exp := types.Expression{Name: "gt", Data: []types.ExpressionArg{
			{DType: "str", Str: "b"},
			{DType: "str", Str: "b"},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("b > a", func(t *testing.T) {
		exp := types.Expression{Name: "gt", Data: []types.ExpressionArg{
			{DType: "str", Str: "b"},
			{DType: "str", Str: "a"},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
}

func TestEvalGTE(t *testing.T) {
	t.Run("2 >= 2", func(t *testing.T) {
		exp := types.Expression{Name: "gte", Data: []types.ExpressionArg{
			{DType: "num", Num: 2},
			{DType: "num", Num: 2},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("2 >= 1", func(t *testing.T) {
		exp := types.Expression{Name: "gte", Data: []types.ExpressionArg{
			{DType: "num", Num: 2},
			{DType: "num", Num: 1},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("1 >= 2", func(t *testing.T) {
		exp := types.Expression{Name: "gte", Data: []types.ExpressionArg{
			{DType: "num", Num: 1},
			{DType: "num", Num: 2},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("a >= b", func(t *testing.T) {
		exp := types.Expression{Name: "gte", Data: []types.ExpressionArg{
			{DType: "str", Str: "a"},
			{DType: "str", Str: "b"},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("b >= b", func(t *testing.T) {
		exp := types.Expression{Name: "gte", Data: []types.ExpressionArg{
			{DType: "str", Str: "b"},
			{DType: "str", Str: "b"},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("b >= a", func(t *testing.T) {
		exp := types.Expression{Name: "gte", Data: []types.ExpressionArg{
			{DType: "str", Str: "b"},
			{DType: "str", Str: "a"},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
}

// Logic operators
func TestEvalAND(t *testing.T) {
	t.Run("0 && 0 ", func(t *testing.T) {
		exp := types.Expression{Name: "and", Data: []types.ExpressionArg{
			{DType: "num", Num: 0},
			{DType: "num", Num: 0},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("1 && 0 ", func(t *testing.T) {
		exp := types.Expression{Name: "and", Data: []types.ExpressionArg{
			{DType: "num", Num: 1},
			{DType: "num", Num: 0},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("0 && 1 ", func(t *testing.T) {
		exp := types.Expression{Name: "and", Data: []types.ExpressionArg{
			{DType: "num", Num: 0},
			{DType: "num", Num: 1},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("1 && 1 ", func(t *testing.T) {
		exp := types.Expression{Name: "and", Data: []types.ExpressionArg{
			{DType: "num", Num: 1},
			{DType: "num", Num: 1},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
}

func TestEvalOR(t *testing.T) {
	t.Run("0 || 0 ", func(t *testing.T) {
		exp := types.Expression{Name: "or", Data: []types.ExpressionArg{
			{DType: "num", Num: 0},
			{DType: "num", Num: 0},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("1 || 0 ", func(t *testing.T) {
		exp := types.Expression{Name: "or", Data: []types.ExpressionArg{
			{DType: "num", Num: 1},
			{DType: "num", Num: 0},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("0 || 1 ", func(t *testing.T) {
		exp := types.Expression{Name: "or", Data: []types.ExpressionArg{
			{DType: "num", Num: 0},
			{DType: "num", Num: 1},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("1 || 1 ", func(t *testing.T) {
		exp := types.Expression{Name: "or", Data: []types.ExpressionArg{
			{DType: "num", Num: 1},
			{DType: "num", Num: 1},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
}

func TestEvalNOT(t *testing.T) {
	t.Run("0", func(t *testing.T) {
		exp := types.Expression{Name: "not", Data: []types.ExpressionArg{
			{DType: "num", Num: 0},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
	t.Run("1", func(t *testing.T) {
		exp := types.Expression{Name: "not", Data: []types.ExpressionArg{
			{DType: "num", Num: 1},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
}

func TestEvalSum(t *testing.T) {

	testAdd := func (expected float64, label string, values ...types.ExpressionArg) {
		
		t.Run(fmt.Sprintf("Sum %s", label), func(t *testing.T) {
			exp := types.Expression{Name: "sum", Data: values}
			EvalContext := EvalContext{}
			ret, err := ExpressionEval(exp, EvalContext)
			if err != nil {
				t.Errorf("unexpected error: %s", err.Error())
				return
			}
			resTS := ret.(float64)
			if resTS != expected {
				t.Errorf("unexpected value: %f - expected ca. %f", ret, expected)
			}
		})
	}

	argNum := func(v float64) types.ExpressionArg {
		return types.ExpressionArg{DType: "num", Num: v}
	}

	argBool := func(v bool) types.ExpressionArg {
		var vN  float64
		if v {
			vN = 1
		} else {
			vN = 0
		}
		return types.ExpressionArg{
			DType: "exp",
			Exp: &types.Expression{Name: "or", Data: []types.ExpressionArg{ argNum(vN), argNum(vN) }, },
		}
	}

	testAdd(1, "0 + 1", argNum(0), argNum(1))
	testAdd(2, "1 + 1", argNum(1), argNum(1) )
	testAdd( 1, "-1 + 2", argNum(-1), argNum(2))
	testAdd(3, "1+1+1", argNum(1), argNum(1), argNum(1))
	testAdd(2, "true + true", argBool(true), argBool(true))
	testAdd(0, "false + false", argBool(false), argBool(false))
	testAdd(1, "true + false", argBool(true), argBool(false))
	testAdd(1, "false + true", argBool(false), argBool(true))
	
}


func TestEvalNeg(t *testing.T) {

	testNeg := func (v1 float64,  expected float64) {
		t.Run(fmt.Sprintf("Negate %f", v1), func(t *testing.T) {
			exp := types.Expression{Name: "neg", Data: []types.ExpressionArg{
				{DType: "num", Num: v1},
			}}
			EvalContext := EvalContext{}
			ret, err := ExpressionEval(exp, EvalContext)
			if err != nil {
				t.Errorf("unexpected error: %s", err.Error())
				return
			}
			resTS := ret.(float64)
			if resTS != expected {
				t.Errorf("unexpected value: %f - expected ca. %f", ret, expected)
			}
		})
	}

	testNeg(0, 0)
	testNeg(1, -1)
	testNeg(-1, 1)
}


func TestEvalTimestampWithOffset(t *testing.T) {
	t.Run("T + 0", func(t *testing.T) {
		exp := types.Expression{Name: "timestampWithOffset", Data: []types.ExpressionArg{
			{DType: "num", Num: 0},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		resTS := int64(ret.(float64))
		if resTS > time.Now().Unix()+1 || resTS < time.Now().Unix()-1 {
			t.Errorf("unexpected value: %d - expected ca. %d", ret, time.Now().Unix()+0)
		}
	})

	t.Run("T + 10", func(t *testing.T) {
		exp := types.Expression{Name: "timestampWithOffset", Data: []types.ExpressionArg{
			{DType: "num", Num: 10},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		resTS := int64(ret.(float64))
		if resTS > time.Now().Unix()+11 || resTS < time.Now().Unix()+9 {
			t.Errorf("unexpected value: %d - expected ca. %d", ret, time.Now().Unix()+10)
		}
	})

	t.Run("T - 10", func(t *testing.T) {
		exp := types.Expression{Name: "timestampWithOffset", Data: []types.ExpressionArg{
			{DType: "num", Num: -10},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		resTS := int64(ret.(float64))
		if resTS < time.Now().Unix()-11 || resTS > time.Now().Unix()-9 {
			t.Errorf("unexpected value: %d - expected ca. %d", ret, time.Now().Unix()-10)
		}
	})

	t.Run("T + No num", func(t *testing.T) {
		exp := types.Expression{Name: "timestampWithOffset", Data: []types.ExpressionArg{
			{DType: "str", Str: "0"},
		}}
		EvalContext := EvalContext{}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Errorf("unexpected lack of error: parameter 1 was not num")
			return
		}
	})

	t.Run("R + 0", func(t *testing.T) {
		r := time.Now().Unix() - 31536000
		exp := types.Expression{Name: "timestampWithOffset", Data: []types.ExpressionArg{
			{DType: "num", Num: 0},
			{DType: "num", Num: float64(r)},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		resTS := int64(ret.(float64))
		if resTS > r+1 || resTS < r-1 {
			t.Errorf("unexpected value: %d - expected ca. %d", ret, r+0)
		}
	})

	t.Run("R + 10", func(t *testing.T) {
		r := time.Now().Unix() - 31536000
		exp := types.Expression{Name: "timestampWithOffset", Data: []types.ExpressionArg{
			{DType: "num", Num: 10},
			{DType: "num", Num: float64(r)},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		resTS := int64(ret.(float64))
		if resTS > r+11 || resTS < r+9 {
			t.Errorf("unexpected value: %d - expected ca. %d", ret, r+10)
		}
	})

	t.Run("R - 10", func(t *testing.T) {
		r := time.Now().Unix() - 31536000
		exp := types.Expression{Name: "timestampWithOffset", Data: []types.ExpressionArg{
			{DType: "num", Num: -10},
			{DType: "num", Num: float64(r)},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		resTS := int64(ret.(float64))
		if resTS > r-9 || resTS < r-11 {
			t.Errorf("unexpected value: %d - expected ca. %d", ret, r-10)
		}
	})

	t.Run("R + No num", func(t *testing.T) {
		r := time.Now().Unix() - 31536000
		exp := types.Expression{Name: "timestampWithOffset", Data: []types.ExpressionArg{
			{DType: "str", Str: "0"},
			{DType: "num", Num: float64(r)},
		}}
		EvalContext := EvalContext{}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Errorf("unexpected lack of error: parameter 1 was not num")
			return
		}
	})

	t.Run("No num + 10", func(t *testing.T) {
		exp := types.Expression{Name: "timestampWithOffset", Data: []types.ExpressionArg{
			{DType: "num", Num: 10},
			{DType: "str", Str: "1"},
		}}
		EvalContext := EvalContext{}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Errorf("unexpected lack of error: parameter 2 was not num")
			return
		}
	})

	t.Run("No num + No num", func(t *testing.T) {
		exp := types.Expression{Name: "timestampWithOffset", Data: []types.ExpressionArg{
			{DType: "str", Str: "0"},
			{DType: "str", Str: "1"},
		}}
		EvalContext := EvalContext{}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Errorf("unexpected lack of error: parameters 1 & 2 were not num")
			return
		}
	})

	t.Run("Valid Exp", func(t *testing.T) {
		exp := types.Expression{Name: "timestampWithOffset", Data: []types.ExpressionArg{
			{
				DType: "exp", Exp: &types.Expression{
					Name: "timestampWithOffset", Data: []types.ExpressionArg{
						{DType: "num", Num: -float64(time.Now().Unix())},
					}},
			}}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		resTS := int64(ret.(float64))
		if resTS-time.Now().Unix() > 1 {
			t.Errorf("unexpected value: %d, expected %d", resTS, time.Now().Unix())
		}
	})

	t.Run("Valid Exp + Valid Exp", func(t *testing.T) {
		exp := types.Expression{Name: "timestampWithOffset", Data: []types.ExpressionArg{
			{DType: "exp", Exp: &types.Expression{
				Name: "timestampWithOffset", Data: []types.ExpressionArg{
					{DType: "num", Num: -float64(time.Now().Unix())},
				}},
			},
			{DType: "exp", Exp: &types.Expression{
				Name: "timestampWithOffset", Data: []types.ExpressionArg{
					{DType: "num", Num: -float64(time.Now().Unix())},
				}},
			},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		resTS := int64(ret.(float64))
		if resTS > 1 {
			t.Errorf("unexpected value: %d, expected %d", resTS, 0)
		}
	})

	t.Run("Not Valid Exp + Valid Exp", func(t *testing.T) {
		exp := types.Expression{Name: "timestampWithOffset", Data: []types.ExpressionArg{
			{DType: "exp", Exp: &types.Expression{
				Name: "or", Data: []types.ExpressionArg{
					{DType: "num", Num: 1},
					{DType: "num", Num: 1},
				}},
			},
			{DType: "exp", Exp: &types.Expression{
				Name: "timestampWithOffset", Data: []types.ExpressionArg{
					{DType: "num", Num: -float64(time.Now().Unix())},
				}},
			},
		}}
		EvalContext := EvalContext{}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Errorf("unexpected lack of error")
			return
		}
	})

	t.Run("Valid Exp + Not Valid Exp", func(t *testing.T) {
		exp := types.Expression{Name: "timestampWithOffset", Data: []types.ExpressionArg{
			{DType: "exp", Exp: &types.Expression{
				Name: "timestampWithOffset", Data: []types.ExpressionArg{
					{DType: "num", Num: -float64(time.Now().Unix())},
				}},
			},
			{DType: "exp", Exp: &types.Expression{
				Name: "or", Data: []types.ExpressionArg{
					{DType: "num", Num: 1},
					{DType: "num", Num: 1},
				}},
			},
		}}
		EvalContext := EvalContext{}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Errorf("unexpected lack of error")
			return
		}
	})
}

func TestEvalGetISOWeekForTs(t *testing.T) {
	t.Run("wrong argument type", func(t *testing.T) {
		exp := types.Expression{Name: "getISOWeekForTs", Data: []types.ExpressionArg{
			{DType: "str", Str: "test"},
		}}
		EvalContext := EvalContext{}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should return type error")
			return
		}
	})
	t.Run("with number", func(t *testing.T) {
		exp := types.Expression{Name: "getISOWeekForTs", Data: []types.ExpressionArg{
			{DType: "num", Num: float64(time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local).Unix())},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		iw := ret.(float64)
		if iw != 1 {
			t.Errorf("unexpected value: %f", iw)
			return
		}
	})

	t.Run("with expression", func(t *testing.T) {
		exp := types.Expression{Name: "getISOWeekForTs", Data: []types.ExpressionArg{
			{DType: "exp", Exp: &types.Expression{
				Name: "timestampWithOffset", Data: []types.ExpressionArg{
					{DType: "num", Num: 0},
				},
			},
			},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		iw := ret.(float64)
		_, ciw := time.Now().ISOWeek()
		if iw != float64(ciw) {
			t.Errorf("unexpected value: %f (should be %d)", iw, ciw)
			return
		}
	})
}

func TestEvalGetTsForNextISOWeek(t *testing.T) {
	t.Run("wrong iso week type", func(t *testing.T) {
		exp := types.Expression{Name: "getTsForNextISOWeek", Data: []types.ExpressionArg{
			{DType: "str", Str: "test"},
		}}
		EvalContext := EvalContext{}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should return type error")
			return
		}
	})

	t.Run("with iso week not in range", func(t *testing.T) {
		exp := types.Expression{Name: "getTsForNextISOWeek", Data: []types.ExpressionArg{
			{DType: "num", Num: 0},
		}}
		EvalContext := EvalContext{}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should return range error")
			return
		}
	})

	t.Run("wrong reference type", func(t *testing.T) {
		exp := types.Expression{Name: "getTsForNextISOWeek", Data: []types.ExpressionArg{
			{DType: "num", Num: 3},
			{DType: "str", Str: "test"},
		}}
		EvalContext := EvalContext{}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should return type error")
			return
		}
	})

	t.Run("without reference", func(t *testing.T) {
		exp := types.Expression{Name: "getTsForNextISOWeek", Data: []types.ExpressionArg{
			{DType: "num", Num: 1},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		ts := ret.(float64)
		tsD := time.Unix(int64(ts), 0)
		refTs := time.Now().AddDate(1, 0, 0)
		// beginning of the year
		refTs = time.Date(refTs.Year(), 1, 1, 0, 0, 0, 0, time.Local)

		y_i, w_i := refTs.ISOWeek()
		y, w := tsD.ISOWeek()
		if y != y_i || w != w_i {
			t.Errorf("unexpected value: %d-%d, expected %d-%d", y, w, y_i, w_i)
		}
	})

	t.Run("with absolute reference", func(t *testing.T) {
		refTs := time.Date(2023, 9, 10, 0, 0, 0, 0, time.Local)
		exp := types.Expression{Name: "getTsForNextISOWeek", Data: []types.ExpressionArg{
			{DType: "num", Num: 1},
			{DType: "num", Num: float64(refTs.Unix())},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		ts := ret.(float64)
		tsD := time.Unix(int64(ts), 0)
		refTs = refTs.AddDate(1, 0, 0)
		// beginning of the year
		refTs = time.Date(refTs.Year(), 1, 1, 0, 0, 0, 0, time.Local)
		y_i, w_i := refTs.ISOWeek()
		y, w := tsD.ISOWeek()
		if y != y_i || w != w_i {
			t.Errorf("unexpected value: %d-%d, expected %d-%d", y, w, y_i, w_i)
		}

	})

	t.Run("with relative reference", func(t *testing.T) {
		exp := types.Expression{Name: "getTsForNextISOWeek", Data: []types.ExpressionArg{
			{DType: "num", Num: 1},
			{DType: "exp", Exp: &types.Expression{
				Name: "timestampWithOffset",
				Data: []types.ExpressionArg{
					{DType: "num", Num: 0},
				},
			}},
		}}
		EvalContext := EvalContext{}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		ts := ret.(float64)
		tsD := time.Unix(int64(ts), 0)
		refTs := time.Now().AddDate(1, 0, 0)
		// beginning of the year
		refTs = time.Date(refTs.Year(), 1, 1, 0, 0, 0, 0, time.Local)

		y_i, w_i := refTs.ISOWeek()
		y, w := tsD.ISOWeek()
		if y != y_i || w != w_i {
			t.Errorf("unexpected value: %d-%d, expected %d-%d", y, w, y_i, w_i)
		}
	})
}

func TestEvalHasMessageTypeAssigned(t *testing.T) {
	t.Run("participant has no messages", func(t *testing.T) {
		exp := types.Expression{Name: "hasMessageTypeAssigned", Data: []types.ExpressionArg{
			{DType: "str", Str: "testMessage"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				Messages: []types.ParticipantMessage{},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		resTS := ret.(bool)
		if resTS {
			t.Errorf("unexpected value: %v", ret)
		}
	})

	t.Run("participant has messages but none that are looked for", func(t *testing.T) {
		exp := types.Expression{Name: "hasMessageTypeAssigned", Data: []types.ExpressionArg{
			{DType: "str", Str: "testMessage"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				Messages: []types.ParticipantMessage{
					{Type: "testMessage2", ScheduledFor: 100},
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		resTS := ret.(bool)
		if resTS {
			t.Errorf("unexpected value: %v", ret)
		}
	})

	t.Run("participant has messages and one is the one looked for", func(t *testing.T) {
		exp := types.Expression{Name: "hasMessageTypeAssigned", Data: []types.ExpressionArg{
			{DType: "str", Str: "testMessage"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				Messages: []types.ParticipantMessage{
					{Type: "testMessage2", ScheduledFor: 100},
					{Type: "testMessage", ScheduledFor: 200},
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		resTS := ret.(bool)
		if !resTS {
			t.Errorf("unexpected value: %v", ret)
		}
	})
}

func TestEvalGenerateRandomNumber(t *testing.T) {
	t.Run("with invalid args", func(t *testing.T) {
		exp := types.Expression{Name: "generateRandomNumber", Data: []types.ExpressionArg{
			{DType: "str", Str: "wrong"},
			{DType: "str", Str: "wrong"},
		}}
		EvalContext := EvalContext{}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should return error")
			return
		}
	})

	t.Run("with not enough args", func(t *testing.T) {
		exp := types.Expression{Name: "generateRandomNumber", Data: []types.ExpressionArg{
			{DType: "num", Num: 10},
		}}
		EvalContext := EvalContext{}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should return error")
			return
		}
	})

	t.Run("with valid args", func(t *testing.T) {
		// repeat 100 times
		for i := 0; i < 100; i++ {
			exp := types.Expression{Name: "generateRandomNumber", Data: []types.ExpressionArg{
				{DType: "num", Num: 10},
				{DType: "num", Num: 20},
			}}
			EvalContext := EvalContext{}
			val, err := ExpressionEval(exp, EvalContext)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			// logger.Debug.Println(val.(float64))
			if val.(float64) < 10 || val.(float64) > 20 {
				t.Errorf("unexpected value: %v", val)
				return
			}
		}
	})
}

func TestEvalParseValueAsNum(t *testing.T) {
	testPState := types.ParticipantState{
		Flags: map[string]string{
			"testKey": "3",
		},
	}

	t.Run("attempt to parse incorrect string", func(t *testing.T) {
		exp := types.Expression{Name: "parseValueAsNum", Data: []types.ExpressionArg{
			{DType: "str", Str: "wrong"},
		}}
		EvalContext := EvalContext{
			ParticipantState: testPState,
		}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should return error")
			return
		}
	})

	t.Run("attempt to parse correct string", func(t *testing.T) {
		exp := types.Expression{Name: "parseValueAsNum", Data: []types.ExpressionArg{
			{DType: "str", Str: "15"},
		}}
		EvalContext := EvalContext{
			ParticipantState: testPState,
		}
		res, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if res != 15.0 {
			t.Errorf("unexpected value: %v", res)
			return
		}
	})

	t.Run("already a number", func(t *testing.T) {
		exp := types.Expression{Name: "parseValueAsNum", Data: []types.ExpressionArg{
			{DType: "num", Num: 65},
		}}
		EvalContext := EvalContext{
			ParticipantState: testPState,
		}
		res, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if res != 65.0 {
			t.Errorf("unexpected value: %v", res)
			return
		}
	})

	t.Run("expression that returns error", func(t *testing.T) {
		exp := types.Expression{Name: "parseValueAsNum", Data: []types.ExpressionArg{
			{DType: "exp", Exp: &types.Expression{Name: "wrong", Data: []types.ExpressionArg{}}},
		}}
		EvalContext := EvalContext{
			ParticipantState: testPState,
		}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Errorf("should return an error: %v", err)
			return
		}
	})

	t.Run("expression that returns number", func(t *testing.T) {
		exp := types.Expression{Name: "parseValueAsNum", Data: []types.ExpressionArg{
			{DType: "exp", Exp: &types.Expression{Name: "timestampWithOffset", Data: []types.ExpressionArg{
				{DType: "num", Num: -10},
			}}},
		}}
		EvalContext := EvalContext{
			ParticipantState: testPState,
		}
		_, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
	})
	t.Run("expression that returns boolean", func(t *testing.T) {
		exp := types.Expression{Name: "parseValueAsNum", Data: []types.ExpressionArg{
			{DType: "exp", Exp: &types.Expression{Name: "gt", Data: []types.ExpressionArg{
				{DType: "num", Num: -10},
				{DType: "num", Num: 10},
			}}},
		}}
		EvalContext := EvalContext{
			ParticipantState: testPState,
		}
		_, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should return an error")
			return
		}
	})

	t.Run("expression that returns string", func(t *testing.T) {
		exp := types.Expression{Name: "parseValueAsNum", Data: []types.ExpressionArg{
			{DType: "exp", Exp: &types.Expression{Name: "getParticipantFlagValue", Data: []types.ExpressionArg{
				{DType: "str", Str: "testKey"},
			}}},
		}}
		EvalContext := EvalContext{
			ParticipantState: testPState,
		}
		res, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if res != 3.0 {
			t.Errorf("unexpected value: %v", res)
			return
		}
	})
}

func TestEvalGetMessageNextTime(t *testing.T) {
	t.Run("participant has no messages", func(t *testing.T) {
		exp := types.Expression{Name: "getMessageNextTime", Data: []types.ExpressionArg{
			{DType: "str", Str: "testMessage"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should return error")
			return
		}
		resTS := ret.(int64)
		if resTS != 0 {
			t.Errorf("unexpected value: %d", ret)
		}
	})

	t.Run("participant has messages but none that are looked for", func(t *testing.T) {
		exp := types.Expression{Name: "getMessageNextTime", Data: []types.ExpressionArg{
			{DType: "str", Str: "testMessage"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				Messages: []types.ParticipantMessage{
					{Type: "testMessage2", ScheduledFor: 100},
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err == nil {
			t.Error("should return error")
			return
		}
		resTS := ret.(int64)
		if resTS != 0 {
			t.Errorf("unexpected value: %d", ret)
		}
	})

	t.Run("participant has messages and one is the one looked for", func(t *testing.T) {
		exp := types.Expression{Name: "getMessageNextTime", Data: []types.ExpressionArg{
			{DType: "str", Str: "testMessage"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				Messages: []types.ParticipantMessage{
					{Type: "testMessage2", ScheduledFor: 50},
					{Type: "testMessage", ScheduledFor: 100},
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		resTS := ret.(int64)
		if resTS != 100 {
			t.Errorf("unexpected value: %d", ret)
		}
	})

	t.Run("participant has messages and two from the specified type", func(t *testing.T) {
		exp := types.Expression{Name: "getMessageNextTime", Data: []types.ExpressionArg{
			{DType: "str", Str: "testMessage"},
		}}
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				Messages: []types.ParticipantMessage{
					{Type: "testMessage1", ScheduledFor: 100},
					{Type: "testMessage", ScheduledFor: 200},
					{Type: "testMessage", ScheduledFor: 400},
				},
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		resTS := ret.(int64)
		if resTS != 200 {
			t.Errorf("unexpected value: %d", ret)
		}
	})
}

func TestNow(t *testing.T) {
	t.Run("testing now", func(t *testing.T) {
		cur := time.Now()
		now := Now()
		if(cur.Sub(now).Abs() > time.Microsecond) {
			t.Errorf("Current time is more than 1 microsecond")
		}
	})

	t.Run("testing change time", func(t *testing.T) {
		cur := time.Unix(1730419200, 0)
		Now = func() time.Time {
			return cur
		}
		now := Now()
		if(cur.Sub(now).Abs() > 0) {
			t.Errorf("Current time is not the time set %s got %s", cur, now)
		}
		Now = time.Now // resetting to current time
	})

	t.Run("testing change time on timestampWithOffset", func(t *testing.T) {
		curTS := int64(1730419200)
		cur := time.Unix(curTS, 0)
		Now = func() time.Time {
			return cur
		}
		exp := types.Expression{Name: "timestampWithOffset", Data: []types.ExpressionArg{
				{DType: "num", Num: -10},
			},
		}
		
		EvalContext := EvalContext{
			ParticipantState: types.ParticipantState{
				StudyStatus: types.PARTICIPANT_STUDY_STATUS_ACTIVE,
			},
		}
		ret, err := ExpressionEval(exp, EvalContext)
		if err != nil {
			t.Error(err)
		} 
		resTS := int64(ret.(float64))
		expTS := curTS - 10
		if( resTS != expTS) {
			t.Errorf("Unexpected timestamp got %d, expecting %d", resTS, expTS)
		}
		Now = time.Now // resetting to current time
	})
}
