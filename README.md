# Quesma Shipper protocol

The sole authority for the wire contract between Quesma Shipper and its control plane.

This repository contains the protocol specification and consumable contract assets—not a runnable
shipper or control plane. Those implementations live in separate repositories and may not be
publicly accessible. You can still inspect, test, and propose changes to the complete wire contract
here; maintainers run the downstream implementation contract tests before release.

The Go module embeds the normative [protocol document](PROTOCOL.md), JSON Schemas, configuration
authority rulebook, and golden fixtures:

```go
import protocol "github.com/QuesmaOrg/shipper-protocol"

raw, err := protocol.FS.ReadFile("schemas/enroll-request.schema.json")
```

The checked-in fixtures are synthetic test vectors. Token-shaped values, signatures, URLs, and the
documented Ed25519 seed are intentionally public and must never be used as production credentials.

## Development

Install the Go version declared in [`go.mod`](go.mod), then run the complete repository gate:

```sh
go test ./...
go vet ./...
go run ./cmd/protocolgen -check
```

See [`CONTRIBUTING.md`](CONTRIBUTING.md) before changing a released contract. Security reports
should follow [`SECURITY.md`](SECURITY.md).

## Versioning and license

Release tags use ordinary Go module versions such as `v0.1.0`. Wire versions are immutable after
release; breaking changes get a new endpoint and schema version. Schema `$id` values use this
repository's public namespace. Treat them as stable identifiers rather than download locations;
published identifiers must not be changed in place.

Licensed under the [Apache License 2.0](LICENSE).
