package ai

import (
	"encoding/json"
	"testing"
)

func TestExportSettingJSONOmitsZeroValues(t *testing.T) {
	data, err := ExportSettingJSON(BaseAgentSetting{
		ProviderID:      1,
		Models:          []string{"m1"},
		ReasoningEffort: ReasoningEffortMedium,
	})
	if err != nil {
		t.Fatalf("ExportSettingJSON error: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}
	if _, ok := got["sysPrompt"]; ok {
		t.Fatalf("zero sysPrompt should be omitted: %s", data)
	}
	if got["providerID"] != float64(1) {
		t.Fatalf("providerID missing: %s", data)
	}
}

func TestImportSettingJSONDoesNotOverwriteWithZeroValues(t *testing.T) {
	base := BaseAgentSetting{
		ProviderID:      1,
		SysPrompt:       "old",
		Models:          []string{"old-model"},
		ReasoningEffort: ReasoningEffortMedium,
	}
	next, err := ImportSettingJSON(base, []byte(`{
		"providerID": 2,
		"sysPrompt": "",
		"models": [],
		"reasoningEffort": ""
	}`))
	if err != nil {
		t.Fatalf("ImportSettingJSON error: %v", err)
	}
	if next.ProviderID != 2 {
		t.Fatalf("expected provider override, got %d", next.ProviderID)
	}
	if next.SysPrompt != "old" {
		t.Fatalf("zero string should not overwrite, got %q", next.SysPrompt)
	}
	if next.Models[0] != "old-model" {
		t.Fatalf("empty slice should not overwrite, got %#v", next.Models)
	}
	if next.ReasoningEffort != ReasoningEffortMedium {
		t.Fatalf("zero reasoning should not overwrite, got %q", next.ReasoningEffort)
	}
}

func TestExportSettingDoc(t *testing.T) {
	docs, err := ExportSettingDoc(BaseAgentSetting{})
	if err != nil {
		t.Fatalf("ExportSettingDoc error: %v", err)
	}
	if len(docs) != 4 {
		t.Fatalf("expected 4 docs, got %d", len(docs))
	}
	if docs[0].Name != "providerID" || docs[0].JSONType != "number" {
		t.Fatalf("unexpected first doc: %#v", docs[0])
	}
	if docs[2].Name != "models" || !docs[2].IsArray {
		t.Fatalf("models doc should be array: %#v", docs[2])
	}
}

func TestSettingRejectsUnsupportedFieldType(t *testing.T) {
	type badSetting struct {
		Nested struct {
			Name string
		}
	}
	if _, err := ExportSettingJSON(badSetting{}); err == nil {
		t.Fatal("expected unsupported field type error")
	}
}
