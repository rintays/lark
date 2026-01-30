package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newAuthCmd(state *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Fetch and cache a tenant access token",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := ensureTenantToken(context.Background(), state)
			if err != nil {
				return err
			}
			payload := map[string]any{
				"tenant_access_token": token,
				"expires_at":          state.Config.TenantAccessTokenExpiresAt,
			}
			return state.Printer.Print(payload, fmt.Sprintf("tenant_access_token: %s", token))
		},
	}
	return cmd
}
