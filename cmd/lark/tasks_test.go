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

func TestTasksCreateCommand(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/open-apis/task/v2/tasks" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("user_id_type"); got != "open_id" {
			t.Fatalf("unexpected user_id_type: %s", got)
		}
		if r.Header.Get("Authorization") != "Bearer tenant-token" {
			t.Fatalf("unexpected authorization: %s", r.Header.Get("Authorization"))
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["summary"] != "Task A" {
			t.Fatalf("unexpected summary: %+v", payload)
		}
		if payload["description"] != "Desc" {
			t.Fatalf("unexpected description: %+v", payload)
		}
		due := payload["due"].(map[string]any)
		if due["timestamp"] != "1700000000000" {
			t.Fatalf("unexpected due timestamp: %+v", due)
		}
		if due["is_all_day"] != true {
			t.Fatalf("unexpected due is_all_day: %+v", due)
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
				"task": map[string]any{
					"guid":    "t1",
					"summary": "Task A",
					"status":  "todo",
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

	cmd := newTasksCmd(state)
	cmd.SetArgs([]string{
		"create",
		"--summary", "Task A",
		"--description", "Desc",
		"--due", "1700000000",
		"--due-all-day",
		"--assignee", "ou_1",
		"--follower", "ou_2",
	})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("tasks create error: %v", err)
	}
	if !strings.Contains(buf.String(), "t1\tTask A") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestTasksListCommand(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/open-apis/task/v2/tasks" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("page_size"); got != "2" {
			t.Fatalf("unexpected page_size: %s", got)
		}
		if got := r.URL.Query().Get("completed"); got != "true" {
			t.Fatalf("unexpected completed: %s", got)
		}
		if got := r.URL.Query().Get("type"); got != "my_tasks" {
			t.Fatalf("unexpected type: %s", got)
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
					{"guid": "t1", "summary": "Task A", "status": "done"},
					{"guid": "t2", "summary": "Task B", "status": "done"},
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
			AppID:                      "app",
			AppSecret:                  "secret",
			BaseURL:                    baseURL,
			TenantAccessToken:          "tenant-token",
			TenantAccessTokenExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
		},
		Printer: output.Printer{Writer: &buf},
	}
	withUserAccount(state.Config, defaultUserAccountName, "user-token", "", time.Now().Add(2*time.Hour).Unix(), "")
	sdkClient, err := larksdk.New(state.Config, larksdk.WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("sdk client error: %v", err)
	}
	state.SDK = sdkClient

	cmd := newTasksCmd(state)
	cmd.SetArgs([]string{
		"list",
		"--limit", "2",
		"--page-size", "2",
		"--completed",
	})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("tasks list error: %v", err)
	}
	if !strings.Contains(buf.String(), "t1\tTask A") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestTasksUpdateCommand(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/open-apis/task/v2/tasks/t1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		updateFields := payload["update_fields"].([]any)
		if len(updateFields) != 2 {
			t.Fatalf("unexpected update_fields: %+v", updateFields)
		}
		taskPayload := payload["task"].(map[string]any)
		if taskPayload["summary"] != "New" {
			t.Fatalf("unexpected summary: %+v", taskPayload)
		}
		if _, ok := taskPayload["due"]; !ok {
			t.Fatalf("expected due to be cleared")
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 0,
			"msg":  "success",
			"data": map[string]any{
				"task": map[string]any{
					"guid":    "t1",
					"summary": "New",
					"status":  "todo",
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

	cmd := newTasksCmd(state)
	cmd.SetArgs([]string{
		"update", "t1",
		"--summary", "New",
		"--clear-due",
	})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("tasks update error: %v", err)
	}
	if !strings.Contains(buf.String(), "t1\tNew") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}
