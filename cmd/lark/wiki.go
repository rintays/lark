package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"lark/internal/larksdk"
)

const maxWikiSearchPageSize = 50

func newWikiCmd(state *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wiki",
		Short: "Manage Wiki resources",
	}
	cmd.AddCommand(newWikiSearchCmd(state))
	return cmd
}

func newWikiSearchCmd(state *appState) *cobra.Command {
	var query string
	var spaceID string
	var nodeID string
	var limit int
	var userAccessToken string
	const userTokenHint = "wiki search requires a user access token; pass --user-access-token, set LARK_USER_ACCESS_TOKEN, or run `lark auth user login`"

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search Wiki nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			if query == "" {
				return errors.New("query is required")
			}
			if limit <= 0 {
				return errors.New("limit must be greater than 0")
			}
			token := userAccessToken
			if token == "" {
				token = os.Getenv("LARK_USER_ACCESS_TOKEN")
			}
			if token == "" {
				var err error
				token, err = ensureUserToken(context.Background(), state)
				if err != nil || token == "" {
					return errors.New(userTokenHint)
				}
			}
			if state.SDK == nil {
				return errors.New("sdk client is required")
			}
			nodes := make([]larksdk.WikiV1Node, 0, limit)
			pageToken := ""
			remaining := limit
			for {
				pageSize := remaining
				if pageSize > maxWikiSearchPageSize {
					pageSize = maxWikiSearchPageSize
				}
				result, err := state.SDK.SearchWikiNodes(context.Background(), token, larksdk.SearchWikiNodesRequest{
					Query:     query,
					SpaceID:   spaceID,
					NodeID:    nodeID,
					PageSize:  pageSize,
					PageToken: pageToken,
				})
				if err != nil {
					return err
				}
				nodes = append(nodes, result.Items...)
				if len(nodes) >= limit || !result.HasMore {
					break
				}
				remaining = limit - len(nodes)
				pageToken = result.PageToken
				if pageToken == "" {
					break
				}
			}
			if len(nodes) > limit {
				nodes = nodes[:limit]
			}
			payload := map[string]any{"nodes": nodes}
			lines := make([]string, 0, len(nodes))
			for _, node := range nodes {
				lines = append(lines, fmt.Sprintf("%s\t%s\t%s", node.NodeID, node.Title, node.URL))
			}
			text := "no nodes found"
			if len(lines) > 0 {
				text = strings.Join(lines, "\n")
			}
			return state.Printer.Print(payload, text)
		},
	}

	cmd.Flags().StringVar(&query, "query", "", "search text")
	cmd.Flags().StringVar(&spaceID, "space-id", "", "space ID filter")
	cmd.Flags().StringVar(&nodeID, "node-id", "", "node ID filter")
	cmd.Flags().IntVar(&limit, "limit", 50, "max number of nodes to return")
	cmd.Flags().StringVar(&userAccessToken, "user-access-token", "", "user access token (OAuth)")
	return cmd
}
