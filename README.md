# carthooks-controlplane-sdk-go

Arcubase internal platform/control-plane Go SDK.

This repository is intentionally separate from `github.com/carthooks/carthooks-sdk-go`.

- `github.com/carthooks/carthooks-sdk-go`
  - external tenant/data-plane SDK
  - public integration surface
- `github.com/carthooks/carthooks-controlplane-sdk-go`
  - internal platform/control-plane SDK
  - trusted internal callers only

The control-plane SDK uses `X-Internal-Auth` and must never be merged back into the external SDK.

## Source Of Truth

The source of truth for request/response contracts is Arcubase upstream:

- `platform/src/vo/control_plane_tenant.go`
- `platform/src/vo/control_plane_service_account.go`
- `platform/src/vo/control_plane_runtime_credential.go`

This SDK is an internal client projection of those contracts. Arcubase upstream owns schema changes first; this repo follows.

## Type Consistency Policy

To keep SDK types aligned with Arcubase upstream, follow these engineering rules:

1. Arcubase upstream `vo` structs are the canonical contract source.
2. Any control-plane request/response field change must land in Arcubase upstream first.
3. The same change must update this SDK in the same rollout window, not later as cleanup.
4. Add or update SDK tests with real JSON envelopes from Arcubase responses.
5. Breaking field changes require an explicit version bump and release note.

Recommended next step:

- generate a dedicated internal control-plane OpenAPI/JSON schema artifact from Arcubase upstream
- regenerate SDK request/response types from that artifact in CI

Until codegen is added, handwritten types in this repo must be treated as a mirrored contract and reviewed against upstream `vo` diffs on every change.
