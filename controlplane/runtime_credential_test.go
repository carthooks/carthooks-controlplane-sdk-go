package controlplane

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIssueServiceAccountRuntimeCredential(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/internal/service-accounts/sa-1/runtime-credential" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("X-Internal-Auth") != "secret" {
			t.Fatalf("missing X-Internal-Auth header")
		}

		var req IssueServiceAccountRuntimeCredentialRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.ArcubaseTenantID != "tenant-1" {
			t.Fatalf("unexpected tenant id: %s", req.ArcubaseTenantID)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"credentialType": "service_account_access_token",
				"accessToken":    "token-1",
				"expiresAt":      time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC).Format(time.RFC3339),
				"scope":          []string{"rows:read"},
			},
		})
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, InternalAuthKey: "secret"})
	resp, err := client.IssueServiceAccountRuntimeCredential(context.Background(), "sa-1", IssueServiceAccountRuntimeCredentialRequest{
		ArcubaseTenantID: "tenant-1",
		RequestedScope:   []string{"rows:read"},
	})
	if err != nil {
		t.Fatalf("IssueServiceAccountRuntimeCredential error: %v", err)
	}
	if resp.AccessToken != "token-1" {
		t.Fatalf("unexpected token: %s", resp.AccessToken)
	}
}

func TestIssueServiceAccountRuntimeCredential_ErrorEnvelope(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{
				"key":     "SERVICE_ACCOUNT_BINDING_NOT_READY",
				"message": "binding not ready",
				"type":    "user",
			},
		})
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, InternalAuthKey: "secret"})
	_, err := client.IssueServiceAccountRuntimeCredential(context.Background(), "sa-1", IssueServiceAccountRuntimeCredentialRequest{
		ArcubaseTenantID: "tenant-1",
		RequestedScope:   []string{"rows:read"},
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Key != "SERVICE_ACCOUNT_BINDING_NOT_READY" {
		t.Fatalf("unexpected error key: %s", apiErr.Key)
	}
}
