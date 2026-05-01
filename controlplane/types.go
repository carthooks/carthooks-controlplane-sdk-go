package controlplane

import "time"

type Config struct {
	BaseURL         string
	InternalAuthKey string
	RequestTimeout  time.Duration
}

type EnsureTenantRequest struct {
	ExternalSource      string `json:"externalSource"`
	ExternalSubjectType string `json:"externalSubjectType"`
	ExternalSubjectID   string `json:"externalSubjectId"`
	DisplayName         string `json:"displayName,omitempty"`
	SlugHint            string `json:"slugHint,omitempty"`
	IdempotencyKey      string `json:"idempotencyKey,omitempty"`
}

type TenantStatusRequest struct {
	ExternalSource      string
	ExternalSubjectType string
	ExternalSubjectID   string
}

type TenantControlPlaneResponse struct {
	ArcubaseTenantID    string `json:"arcubaseTenantId"`
	ArcubaseTenantSlug  string `json:"arcubaseTenantSlug"`
	ArcubaseInstanceID  string `json:"arcubaseInstanceId"`
	TenantBindingStatus string `json:"tenantBindingStatus"`
	Created             bool   `json:"created"`
}

type EnsureServiceAccountRequest struct {
	ArcubaseTenantID    string `json:"arcubaseTenantId"`
	ExternalSource      string `json:"externalSource"`
	ExternalSubjectType string `json:"externalSubjectType"`
	ExternalSubjectID   string `json:"externalSubjectId"`
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
	ArcubaseTenantID string   `json:"arcubaseTenantId"`
	RequestedScope   []string `json:"requestedScope"`
}

type ServiceAccountRuntimeCredentialResponse struct {
	CredentialType string    `json:"credentialType"`
	AccessToken    string    `json:"accessToken"`
	ExpiresAt      time.Time `json:"expiresAt"`
	Scope          []string  `json:"scope"`
}

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
