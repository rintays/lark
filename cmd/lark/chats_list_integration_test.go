package main

import (
	"bytes"
	"encoding/json"
	"testing"

	"lark/internal/testutil"
)

func TestChatsListIntegration(t *testing.T) {
	testutil.RequireIntegration(t)

	var buf bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"--json", "chats", "list", "--limit", "1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("chats list failed: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json output: %v; out=%q", err, buf.String())
	}
	chats, ok := payload["chats"]
	if !ok {
		t.Fatalf("expected chats in output, got: %v", payload)
	}
	if _, ok := chats.([]any); !ok {
		t.Fatalf("expected chats to be an array, got: %T (%v)", chats, chats)
	}
}
