# Contributing

Thanks for helping improve the Quesma Shipper protocol. This repository is the contract between
separately maintained client and server implementations, so even small edits can have compatibility
or security consequences. Participation is governed by our [Code of Conduct](CODE_OF_CONDUCT.md).

## Before proposing a change

- Open an issue before making a breaking change or adding a protocol capability. Describe the use
  case and the compatibility impact rather than only the proposed wire shape.
- Never change the meaning or accepted shape of a released wire version. Add a new endpoint and
  schema version for a breaking change.
- Treat schema `$id` values as stable identifiers. Do not change identifiers published in a release;
  introduce a new protocol version and identifier when a breaking change is required.
- Do not include production payloads, credentials, hostnames, customer data, or presigned URLs in
  issues, pull requests, or fixtures. All fixtures must be minimal synthetic examples.

## Making a change

Keep the normative document, schemas, and fixtures in sync. If you edit `authority.json`, regenerate
its table in `PROTOCOL.md`:

```sh
go run ./cmd/protocolgen
```

Run the same checks as CI before opening a pull request:

```sh
go test ./...
go vet ./...
go run ./cmd/protocolgen -check
```

Explain compatibility and security implications in the pull request. Maintainers will run the
private Shipper and control-plane integration suites; external contributors are not expected to have
access to those repositories.

Unless explicitly stated otherwise, contributions submitted for inclusion are licensed under the
repository's [Apache License 2.0](LICENSE), as described in section 5 of that license.
