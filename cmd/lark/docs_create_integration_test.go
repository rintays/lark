package main

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"

	"lark/internal/testutil"
)

func TestDocsCreateIntegration(t *testing.T) {
	testutil.RequireIntegration(t)

	folderToken := os.Getenv("LARK_TEST_FOLDER_TOKEN")
	if folderToken == "" {
		t.Skip("missing LARK_TEST_FOLDER_TOKEN")
	}

	title := "clawdbot-it " + time.Now().Format("20060102-150405")

	var buf bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"--json", "docs", "create", "--title", title, "--folder-id", folderToken})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("docs create failed: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json output: %v; out=%q", err, buf.String())
	}
	doc, ok := payload["document"].(map[string]any)
	if !ok {
		t.Fatalf("expected document object, got: %v", payload)
	}
	id, _ := doc["document_id"].(string)
	if id == "" {
		// fallback in case JSON keys change
		id, _ = doc["documentId"].(string)
	}
	if id == "" {
		t.Fatalf("expected non-empty document id, got: %v", doc)
	}
}
