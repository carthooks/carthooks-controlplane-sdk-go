package controlplane

import "time"

type Config struct {
	BaseURL         string
	InternalAuthKey string
	RequestTimeout  time.Duration
}

type SyncServiceAccountProfileRequest = EnsureServiceAccountRequest

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
