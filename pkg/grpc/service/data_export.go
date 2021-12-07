package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/coneno/logger"
	"github.com/influenzanet/go-utils/pkg/api_types"
	"github.com/influenzanet/go-utils/pkg/constants"
	"github.com/influenzanet/go-utils/pkg/token_checks"
	"github.com/influenzanet/study-service/pkg/api"
	"github.com/influenzanet/study-service/pkg/exporter"
	"github.com/influenzanet/study-service/pkg/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	loggingAPI "github.com/influenzanet/logging-service/pkg/api"
)

const chunkSize = 64 * 1024 // 64 KiB

type ResponseFormat int

const (
	FLAT_JSON       ResponseFormat = iota
	WIDE_FORMAT_CSV ResponseFormat = iota
	LONG_FORMAT_CSV ResponseFormat = iota
)

func (s *studyServiceServer) GetStudyResponseStatistics(ctx context.Context, req *api.SurveyResponseQuery) (*api.StudyResponseStatistics, error) {
	if req == nil || token_checks.IsTokenEmpty(req.Token) || req.StudyKey == "" {
		return nil, s.missingArgumentError()
	}

	if !(token_checks.CheckRoleInToken(req.Token, constants.USER_ROLE_ADMIN) &&
		token_checks.CheckRoleInToken(req.Token, constants.USER_ROLE_RESEARCHER)) {
		err := s.HasRoleInStudy(req.Token.InstanceId, req.StudyKey, req.Token.Id, []string{
			types.STUDY_ROLE_OWNER,
			types.STUDY_ROLE_MAINTAINER,
			"analyst"})
		if err != nil {
			s.SaveLogEvent(req.Token.InstanceId, req.Token.Id, loggingAPI.LogEventType_SECURITY, constants.LOG_EVENT_DOWNLOAD_RESPONSES, "Statistics: permission denied for "+req.StudyKey)
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	keys, err := s.studyDBservice.GetSurveyResponseKeys(req.Token.InstanceId, req.StudyKey, req.From, req.Until)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &api.StudyResponseStatistics{
		SurveyResponseCounts: map[string]int64{},
	}

	for _, k := range keys {
		count, err := s.studyDBservice.CountSurveyResponsesByKey(req.Token.InstanceId, req.StudyKey, k, req.From, req.Until)
		if err != nil {
			logger.Error.Printf("unexpected error: %v", err)
			continue
		}
		resp.SurveyResponseCounts[k] = count
	}
	s.SaveLogEvent(req.Token.InstanceId, req.Token.Id, loggingAPI.LogEventType_LOG, constants.LOG_EVENT_DOWNLOAD_RESPONSES, "statistics: "+req.StudyKey)
	return resp, nil
}

func (s *studyServiceServer) StreamStudyResponses(req *api.SurveyResponseQuery, stream api.StudyServiceApi_StreamStudyResponsesServer) error {
	if req == nil || token_checks.IsTokenEmpty(req.Token) || req.StudyKey == "" {
		return s.missingArgumentError()
	}

	if !(token_checks.CheckRoleInToken(req.Token, constants.USER_ROLE_ADMIN) &&
		token_checks.CheckRoleInToken(req.Token, constants.USER_ROLE_RESEARCHER)) {
		err := s.HasRoleInStudy(req.Token.InstanceId, req.StudyKey, req.Token.Id, []string{
			types.STUDY_ROLE_OWNER,
			types.STUDY_ROLE_MAINTAINER,
		})
		if err != nil {
			s.SaveLogEvent(req.Token.InstanceId, req.Token.Id, loggingAPI.LogEventType_SECURITY, constants.LOG_EVENT_DOWNLOAD_RESPONSES, "permission denied for "+req.StudyKey)
			return status.Error(codes.Internal, err.Error())
		}
	}

	sendResponseOverGrpc := func(instanceID string, studyKey string, response types.SurveyResponse, args ...interface{}) error {
		if len(args) != 1 {
			return errors.New("StreamStudyResponses callback: unexpected number of args")
		}
		stream, ok := args[0].(api.StudyServiceApi_StreamStudyResponsesServer)
		if !ok {
			return errors.New(("StreamStudyResponses callback: can't parse stream"))
		}

		if err := stream.Send(response.ToAPI()); err != nil {
			return err
		}
		return nil
	}

	err := s.studyDBservice.PerformActionForSurveyResponses(req.Token.InstanceId, req.StudyKey, req.SurveyKey, req.From, req.Until,
		sendResponseOverGrpc, stream)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	s.SaveLogEvent(req.Token.InstanceId, req.Token.Id, loggingAPI.LogEventType_LOG, constants.LOG_EVENT_DOWNLOAD_RESPONSES, req.StudyKey)
	return nil
}

func (s *studyServiceServer) HasAccessToDownload(t *api_types.TokenInfos, studyKey string) error {
	// Admin and service account has access to studies:
	if token_checks.CheckIfAnyRolesInToken(t, []string{
		constants.USER_ROLE_ADMIN,
		constants.USER_ROLE_SERVICE_ACCOUNT,
	}) {
		return nil
	}

	// Should have at least researcher access:
	if !token_checks.CheckRoleInToken(t, constants.USER_ROLE_RESEARCHER) {
		return fmt.Errorf("unexpected roles %v", t.Payload["roles"])
	}

	// Should have access to the specific study as well
	err := s.HasRoleInStudy(t.InstanceId, studyKey, t.Id, []string{
		types.STUDY_ROLE_OWNER,
		types.STUDY_ROLE_MAINTAINER,
	})
	return err
}

// TODO: Test GetResponsesFlatJSON
func (s *studyServiceServer) GetResponsesFlatJSON(req *api.ResponseExportQuery, stream api.StudyServiceApi_GetResponsesFlatJSONServer) error {
	buf, err := s.getResponseExportBuffer(req, FLAT_JSON)

	if err != nil {
		return err
	}

	return StreamFile(stream, buf)
}

// TODO: Test GetResponsesWideFormatCSV
func (s *studyServiceServer) GetResponsesWideFormatCSV(req *api.ResponseExportQuery, stream api.StudyServiceApi_GetResponsesWideFormatCSVServer) error {
	buf, err := s.getResponseExportBuffer(req, WIDE_FORMAT_CSV)

	if err != nil {
		return err
	}

	return StreamFile(stream, buf)
}

// TODO: Test GetResponsesLongFormatCSV
func (s *studyServiceServer) GetResponsesLongFormatCSV(req *api.ResponseExportQuery, stream api.StudyServiceApi_GetResponsesLongFormatCSVServer) error {
	buf, err := s.getResponseExportBuffer(req, LONG_FORMAT_CSV)

	if err != nil {
		return err
	}

	return StreamFile(stream, buf)
}

// TODO: Test GetSurveyInfoPreviewCSV
func (s *studyServiceServer) GetSurveyInfoPreviewCSV(req *api.SurveyInfoExportQuery, stream api.StudyServiceApi_GetSurveyInfoPreviewCSVServer) error {
	responseExporter, err := s.getResponseExporterSurveyInfo(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = responseExporter.GetSurveyInfoCSV(buf)
	if err != nil {
		logger.Info.Println(err)
		return status.Error(codes.Internal, err.Error())
	}

	return StreamFile(stream, buf)
}

// TODO: Test GetSurveyInfoPreview
func (s *studyServiceServer) GetSurveyInfoPreview(ctx context.Context, req *api.SurveyInfoExportQuery) (*api.SurveyInfoExport, error) {
	responseExporter, err := s.getResponseExporterSurveyInfo(req)
	if err != nil {
		return nil, err
	}

	resp := responseExporter.GetSurveyVersionDefs()
	versions := make([]*api.SurveyVersionPreview, len(resp))
	for i, v := range resp {
		versions[i] = v.ToAPI()
	}

	return &api.SurveyInfoExport{
		Key:      req.SurveyKey,
		Versions: versions,
	}, nil
}

type StreamObj interface {
	Send(*api.Chunk) error
}

func StreamFile(stream StreamObj, buf *bytes.Buffer) error {
	chnk := &api.Chunk{}

	for currentByte := 0; currentByte < len(buf.Bytes()); currentByte += chunkSize {
		if currentByte+chunkSize > len(buf.Bytes()) {
			chnk.Chunk = buf.Bytes()[currentByte:len(buf.Bytes())]
		} else {
			chnk.Chunk = buf.Bytes()[currentByte : currentByte+chunkSize]
		}

		if err := stream.Send(chnk); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

func (s *studyServiceServer) getResponseExportBuffer(req *api.ResponseExportQuery, fmt ResponseFormat) (*bytes.Buffer, error) {
	responseExporter, err := s.getResponseExporterResponseExport(req)
	if err != nil {
		return nil, err
	}

	// Download responses
	err = s.studyDBservice.PerformActionForSurveyResponses(
		req.Token.InstanceId, req.StudyKey, req.SurveyKey,
		req.From, req.Until, func(instanceID, studyKey string, response types.SurveyResponse, args ...interface{}) error {
			if len(args) != 1 {
				return errors.New("[getResponseExportBuffer]: wrong DB method argument")
			}
			rExp, ok := args[0].(*exporter.ResponseExporter)
			if !ok {
				return errors.New("[getResponseExportBuffer]: wrong DB method argument")
			}
			return rExp.AddResponse(response.ToAPI())
		},
		responseExporter,
	)
	if err != nil {
		logger.Info.Print(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	buf := new(bytes.Buffer)
	includeMeta := &exporter.IncludeMeta{
		Postion:        req.IncludeMeta.Position,
		ItemVersion:    req.IncludeMeta.ItemVersion,
		InitTimes:      req.IncludeMeta.InitTimes,
		ResponsedTimes: req.IncludeMeta.ResponsedTimes,
		DisplayedTimes: req.IncludeMeta.DisplayedTimes,
	}

	switch fmt {
	case FLAT_JSON:
		err = responseExporter.GetResponsesJSON(buf, includeMeta)
	case WIDE_FORMAT_CSV:
		err = responseExporter.GetResponsesCSV(buf, includeMeta)
	case LONG_FORMAT_CSV:
		err = responseExporter.GetResponsesLongFormatCSV(buf, includeMeta)
	default:
		return nil, status.Error(codes.Internal, errors.New("[getResponseExportBuffer]: wrong response format").Error())
	}

	if err != nil {
		logger.Info.Println(err)
		return nil, err
	}
	return buf, nil
}

func (s *studyServiceServer) getResponseExporterSurveyInfo(req *api.SurveyInfoExportQuery) (*exporter.ResponseExporter, error) {
	if req == nil {
		return nil, s.missingArgumentError()
	}
	return s.getResponseExporter(req.Token, req.StudyKey, req.SurveyKey, req.PreviewLanguage, req.ShortQuestionKeys, "ignored")
}

func (s *studyServiceServer) getResponseExporterResponseExport(req *api.ResponseExportQuery) (*exporter.ResponseExporter, error) {
	if req == nil {
		return nil, s.missingArgumentError()
	}
	return s.getResponseExporter(req.Token, req.StudyKey, req.SurveyKey, "ignored", req.ShortQuestionKeys, req.Separator)
}

func (s *studyServiceServer) getResponseExporter(
	token *api_types.TokenInfos,
	studyKey string,
	surveyKey string,
	previewLanguage string,
	shortQuestionKeys bool,
	separator string) (*exporter.ResponseExporter, error) {
	if token_checks.IsTokenEmpty(token) || studyKey == "" {
		return nil, s.missingArgumentError()
	}

	if err := s.HasAccessToDownload(token, studyKey); err != nil {
		s.SaveLogEvent(token.InstanceId, token.Id, loggingAPI.LogEventType_SECURITY, constants.LOG_EVENT_GET_SURVEY_DEF, "Permission denied for "+studyKey)
		return nil, status.Error(codes.Internal, err.Error())
	}

	surveyDef, err := s.studyDBservice.FindSurveyDef(token.InstanceId, studyKey, surveyKey)
	if err != nil {
		logger.Info.Println(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	responseExporter, err := exporter.NewResponseExporter(
		surveyDef.ToAPI(),
		previewLanguage,
		shortQuestionKeys,
		separator,
	)
	if err != nil {
		logger.Info.Println(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return responseExporter, nil
}

func (s *studyServiceServer) missingArgumentError() error {
	return status.Error(codes.InvalidArgument, "missing argument")
}
