package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"lark/internal/config"
	"lark/internal/larksdk"
	"lark/internal/output"
	"lark/internal/testutil"
)

func TestTasklistsCreateCommand(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/open-apis/task/v2/tasklists" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["name"] != "List A" {
			t.Fatalf("unexpected name: %+v", payload)
		}
		members := payload["members"].([]any)
		if len(members) != 2 {
			t.Fatalf("unexpected members: %+v", members)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 0,
			"msg":  "success",
			"data": map[string]any{
				"tasklist": map[string]any{
					"guid": "l1",
					"name": "List A",
					"url":  "https://applink",
				},
			},
		})
	})
	httpClient, baseURL := testutil.NewTestClient(handler)

	var buf bytes.Buffer
	state := &appState{
		TokenType: "tenant",
		Config: &config.Config{
			AppID:                      "app",
			AppSecret:                  "secret",
			BaseURL:                    baseURL,
			TenantAccessToken:          "tenant-token",
			TenantAccessTokenExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
		},
		Printer: output.Printer{Writer: &buf},
	}
	sdkClient, err := larksdk.New(state.Config, larksdk.WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("sdk client error: %v", err)
	}
	state.SDK = sdkClient

	cmd := newTasklistsCmd(state)
	cmd.SetArgs([]string{
		"create",
		"--name", "List A",
		"--editor", "ou_1",
		"--viewer", "ou_2",
	})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("tasklists create error: %v", err)
	}
	if !strings.Contains(buf.String(), "l1\tList A") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestTasklistsUpdateCommand(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/open-apis/task/v2/tasklists/l1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		updateFields := payload["update_fields"].([]any)
		if len(updateFields) != 1 || updateFields[0] != "name" {
			t.Fatalf("unexpected update_fields: %+v", updateFields)
		}
		tasklistPayload := payload["tasklist"].(map[string]any)
		if tasklistPayload["name"] != "List B" {
			t.Fatalf("unexpected tasklist payload: %+v", tasklistPayload)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 0,
			"msg":  "success",
			"data": map[string]any{
				"tasklist": map[string]any{
					"guid": "l1",
					"name": "List B",
					"url":  "https://applink",
				},
			},
		})
	})
	httpClient, baseURL := testutil.NewTestClient(handler)

	var buf bytes.Buffer
	state := &appState{
		TokenType: "tenant",
		Config: &config.Config{
			AppID:                      "app",
			AppSecret:                  "secret",
			BaseURL:                    baseURL,
			TenantAccessToken:          "tenant-token",
			TenantAccessTokenExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
		},
		Printer: output.Printer{Writer: &buf},
	}
	sdkClient, err := larksdk.New(state.Config, larksdk.WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("sdk client error: %v", err)
	}
	state.SDK = sdkClient

	cmd := newTasklistsCmd(state)
	cmd.SetArgs([]string{
		"update", "l1",
		"--name", "List B",
	})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("tasklists update error: %v", err)
	}
	if !strings.Contains(buf.String(), "l1\tList B") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}
