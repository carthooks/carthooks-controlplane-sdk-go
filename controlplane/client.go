package controlplane

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type traceIDContextKey struct{}

func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	traceID = strings.TrimSpace(traceID)
	if traceID == "" {
		return ctx
	}
	return context.WithValue(ctx, traceIDContextKey{}, traceID)
}

func TraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	traceID, _ := ctx.Value(traceIDContextKey{}).(string)
	return strings.TrimSpace(traceID)
}

type Client struct {
	baseURL         string
	internalAuthKey string
	httpClient      *http.Client
}

func NewClient(cfg Config) *Client {
	timeout := cfg.RequestTimeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &Client{
		baseURL:         strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/"),
		internalAuthKey: strings.TrimSpace(cfg.InternalAuthKey),
		httpClient:      &http.Client{Timeout: timeout},
	}
}

func (c *Client) EnsureTenant(ctx context.Context, req EnsureTenantRequest) (*TenantControlPlaneResponse, error) {
	return doJSON[EnsureTenantRequest, TenantControlPlaneResponse](ctx, c, http.MethodPost, "/api/internal/tenants/ensure", req, nil)
}

func (c *Client) GetTenantStatus(ctx context.Context, req TenantStatusRequest) (*TenantControlPlaneResponse, error) {
	params := url.Values{}
	params.Set("externalSource", strings.TrimSpace(req.ExternalSource))
	params.Set("externalSubjectType", strings.TrimSpace(req.ExternalSubjectType))
	params.Set("externalSubjectId", strings.TrimSpace(req.ExternalSubjectID))
	return doJSON[struct{}, TenantControlPlaneResponse](ctx, c, http.MethodGet, "/api/internal/tenants/status", struct{}{}, params)
}

func (c *Client) EnsureServiceAccount(ctx context.Context, req EnsureServiceAccountRequest) (*ServiceAccountControlPlaneResponse, error) {
	return doJSON[EnsureServiceAccountRequest, ServiceAccountControlPlaneResponse](ctx, c, http.MethodPost, "/api/internal/service-accounts/ensure", req, nil)
}

func (c *Client) SyncServiceAccountProfile(ctx context.Context, serviceAccountID string, req SyncServiceAccountProfileRequest) (*ServiceAccountControlPlaneResponse, error) {
	path := fmt.Sprintf("/api/internal/service-accounts/%s/profile/sync", url.PathEscape(strings.TrimSpace(serviceAccountID)))
	return doJSON[SyncServiceAccountProfileRequest, ServiceAccountControlPlaneResponse](ctx, c, http.MethodPost, path, req, nil)
}

func (c *Client) IssueServiceAccountRuntimeCredential(ctx context.Context, serviceAccountID string, req IssueServiceAccountRuntimeCredentialRequest) (*ServiceAccountRuntimeCredentialResponse, error) {
	path := fmt.Sprintf("/api/internal/service-accounts/%s/runtime-credential", url.PathEscape(strings.TrimSpace(serviceAccountID)))
	return doJSON[IssueServiceAccountRuntimeCredentialRequest, ServiceAccountRuntimeCredentialResponse](ctx, c, http.MethodPost, path, req, nil)
}

func (c *Client) IssueRuntimeSession(ctx context.Context, req IssueRuntimeSessionRequest) (*RuntimeSessionResponse, error) {
	return doJSON[IssueRuntimeSessionRequest, RuntimeSessionResponse](ctx, c, http.MethodPost, "/api/internal/runtime-sessions/issue", req, nil)
}

func (c *Client) EnsureTenantUser(ctx context.Context, req EnsureTenantUserRequest) (*TenantUserControlPlaneResponse, error) {
	return doJSON[EnsureTenantUserRequest, TenantUserControlPlaneResponse](ctx, c, http.MethodPost, "/api/internal/tenant-users/ensure", req, nil)
}

func (c *Client) ListTenantDepartments(ctx context.Context, arcubaseTenantID string) (*ListTenantDepartmentsResponse, error) {
	path := fmt.Sprintf("/api/internal/tenants/%s/departments/tree", url.PathEscape(strings.TrimSpace(arcubaseTenantID)))
	return doJSON[struct{}, ListTenantDepartmentsResponse](ctx, c, http.MethodGet, path, struct{}{}, nil)
}

func (c *Client) GetTenantDepartmentPaths(ctx context.Context, req GetTenantDepartmentPathsRequest) (*GetTenantDepartmentPathsResponse, error) {
	path := fmt.Sprintf("/api/internal/tenants/%s/departments/paths", url.PathEscape(strings.TrimSpace(req.ArcubaseTenantID)))
	return doJSON[GetTenantDepartmentPathsRequest, GetTenantDepartmentPathsResponse](ctx, c, http.MethodPost, path, req, nil)
}

func (c *Client) CreateTenantDepartment(ctx context.Context, arcubaseTenantID string, req CreateTenantDepartmentRequest) (*TenantDepartmentItem, error) {
	path := fmt.Sprintf("/api/internal/tenants/%s/departments", url.PathEscape(strings.TrimSpace(arcubaseTenantID)))
	return doJSON[CreateTenantDepartmentRequest, TenantDepartmentItem](ctx, c, http.MethodPost, path, req, nil)
}

func (c *Client) RenameTenantDepartment(ctx context.Context, arcubaseTenantID, departmentID string, req RenameTenantDepartmentRequest) (*TenantDepartmentItem, error) {
	path := fmt.Sprintf("/api/internal/tenants/%s/departments/%s/rename", url.PathEscape(strings.TrimSpace(arcubaseTenantID)), url.PathEscape(strings.TrimSpace(departmentID)))
	return doJSON[RenameTenantDepartmentRequest, TenantDepartmentItem](ctx, c, http.MethodPut, path, req, nil)
}

func (c *Client) DeleteTenantDepartment(ctx context.Context, arcubaseTenantID, departmentID string) error {
	path := fmt.Sprintf("/api/internal/tenants/%s/departments/%s", url.PathEscape(strings.TrimSpace(arcubaseTenantID)), url.PathEscape(strings.TrimSpace(departmentID)))
	_, err := doJSON[struct{}, bool](ctx, c, http.MethodDelete, path, struct{}{}, nil)
	return err
}

func (c *Client) RelocateTenantDepartment(ctx context.Context, arcubaseTenantID, departmentID string, req RelocateTenantDepartmentRequest) (*TenantDepartmentItem, error) {
	path := fmt.Sprintf("/api/internal/tenants/%s/departments/%s/relocate", url.PathEscape(strings.TrimSpace(arcubaseTenantID)), url.PathEscape(strings.TrimSpace(departmentID)))
	return doJSON[RelocateTenantDepartmentRequest, TenantDepartmentItem](ctx, c, http.MethodPut, path, req, nil)
}

func (c *Client) GetTenantUserDepartments(ctx context.Context, tenantUserID string, arcubaseTenantID string) (*TenantUserDepartmentsResponse, error) {
	params := url.Values{}
	params.Set("arcubaseTenantId", strings.TrimSpace(arcubaseTenantID))
	path := fmt.Sprintf("/api/internal/tenant-users/%s/departments", url.PathEscape(strings.TrimSpace(tenantUserID)))
	return doJSON[struct{}, TenantUserDepartmentsResponse](ctx, c, http.MethodGet, path, struct{}{}, params)
}

func (c *Client) UpdateTenantUserDepartments(ctx context.Context, tenantUserID string, req UpdateTenantUserDepartmentsRequest) (*TenantUserDepartmentsResponse, error) {
	path := fmt.Sprintf("/api/internal/tenant-users/%s/departments", url.PathEscape(strings.TrimSpace(tenantUserID)))
	return doJSON[UpdateTenantUserDepartmentsRequest, TenantUserDepartmentsResponse](ctx, c, http.MethodPut, path, req, nil)
}

func doJSON[Req any, Resp any](ctx context.Context, c *Client, method, path string, reqBody Req, params url.Values) (*Resp, error) {
	fullURL, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return nil, err
	}
	if params != nil && len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	var bodyReader *bytes.Reader
	if method == http.MethodGet {
		bodyReader = bytes.NewReader(nil)
	} else {
		body, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(body)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Accept", "application/json")
	if method != http.MethodGet {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	if c.internalAuthKey != "" {
		httpReq.Header.Set("X-Internal-Auth", c.internalAuthKey)
	}
	if traceID := TraceIDFromContext(ctx); traceID != "" {
		httpReq.Header.Set("X-Trace-Id", traceID)
		httpReq.Header.Set("X-Request-Id", traceID)
	}

	return decodeResponse[Resp](c.httpClient.Do(httpReq))
}

func decodeResponse[T any](resp *http.Response, err error) (*T, error) {
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var envelope struct {
		Data  json.RawMessage `json:"data"`
		Error *envelopeError  `json:"error"`
	}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&envelope); decodeErr != nil {
		return nil, decodeErr
	}

	if envelope.Error != nil {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		apiErr.Code = envelope.Error.Code
		apiErr.Key = envelope.Error.Key
		apiErr.Message = envelope.Error.Message
		apiErr.Type = envelope.Error.Type
		if apiErr.Code == "" && apiErr.Key == "" && apiErr.Message == "" {
			apiErr.Message = fmt.Sprintf("control-plane response status %d", resp.StatusCode)
		}
		return nil, apiErr
	}

	if resp.StatusCode != http.StatusOK {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		if apiErr.Code == "" && apiErr.Key == "" && apiErr.Message == "" {
			apiErr.Message = fmt.Sprintf("control-plane response status %d", resp.StatusCode)
		}
		return nil, apiErr
	}

	var out T
	if len(envelope.Data) == 0 || string(envelope.Data) == "null" {
		return &out, nil
	}
	if err := json.Unmarshal(envelope.Data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
