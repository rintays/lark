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

func TestTasklistsListArgs(t *testing.T) {
	state := &appState{
		Config: &config.Config{},
	}
	cmd := newTasklistsCmd(state)
	cmd.SetArgs([]string{"list", "extra"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTasklistsListLimitClamp(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/open-apis/task/v2/tasklists" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("page_size") != "500" {
			t.Fatalf("unexpected page_size: %s", r.URL.Query().Get("page_size"))
		}
		if r.Header.Get("Authorization") != "Bearer user-token" {
			t.Fatalf("unexpected authorization: %s", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 0,
			"msg":  "success",
			"data": map[string]any{
				"items": []map[string]any{
					{
						"guid":       "l1",
						"name":       "List A",
						"updated_at": "2025-01-01T00:00:00Z",
						"owner": map[string]any{
							"id": "ou_1",
						},
					},
				},
				"has_more": false,
			},
		})
	})
	httpClient, baseURL := testutil.NewTestClient(handler)

	var buf bytes.Buffer
	state := &appState{
		TokenType: "user",
		Config: &config.Config{
			AppID:     "app",
			AppSecret: "secret",
			BaseURL:   baseURL,
			UserAccounts: map[string]*config.UserAccount{
				defaultUserAccountName: {
					UserAccessToken:          "user-token",
					UserAccessTokenExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
				},
			},
		},
		Printer: output.Printer{Writer: &buf},
	}
	sdkClient, err := larksdk.New(state.Config, larksdk.WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("sdk client error: %v", err)
	}
	state.SDK = sdkClient

	cmd := newTasklistsCmd(state)
	cmd.SetArgs([]string{"list", "--limit", "600"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("tasklists list error: %v", err)
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
