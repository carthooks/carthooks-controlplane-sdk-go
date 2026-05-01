ARCUBASE_ROOT ?=

.PHONY: test sync-contract validate-contract

test:
	go test ./...

sync-contract:
	./scripts/sync-contract.sh "$(ARCUBASE_ROOT)"

validate-contract:
	./scripts/validate-contract.sh "$(ARCUBASE_ROOT)"
