package controlplane

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEnsureTenant(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/internal/tenants/ensure" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("X-Internal-Auth") != "secret" {
			t.Fatalf("missing X-Internal-Auth header")
		}
		if r.Header.Get("X-Trace-Id") != "trace_test_1" {
			t.Fatalf("missing X-Trace-Id header")
		}
		if r.Header.Get("X-Request-Id") != "trace_test_1" {
			t.Fatalf("missing X-Request-Id header")
		}

		var req EnsureTenantRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.ExternalSubjectID != "team-1" {
			t.Fatalf("unexpected externalSubjectId: %s", req.ExternalSubjectID)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"arcubaseTenantId":    "tenant-1",
				"tenantBindingStatus": "ready",
				"created":             true,
			},
		})
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, InternalAuthKey: "secret"})
	resp, err := client.EnsureTenant(ContextWithTraceID(context.Background(), "trace_test_1"), EnsureTenantRequest{
		ExternalSource:      "botworks",
		ExternalSubjectType: "team",
		ExternalSubjectID:   "team-1",
	})
	if err != nil {
		t.Fatalf("EnsureTenant error: %v", err)
	}
	if resp.ArcubaseTenantID != "tenant-1" {
		t.Fatalf("unexpected tenant id: %s", resp.ArcubaseTenantID)
	}
}

func TestGetTenantStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Query().Get("externalSource") != "botworks" {
			t.Fatalf("unexpected externalSource")
		}
		if r.Header.Get("X-Internal-Auth") != "secret" {
			t.Fatalf("missing X-Internal-Auth header")
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"arcubaseTenantId":    "tenant-1",
				"tenantBindingStatus": "ready",
			},
		})
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, InternalAuthKey: "secret"})
	resp, err := client.GetTenantStatus(context.Background(), TenantStatusRequest{
		ExternalSource:      "botworks",
		ExternalSubjectType: "team",
		ExternalSubjectID:   "team-1",
	})
	if err != nil {
		t.Fatalf("GetTenantStatus error: %v", err)
	}
	if resp.TenantBindingStatus != "ready" {
		t.Fatalf("unexpected status: %s", resp.TenantBindingStatus)
	}
}

func TestEnsureTenant_ReturnsEnvelopeErrorEvenWhenHTTP200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{
				"code":    "CONTROL_PLANE_UNAUTHORIZED",
				"message": "Control-plane unauthorized",
				"type":    "user",
				"key":     "CONTROL_PLANE_UNAUTHORIZED",
			},
		})
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, InternalAuthKey: "secret"})
	resp, err := client.EnsureTenant(context.Background(), EnsureTenantRequest{
		ExternalSource:      "botworks",
		ExternalSubjectType: "team",
		ExternalSubjectID:   "team-1",
	})
	if err == nil {
		t.Fatalf("expected error, got nil response=%+v", resp)
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Code != "CONTROL_PLANE_UNAUTHORIZED" {
		t.Fatalf("unexpected code: %s", apiErr.Code)
	}
}

func TestEnsureServiceAccount_SendsExplicitNonAdminFlag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/internal/service-accounts/ensure" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		value, ok := body["isAdmin"].(bool)
		if !ok {
			t.Fatalf("expected isAdmin boolean in request body, got %#v", body["isAdmin"])
		}
		if value {
			t.Fatalf("expected isAdmin=false")
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"serviceAccountId":            "svc-1",
				"serviceAccountBindingStatus": "ready",
				"arcubaseTenantId":            "tenant-1",
				"tenantUserId":                "tu-1",
				"tenantUserNumericId":         123,
				"externalSource":              "botworks",
				"externalSubjectType":         "digiemployee",
				"externalSubjectId":           "de-1",
			},
		})
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, InternalAuthKey: "secret"})
	resp, err := client.EnsureServiceAccount(context.Background(), EnsureServiceAccountRequest{
		ArcubaseTenantID:    "tenant-1",
		ExternalSource:      "botworks",
		ExternalSubjectType: "digiemployee",
		ExternalSubjectID:   "de-1",
		IsAdmin:             false,
	})
	if err != nil {
		t.Fatalf("EnsureServiceAccount error: %v", err)
	}
	if resp.ServiceAccountID != "svc-1" {
		t.Fatalf("unexpected service account id: %s", resp.ServiceAccountID)
	}
	if resp.TenantUserNumericID != 123 {
		t.Fatalf("unexpected tenant user numeric id: %d", resp.TenantUserNumericID)
	}
}
