# carthooks-controlplane-sdk-go

Arcubase internal platform/control-plane Go SDK.

This repository is intentionally separate from `github.com/carthooks/carthooks-sdk-go`.

- `carthooks-sdk-go`: external tenant/data-plane SDK
- `carthooks-controlplane-sdk-go`: internal platform/control-plane SDK

The control-plane SDK is only for trusted internal callers and uses `X-Internal-Auth`.
