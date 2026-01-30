package larksdk

import (
	"context"
	"errors"
	"fmt"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkwiki "github.com/larksuite/oapi-sdk-go/v3/service/wiki/v2"
)

func (c *Client) ListWikiSpaces(ctx context.Context, token string, req ListWikiSpacesRequest) (ListWikiSpacesResult, error) {
	if !c.available() {
		return ListWikiSpacesResult{}, ErrUnavailable
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return ListWikiSpacesResult{}, errors.New("tenant access token is required")
	}

	builder := larkwiki.NewListSpaceReqBuilder()
	if req.PageSize > 0 {
		builder.PageSize(req.PageSize)
	}
	if req.PageToken != "" {
		builder.PageToken(req.PageToken)
	}

	resp, err := c.sdk.Wiki.V2.Space.List(ctx, builder.Build(), larkcore.WithTenantAccessToken(tenantToken))
	if err != nil {
		return ListWikiSpacesResult{}, err
	}
	if resp == nil {
		return ListWikiSpacesResult{}, errors.New("list wiki spaces failed: empty response")
	}
	if !resp.Success() {
		return ListWikiSpacesResult{}, fmt.Errorf("list wiki spaces failed: %s", resp.Msg)
	}

	result := ListWikiSpacesResult{}
	if resp.Data != nil {
		if resp.Data.Items != nil {
			result.Items = make([]WikiSpace, 0, len(resp.Data.Items))
			for _, space := range resp.Data.Items {
				result.Items = append(result.Items, mapWikiSpace(space))
			}
		}
		if resp.Data.PageToken != nil {
			result.PageToken = *resp.Data.PageToken
		}
		if resp.Data.HasMore != nil {
			result.HasMore = *resp.Data.HasMore
		}
	}
	return result, nil
}

func (c *Client) GetWikiSpace(ctx context.Context, token string, req GetWikiSpaceRequest) (WikiSpace, error) {
	if !c.available() {
		return WikiSpace{}, ErrUnavailable
	}
	if req.SpaceID == "" {
		return WikiSpace{}, errors.New("space id is required")
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return WikiSpace{}, errors.New("tenant access token is required")
	}

	builder := larkwiki.NewGetSpaceReqBuilder().SpaceId(req.SpaceID)
	if req.Lang != "" {
		builder.Lang(req.Lang)
	}
	resp, err := c.sdk.Wiki.V2.Space.Get(ctx, builder.Build(), larkcore.WithTenantAccessToken(tenantToken))
	if err != nil {
		return WikiSpace{}, err
	}
	if resp == nil {
		return WikiSpace{}, errors.New("get wiki space failed: empty response")
	}
	if !resp.Success() {
		return WikiSpace{}, fmt.Errorf("get wiki space failed: %s", resp.Msg)
	}
	if resp.Data == nil || resp.Data.Space == nil {
		return WikiSpace{}, nil
	}
	return mapWikiSpace(resp.Data.Space), nil
}

func (c *Client) GetWikiNode(ctx context.Context, token string, req GetWikiNodeRequest) (WikiNode, error) {
	if !c.available() {
		return WikiNode{}, ErrUnavailable
	}
	if req.Token == "" {
		return WikiNode{}, errors.New("node token is required")
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return WikiNode{}, errors.New("tenant access token is required")
	}

	builder := larkwiki.NewGetNodeSpaceReqBuilder().Token(req.Token)
	if req.ObjType != "" {
		builder.ObjType(req.ObjType)
	}
	resp, err := c.sdk.Wiki.V2.Space.GetNode(ctx, builder.Build(), larkcore.WithTenantAccessToken(tenantToken))
	if err != nil {
		return WikiNode{}, err
	}
	if resp == nil {
		return WikiNode{}, errors.New("get wiki node failed: empty response")
	}
	if !resp.Success() {
		return WikiNode{}, fmt.Errorf("get wiki node failed: %s", resp.Msg)
	}
	if resp.Data == nil || resp.Data.Node == nil {
		return WikiNode{}, nil
	}
	return mapWikiNode(resp.Data.Node), nil
}

func (c *Client) ListWikiNodes(ctx context.Context, token string, req ListWikiNodesRequest) (ListWikiNodesResult, error) {
	if !c.available() {
		return ListWikiNodesResult{}, ErrUnavailable
	}
	if req.SpaceID == "" {
		return ListWikiNodesResult{}, errors.New("space id is required")
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return ListWikiNodesResult{}, errors.New("tenant access token is required")
	}

	builder := larkwiki.NewListSpaceNodeReqBuilder().SpaceId(req.SpaceID)
	if req.ParentNodeToken != "" {
		builder.ParentNodeToken(req.ParentNodeToken)
	}
	if req.PageSize > 0 {
		builder.PageSize(req.PageSize)
	}
	if req.PageToken != "" {
		builder.PageToken(req.PageToken)
	}

	resp, err := c.sdk.Wiki.V2.SpaceNode.List(ctx, builder.Build(), larkcore.WithTenantAccessToken(tenantToken))
	if err != nil {
		return ListWikiNodesResult{}, err
	}
	if resp == nil {
		return ListWikiNodesResult{}, errors.New("list wiki nodes failed: empty response")
	}
	if !resp.Success() {
		return ListWikiNodesResult{}, fmt.Errorf("list wiki nodes failed: %s", resp.Msg)
	}

	result := ListWikiNodesResult{}
	if resp.Data != nil {
		if resp.Data.Items != nil {
			result.Items = make([]WikiNode, 0, len(resp.Data.Items))
			for _, node := range resp.Data.Items {
				result.Items = append(result.Items, mapWikiNode(node))
			}
		}
		if resp.Data.PageToken != nil {
			result.PageToken = *resp.Data.PageToken
		}
		if resp.Data.HasMore != nil {
			result.HasMore = *resp.Data.HasMore
		}
	}
	return result, nil
}

func (c *Client) ListWikiSpaceMembers(ctx context.Context, token string, req ListWikiSpaceMembersRequest) (ListWikiSpaceMembersResult, error) {
	if !c.available() {
		return ListWikiSpaceMembersResult{}, ErrUnavailable
	}
	if req.SpaceID == "" {
		return ListWikiSpaceMembersResult{}, errors.New("space id is required")
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return ListWikiSpaceMembersResult{}, errors.New("tenant access token is required")
	}

	builder := larkwiki.NewListSpaceMemberReqBuilder().SpaceId(req.SpaceID)
	if req.PageSize > 0 {
		builder.PageSize(req.PageSize)
	}
	if req.PageToken != "" {
		builder.PageToken(req.PageToken)
	}

	resp, err := c.sdk.Wiki.V2.SpaceMember.List(ctx, builder.Build(), larkcore.WithTenantAccessToken(tenantToken))
	if err != nil {
		return ListWikiSpaceMembersResult{}, err
	}
	if resp == nil {
		return ListWikiSpaceMembersResult{}, errors.New("list wiki space members failed: empty response")
	}
	if !resp.Success() {
		return ListWikiSpaceMembersResult{}, fmt.Errorf("list wiki space members failed: %s", resp.Msg)
	}

	result := ListWikiSpaceMembersResult{}
	if resp.Data != nil {
		if resp.Data.Members != nil {
			result.Members = make([]WikiMember, 0, len(resp.Data.Members))
			for _, member := range resp.Data.Members {
				result.Members = append(result.Members, mapWikiMember(member))
			}
		}
		if resp.Data.PageToken != nil {
			result.PageToken = *resp.Data.PageToken
		}
		if resp.Data.HasMore != nil {
			result.HasMore = *resp.Data.HasMore
		}
	}
	return result, nil
}

func (c *Client) CreateWikiSpaceMember(ctx context.Context, token string, req CreateWikiSpaceMemberRequest) (WikiMember, error) {
	if !c.available() {
		return WikiMember{}, ErrUnavailable
	}
	if req.SpaceID == "" {
		return WikiMember{}, errors.New("space id is required")
	}
	if req.MemberType == "" || req.MemberID == "" {
		return WikiMember{}, errors.New("member type and member id are required")
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return WikiMember{}, errors.New("tenant access token is required")
	}

	memberBuilder := larkwiki.NewMemberBuilder().MemberType(req.MemberType).MemberId(req.MemberID)
	if req.MemberRole != "" {
		memberBuilder.MemberRole(req.MemberRole)
	}
	if req.Type != "" {
		memberBuilder.Type(req.Type)
	}

	builder := larkwiki.NewCreateSpaceMemberReqBuilder().
		SpaceId(req.SpaceID).
		NeedNotification(req.NeedNotification).
		Member(memberBuilder.Build())

	resp, err := c.sdk.Wiki.V2.SpaceMember.Create(ctx, builder.Build(), larkcore.WithTenantAccessToken(tenantToken))
	if err != nil {
		return WikiMember{}, err
	}
	if resp == nil {
		return WikiMember{}, errors.New("create wiki space member failed: empty response")
	}
	if !resp.Success() {
		return WikiMember{}, fmt.Errorf("create wiki space member failed: %s", resp.Msg)
	}
	if resp.Data == nil || resp.Data.Member == nil {
		return WikiMember{}, nil
	}
	return mapWikiMember(resp.Data.Member), nil
}

func (c *Client) DeleteWikiSpaceMember(ctx context.Context, token string, req DeleteWikiSpaceMemberRequest) (WikiMember, error) {
	if !c.available() {
		return WikiMember{}, ErrUnavailable
	}
	if req.SpaceID == "" {
		return WikiMember{}, errors.New("space id is required")
	}
	if req.MemberType == "" || req.MemberID == "" {
		return WikiMember{}, errors.New("member type and member id are required")
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return WikiMember{}, errors.New("tenant access token is required")
	}

	memberBuilder := larkwiki.NewMemberBuilder().MemberType(req.MemberType).MemberId(req.MemberID)
	if req.MemberRole != "" {
		memberBuilder.MemberRole(req.MemberRole)
	}
	if req.Type != "" {
		memberBuilder.Type(req.Type)
	}

	builder := larkwiki.NewDeleteSpaceMemberReqBuilder().
		SpaceId(req.SpaceID).
		MemberId(req.MemberID).
		Member(memberBuilder.Build())

	resp, err := c.sdk.Wiki.V2.SpaceMember.Delete(ctx, builder.Build(), larkcore.WithTenantAccessToken(tenantToken))
	if err != nil {
		return WikiMember{}, err
	}
	if resp == nil {
		return WikiMember{}, errors.New("delete wiki space member failed: empty response")
	}
	if !resp.Success() {
		return WikiMember{}, fmt.Errorf("delete wiki space member failed: %s", resp.Msg)
	}
	if resp.Data == nil || resp.Data.Member == nil {
		return WikiMember{}, nil
	}
	return mapWikiMember(resp.Data.Member), nil
}

func (c *Client) GetWikiTask(ctx context.Context, token string, req GetWikiTaskRequest) (WikiTask, error) {
	if !c.available() {
		return WikiTask{}, ErrUnavailable
	}
	if req.TaskID == "" {
		return WikiTask{}, errors.New("task id is required")
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return WikiTask{}, errors.New("tenant access token is required")
	}

	builder := larkwiki.NewGetTaskReqBuilder().TaskId(req.TaskID)
	if req.TaskType != "" {
		builder.TaskType(req.TaskType)
	}
	resp, err := c.sdk.Wiki.V2.Task.Get(ctx, builder.Build(), larkcore.WithTenantAccessToken(tenantToken))
	if err != nil {
		return WikiTask{}, err
	}
	if resp == nil {
		return WikiTask{}, errors.New("get wiki task failed: empty response")
	}
	if !resp.Success() {
		return WikiTask{}, fmt.Errorf("get wiki task failed: %s", resp.Msg)
	}
	if resp.Data == nil || resp.Data.Task == nil {
		return WikiTask{}, nil
	}
	return mapWikiTask(resp.Data.Task), nil
}

func mapWikiSpace(space *larkwiki.Space) WikiSpace {
	if space == nil {
		return WikiSpace{}
	}
	result := WikiSpace{}
	if space.SpaceId != nil {
		result.SpaceID = *space.SpaceId
	}
	if space.Name != nil {
		result.Name = *space.Name
	}
	if space.Description != nil {
		result.Description = *space.Description
	}
	if space.SpaceType != nil {
		result.SpaceType = *space.SpaceType
	}
	if space.Visibility != nil {
		result.Visibility = *space.Visibility
	}
	if space.OpenSharing != nil {
		result.OpenSharing = *space.OpenSharing
	}
	return result
}

func mapWikiNode(node *larkwiki.Node) WikiNode {
	if node == nil {
		return WikiNode{}
	}
	result := WikiNode{}
	if node.SpaceId != nil {
		result.SpaceID = *node.SpaceId
	}
	if node.NodeToken != nil {
		result.NodeToken = *node.NodeToken
	}
	if node.ObjToken != nil {
		result.ObjToken = *node.ObjToken
	}
	if node.ObjType != nil {
		result.ObjType = *node.ObjType
	}
	if node.ParentNodeToken != nil {
		result.ParentNodeToken = *node.ParentNodeToken
	}
	if node.NodeType != nil {
		result.NodeType = *node.NodeType
	}
	if node.OriginNodeToken != nil {
		result.OriginNodeToken = *node.OriginNodeToken
	}
	if node.OriginSpaceId != nil {
		result.OriginSpaceID = *node.OriginSpaceId
	}
	if node.HasChild != nil {
		result.HasChild = *node.HasChild
	}
	if node.Title != nil {
		result.Title = *node.Title
	}
	if node.ObjCreateTime != nil {
		result.ObjCreateTime = *node.ObjCreateTime
	}
	if node.ObjEditTime != nil {
		result.ObjEditTime = *node.ObjEditTime
	}
	if node.NodeCreateTime != nil {
		result.NodeCreateTime = *node.NodeCreateTime
	}
	if node.Creator != nil {
		result.Creator = *node.Creator
	}
	if node.Owner != nil {
		result.Owner = *node.Owner
	}
	if node.NodeCreator != nil {
		result.NodeCreator = *node.NodeCreator
	}
	return result
}

func mapWikiMember(member *larkwiki.Member) WikiMember {
	if member == nil {
		return WikiMember{}
	}
	result := WikiMember{}
	if member.MemberType != nil {
		result.MemberType = *member.MemberType
	}
	if member.MemberId != nil {
		result.MemberID = *member.MemberId
	}
	if member.MemberRole != nil {
		result.MemberRole = *member.MemberRole
	}
	if member.Type != nil {
		result.Type = *member.Type
	}
	return result
}

func mapWikiTask(task *larkwiki.TaskResult) WikiTask {
	if task == nil {
		return WikiTask{}
	}
	result := WikiTask{}
	if task.TaskId != nil {
		result.TaskID = *task.TaskId
	}
	if task.MoveResult != nil {
		result.MoveResults = make([]WikiMoveResult, 0, len(task.MoveResult))
		for _, move := range task.MoveResult {
			result.MoveResults = append(result.MoveResults, mapWikiMoveResult(move))
		}
	}
	return result
}

func mapWikiMoveResult(move *larkwiki.MoveResult) WikiMoveResult {
	if move == nil {
		return WikiMoveResult{}
	}
	result := WikiMoveResult{}
	if move.Node != nil {
		result.Node = mapWikiNode(move.Node)
	}
	if move.Status != nil {
		result.Status = *move.Status
	}
	if move.StatusMsg != nil {
		result.StatusMsg = *move.StatusMsg
	}
	return result
}
