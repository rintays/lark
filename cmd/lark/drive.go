package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"lark/internal/larkapi"
)

const maxDrivePageSize = 200

func newDriveCmd(state *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drive",
		Short: "Manage Drive files",
	}
	cmd.AddCommand(newDriveListCmd(state))
	return cmd
}

func newDriveListCmd(state *appState) *cobra.Command {
	var folderID string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List files in a Drive folder",
		RunE: func(cmd *cobra.Command, args []string) error {
			if limit <= 0 {
				return errors.New("limit must be greater than 0")
			}
			token, err := ensureTenantToken(context.Background(), state)
			if err != nil {
				return err
			}
			files := make([]larkapi.DriveFile, 0, limit)
			pageToken := ""
			remaining := limit
			for {
				pageSize := remaining
				if pageSize > maxDrivePageSize {
					pageSize = maxDrivePageSize
				}
				result, err := state.Client.ListDriveFiles(context.Background(), token, larkapi.ListDriveFilesRequest{
					FolderToken: folderID,
					PageSize:    pageSize,
					PageToken:   pageToken,
				})
				if err != nil {
					return err
				}
				files = append(files, result.Files...)
				if len(files) >= limit || !result.HasMore {
					break
				}
				remaining = limit - len(files)
				pageToken = result.PageToken
				if pageToken == "" {
					break
				}
			}
			if len(files) > limit {
				files = files[:limit]
			}
			payload := map[string]any{"files": files}
			lines := make([]string, 0, len(files))
			for _, file := range files {
				lines = append(lines, fmt.Sprintf("%s\t%s\t%s\t%s", file.Token, file.Name, file.FileType, file.URL))
			}
			text := "no files found"
			if len(lines) > 0 {
				text = strings.Join(lines, "\n")
			}
			return state.Printer.Print(payload, text)
		},
	}

	cmd.Flags().StringVar(&folderID, "folder-id", "", "Drive folder token (default: root)")
	cmd.Flags().IntVar(&limit, "limit", 50, "max number of files to return")
	return cmd
}
