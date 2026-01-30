package larksdk

import (
	"context"
	"errors"
	"fmt"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"

	"lark/internal/larkapi"
)

func (c *Client) ListDriveFiles(ctx context.Context, token string, req larkapi.ListDriveFilesRequest) (larkapi.ListDriveFilesResult, error) {
	if !c.available() {
		return larkapi.ListDriveFilesResult{}, ErrUnavailable
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return larkapi.ListDriveFilesResult{}, errors.New("tenant access token is required")
	}

	builder := larkdrive.NewListFileReqBuilder()
	if req.FolderToken != "" {
		builder.FolderToken(req.FolderToken)
	}
	if req.PageSize > 0 {
		builder.PageSize(req.PageSize)
	}
	if req.PageToken != "" {
		builder.PageToken(req.PageToken)
	}

	resp, err := c.sdk.Drive.V1.File.List(ctx, builder.Build(), larkcore.WithTenantAccessToken(tenantToken))
	if err != nil {
		return larkapi.ListDriveFilesResult{}, err
	}
	if resp == nil {
		return larkapi.ListDriveFilesResult{}, errors.New("list drive files failed: empty response")
	}
	if !resp.Success() {
		return larkapi.ListDriveFilesResult{}, fmt.Errorf("list drive files failed: %s", resp.Msg)
	}

	result := larkapi.ListDriveFilesResult{}
	if resp.Data != nil {
		if resp.Data.Files != nil {
			result.Files = make([]larkapi.DriveFile, 0, len(resp.Data.Files))
			for _, file := range resp.Data.Files {
				result.Files = append(result.Files, mapDriveFile(file))
			}
		}
		if resp.Data.NextPageToken != nil {
			result.PageToken = *resp.Data.NextPageToken
		}
		if resp.Data.HasMore != nil {
			result.HasMore = *resp.Data.HasMore
		}
	}
	return result, nil
}

func mapDriveFile(file *larkdrive.File) larkapi.DriveFile {
	if file == nil {
		return larkapi.DriveFile{}
	}
	result := larkapi.DriveFile{}
	if file.Token != nil {
		result.Token = *file.Token
	}
	if file.Name != nil {
		result.Name = *file.Name
	}
	if file.Type != nil {
		result.FileType = *file.Type
	}
	if file.Url != nil {
		result.URL = *file.Url
	}
	if file.ParentToken != nil {
		result.ParentID = *file.ParentToken
	}
	if file.OwnerId != nil {
		result.OwnerID = *file.OwnerId
	}
	return result
}
