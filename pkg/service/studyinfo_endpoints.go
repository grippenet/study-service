package service

import (
	"context"
	"log"

	"github.com/influenzanet/study-service/pkg/api"
	"github.com/influenzanet/study-service/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *studyServiceServer) GetStudiesForUser(ctx context.Context, req *api.GetStudiesForUserReq) (*api.Studies, error) {
	if req == nil || utils.IsTokenEmpty(req.Token) {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	// for every profile form the token
	profileIDs := []string{req.Token.ProfilId}
	profileIDs = append(profileIDs, req.Token.OtherProfileIds...)

	studies, err := s.studyDBservice.GetStudiesByStatus(req.Token.InstanceId, "", false)
	if err != nil {
		log.Printf("GetStudiesForUser.GetStudiesByStatus: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &api.Studies{Studies: []*api.Study{}}
	for _, study := range studies {

		for _, profileID := range profileIDs {
			// ParticipantID
			participantID, err := utils.ProfileIDtoParticipantID(profileID, s.StudyGlobalSecret, study.SecretKey)
			if err != nil {
				continue
			}

			pState, err := s.studyDBservice.FindParticipantState(req.Token.InstanceId, study.Key, participantID)
			if err != nil {
				// user not in the study
				continue
			}

			if pState.StudyStatus != "active" {
				continue
			}

			// at least one profile in the study:
			resp.Studies = append(resp.Studies, &api.Study{
				Key:    study.Key,
				Status: study.Status,
				Props:  study.Props.ToAPI(),
			})
			break
		}
	}

	return resp, nil
}

func (s *studyServiceServer) GetActiveStudies(ctx context.Context, req *api.TokenInfos) (*api.Studies, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (s *studyServiceServer) HasParticipantStateWithCondition(ctx context.Context, req *api.ProfilesWithConditionReq) (*api.AssignedSurveys, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}
