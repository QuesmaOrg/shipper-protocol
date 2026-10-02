package protocol

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// ValidateConfigDocument validates a JSON representation of a served YAML document
// against the canonical authoring schema. Pass the decoded object, not the base64
// /v1/config response envelope. The empty object is a valid layer of defaults.
//
// This enforces the schema's custom Go duration formats as well as its shape.
// It does not replace client resolution: source IDs, rule packs, root ceilings,
// recipient checksums, and interactions with machine-local layers remain client
// responsibilities. Clients may accept older documents outside this authoring
// profile; this function is intended for control-plane writes, not client reads.
func ValidateConfigDocument(rawJSON []byte) error {
	schema, err := configDocumentSchema()
	if err != nil {
		return err
	}
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(rawJSON))
	if err != nil {
		return fmt.Errorf("config document: %w", err)
	}
	return schema.Validate(doc)
}

var configDocumentSchema = sync.OnceValues(func() (*jsonschema.Schema, error) {
	const name = "schemas/config-document.schema.json"
	raw, err := FS.ReadFile(name)
	if err != nil {
		return nil, err
	}
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	compiler := jsonschema.NewCompiler()
	compiler.AssertFormat()
	compiler.RegisterFormat(&jsonschema.Format{
		Name:     "go-positive-duration",
		Validate: durationAtLeast(time.Nanosecond),
	})
	compiler.RegisterFormat(&jsonschema.Format{
		Name:     "go-schedule-duration",
		Validate: durationAtLeast(time.Minute),
	})
	if err := compiler.AddResource(name, doc); err != nil {
		return nil, err
	}
	return compiler.Compile(name)
})

func durationAtLeast(minimum time.Duration) func(any) error {
	return func(value any) error {
		s, ok := value.(string)
		if !ok {
			return nil // The schema's type keyword handles non-strings.
		}
		duration, err := time.ParseDuration(s)
		if err != nil {
			return err
		}
		if duration < minimum {
			return fmt.Errorf("must be at least %s", minimum)
		}
		return nil
	}
}
