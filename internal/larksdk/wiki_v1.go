package larksdk

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkwikiv1 "github.com/larksuite/oapi-sdk-go/v3/service/wiki/v1"
)

type searchWikiNodesResponse struct {
	*larkcore.ApiResp `json:"-"`
	larkcore.CodeError
	Data *searchWikiNodesResponseData `json:"data"`
}

type searchWikiNodesResponseData struct {
	Items     []*larkwikiv1.Node `json:"items"`
	PageToken *string            `json:"page_token"`
	HasMore   *bool              `json:"has_more"`
}

func (r *searchWikiNodesResponse) Success() bool {
	return r.Code == 0
}

func (c *Client) SearchWikiNodes(ctx context.Context, token string, req SearchWikiNodesRequest) (SearchWikiNodesResult, error) {
	if !c.available() || c.coreConfig == nil {
		return SearchWikiNodesResult{}, ErrUnavailable
	}
	if token == "" {
		return SearchWikiNodesResult{}, errors.New("user access token is required")
	}
	if req.Query == "" {
		return SearchWikiNodesResult{}, errors.New("query is required")
	}

	payload := map[string]any{
		"query": req.Query,
	}
	if req.SpaceID != "" {
		payload["space_id"] = req.SpaceID
	}
	if req.NodeID != "" {
		payload["node_id"] = req.NodeID
	}

	apiReq := &larkcore.ApiReq{
		ApiPath:                   "/open-apis/wiki/v1/nodes/search",
		HttpMethod:                http.MethodPost,
		PathParams:                larkcore.PathParams{},
		QueryParams:               larkcore.QueryParams{},
		Body:                      payload,
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeUser},
	}
	if req.PageSize > 0 {
		apiReq.QueryParams.Set("page_size", fmt.Sprint(req.PageSize))
	}
	if req.PageToken != "" {
		apiReq.QueryParams.Set("page_token", req.PageToken)
	}

	apiResp, err := larkcore.Request(ctx, apiReq, c.coreConfig, larkcore.WithUserAccessToken(token))
	if err != nil {
		return SearchWikiNodesResult{}, err
	}
	if apiResp == nil {
		return SearchWikiNodesResult{}, errors.New("search wiki nodes failed: empty response")
	}
	resp := &searchWikiNodesResponse{ApiResp: apiResp}
	if err := apiResp.JSONUnmarshalBody(resp, c.coreConfig); err != nil {
		return SearchWikiNodesResult{}, err
	}
	if !resp.Success() {
		return SearchWikiNodesResult{}, fmt.Errorf("search wiki nodes failed: %s", resp.Msg)
	}

	result := SearchWikiNodesResult{}
	if resp.Data != nil {
		if resp.Data.Items != nil {
			result.Items = make([]WikiV1Node, 0, len(resp.Data.Items))
			for _, node := range resp.Data.Items {
				result.Items = append(result.Items, mapWikiV1Node(node))
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

func mapWikiV1Node(node *larkwikiv1.Node) WikiV1Node {
	if node == nil {
		return WikiV1Node{}
	}
	result := WikiV1Node{}
	if node.NodeId != nil {
		result.NodeID = *node.NodeId
	}
	if node.SpaceId != nil {
		result.SpaceID = *node.SpaceId
	}
	if node.ParentId != nil {
		result.ParentID = *node.ParentId
	}
	if node.ObjType != nil {
		result.ObjType = *node.ObjType
	}
	if node.Title != nil {
		result.Title = *node.Title
	}
	if node.Url != nil {
		result.URL = *node.Url
	}
	if node.Icon != nil {
		result.Icon = *node.Icon
	}
	if node.AreaId != nil {
		result.AreaID = *node.AreaId
	}
	if node.SortId != nil {
		result.SortID = *node.SortId
	}
	if node.Domain != nil {
		result.Domain = *node.Domain
	}
	if node.ObjToken != nil {
		result.ObjToken = *node.ObjToken
	}
	if node.CreateTime != nil {
		result.CreateTime = *node.CreateTime
	}
	if node.UpdateTime != nil {
		result.UpdateTime = *node.UpdateTime
	}
	if node.DeleteTime != nil {
		result.DeleteTime = *node.DeleteTime
	}
	if node.ChildNum != nil {
		result.ChildNum = *node.ChildNum
	}
	if node.Version != nil {
		result.Version = *node.Version
	}
	return result
}
