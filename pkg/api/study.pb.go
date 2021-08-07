// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.22.0-devel
// 	protoc        v3.11.4
// source: study_service/study.proto

package api

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Study struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                        string          `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`   // db id
	Key                       string          `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"` // user defined unique study identifier
	SecretKey                 string          `protobuf:"bytes,3,opt,name=secret_key,json=secretKey,proto3" json:"secret_key,omitempty"`
	Status                    string          `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	Props                     *Study_Props    `protobuf:"bytes,5,opt,name=props,proto3" json:"props,omitempty"`
	Rules                     []*Expression   `protobuf:"bytes,6,rep,name=rules,proto3" json:"rules,omitempty"`
	Members                   []*Study_Member `protobuf:"bytes,7,rep,name=members,proto3" json:"members,omitempty"`
	Stats                     *Study_Stats    `protobuf:"bytes,8,opt,name=stats,proto3" json:"stats,omitempty"`
	ParticipantFileUploadRule *Expression     `protobuf:"bytes,9,opt,name=participant_file_upload_rule,json=participantFileUploadRule,proto3" json:"participant_file_upload_rule,omitempty"`
}

func (x *Study) Reset() {
	*x = Study{}
	if protoimpl.UnsafeEnabled {
		mi := &file_study_service_study_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Study) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Study) ProtoMessage() {}

func (x *Study) ProtoReflect() protoreflect.Message {
	mi := &file_study_service_study_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Study.ProtoReflect.Descriptor instead.
func (*Study) Descriptor() ([]byte, []int) {
	return file_study_service_study_proto_rawDescGZIP(), []int{0}
}

func (x *Study) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Study) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Study) GetSecretKey() string {
	if x != nil {
		return x.SecretKey
	}
	return ""
}

func (x *Study) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Study) GetProps() *Study_Props {
	if x != nil {
		return x.Props
	}
	return nil
}

func (x *Study) GetRules() []*Expression {
	if x != nil {
		return x.Rules
	}
	return nil
}

func (x *Study) GetMembers() []*Study_Member {
	if x != nil {
		return x.Members
	}
	return nil
}

func (x *Study) GetStats() *Study_Stats {
	if x != nil {
		return x.Stats
	}
	return nil
}

func (x *Study) GetParticipantFileUploadRule() *Expression {
	if x != nil {
		return x.ParticipantFileUploadRule
	}
	return nil
}

type StudyForUser struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key        string       `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"` // user defined unique study identifier
	Status     string       `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Props      *Study_Props `protobuf:"bytes,3,opt,name=props,proto3" json:"props,omitempty"`
	Stats      *Study_Stats `protobuf:"bytes,4,opt,name=stats,proto3" json:"stats,omitempty"`
	ProfileIds []string     `protobuf:"bytes,5,rep,name=profile_ids,json=profileIds,proto3" json:"profile_ids,omitempty"`
}

func (x *StudyForUser) Reset() {
	*x = StudyForUser{}
	if protoimpl.UnsafeEnabled {
		mi := &file_study_service_study_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StudyForUser) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StudyForUser) ProtoMessage() {}

func (x *StudyForUser) ProtoReflect() protoreflect.Message {
	mi := &file_study_service_study_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StudyForUser.ProtoReflect.Descriptor instead.
func (*StudyForUser) Descriptor() ([]byte, []int) {
	return file_study_service_study_proto_rawDescGZIP(), []int{1}
}

func (x *StudyForUser) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *StudyForUser) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *StudyForUser) GetProps() *Study_Props {
	if x != nil {
		return x.Props
	}
	return nil
}

func (x *StudyForUser) GetStats() *Study_Stats {
	if x != nil {
		return x.Stats
	}
	return nil
}

func (x *StudyForUser) GetProfileIds() []string {
	if x != nil {
		return x.ProfileIds
	}
	return nil
}

type Tag struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Label []*LocalisedObject `protobuf:"bytes,1,rep,name=label,proto3" json:"label,omitempty"`
}

func (x *Tag) Reset() {
	*x = Tag{}
	if protoimpl.UnsafeEnabled {
		mi := &file_study_service_study_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tag) ProtoMessage() {}

func (x *Tag) ProtoReflect() protoreflect.Message {
	mi := &file_study_service_study_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tag.ProtoReflect.Descriptor instead.
func (*Tag) Descriptor() ([]byte, []int) {
	return file_study_service_study_proto_rawDescGZIP(), []int{2}
}

func (x *Tag) GetLabel() []*LocalisedObject {
	if x != nil {
		return x.Label
	}
	return nil
}

type AssignedSurvey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SurveyKey  string `protobuf:"bytes,1,opt,name=survey_key,json=surveyKey,proto3" json:"survey_key,omitempty"`
	ValidFrom  int64  `protobuf:"varint,2,opt,name=valid_from,json=validFrom,proto3" json:"valid_from,omitempty"`
	ValidUntil int64  `protobuf:"varint,3,opt,name=valid_until,json=validUntil,proto3" json:"valid_until,omitempty"`
	StudyKey   string `protobuf:"bytes,4,opt,name=study_key,json=studyKey,proto3" json:"study_key,omitempty"`
	Category   string `protobuf:"bytes,5,opt,name=category,proto3" json:"category,omitempty"`
	ProfileId  string `protobuf:"bytes,6,opt,name=profile_id,json=profileId,proto3" json:"profile_id,omitempty"`
}

func (x *AssignedSurvey) Reset() {
	*x = AssignedSurvey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_study_service_study_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssignedSurvey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssignedSurvey) ProtoMessage() {}

func (x *AssignedSurvey) ProtoReflect() protoreflect.Message {
	mi := &file_study_service_study_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssignedSurvey.ProtoReflect.Descriptor instead.
func (*AssignedSurvey) Descriptor() ([]byte, []int) {
	return file_study_service_study_proto_rawDescGZIP(), []int{3}
}

func (x *AssignedSurvey) GetSurveyKey() string {
	if x != nil {
		return x.SurveyKey
	}
	return ""
}

func (x *AssignedSurvey) GetValidFrom() int64 {
	if x != nil {
		return x.ValidFrom
	}
	return 0
}

func (x *AssignedSurvey) GetValidUntil() int64 {
	if x != nil {
		return x.ValidUntil
	}
	return 0
}

func (x *AssignedSurvey) GetStudyKey() string {
	if x != nil {
		return x.StudyKey
	}
	return ""
}

func (x *AssignedSurvey) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *AssignedSurvey) GetProfileId() string {
	if x != nil {
		return x.ProfileId
	}
	return ""
}

type SurveyInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StudyKey        string             `protobuf:"bytes,1,opt,name=study_key,json=studyKey,proto3" json:"study_key,omitempty"`
	SurveyKey       string             `protobuf:"bytes,2,opt,name=survey_key,json=surveyKey,proto3" json:"survey_key,omitempty"`
	Name            []*LocalisedObject `protobuf:"bytes,3,rep,name=name,proto3" json:"name,omitempty"`
	Description     []*LocalisedObject `protobuf:"bytes,4,rep,name=description,proto3" json:"description,omitempty"`
	TypicalDuration []*LocalisedObject `protobuf:"bytes,5,rep,name=typical_duration,json=typicalDuration,proto3" json:"typical_duration,omitempty"`
}

func (x *SurveyInfo) Reset() {
	*x = SurveyInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_study_service_study_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SurveyInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SurveyInfo) ProtoMessage() {}

func (x *SurveyInfo) ProtoReflect() protoreflect.Message {
	mi := &file_study_service_study_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SurveyInfo.ProtoReflect.Descriptor instead.
func (*SurveyInfo) Descriptor() ([]byte, []int) {
	return file_study_service_study_proto_rawDescGZIP(), []int{4}
}

func (x *SurveyInfo) GetStudyKey() string {
	if x != nil {
		return x.StudyKey
	}
	return ""
}

func (x *SurveyInfo) GetSurveyKey() string {
	if x != nil {
		return x.SurveyKey
	}
	return ""
}

func (x *SurveyInfo) GetName() []*LocalisedObject {
	if x != nil {
		return x.Name
	}
	return nil
}

func (x *SurveyInfo) GetDescription() []*LocalisedObject {
	if x != nil {
		return x.Description
	}
	return nil
}

func (x *SurveyInfo) GetTypicalDuration() []*LocalisedObject {
	if x != nil {
		return x.TypicalDuration
	}
	return nil
}

type AssignedSurveys struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Surveys     []*AssignedSurvey `protobuf:"bytes,1,rep,name=surveys,proto3" json:"surveys,omitempty"`
	SurveyInfos []*SurveyInfo     `protobuf:"bytes,2,rep,name=survey_infos,json=surveyInfos,proto3" json:"survey_infos,omitempty"`
}

func (x *AssignedSurveys) Reset() {
	*x = AssignedSurveys{}
	if protoimpl.UnsafeEnabled {
		mi := &file_study_service_study_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssignedSurveys) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssignedSurveys) ProtoMessage() {}

func (x *AssignedSurveys) ProtoReflect() protoreflect.Message {
	mi := &file_study_service_study_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssignedSurveys.ProtoReflect.Descriptor instead.
func (*AssignedSurveys) Descriptor() ([]byte, []int) {
	return file_study_service_study_proto_rawDescGZIP(), []int{5}
}

func (x *AssignedSurveys) GetSurveys() []*AssignedSurvey {
	if x != nil {
		return x.Surveys
	}
	return nil
}

func (x *AssignedSurveys) GetSurveyInfos() []*SurveyInfo {
	if x != nil {
		return x.SurveyInfos
	}
	return nil
}

type Study_Props struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name               []*LocalisedObject `protobuf:"bytes,1,rep,name=name,proto3" json:"name,omitempty"`
	Description        []*LocalisedObject `protobuf:"bytes,2,rep,name=description,proto3" json:"description,omitempty"`
	Tags               []*Tag             `protobuf:"bytes,3,rep,name=tags,proto3" json:"tags,omitempty"`
	StartDate          int64              `protobuf:"varint,4,opt,name=start_date,json=startDate,proto3" json:"start_date,omitempty"`
	EndDate            int64              `protobuf:"varint,5,opt,name=end_date,json=endDate,proto3" json:"end_date,omitempty"`
	SystemDefaultStudy bool               `protobuf:"varint,6,opt,name=system_default_study,json=systemDefaultStudy,proto3" json:"system_default_study,omitempty"`
}

func (x *Study_Props) Reset() {
	*x = Study_Props{}
	if protoimpl.UnsafeEnabled {
		mi := &file_study_service_study_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Study_Props) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Study_Props) ProtoMessage() {}

func (x *Study_Props) ProtoReflect() protoreflect.Message {
	mi := &file_study_service_study_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Study_Props.ProtoReflect.Descriptor instead.
func (*Study_Props) Descriptor() ([]byte, []int) {
	return file_study_service_study_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Study_Props) GetName() []*LocalisedObject {
	if x != nil {
		return x.Name
	}
	return nil
}

func (x *Study_Props) GetDescription() []*LocalisedObject {
	if x != nil {
		return x.Description
	}
	return nil
}

func (x *Study_Props) GetTags() []*Tag {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *Study_Props) GetStartDate() int64 {
	if x != nil {
		return x.StartDate
	}
	return 0
}

func (x *Study_Props) GetEndDate() int64 {
	if x != nil {
		return x.EndDate
	}
	return 0
}

func (x *Study_Props) GetSystemDefaultStudy() bool {
	if x != nil {
		return x.SystemDefaultStudy
	}
	return false
}

type Study_Member struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId   string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Role     string `protobuf:"bytes,2,opt,name=role,proto3" json:"role,omitempty"`
	Username string `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty"`
}

func (x *Study_Member) Reset() {
	*x = Study_Member{}
	if protoimpl.UnsafeEnabled {
		mi := &file_study_service_study_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Study_Member) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Study_Member) ProtoMessage() {}

func (x *Study_Member) ProtoReflect() protoreflect.Message {
	mi := &file_study_service_study_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Study_Member.ProtoReflect.Descriptor instead.
func (*Study_Member) Descriptor() ([]byte, []int) {
	return file_study_service_study_proto_rawDescGZIP(), []int{0, 1}
}

func (x *Study_Member) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Study_Member) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *Study_Member) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

type Study_Stats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ParticipantCount int64 `protobuf:"varint,1,opt,name=participant_count,json=participantCount,proto3" json:"participant_count,omitempty"`
	ResponseCount    int64 `protobuf:"varint,2,opt,name=response_count,json=responseCount,proto3" json:"response_count,omitempty"`
}

func (x *Study_Stats) Reset() {
	*x = Study_Stats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_study_service_study_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Study_Stats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Study_Stats) ProtoMessage() {}

func (x *Study_Stats) ProtoReflect() protoreflect.Message {
	mi := &file_study_service_study_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Study_Stats.ProtoReflect.Descriptor instead.
func (*Study_Stats) Descriptor() ([]byte, []int) {
	return file_study_service_study_proto_rawDescGZIP(), []int{0, 2}
}

func (x *Study_Stats) GetParticipantCount() int64 {
	if x != nil {
		return x.ParticipantCount
	}
	return 0
}

func (x *Study_Stats) GetResponseCount() int64 {
	if x != nil {
		return x.ResponseCount
	}
	return 0
}

var File_study_service_study_proto protoreflect.FileDescriptor

var file_study_service_study_proto_rawDesc = []byte{
	0x0a, 0x19, 0x73, 0x74, 0x75, 0x64, 0x79, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f,
	0x73, 0x74, 0x75, 0x64, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x69, 0x6e, 0x66,
	0x6c, 0x75, 0x65, 0x6e, 0x7a, 0x61, 0x6e, 0x65, 0x74, 0x2e, 0x73, 0x74, 0x75, 0x64, 0x79, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x1e, 0x73, 0x74, 0x75, 0x64, 0x79, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x65, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x73, 0x74, 0x75, 0x64, 0x79, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xb4, 0x07, 0x0a, 0x05, 0x53, 0x74, 0x75, 0x64, 0x79, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x3d, 0x0a, 0x05, 0x70, 0x72, 0x6f, 0x70, 0x73, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x69, 0x6e, 0x66, 0x6c, 0x75, 0x65, 0x6e, 0x7a,
	0x61, 0x6e, 0x65, 0x74, 0x2e, 0x73, 0x74, 0x75, 0x64, 0x79, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x53, 0x74, 0x75, 0x64, 0x79, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x52, 0x05,
	0x70, 0x72, 0x6f, 0x70, 0x73, 0x12, 0x3c, 0x0a, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x06,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x69, 0x6e, 0x66, 0x6c, 0x75, 0x65, 0x6e, 0x7a, 0x61,
	0x6e, 0x65, 0x74, 0x2e, 0x73, 0x74, 0x75, 0x64, 0x79, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x45, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x72, 0x75,
	0x6c, 0x65, 0x73, 0x12, 0x42, 0x0a, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18, 0x07,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x69, 0x6e, 0x66, 0x6c, 0x75, 0x65, 0x6e, 0x7a, 0x61,
	0x6e, 0x65, 0x74, 0x2e, 0x73, 0x74, 0x75, 0x64, 0x79, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x53, 0x74, 0x75, 0x64, 0x79, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x07,
	0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x12, 0x3d, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x73,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x69, 0x6e, 0x66, 0x6c, 0x75, 0x65, 0x6e,
	0x7a, 0x61, 0x6e, 0x65, 0x74, 0x2e, 0x73, 0x74, 0x75, 0x64, 0x79, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x53, 0x74, 0x75, 0x64, 0x79, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52,
	0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x12, 0x67, 0x0a, 0x1c, 0x70, 0x61, 0x72, 0x74, 0x69, 0x63,
	0x69, 0x70, 0x61, 0x6e, 0x74, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x75, 0x70, 0x6c, 0x6f, 0x61,
	0x64, 0x5f, 0x72, 0x75, 0x6c, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x69,
	0x6e, 0x66, 0x6c, 0x75, 0x65, 0x6e, 0x7a, 0x61, 0x6e, 0x65, 0x74, 0x2e, 0x73, 0x74, 0x75, 0x64,
	0x79, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x45, 0x78, 0x70, 0x72, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x19, 0x70, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e,
	0x74, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x75, 0x6c, 0x65, 0x1a,
	0xb8, 0x02, 0x0a, 0x05, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x12, 0x3f, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x69, 0x6e, 0x66, 0x6c, 0x75, 0x65,
	0x6e, 0x7a, 0x61, 0x6e, 0x65, 0x74, 0x2e, 0x73, 0x74, 0x75, 0x64, 0x79, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x69, 0x73, 0x65, 0x64, 0x4f, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x4d, 0x0a, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x2b, 0x2e, 0x69, 0x6e, 0x66, 0x6c, 0x75, 0x65, 0x6e, 0x7a, 0x61, 0x6e, 0x65, 0x74, 0x2e, 0x73,
	0x74, 0x75, 0x64, 0x79, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x6f, 0x63,
	0x61, 0x6c, 0x69, 0x73, 0x65, 0x64, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x33, 0x0a, 0x04, 0x74, 0x61, 0x67,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x69, 0x6e, 0x66, 0x6c, 0x75, 0x65,
	0x6e, 0x7a, 0x61, 0x6e, 0x65, 0x74, 0x2e, 0x73, 0x74, 0x75, 0x64, 0x79, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x54, 0x61, 0x67, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12, 0x1d,
	0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x44, 0x61, 0x74, 0x65, 0x12, 0x19, 0x0a,
	0x08, 0x65, 0x6e, 0x64, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x07, 0x65, 0x6e, 0x64, 0x44, 0x61, 0x74, 0x65, 0x12, 0x30, 0x0a, 0x14, 0x73, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x5f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x73, 0x74, 0x75, 0x64, 0x79,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x12, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x44, 0x65,
	0x66, 0x61, 0x75, 0x6c, 0x74, 0x53, 0x74, 0x75, 0x64, 0x79, 0x1a, 0x51, 0x0a, 0x06, 0x4d, 0x65,
	0x6d, 0x62, 0x65, 0x72, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x72, 0x6f, 0x6c,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x1a, 0x5b, 0x0a,
	0x05, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x2b, 0x0a, 0x11, 0x70, 0x61, 0x72, 0x74, 0x69, 0x63,
	0x69, 0x70, 0x61, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x10, 0x70, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0xd7, 0x01, 0x0a, 0x0c, 0x53,
	0x74, 0x75, 0x64, 0x79, 0x46, 0x6f, 0x72, 0x55, 0x73, 0x65, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x3d, 0x0a, 0x05, 0x70, 0x72, 0x6f, 0x70, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x69, 0x6e, 0x66, 0x6c, 0x75, 0x65, 0x6e, 0x7a, 0x61,
	0x6e, 0x65, 0x74, 0x2e, 0x73, 0x74, 0x75, 0x64, 0x79, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x53, 0x74, 0x75, 0x64, 0x79, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x52, 0x05, 0x70,
	0x72, 0x6f, 0x70, 0x73, 0x12, 0x3d, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x69, 0x6e, 0x66, 0x6c, 0x75, 0x65, 0x6e, 0x7a, 0x61, 0x6e,
	0x65, 0x74, 0x2e, 0x73, 0x74, 0x75, 0x64, 0x79, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x53, 0x74, 0x75, 0x64, 0x79, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x05, 0x73, 0x74,
	0x61, 0x74, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69,
	0x64, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x49, 0x64, 0x73, 0x22, 0x48, 0x0a, 0x03, 0x54, 0x61, 0x67, 0x12, 0x41, 0x0a, 0x05, 0x6c,
	0x61, 0x62, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x69, 0x6e, 0x66,
	0x6c, 0x75, 0x65, 0x6e, 0x7a, 0x61, 0x6e, 0x65, 0x74, 0x2e, 0x73, 0x74, 0x75, 0x64, 0x79, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x69, 0x73, 0x65,
	0x64, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x22, 0xc7,
	0x01, 0x0a, 0x0e, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x53, 0x75, 0x72, 0x76, 0x65,
	0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x5f, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x4b, 0x65, 0x79,
	0x12, 0x1d, 0x0a, 0x0a, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x5f, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x46, 0x72, 0x6f, 0x6d, 0x12,
	0x1f, 0x0a, 0x0b, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x5f, 0x75, 0x6e, 0x74, 0x69, 0x6c, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x55, 0x6e, 0x74, 0x69, 0x6c,
	0x12, 0x1b, 0x0a, 0x09, 0x73, 0x74, 0x75, 0x64, 0x79, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x74, 0x75, 0x64, 0x79, 0x4b, 0x65, 0x79, 0x12, 0x1a, 0x0a,
	0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70,
	0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x22, 0xb0, 0x02, 0x0a, 0x0a, 0x53, 0x75, 0x72,
	0x76, 0x65, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x74, 0x75, 0x64, 0x79,
	0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x74, 0x75, 0x64,
	0x79, 0x4b, 0x65, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x5f, 0x6b,
	0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79,
	0x4b, 0x65, 0x79, 0x12, 0x3f, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x2b, 0x2e, 0x69, 0x6e, 0x66, 0x6c, 0x75, 0x65, 0x6e, 0x7a, 0x61, 0x6e, 0x65, 0x74,
	0x2e, 0x73, 0x74, 0x75, 0x64, 0x79, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c,
	0x6f, 0x63, 0x61, 0x6c, 0x69, 0x73, 0x65, 0x64, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x4d, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x69, 0x6e, 0x66, 0x6c,
	0x75, 0x65, 0x6e, 0x7a, 0x61, 0x6e, 0x65, 0x74, 0x2e, 0x73, 0x74, 0x75, 0x64, 0x79, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x69, 0x73, 0x65, 0x64,
	0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x56, 0x0a, 0x10, 0x74, 0x79, 0x70, 0x69, 0x63, 0x61, 0x6c, 0x5f, 0x64,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e,
	0x69, 0x6e, 0x66, 0x6c, 0x75, 0x65, 0x6e, 0x7a, 0x61, 0x6e, 0x65, 0x74, 0x2e, 0x73, 0x74, 0x75,
	0x64, 0x79, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x6c,
	0x69, 0x73, 0x65, 0x64, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x0f, 0x74, 0x79, 0x70, 0x69,
	0x63, 0x61, 0x6c, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xa2, 0x01, 0x0a, 0x0f,
	0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x53, 0x75, 0x72, 0x76, 0x65, 0x79, 0x73, 0x12,
	0x44, 0x0a, 0x07, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x2a, 0x2e, 0x69, 0x6e, 0x66, 0x6c, 0x75, 0x65, 0x6e, 0x7a, 0x61, 0x6e, 0x65, 0x74, 0x2e,
	0x73, 0x74, 0x75, 0x64, 0x79, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x41, 0x73,
	0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x53, 0x75, 0x72, 0x76, 0x65, 0x79, 0x52, 0x07, 0x73, 0x75,
	0x72, 0x76, 0x65, 0x79, 0x73, 0x12, 0x49, 0x0a, 0x0c, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x5f,
	0x69, 0x6e, 0x66, 0x6f, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x69, 0x6e,
	0x66, 0x6c, 0x75, 0x65, 0x6e, 0x7a, 0x61, 0x6e, 0x65, 0x74, 0x2e, 0x73, 0x74, 0x75, 0x64, 0x79,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x75, 0x72, 0x76, 0x65, 0x79, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x0b, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x73,
	0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69,
	0x6e, 0x66, 0x6c, 0x75, 0x65, 0x6e, 0x7a, 0x61, 0x6e, 0x65, 0x74, 0x2f, 0x73, 0x74, 0x75, 0x64,
	0x79, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70,
	0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_study_service_study_proto_rawDescOnce sync.Once
	file_study_service_study_proto_rawDescData = file_study_service_study_proto_rawDesc
)

func file_study_service_study_proto_rawDescGZIP() []byte {
	file_study_service_study_proto_rawDescOnce.Do(func() {
		file_study_service_study_proto_rawDescData = protoimpl.X.CompressGZIP(file_study_service_study_proto_rawDescData)
	})
	return file_study_service_study_proto_rawDescData
}

var file_study_service_study_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_study_service_study_proto_goTypes = []interface{}{
	(*Study)(nil),           // 0: influenzanet.study_service.Study
	(*StudyForUser)(nil),    // 1: influenzanet.study_service.StudyForUser
	(*Tag)(nil),             // 2: influenzanet.study_service.Tag
	(*AssignedSurvey)(nil),  // 3: influenzanet.study_service.AssignedSurvey
	(*SurveyInfo)(nil),      // 4: influenzanet.study_service.SurveyInfo
	(*AssignedSurveys)(nil), // 5: influenzanet.study_service.AssignedSurveys
	(*Study_Props)(nil),     // 6: influenzanet.study_service.Study.Props
	(*Study_Member)(nil),    // 7: influenzanet.study_service.Study.Member
	(*Study_Stats)(nil),     // 8: influenzanet.study_service.Study.Stats
	(*Expression)(nil),      // 9: influenzanet.study_service.Expression
	(*LocalisedObject)(nil), // 10: influenzanet.study_service.LocalisedObject
}
var file_study_service_study_proto_depIdxs = []int32{
	6,  // 0: influenzanet.study_service.Study.props:type_name -> influenzanet.study_service.Study.Props
	9,  // 1: influenzanet.study_service.Study.rules:type_name -> influenzanet.study_service.Expression
	7,  // 2: influenzanet.study_service.Study.members:type_name -> influenzanet.study_service.Study.Member
	8,  // 3: influenzanet.study_service.Study.stats:type_name -> influenzanet.study_service.Study.Stats
	9,  // 4: influenzanet.study_service.Study.participant_file_upload_rule:type_name -> influenzanet.study_service.Expression
	6,  // 5: influenzanet.study_service.StudyForUser.props:type_name -> influenzanet.study_service.Study.Props
	8,  // 6: influenzanet.study_service.StudyForUser.stats:type_name -> influenzanet.study_service.Study.Stats
	10, // 7: influenzanet.study_service.Tag.label:type_name -> influenzanet.study_service.LocalisedObject
	10, // 8: influenzanet.study_service.SurveyInfo.name:type_name -> influenzanet.study_service.LocalisedObject
	10, // 9: influenzanet.study_service.SurveyInfo.description:type_name -> influenzanet.study_service.LocalisedObject
	10, // 10: influenzanet.study_service.SurveyInfo.typical_duration:type_name -> influenzanet.study_service.LocalisedObject
	3,  // 11: influenzanet.study_service.AssignedSurveys.surveys:type_name -> influenzanet.study_service.AssignedSurvey
	4,  // 12: influenzanet.study_service.AssignedSurveys.survey_infos:type_name -> influenzanet.study_service.SurveyInfo
	10, // 13: influenzanet.study_service.Study.Props.name:type_name -> influenzanet.study_service.LocalisedObject
	10, // 14: influenzanet.study_service.Study.Props.description:type_name -> influenzanet.study_service.LocalisedObject
	2,  // 15: influenzanet.study_service.Study.Props.tags:type_name -> influenzanet.study_service.Tag
	16, // [16:16] is the sub-list for method output_type
	16, // [16:16] is the sub-list for method input_type
	16, // [16:16] is the sub-list for extension type_name
	16, // [16:16] is the sub-list for extension extendee
	0,  // [0:16] is the sub-list for field type_name
}

func init() { file_study_service_study_proto_init() }
func file_study_service_study_proto_init() {
	if File_study_service_study_proto != nil {
		return
	}
	file_study_service_expression_proto_init()
	file_study_service_survey_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_study_service_study_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Study); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_study_service_study_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StudyForUser); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_study_service_study_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tag); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_study_service_study_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssignedSurvey); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_study_service_study_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SurveyInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_study_service_study_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssignedSurveys); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_study_service_study_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Study_Props); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_study_service_study_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Study_Member); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_study_service_study_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Study_Stats); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_study_service_study_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_study_service_study_proto_goTypes,
		DependencyIndexes: file_study_service_study_proto_depIdxs,
		MessageInfos:      file_study_service_study_proto_msgTypes,
	}.Build()
	File_study_service_study_proto = out.File
	file_study_service_study_proto_rawDesc = nil
	file_study_service_study_proto_goTypes = nil
	file_study_service_study_proto_depIdxs = nil
}
