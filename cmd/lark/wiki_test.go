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

func TestWikiSearchCommand(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/open-apis/wiki/v1/nodes/search" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("page_size") != "2" {
			t.Fatalf("unexpected page_size: %s", r.URL.Query().Get("page_size"))
		}
		if r.Header.Get("Authorization") != "Bearer user-token" {
			t.Fatalf("unexpected authorization: %s", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["query"] != "design" {
			t.Fatalf("unexpected query: %+v", payload)
		}
		if payload["space_id"] != "space" {
			t.Fatalf("unexpected space_id: %+v", payload)
		}
		if payload["node_id"] != "node" {
			t.Fatalf("unexpected node_id: %+v", payload)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 0,
			"msg":  "ok",
			"data": map[string]any{
				"items":    []map[string]any{{"node_id": "n1", "title": "Design", "url": "https://example.com/wiki"}},
				"has_more": false,
			},
		})
	})
	httpClient, baseURL := testutil.NewTestClient(handler)

	var buf bytes.Buffer
	state := &appState{
		Config: &config.Config{
			AppID:                    "app",
			AppSecret:                "secret",
			BaseURL:                  baseURL,
			UserAccessToken:          "user-token",
			UserAccessTokenExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
		},
		Printer: output.Printer{Writer: &buf},
	}
	sdkClient, err := larksdk.New(state.Config, larksdk.WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("sdk client error: %v", err)
	}
	state.SDK = sdkClient

	cmd := newWikiCmd(state)
	cmd.SetArgs([]string{"search", "--query", "design", "--space-id", "space", "--node-id", "node", "--limit", "2"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("wiki search error: %v", err)
	}

	if !strings.Contains(buf.String(), "n1\tDesign\thttps://example.com/wiki") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}
