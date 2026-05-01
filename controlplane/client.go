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
