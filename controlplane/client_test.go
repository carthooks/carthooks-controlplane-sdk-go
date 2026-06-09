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
				"tenantUserId":                "123",
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
	if resp.TenantUserID != "123" {
		t.Fatalf("unexpected tenant user id: %s", resp.TenantUserID)
	}
}

func TestEnsureTenantUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/internal/tenant-users/ensure" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var req EnsureTenantUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.ArcubaseTenantID != "tenant-1" {
			t.Fatalf("unexpected arcubaseTenantId: %s", req.ArcubaseTenantID)
		}
		if req.ExternalSource != "botworks" {
			t.Fatalf("unexpected externalSource: %s", req.ExternalSource)
		}
		if req.ExternalSubjectType != "human_member" {
			t.Fatalf("unexpected externalSubjectType: %s", req.ExternalSubjectType)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"arcubaseTenantId":    "tenant-1",
				"tenantUserId":        "2188889901",
				"bindingStatus":       "ready",
				"created":             true,
				"externalSource":      "botworks",
				"externalSubjectType": "human_member",
				"externalSubjectId":   "u_1",
				"displayName":         "Alice",
				"email":               "alice@example.com",
			},
		})
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, InternalAuthKey: "secret"})
	resp, err := client.EnsureTenantUser(context.Background(), EnsureTenantUserRequest{
		ArcubaseTenantID:    "tenant-1",
		ExternalSource:      "botworks",
		ExternalSubjectType: "human_member",
		ExternalSubjectID:   "u_1",
		DisplayName:         "Alice",
		Email:               "alice@example.com",
	})
	if err != nil {
		t.Fatalf("EnsureTenantUser error: %v", err)
	}
	if resp.TenantUserID != "2188889901" {
		t.Fatalf("unexpected tenant user id: %s", resp.TenantUserID)
	}
	if !resp.Created {
		t.Fatalf("expected created=true")
	}
}

func TestTenantDepartmentClientMethods(t *testing.T) {
	var calls []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		switch r.Method + " " + r.URL.Path {
		case "GET /api/internal/tenants/tenant-1/departments/tree":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"items": []map[string]any{{
						"id":             "1",
						"name":           "Root",
						"parentId":       "0",
						"depth":          1,
						"path":           "/0",
						"childDeptCount": 1,
					}},
				},
			})
		case "POST /api/internal/tenants/tenant-1/departments/paths":
			var req GetTenantDepartmentPathsRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode paths request: %v", err)
			}
			if len(req.DepartmentIDs) != 1 || req.DepartmentIDs[0] != "2" {
				t.Fatalf("unexpected departmentIds: %#v", req.DepartmentIDs)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"paths": [][]map[string]any{{
						{"id": "1", "name": "Root", "parentId": "0"},
						{"id": "2", "name": "Sales", "parentId": "1"},
					}},
				},
			})
		case "POST /api/internal/tenants/tenant-1/departments":
			var req CreateTenantDepartmentRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode create request: %v", err)
			}
			if req.ParentID != "1" || req.Name != "Sales" {
				t.Fatalf("unexpected create request: %#v", req)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"id": "2", "name": "Sales", "parentId": "1"},
			})
		case "PUT /api/internal/tenants/tenant-1/departments/2/rename":
			var req RenameTenantDepartmentRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode rename request: %v", err)
			}
			if req.Name != "Revenue" {
				t.Fatalf("unexpected rename request: %#v", req)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"id": "2", "name": "Revenue", "parentId": "1"},
			})
		case "DELETE /api/internal/tenants/tenant-1/departments/2":
			_ = json.NewEncoder(w).Encode(map[string]any{"data": true})
		case "PUT /api/internal/tenants/tenant-1/departments/2/relocate":
			var req RelocateTenantDepartmentRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode relocate request: %v", err)
			}
			if req.TargetID != "1" || req.Type != "inner" {
				t.Fatalf("unexpected relocate request: %#v", req)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"id": "2", "name": "Revenue", "parentId": "1"},
			})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, InternalAuthKey: "secret"})

	tree, err := client.ListTenantDepartments(context.Background(), "tenant-1")
	if err != nil {
		t.Fatalf("ListTenantDepartments error: %v", err)
	}
	if len(tree.Items) != 1 || tree.Items[0].ID != "1" {
		t.Fatalf("unexpected tree response: %#v", tree)
	}

	paths, err := client.GetTenantDepartmentPaths(context.Background(), GetTenantDepartmentPathsRequest{
		ArcubaseTenantID: "tenant-1",
		DepartmentIDs:    []string{"2"},
	})
	if err != nil {
		t.Fatalf("GetTenantDepartmentPaths error: %v", err)
	}
	if len(paths.Paths) != 1 || len(paths.Paths[0]) != 2 {
		t.Fatalf("unexpected paths response: %#v", paths)
	}

	created, err := client.CreateTenantDepartment(context.Background(), "tenant-1", CreateTenantDepartmentRequest{
		ParentID: "1",
		Name:     "Sales",
	})
	if err != nil {
		t.Fatalf("CreateTenantDepartment error: %v", err)
	}
	if created.ID != "2" {
		t.Fatalf("unexpected created department: %#v", created)
	}

	renamed, err := client.RenameTenantDepartment(context.Background(), "tenant-1", "2", RenameTenantDepartmentRequest{Name: "Revenue"})
	if err != nil {
		t.Fatalf("RenameTenantDepartment error: %v", err)
	}
	if renamed.Name != "Revenue" {
		t.Fatalf("unexpected renamed department: %#v", renamed)
	}

	if err := client.DeleteTenantDepartment(context.Background(), "tenant-1", "2"); err != nil {
		t.Fatalf("DeleteTenantDepartment error: %v", err)
	}

	relocated, err := client.RelocateTenantDepartment(context.Background(), "tenant-1", "2", RelocateTenantDepartmentRequest{
		TargetID: "1",
		Type:     "inner",
	})
	if err != nil {
		t.Fatalf("RelocateTenantDepartment error: %v", err)
	}
	if relocated.ParentID != "1" {
		t.Fatalf("unexpected relocated department: %#v", relocated)
	}

	expected := []string{
		"GET /api/internal/tenants/tenant-1/departments/tree",
		"POST /api/internal/tenants/tenant-1/departments/paths",
		"POST /api/internal/tenants/tenant-1/departments",
		"PUT /api/internal/tenants/tenant-1/departments/2/rename",
		"DELETE /api/internal/tenants/tenant-1/departments/2",
		"PUT /api/internal/tenants/tenant-1/departments/2/relocate",
	}
	if len(calls) != len(expected) {
		t.Fatalf("unexpected calls: %#v", calls)
	}
	for i := range expected {
		if calls[i] != expected[i] {
			t.Fatalf("call %d: expected %s, got %s", i, expected[i], calls[i])
		}
	}
}

func TestTenantUserDepartmentsClientMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method + " " + r.URL.Path {
		case "GET /api/internal/tenant-users/2188889901/departments":
			if r.URL.Query().Get("arcubaseTenantId") != "tenant-1" {
				t.Fatalf("missing arcubaseTenantId query")
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"arcubaseTenantId": "tenant-1",
					"tenantUserId":     "2188889901",
					"departmentIds":    []string{"1", "2"},
				},
			})
		case "PUT /api/internal/tenant-users/2188889901/departments":
			var req UpdateTenantUserDepartmentsRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode update request: %v", err)
			}
			if req.ArcubaseTenantID != "tenant-1" || len(req.DepartmentIDs) != 1 || req.DepartmentIDs[0] != "2" {
				t.Fatalf("unexpected update request: %#v", req)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"arcubaseTenantId": "tenant-1",
					"tenantUserId":     "2188889901",
					"departmentIds":    []string{"2"},
				},
			})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, InternalAuthKey: "secret"})
	current, err := client.GetTenantUserDepartments(context.Background(), "2188889901", "tenant-1")
	if err != nil {
		t.Fatalf("GetTenantUserDepartments error: %v", err)
	}
	if len(current.DepartmentIDs) != 2 {
		t.Fatalf("unexpected department IDs: %#v", current.DepartmentIDs)
	}

	updated, err := client.UpdateTenantUserDepartments(context.Background(), "2188889901", UpdateTenantUserDepartmentsRequest{
		ArcubaseTenantID: "tenant-1",
		DepartmentIDs:    []string{"2"},
	})
	if err != nil {
		t.Fatalf("UpdateTenantUserDepartments error: %v", err)
	}
	if len(updated.DepartmentIDs) != 1 || updated.DepartmentIDs[0] != "2" {
		t.Fatalf("unexpected updated department IDs: %#v", updated.DepartmentIDs)
	}
}
