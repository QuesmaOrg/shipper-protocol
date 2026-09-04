# Changelog

All notable changes to this module are recorded here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/). Versions follow Go module semantic
versioning. Wire versions (`/v1/`, `/v2/`) are immutable and are not tracked by this file.

## [Unreleased]

### Added

- Initial public release of the protocol module: PROTOCOL.md, JSON Schemas for `/v1/enroll`,
  `/v1/config`, and `/v2/uploads/authorize`, the configuration authority rulebook, and golden
  fixtures. All assets are embedded in `protocol.FS`.

### Changed

- `/v1/enroll` now requires a lowercase `install_id`. The v2 object-key grammar has always
  required the lowercase spelling in the `install=` segment, so an uppercase enrollment could
  never authorize an upload. A golden rejection fixture pins the rule.
