package protocol_test

import (
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

func TestConfigDocumentValidation(t *testing.T) {
	for _, raw := range []string{`{}`, `{"mode":{"schedule":"1m"},"drain_deadline":"1ns"}`} {
		if err := protocol.ValidateConfigDocument([]byte(raw)); err != nil {
			t.Errorf("rejected %s: %v", raw, err)
		}
	}
	for _, raw := range []string{
		`{`, `{} {}`, `null`,
		`{"max_files_per_run":1.5}`, `{"max_files_per_run":null}`,
		`{"drain_deadline":"0s"}`, `{"drain_deadline":"-1s"}`, `{"drain_deadline":"999999999999999999h"}`,
		`{"mode":{"schedule":"30s"}}`, `{"mode":{"schedule":"forever"}}`,
		`{"sources":"hello"}`, `{"sources":[{}]}`, `{"sources":[{"id":"example","enabled":"yes"}]}`,
		`{"unknown":true}`, `{"mode":{"unknown":true}}`,
	} {
		if err := protocol.ValidateConfigDocument([]byte(raw)); err == nil {
			t.Errorf("accepted %s", raw)
		}
	}
}
