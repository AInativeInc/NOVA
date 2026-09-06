# NOVA

NOVA now includes scaffolding for:
- `likeness-manager` (real-person digital double + consent lifecycle)
- `character-ownership` (multi-creator ownership splits)
- `licensing-engine` (consent/usage validation)
- `royalty-service` (automatic royalty split calculations)

## Key Artifacts

- gRPC contracts: `/api/proto/*.proto`
- Clean-architecture domain/usecase skeletons: `/internal/*`
- GWAV integration types + adapter contracts: `/pkg/gwav`, `/internal/licensingengine/adapter`
- Postgres migration (consent, ownership, royalties, immutable ledger):
  `/deploy/migrations/0001_likeness_ownership_royalty.sql`
- Kubernetes and Helm deployment templates: `/deploy/kubernetes`, `/deploy/helm/nova`
- Example runtime config: `/configs/example.yaml`

## Testing

Run targeted tests:

```bash
go test ./internal/licensingengine/usecase ./internal/likenessmanager/usecase ./internal/characterownership/usecase ./internal/royaltyservice/usecase
```
