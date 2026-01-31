package authregistry

import (
	"fmt"
)

// RequiredUserScopesFromServices returns the stable, de-duped union of
// RequiredUserScopes declared by the given services.
//
// Services with no RequiredUserScopes are allowed and simply contribute nothing
// to the union.
func RequiredUserScopesFromServices(services []string) ([]string, error) {
	scopes, _, err := RequiredUserScopesFromServicesReport(services)
	return scopes, err
}

// RequiredUserScopesFromServicesReport returns the stable, de-duped union of
// RequiredUserScopes declared by the given services.
//
// It also returns a stable-sorted list of services that require a user token
// but do not yet declare RequiredUserScopes (unknown vs explicitly empty).
func RequiredUserScopesFromServicesReport(services []string) (scopes []string, undeclared []string, err error) {
	services = normalizeServices(services)
	var missing []string
	for _, name := range services {
		def, ok := Registry[name]
		if !ok {
			return nil, nil, fmt.Errorf("unknown service %q", name)
		}
		requiresUser := false
		for _, tt := range def.TokenTypes {
			if tt == TokenUser {
				requiresUser = true
				break
			}
		}
		if requiresUser && def.RequiredUserScopes == nil {
			missing = append(missing, name)
		}
		scopes = append(scopes, def.RequiredUserScopes...)
	}
	return uniqueSorted(scopes), uniqueSorted(missing), nil
}

// RequiresOfflineFromServices reports whether any of the given services declares
// RequiresOffline.
func RequiresOfflineFromServices(services []string) (bool, error) {
	services = normalizeServices(services)
	for _, name := range services {
		def, ok := Registry[name]
		if !ok {
			return false, fmt.Errorf("unknown service %q", name)
		}
		if def.RequiresOffline {
			return true, nil
		}
	}
	return false, nil
}
