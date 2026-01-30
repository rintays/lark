package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"lark/internal/config"
	"lark/internal/larkapi"
	"lark/internal/output"
	"lark/internal/testutil"
)

func TestDriveListCommand(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/open-apis/drive/v1/files" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("page_size") != "2" {
			t.Fatalf("unexpected page_size: %s", r.URL.Query().Get("page_size"))
		}
		if r.URL.Query().Get("folder_token") != "root" {
			t.Fatalf("unexpected folder_token: %s", r.URL.Query().Get("folder_token"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 0,
			"msg":  "ok",
			"data": map[string]any{
				"files":    []map[string]any{{"token": "f1", "name": "Doc", "type": "docx", "url": "https://example.com/doc"}},
				"has_more": false,
			},
		})
	})
	httpClient, baseURL := testutil.NewTestClient(handler)

	var buf bytes.Buffer
	state := &appState{
		Config: &config.Config{
			AppID:                      "app",
			AppSecret:                  "secret",
			BaseURL:                    baseURL,
			TenantAccessToken:          "token",
			TenantAccessTokenExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
		},
		Printer: output.Printer{Writer: &buf},
		Client:  &larkapi.Client{BaseURL: baseURL, HTTPClient: httpClient},
	}

	cmd := newDriveCmd(state)
	cmd.SetArgs([]string{"list", "--folder-id", "root", "--limit", "2"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("drive list error: %v", err)
	}

	if !strings.Contains(buf.String(), "f1\tDoc\tdocx\thttps://example.com/doc") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestDriveSearchCommand(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/open-apis/drive/v1/files/search" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["query"] != "budget" {
			t.Fatalf("unexpected query: %+v", payload)
		}
		if payload["page_size"].(float64) != 2 {
			t.Fatalf("unexpected page_size: %+v", payload)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 0,
			"msg":  "ok",
			"data": map[string]any{
				"files":    []map[string]any{{"token": "f2", "name": "Budget", "type": "sheet", "url": "https://example.com/sheet"}},
				"has_more": false,
			},
		})
	})
	httpClient, baseURL := testutil.NewTestClient(handler)

	var buf bytes.Buffer
	state := &appState{
		Config: &config.Config{
			AppID:                      "app",
			AppSecret:                  "secret",
			BaseURL:                    baseURL,
			TenantAccessToken:          "token",
			TenantAccessTokenExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
		},
		Printer: output.Printer{Writer: &buf},
		Client:  &larkapi.Client{BaseURL: baseURL, HTTPClient: httpClient},
	}

	cmd := newDriveCmd(state)
	cmd.SetArgs([]string{"search", "--query", "budget", "--limit", "2"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("drive search error: %v", err)
	}

	if !strings.Contains(buf.String(), "f2\tBudget\tsheet\thttps://example.com/sheet") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestDriveGetCommand(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/open-apis/drive/v1/files/f1":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"code": 0,
				"msg":  "ok",
				"data": map[string]any{
					"file": map[string]any{"token": "f1", "name": "Doc", "type": "docx", "url": "https://example.com/doc"},
				},
			})
		case "/open-apis/drive/v1/files/f2":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"code": 0,
				"msg":  "ok",
				"data": map[string]any{
					"file": map[string]any{"token": "f2", "name": "Sheet", "type": "sheet", "url": "https://example.com/sheet"},
				},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	})
	httpClient, baseURL := testutil.NewTestClient(handler)

	cases := []struct {
		name     string
		args     []string
		expected string
	}{
		{name: "arg", args: []string{"get", "f1"}, expected: "f1\tDoc\tdocx\thttps://example.com/doc"},
		{name: "flag", args: []string{"get", "--file-token", "f2"}, expected: "f2\tSheet\tsheet\thttps://example.com/sheet"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			state := &appState{
				Config: &config.Config{
					AppID:                      "app",
					AppSecret:                  "secret",
					BaseURL:                    baseURL,
					TenantAccessToken:          "token",
					TenantAccessTokenExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
				},
				Printer: output.Printer{Writer: &buf},
				Client:  &larkapi.Client{BaseURL: baseURL, HTTPClient: httpClient},
			}

			cmd := newDriveCmd(state)
			cmd.SetArgs(tc.args)
			if err := cmd.Execute(); err != nil {
				t.Fatalf("drive get error: %v", err)
			}

			if !strings.Contains(buf.String(), tc.expected) {
				t.Fatalf("unexpected output: %q", buf.String())
			}
		})
	}
}

func TestDriveURLsCommand(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/open-apis/drive/v1/files/f1":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"code": 0,
				"msg":  "ok",
				"data": map[string]any{
					"file": map[string]any{"token": "f1", "name": "Doc", "url": "https://example.com/doc"},
				},
			})
		case "/open-apis/drive/v1/files/f2":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"code": 0,
				"msg":  "ok",
				"data": map[string]any{
					"file": map[string]any{"token": "f2", "name": "Sheet", "url": "https://example.com/sheet"},
				},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	})
	httpClient, baseURL := testutil.NewTestClient(handler)

	var buf bytes.Buffer
	state := &appState{
		Config: &config.Config{
			AppID:                      "app",
			AppSecret:                  "secret",
			BaseURL:                    baseURL,
			TenantAccessToken:          "token",
			TenantAccessTokenExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
		},
		Printer: output.Printer{Writer: &buf},
		Client:  &larkapi.Client{BaseURL: baseURL, HTTPClient: httpClient},
	}

	cmd := newDriveCmd(state)
	cmd.SetArgs([]string{"urls", "f1", "f2"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("drive urls error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "f1\thttps://example.com/doc\tDoc") {
		t.Fatalf("unexpected output: %q", output)
	}
	if !strings.Contains(output, "f2\thttps://example.com/sheet\tSheet") {
		t.Fatalf("unexpected output: %q", output)
	}
}
