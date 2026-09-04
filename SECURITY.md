# Security policy

## Reporting a vulnerability

Do not open a public issue for a suspected vulnerability or include real credentials, customer data,
presigned URLs, or production payloads in a report. Email `contact@quesma.com` with the subject
`Security: shipper-protocol` and include:

- the affected protocol or module version;
- the impact and conditions required to reproduce it;
- minimal reproduction steps using synthetic data; and
- any suggested mitigation, if known.

We will acknowledge the report and coordinate disclosure after assessing it. Please allow time for
both protocol peers to be patched before publishing details.

## Supported versions

Security fixes are made on the latest released module version. The default branch is development
code and is not a supported release. Because wire versions are immutable after release, a fix that
cannot be made compatibly will use a new endpoint and schema version.
