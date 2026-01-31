package main

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"

	"lark/internal/testutil"
)

func TestMailSendIntegration(t *testing.T) {
	testutil.RequireIntegration(t)

	to := os.Getenv("LARK_TEST_MAIL_TO")
	if to == "" {
		t.Skip("missing LARK_TEST_MAIL_TO")
	}

	// Requires a user token. We intentionally rely on the CLI's normal resolution:
	// --user-access-token / env LARK_USER_ACCESS_TOKEN / stored token from `lark auth user login`.

	subject := "clawdbot-it " + time.Now().Format("20060102-150405")

	var buf bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"--json", "mail", "send", "--to", to, "--subject", subject, "--text", "ping"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("mail send failed: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json output: %v; out=%q", err, buf.String())
	}
	id, _ := payload["message_id"].(string)
	if id == "" {
		t.Fatalf("expected non-empty message_id, got: %v", payload)
	}
}
