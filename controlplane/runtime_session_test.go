package controlplane

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIssueRuntimeSession(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/internal/runtime-sessions/issue" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("X-Internal-Auth") != "secret" {
			t.Fatalf("missing X-Internal-Auth header")
		}

		var req IssueRuntimeSessionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.SubjectType != "system_admin" {
			t.Fatalf("unexpected subject type: %s", req.SubjectType)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"credentialType": "tenant_access_token",
				"accessToken":    "token-1",
				"expiresAt":      time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC).Format(time.RFC3339),
				"scope":          []string{"admin:runtime"},
				"tenantUserId":   "tu-1",
				"subjectType":    "system_admin",
				"subjectId":      "admin-1",
			},
		})
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, InternalAuthKey: "secret"})
	resp, err := client.IssueRuntimeSession(context.Background(), IssueRuntimeSessionRequest{
		ArcubaseTenantID: "tenant-1",
		SubjectType:      "system_admin",
		RequestedScope:   []string{"admin:runtime"},
	})
	if err != nil {
		t.Fatalf("IssueRuntimeSession error: %v", err)
	}
	if resp.AccessToken != "token-1" {
		t.Fatalf("unexpected token: %s", resp.AccessToken)
	}
	if resp.TenantUserID != "tu-1" {
		t.Fatalf("unexpected tenant user id: %s", resp.TenantUserID)
	}
}
