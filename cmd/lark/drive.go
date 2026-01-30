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
	cmd.AddCommand(newDriveSearchCmd(state))
	cmd.AddCommand(newDriveGetCmd(state))
	cmd.AddCommand(newDriveURLsCmd(state))
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

func newDriveSearchCmd(state *appState) *cobra.Command {
	var query string
	var limit int

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search Drive files by text",
		RunE: func(cmd *cobra.Command, args []string) error {
			if query == "" {
				return errors.New("query is required")
			}
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
				result, err := state.Client.SearchDriveFiles(context.Background(), token, larkapi.SearchDriveFilesRequest{
					Query:     query,
					PageSize:  pageSize,
					PageToken: pageToken,
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

	cmd.Flags().StringVar(&query, "query", "", "search text")
	cmd.Flags().IntVar(&limit, "limit", 50, "max number of files to return")
	return cmd
}

func newDriveGetCmd(state *appState) *cobra.Command {
	var fileToken string

	cmd := &cobra.Command{
		Use:   "get <file-token>",
		Short: "Get Drive file metadata",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				if fileToken != "" && fileToken != args[0] {
					return errors.New("file-token provided twice")
				}
				fileToken = args[0]
			}
			if fileToken == "" {
				return errors.New("file-token is required")
			}
			token, err := ensureTenantToken(context.Background(), state)
			if err != nil {
				return err
			}
			file, err := state.Client.GetDriveFileMetadata(context.Background(), token, fileToken)
			if err != nil {
				return err
			}
			payload := map[string]any{"file": file}
			text := fmt.Sprintf("%s\t%s\t%s\t%s", file.Token, file.Name, file.FileType, file.URL)
			return state.Printer.Print(payload, text)
		},
	}

	cmd.Flags().StringVar(&fileToken, "file-token", "", "Drive file token")
	return cmd
}

func newDriveURLsCmd(state *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "urls <file-id> [file-id...]",
		Short: "Print web URLs for Drive file IDs",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := ensureTenantToken(context.Background(), state)
			if err != nil {
				return err
			}
			files := make([]larkapi.DriveFile, 0, len(args))
			for _, fileID := range args {
				file, err := state.Client.GetDriveFile(context.Background(), token, fileID)
				if err != nil {
					return err
				}
				files = append(files, file)
			}
			payload := map[string]any{"files": files}
			lines := make([]string, 0, len(files))
			for _, file := range files {
				lines = append(lines, fmt.Sprintf("%s\t%s\t%s", file.Token, file.URL, file.Name))
			}
			return state.Printer.Print(payload, strings.Join(lines, "\n"))
		},
	}

	return cmd
}
