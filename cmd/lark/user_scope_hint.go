package main

import (
	"fmt"
	"regexp"
	"strings"
)

var scopeBracketPattern = regexp.MustCompile(`\[(.*?)\]`)

func withUserScopeHint(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	if strings.Contains(msg, "Re-authorize with:") {
		return err
	}
	scopes := extractScopesFromErrorMessage(msg)
	if len(scopes) == 0 {
		return err
	}
	scopes = ensureOfflineAccess(scopes)
	scopeArg := strings.Join(scopes, " ")
	hint := fmt.Sprintf("Missing user OAuth scopes: %s.\nRe-authorize with:\n  lark auth user login --scopes %q --force-consent", strings.Join(scopes, ", "), scopeArg)
	return fmt.Errorf("%s\n%s", msg, hint)
}

func extractScopesFromErrorMessage(msg string) []string {
	matches := scopeBracketPattern.FindAllStringSubmatch(msg, -1)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		scopes := normalizeScopes(parseScopeList(match[1]))
		if len(scopes) == 0 {
			continue
		}
		if containsScopeToken(scopes) {
			return scopes
		}
	}
	return nil
}

func containsScopeToken(scopes []string) bool {
	for _, scope := range scopes {
		if strings.Contains(scope, ":") {
			return true
		}
	}
	return false
}
