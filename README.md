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

The Arcubase Server control-plane HTTP contract is the source of truth.
This SDK is generated from the Arcubase control-plane Swagger artifact and
must not define request/response contract types by hand.

Arcubase upstream integration points:

- `platform/src/vo/control_plane_tenant.go`
- `platform/src/vo/control_plane_service_account.go`
- `platform/src/vo/control_plane_runtime_credential.go`
- `platform/types/controlplane/controlplane_swagger.json`

That means:

- Go request/response contract types live in Arcubase upstream.
- Arcubase upstream generates `platform/types/controlplane/controlplane_swagger.json`.
- This repo syncs that artifact into `contracts/controlplane_swagger.json`.
- `controlplane/types_generated.go` and `controlplane/client_generated.go` are generated from the synced contract.
- Hand-written SDK files may only contain transport/config/error infrastructure.

## Type Consistency Policy

To keep SDK types aligned with Arcubase upstream, follow these engineering rules:

1. Request/response structs are edited in Arcubase upstream, not in this SDK.
2. Arcubase upstream must regenerate `platform/types/controlplane/controlplane_swagger.json` after contract changes.
3. This repo must sync `contracts/controlplane_swagger.json` from upstream in the same rollout window.
4. This repo must run `make generate` after syncing the contract.
5. `make validate-contract` must fail if upstream artifact and local snapshot drift.
6. Breaking field changes require an explicit version bump and release note.

## Contract Sync

When working in the Arcubase monorepo:

```bash
make sync-contract ARCUBASE_ROOT=/path/to/arcubase
make generate
make validate-contract ARCUBASE_ROOT=/path/to/arcubase
```

The validation target compares:

- upstream: `platform/types/controlplane/controlplane_swagger.json`
- local snapshot: `contracts/controlplane_swagger.json`

If they differ, the SDK is out of sync and must not be released.

## Generated Files

Generated files:

- `controlplane/types_generated.go`
- `controlplane/client_generated.go`

Do not edit generated files by hand. Change Arcubase upstream contract, sync the
contract snapshot, and regenerate.
