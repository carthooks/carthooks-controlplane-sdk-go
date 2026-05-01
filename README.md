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

The contract now has two explicit layers:

1. Shared Go request/response structs live in this repository.
2. Arcubase `platform` imports those types and generates the HTTP contract artifact.

Arcubase upstream integration points:

- `platform/src/vo/control_plane_tenant.go`
- `platform/src/vo/control_plane_service_account.go`
- `platform/src/vo/control_plane_runtime_credential.go`
- `platform/types/controlplane/controlplane_swagger.json`

That means:

- Go type shape is centralized in this repo
- HTTP contract artifact is generated in Arcubase upstream
- this repo keeps a checked-in snapshot of that artifact under `contracts/`

## Type Consistency Policy

To keep SDK types aligned with Arcubase upstream, follow these engineering rules:

1. Shared Go request/response structs are edited here.
2. Arcubase upstream must import those structs instead of redefining them.
3. Arcubase upstream must regenerate `platform/types/controlplane/controlplane_swagger.json` after contract changes.
4. This repo must sync `contracts/controlplane_swagger.json` from upstream in the same rollout window.
5. `make validate-contract` must fail if upstream artifact and local snapshot drift.
6. Breaking field changes require an explicit version bump and release note.

## Contract Sync

When working in the Arcubase monorepo:

```bash
make sync-contract ARCUBASE_ROOT=/path/to/arcubase
make validate-contract ARCUBASE_ROOT=/path/to/arcubase
```

The validation target compares:

- upstream: `platform/types/controlplane/controlplane_swagger.json`
- local snapshot: `contracts/controlplane_swagger.json`

If they differ, the SDK is out of sync and must not be released.
