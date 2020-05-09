package main

import (
	"context"

	"github.com/influenzanet/study-service/api"
	"github.com/influenzanet/study-service/models"
	"github.com/influenzanet/study-service/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *studyServiceServer) CreateNewStudy(ctx context.Context, req *api.NewStudyRequest) (*api.Study, error) {
	if req == nil || utils.IsTokenEmpty(req.Token) || req.Study == nil {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	if !utils.CheckIfAnyRolesInToken(req.Token, []string{"RESEARCHER", "ADMIN"}) {
		return nil, status.Error(codes.Unauthenticated, "not authorized to create a study")
	}

	study := models.StudyFromAPI(req.Study)
	study.Members = []models.StudyMember{
		{
			Role:     "owner",
			UserID:   req.Token.Id,
			UserName: utils.GetUsernameFromToken(req.Token),
		},
	}

	cStudy, err := createStudyInDB(req.Token.InstanceId, study)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return cStudy.ToAPI(), nil
}

func (s *studyServiceServer) SaveSurveyToStudy(ctx context.Context, req *api.AddSurveyReq) (*api.Survey, error) {
	if req == nil || utils.IsTokenEmpty(req.Token) || req.StudyKey == "" || req.Survey == nil {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	members, err := getStudyMembers(req.Token.InstanceId, req.StudyKey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !utils.CheckIfMember(req.Token.Id, members, []string{"maintainer", "owner"}) {
		return nil, status.Error(codes.Unauthenticated, "not authorized to access this study")
	}

	newSurvey := models.SurveyFromAPI(req.Survey)
	createdSurvey, err := saveSurveyToDB(req.Token.InstanceId, req.StudyKey, newSurvey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return createdSurvey.ToAPI(), nil
}

func (s *studyServiceServer) RemoveSurveyFromStudy(ctx context.Context, req *api.SurveyReferenceRequest) (*api.Status, error) {
	if req == nil || utils.IsTokenEmpty(req.Token) || req.StudyKey == "" || req.SurveyKey == "" {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	members, err := getStudyMembers(req.Token.InstanceId, req.StudyKey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !utils.CheckIfMember(req.Token.Id, members, []string{"maintainer", "owner"}) {
		return nil, status.Error(codes.Unauthenticated, "not authorized to access this study")
	}
	err = removeSurveyFromStudyDB(req.Token.InstanceId, req.StudyKey, req.SurveyKey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.Status{
		Msg: "survey removed",
	}, nil
}

func (s *studyServiceServer) GetStudySurveyInfos(ctx context.Context, req *api.StudyReferenceReq) (*api.SurveyInfoResp, error) {
	if req == nil || utils.IsTokenEmpty(req.Token) || req.StudyKey == "" {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}
	surveys, err := findAllSurveyDefsForStudyDB(req.Token.InstanceId, req.StudyKey, true)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	infos := make([]*api.SurveyInfoResp_SurveyInfo, len(surveys))
	for i, s := range surveys {
		apiS := s.ToAPI()
		infos[i] = &api.SurveyInfoResp_SurveyInfo{
			Key:         s.Current.SurveyDefinition.Key,
			Name:        apiS.Name,
			Description: apiS.Description,
		}
	}

	resp := api.SurveyInfoResp{
		Infos: infos,
	}
	return &resp, nil
}