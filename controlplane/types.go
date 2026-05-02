package controlplane

import "time"

type Config struct {
	BaseURL         string
	InternalAuthKey string
	RequestTimeout  time.Duration
}

type EnsureTenantRequest struct {
	ExternalSource      string `json:"externalSource" binding:"required"`
	ExternalSubjectType string `json:"externalSubjectType" binding:"required"`
	ExternalSubjectID   string `json:"externalSubjectId" binding:"required"`
	DisplayName         string `json:"displayName,omitempty"`
	SlugHint            string `json:"slugHint,omitempty"`
	IdempotencyKey      string `json:"idempotencyKey,omitempty"`
}

type TenantStatusRequest struct {
	ExternalSource      string `json:"externalSource" form:"externalSource" binding:"required"`
	ExternalSubjectType string `json:"externalSubjectType" form:"externalSubjectType" binding:"required"`
	ExternalSubjectID   string `json:"externalSubjectId" form:"externalSubjectId" binding:"required"`
}

type TenantControlPlaneResponse struct {
	ArcubaseTenantID    string `json:"arcubaseTenantId"`
	ArcubaseTenantSlug  string `json:"arcubaseTenantSlug"`
	ArcubaseInstanceID  string `json:"arcubaseInstanceId"`
	TenantBindingStatus string `json:"tenantBindingStatus"`
	Created             bool   `json:"created"`
}

type EnsureServiceAccountRequest struct {
	ArcubaseTenantID    string `json:"arcubaseTenantId" binding:"required"`
	ExternalSource      string `json:"externalSource" binding:"required"`
	ExternalSubjectType string `json:"externalSubjectType" binding:"required"`
	ExternalSubjectID   string `json:"externalSubjectId" binding:"required"`
	DisplayName         string `json:"displayName,omitempty"`
	AvatarURL           string `json:"avatarUrl,omitempty"`
	Bio                 string `json:"bio,omitempty"`
	IdempotencyKey      string `json:"idempotencyKey,omitempty"`
}

type SyncServiceAccountProfileRequest = EnsureServiceAccountRequest

type ServiceAccountControlPlaneResponse struct {
	ServiceAccountID            string     `json:"serviceAccountId"`
	ServiceAccountBindingStatus string     `json:"serviceAccountBindingStatus"`
	ProfileSynced               bool       `json:"profileSynced"`
	Created                     bool       `json:"created"`
	ArcubaseTenantID            string     `json:"arcubaseTenantId"`
	TenantUserID                string     `json:"tenantUserId"`
	ExternalSource              string     `json:"externalSource"`
	ExternalSubjectType         string     `json:"externalSubjectType"`
	ExternalSubjectID           string     `json:"externalSubjectId"`
	DisplayName                 string     `json:"displayName"`
	AvatarURL                   string     `json:"avatarUrl"`
	Bio                         string     `json:"bio"`
	ProfileSyncedAt             *time.Time `json:"profileSyncedAt"`
}

type IssueServiceAccountRuntimeCredentialRequest struct {
	ArcubaseTenantID string   `json:"arcubaseTenantId" binding:"required"`
	RequestedScope   []string `json:"requestedScope" binding:"required"`
}

type ServiceAccountRuntimeCredentialResponse struct {
	CredentialType string    `json:"credentialType"`
	AccessToken    string    `json:"accessToken"`
	ExpiresAt      time.Time `json:"expiresAt"`
	Scope          []string  `json:"scope"`
	TenantUserID   string    `json:"tenantUserId"`
	SubjectType    string    `json:"subjectType,omitempty"`
	SubjectID      string    `json:"subjectId,omitempty"`
}

type IssueRuntimeSessionRequest struct {
	ArcubaseTenantID string   `json:"arcubaseTenantId" binding:"required"`
	SubjectType      string   `json:"subjectType" binding:"required"`
	SubjectID        string   `json:"subjectId,omitempty"`
	RequestedScope   []string `json:"requestedScope" binding:"required"`
}

type RuntimeSessionResponse = ServiceAccountRuntimeCredentialResponse

type envelopeError struct {
	Code    string `json:"code"`
	Key     string `json:"key"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

type APIError struct {
	StatusCode int
	Code       string
	Key        string
	Message    string
	Type       string
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	if e.Key != "" {
		return e.Key
	}
	return e.Code
}
