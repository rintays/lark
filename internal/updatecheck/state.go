package updatecheck

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// State is persisted to disk to throttle update checks and avoid repeatedly
// prompting users.
//
// Note: keep this stable; it's user-facing state.
type State struct {
	LastCheckedAtUnix int64  `json:"last_checked_at_unix"`
	LatestVersion     string `json:"latest_version,omitempty"`
	LatestPrerelease  bool   `json:"latest_prerelease,omitempty"`

	// PendingLevel is set by the async checker when a newer version exists.
	// Values: "" | "patch" | "minor" | "major".
	PendingLevel string `json:"pending_level,omitempty"`
	PendingFrom  string `json:"pending_from,omitempty"`
	PendingTo    string `json:"pending_to,omitempty"`

	// LastPromptedVersion is used to avoid repeatedly intercepting commands for
	// the same upgrade.
	LastPromptedVersion string `json:"last_prompted_version,omitempty"`
	// SuppressUntilUnix suppresses prompts until this timestamp (unix seconds).
	SuppressUntilUnix int64 `json:"suppress_until_unix,omitempty"`
}

func (s *State) Suppressed(now time.Time) bool {
	if s == nil {
		return false
	}
	if s.SuppressUntilUnix <= 0 {
		return false
	}
	return now.Unix() < s.SuppressUntilUnix
}

func (s *State) ShouldCheck(now time.Time, interval time.Duration) bool {
	if s == nil {
		return true
	}
	if s.LastCheckedAtUnix <= 0 {
		return true
	}
	last := time.Unix(s.LastCheckedAtUnix, 0)
	return now.Sub(last) >= interval
}

func pathForConfig(configPath string) string {
	dir := filepath.Dir(configPath)
	return filepath.Join(dir, "update_state.json")
}

func Load(configPath string) (*State, error) {
	p := pathForConfig(configPath)
	b, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &State{}, nil
		}
		return nil, err
	}
	var s State
	if err := json.Unmarshal(b, &s); err != nil {
		// If corrupted, don't hard-fail the CLI.
		return &State{}, nil
	}
	return &s, nil
}

func Save(configPath string, s *State) error {
	if s == nil {
		return nil
	}
	p := pathForConfig(configPath)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, b, 0o644)
}
