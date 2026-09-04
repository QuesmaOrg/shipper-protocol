# Quesma Shipper protocol

The sole authority for the wire contract between
[Quesma Shipper](https://github.com/QuesmaOrg/quesma-shipper) and its control plane.

The Go module embeds the normative [protocol document](PROTOCOL.md), JSON Schemas, configuration
authority rulebook, and golden fixtures:

```go
import protocol "github.com/QuesmaOrg/shipper-protocol"

raw, err := fs.ReadFile(protocol.FS, "schemas/enroll-request.schema.json")
```

Run the complete repository gate with:

```sh
go test ./...
go vet ./...
go run ./cmd/protocolgen -check
```

Release tags use ordinary module versions such as `v0.1.0`. Existing historical schema `$id`
values are deliberately preserved as stable identifiers; they are not download locations.
