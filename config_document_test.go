package protocol_test

import (
	"encoding/json"
	"io/fs"
	"path"
	"strings"
	"testing"

	protocol "github.com/QuesmaOrg/shipper-protocol"
)

func TestConfigDocumentFixtures(t *testing.T) {
	paths, err := fs.Glob(protocol.FS, "fixtures/v1/config-document/*.json")
	if err != nil || len(paths) == 0 {
		t.Fatalf("config document fixtures: %v (%d files)", err, len(paths))
	}
	for _, assetPath := range paths {
		t.Run(path.Base(assetPath), func(t *testing.T) {
			raw, err := protocol.FS.ReadFile(assetPath)
			if err != nil {
				t.Fatal(err)
			}
			err = protocol.ValidateConfigDocument(raw)
			bad := strings.HasPrefix(path.Base(assetPath), "bad-")
			if bad && err == nil {
				t.Fatal("invalid authored document accepted")
			}
			if !bad && err != nil {
				t.Fatalf("valid authored document rejected: %v", err)
			}
		})
	}
}

func TestConfigDocumentRejectsMalformedJSON(t *testing.T) {
	for _, raw := range []string{"", "{", "{} {}"} {
		if err := protocol.ValidateConfigDocument([]byte(raw)); err == nil {
			t.Errorf("accepted %q", raw)
		}
	}
}

// Every rulebook row must be considered when changing the authoring schema. The
// probes distinguish portable authority checks from checks only a client with
// its compiled catalog and local configuration can perform.
func TestConfigDocumentAuthority(t *testing.T) {
	type probe struct {
		class  string
		accept string
		reject string
	}
	probes := map[string]probe{
		"issued_at / org":   {"envelope", `{"issued_at":"2026-10-01T12:00:00Z","org":"example"}`, `{"org":false}`},
		"config_version":    {"set", `{"config_version":1}`, `{"config_version":2}`},
		"mode.schedule":     {"set", `{"mode":{"schedule":"1m"}}`, `{"mode":{"schedule":"59s"}}`},
		"max_files_per_run": {"set", `{"max_files_per_run":512}`, `{"max_files_per_run":"512"}`},
		"drain_deadline":    {"set", `{"drain_deadline":"1h"}`, `{"drain_deadline":"forever"}`},
		"state_dir":         {"refused", `{}`, `{"state_dir":"/tmp/state"}`},
		"upload_targets":    {"refused", `{}`, `{"upload_targets":[{"origin":"https://example.test"}]}`},
		"send.sink / bucket / prefix / region / path": {"ignored", `{"send":{"sink":"s3","bucket":"ignored","region":"ignored","prefix":"ignored","path":"ignored"}}`, `{"send":false}`},
		"scrub.rule_packs":                            {"union", `{"scrub":{"rule_packs":["gitleaks-core"]}}`, `{"scrub":{"rule_packs":[false]}}`},
		"scrub.secret_key_names":                      {"union", `{"scrub":{"secret_key_names":["EXAMPLE_TOKEN"]}}`, `{"scrub":{"secret_key_names":false}}`},
		"structural_exempt":                           {"set", `{"structural_exempt":{"claude-code-transcripts":["example"]}}`, `{"structural_exempt":{"claude-code-transcripts":false}}`},
		"encryption.additional_recipients":            {"union", `{"encryption":{"additional_recipients":["age1cpx4grz9j4fkn36cfurggwcg4l0da5fyqadl8fwagtcwy55gt44qlclfa5"]}}`, `{"encryption":{"additional_recipients":["not-an-age-key"]}}`},
		"encryption.include_install_recipient":        {"set", `{"encryption":{"include_install_recipient":true}}`, `{"encryption":{"include_install_recipient":false}}`},
		"sources[].enabled":                           {"set", `{"sources":[{"id":"claude-code-transcripts","enabled":false}]}`, `{"sources":[{"id":"claude-code-transcripts","enabled":"false"}]}`},
		"sources[].roots / include / exclude":         {"set", `{"sources":[{"id":"claude-code-transcripts","roots":["~/.claude/projects"],"include":["*.jsonl"],"exclude":["private/**"]}]}`, `{"sources":[{"id":"claude-code-transcripts","roots":false}]}`},
		"sources[].max_file_bytes":                    {"set", `{"sources":[{"id":"claude-code-transcripts","max_file_bytes":1024}]}`, `{"sources":[{"id":"claude-code-transcripts","max_file_bytes":"1024"}]}`},
		"sources[].enrichers{}":                       {"off-only", `{"sources":[{"id":"example","enrichers":{"example":false}}]}`, `{"sources":[{"id":"example","enrichers":{"example":true}}]}`},
		// Deliberately accepted by the portable schema: only the client knows
		// which source IDs its compiled catalog contains.
		"a source id not in the compiled catalog": {"refused", `{"sources":[{"id":"build-specific-source"}]}`, ""},
		"crash_report.enabled / dsn":              {"ignored", `{"crash_report":{"enabled":false,"dsn":"ignored"}}`, `{"crash_report":false}`},
		"autoupdate.enabled":                      {"off-only", `{"autoupdate":{"enabled":false}}`, `{"autoupdate":{"enabled":true}}`},
		"telemetry_endpoint":                      {"envelope", `{"telemetry_endpoint":"/v1/telemetry"}`, `{"telemetry_endpoint":"https://example.test/v1/telemetry"}`},
	}
	raw, err := protocol.FS.ReadFile("authority.json")
	if err != nil {
		t.Fatal(err)
	}
	var authority struct {
		Fields []struct {
			Field string `json:"field"`
			Class string `json:"server_may"`
		} `json:"fields"`
	}
	if err := json.Unmarshal(raw, &authority); err != nil {
		t.Fatal(err)
	}
	for _, field := range authority.Fields {
		t.Run(field.Field, func(t *testing.T) {
			p, ok := probes[field.Field]
			if !ok {
				t.Fatal("authority row needs an authoring-schema probe")
			}
			if p.class != field.Class {
				t.Fatalf("authority changed from %s to %s; review the authoring schema", p.class, field.Class)
			}
			if err := protocol.ValidateConfigDocument([]byte(p.accept)); err != nil {
				t.Fatalf("authority-compatible shape rejected: %v", err)
			}
			if p.reject != "" && protocol.ValidateConfigDocument([]byte(p.reject)) == nil {
				t.Fatal("authority-incompatible shape accepted")
			}
		})
		delete(probes, field.Field)
	}
	if len(probes) != 0 {
		t.Fatalf("schema probes without authority rows: %v", probes)
	}
}
