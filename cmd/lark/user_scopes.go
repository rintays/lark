package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type userOAuthScopeOptions struct {
	Scopes        string
	ScopesSet     bool
	Services      []string
	ServicesSet   bool
	Readonly      bool
	DriveScope    string
	DriveScopeSet bool
}

type serviceScopeSet struct {
	Full     []string
	Readonly []string
}

var userServiceScopes = map[string]serviceScopeSet{
	"drive":  {Full: []string{"drive:drive"}, Readonly: []string{"drive:drive:readonly"}},
	"docs":   {Full: []string{"drive:drive"}, Readonly: []string{"drive:drive:readonly"}},
	"sheets": {Full: []string{"drive:drive"}, Readonly: []string{"drive:drive:readonly"}},
}

var userServiceAliases = map[string][]string{
	"all":  {"drive", "docs", "sheets"},
	"user": {"drive", "docs", "sheets"},
}

var defaultUserServices = []string{"drive"}

func resolveUserOAuthScopes(state *appState, opts userOAuthScopeOptions) ([]string, string, error) {
	if opts.ScopesSet {
		scopes := normalizeScopes(parseScopeList(opts.Scopes))
		if len(scopes) == 0 {
			return nil, "", errors.New("scopes must not be empty")
		}
		scopes = ensureOfflineAccess(scopes)
		return scopes, "flag", nil
	}

	if opts.ServicesSet || opts.Readonly || opts.DriveScopeSet {
		services := opts.Services
		if len(services) == 0 {
			services = defaultUserServices
		}
		scopes, err := scopesFromServices(services, opts.Readonly, opts.DriveScope)
		if err != nil {
			return nil, "", err
		}
		return ensureOfflineAccess(scopes), "services", nil
	}

	if state != nil && state.Config != nil && len(state.Config.UserScopes) > 0 {
		scopes := normalizeScopes(state.Config.UserScopes)
		return ensureOfflineAccess(scopes), "config", nil
	}

	return []string{defaultUserOAuthScope}, "default", nil
}

func scopesFromServices(services []string, readonly bool, driveScope string) ([]string, error) {
	driveScope = strings.ToLower(strings.TrimSpace(driveScope))
	if driveScope != "" {
		switch driveScope {
		case "full", "readonly":
		case "file":
			return nil, errors.New("drive-scope file is not supported; use full or readonly")
		default:
			return nil, fmt.Errorf("invalid drive-scope %q (use full or readonly)", driveScope)
		}
	}
	if readonly {
		if driveScope != "" {
			return nil, errors.New("drive-scope cannot be combined with --readonly")
		}
		driveScope = "readonly"
	}
	if driveScope == "" {
		driveScope = "full"
	}

	var scopes []string
	for _, svc := range expandServiceAliases(services) {
		set, ok := userServiceScopes[svc]
		if !ok {
			return nil, fmt.Errorf("unknown service %q (use `lark auth user services` to list supported services)", svc)
		}
		if driveScope == "readonly" {
			scopes = append(scopes, set.Readonly...)
		} else {
			scopes = append(scopes, set.Full...)
		}
	}
	return normalizeScopes(scopes), nil
}

func parseScopeList(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	fields := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\t' || r == ' '
	})
	return normalizeScopes(fields)
}

func normalizeScopes(scopes []string) []string {
	seen := make(map[string]struct{}, len(scopes))
	out := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		scope = strings.TrimSpace(scope)
		if scope == "" {
			continue
		}
		if _, ok := seen[scope]; ok {
			continue
		}
		seen[scope] = struct{}{}
		out = append(out, scope)
	}
	return out
}

func ensureOfflineAccess(scopes []string) []string {
	for _, scope := range scopes {
		if scope == defaultUserOAuthScope {
			return scopes
		}
	}
	return append([]string{defaultUserOAuthScope}, scopes...)
}

func joinScopes(scopes []string) string {
	return strings.Join(normalizeScopes(scopes), " ")
}

func parseServicesList(raw []string) []string {
	parts := make([]string, 0, len(raw))
	for _, entry := range raw {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		for _, part := range strings.FieldsFunc(entry, func(r rune) bool { return r == ',' || r == ' ' || r == '\t' || r == '\n' }) {
			part = strings.ToLower(strings.TrimSpace(part))
			if part != "" {
				parts = append(parts, part)
			}
		}
	}
	return normalizeServices(parts)
}

func normalizeServices(services []string) []string {
	seen := make(map[string]struct{}, len(services))
	out := make([]string, 0, len(services))
	for _, svc := range services {
		if svc == "" {
			continue
		}
		svc = strings.ToLower(strings.TrimSpace(svc))
		if svc == "" {
			continue
		}
		if _, ok := seen[svc]; ok {
			continue
		}
		seen[svc] = struct{}{}
		out = append(out, svc)
	}
	return out
}

func expandServiceAliases(services []string) []string {
	expanded := make([]string, 0, len(services))
	for _, svc := range services {
		if alias, ok := userServiceAliases[svc]; ok {
			expanded = append(expanded, alias...)
			continue
		}
		expanded = append(expanded, svc)
	}
	return normalizeServices(expanded)
}

func listUserServices() []string {
	services := make([]string, 0, len(userServiceScopes))
	for svc := range userServiceScopes {
		services = append(services, svc)
	}
	sort.Strings(services)
	return services
}
