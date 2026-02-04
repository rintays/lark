package larksdk

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

type TaskTime struct {
	Timestamp string `json:"timestamp,omitempty"`
	IsAllDay  *bool  `json:"is_all_day,omitempty"`
}

type TaskMember struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	Role string `json:"role,omitempty"`
	Name string `json:"name,omitempty"`
}

type TaskReminder struct {
	ID                 string `json:"id,omitempty"`
	RelativeFireMinute int    `json:"relative_fire_minute,omitempty"`
}

type TaskInTasklistInfo struct {
	TasklistGUID string `json:"tasklist_guid,omitempty"`
	SectionGUID  string `json:"section_guid,omitempty"`
}

type Task struct {
	GUID           string               `json:"guid,omitempty"`
	Summary        string               `json:"summary,omitempty"`
	Description    string               `json:"description,omitempty"`
	Due            *TaskTime            `json:"due,omitempty"`
	Start          *TaskTime            `json:"start,omitempty"`
	CompletedAt    string               `json:"completed_at,omitempty"`
	Creator        *TaskMember          `json:"creator,omitempty"`
	Members        []TaskMember         `json:"members,omitempty"`
	Reminders      []TaskReminder       `json:"reminders,omitempty"`
	Tasklists      []TaskInTasklistInfo `json:"tasklists,omitempty"`
	RepeatRule     string               `json:"repeat_rule,omitempty"`
	Mode           int                  `json:"mode,omitempty"`
	IsMilestone    *bool                `json:"is_milestone,omitempty"`
	Status         string               `json:"status,omitempty"`
	URL            string               `json:"url,omitempty"`
	TaskID         string               `json:"task_id,omitempty"`
	ParentTaskGUID string               `json:"parent_task_guid,omitempty"`
	Extra          string               `json:"extra,omitempty"`
	CreatedAt      string               `json:"created_at,omitempty"`
	UpdatedAt      string               `json:"updated_at,omitempty"`
}

type TaskList struct {
	GUID      string       `json:"guid,omitempty"`
	Name      string       `json:"name,omitempty"`
	Creator   *TaskMember  `json:"creator,omitempty"`
	Owner     *TaskMember  `json:"owner,omitempty"`
	Members   []TaskMember `json:"members,omitempty"`
	URL       string       `json:"url,omitempty"`
	CreatedAt string       `json:"created_at,omitempty"`
	UpdatedAt string       `json:"updated_at,omitempty"`
}

type CreateTaskRequest struct {
	Summary     string
	Description string
	Due         *TaskTime
	Start       *TaskTime
	Members     []TaskMember
	Tasklists   []TaskInTasklistInfo
	ClientToken string
	CompletedAt *string
	Extra       string
	RepeatRule  string
	Mode        *int
	IsMilestone *bool
	UserIDType  string
}

type UpdateTaskRequest struct {
	TaskGUID     string
	Task         map[string]any
	UpdateFields []string
	UserIDType   string
}

type GetTaskRequest struct {
	TaskGUID   string
	UserIDType string
}

type ListTasksRequest struct {
	PageSize   int
	PageToken  string
	Completed  *bool
	Type       string
	UserIDType string
}

type ListTasksResult struct {
	Items     []Task
	PageToken string
	HasMore   bool
}

type ListTasklistsRequest struct {
	PageSize   int
	PageToken  string
	UserIDType string
}

type ListTasklistsResult struct {
	Items     []TaskList
	PageToken string
	HasMore   bool
}

type CreateTasklistRequest struct {
	Name       string
	Members    []TaskMember
	UserIDType string
}

type UpdateTasklistRequest struct {
	TasklistGUID      string
	Tasklist          map[string]any
	UpdateFields      []string
	OriginOwnerToRole string
	UserIDType        string
}

type GetTasklistRequest struct {
	TasklistGUID string
	UserIDType   string
}

type createTaskResponse struct {
	*larkcore.ApiResp `json:"-"`
	larkcore.CodeError
	Data *createTaskResponseData `json:"data"`
}

type createTaskResponseData struct {
	Task Task `json:"task"`
}

func (r *createTaskResponse) Success() bool { return r.Code == 0 }

type updateTaskResponse struct {
	*larkcore.ApiResp `json:"-"`
	larkcore.CodeError
	Data *updateTaskResponseData `json:"data"`
}

type updateTaskResponseData struct {
	Task Task `json:"task"`
}

func (r *updateTaskResponse) Success() bool { return r.Code == 0 }

type getTaskResponse struct {
	*larkcore.ApiResp `json:"-"`
	larkcore.CodeError
	Data *getTaskResponseData `json:"data"`
}

type getTaskResponseData struct {
	Task Task `json:"task"`
}

func (r *getTaskResponse) Success() bool { return r.Code == 0 }

type listTasksResponse struct {
	*larkcore.ApiResp `json:"-"`
	larkcore.CodeError
	Data *listTasksResponseData `json:"data"`
}

type listTasksResponseData struct {
	Items     []Task  `json:"items"`
	PageToken *string `json:"page_token"`
	HasMore   *bool   `json:"has_more"`
}

func (r *listTasksResponse) Success() bool { return r.Code == 0 }

type listTasklistsResponse struct {
	*larkcore.ApiResp `json:"-"`
	larkcore.CodeError
	Data *listTasklistsResponseData `json:"data"`
}

type listTasklistsResponseData struct {
	Items     []TaskList `json:"items"`
	Tasklists []TaskList `json:"tasklists"`
	PageToken *string    `json:"page_token"`
	HasMore   *bool      `json:"has_more"`
}

func (r *listTasklistsResponse) Success() bool { return r.Code == 0 }

type deleteTaskResponse struct {
	*larkcore.ApiResp `json:"-"`
	larkcore.CodeError
}

func (r *deleteTaskResponse) Success() bool { return r.Code == 0 }

type createTasklistResponse struct {
	*larkcore.ApiResp `json:"-"`
	larkcore.CodeError
	Data *createTasklistResponseData `json:"data"`
}

type createTasklistResponseData struct {
	Tasklist TaskList `json:"tasklist"`
}

func (r *createTasklistResponse) Success() bool { return r.Code == 0 }

type updateTasklistResponse struct {
	*larkcore.ApiResp `json:"-"`
	larkcore.CodeError
	Data *updateTasklistResponseData `json:"data"`
}

type updateTasklistResponseData struct {
	Tasklist TaskList `json:"tasklist"`
}

func (r *updateTasklistResponse) Success() bool { return r.Code == 0 }

type getTasklistResponse struct {
	*larkcore.ApiResp `json:"-"`
	larkcore.CodeError
	Data *getTasklistResponseData `json:"data"`
}

type getTasklistResponseData struct {
	Tasklist TaskList `json:"tasklist"`
}

func (r *getTasklistResponse) Success() bool { return r.Code == 0 }

type deleteTasklistResponse struct {
	*larkcore.ApiResp `json:"-"`
	larkcore.CodeError
}

func (r *deleteTasklistResponse) Success() bool { return r.Code == 0 }

func (c *Client) CreateTask(ctx context.Context, token string, req CreateTaskRequest) (Task, error) {
	if !c.available() || c.coreConfig == nil {
		return Task{}, ErrUnavailable
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return Task{}, errors.New("tenant access token is required")
	}
	return c.createTask(ctx, req, larkcore.WithTenantAccessToken(tenantToken))
}

func (c *Client) CreateTaskWithUserToken(ctx context.Context, userAccessToken string, req CreateTaskRequest) (Task, error) {
	if !c.available() || c.coreConfig == nil {
		return Task{}, ErrUnavailable
	}
	userAccessToken = strings.TrimSpace(userAccessToken)
	if userAccessToken == "" {
		return Task{}, errors.New("user access token is required")
	}
	return c.createTask(ctx, req, larkcore.WithUserAccessToken(userAccessToken))
}

func (c *Client) createTask(ctx context.Context, req CreateTaskRequest, option larkcore.RequestOptionFunc) (Task, error) {
	if strings.TrimSpace(req.Summary) == "" {
		return Task{}, errors.New("summary is required")
	}

	payload := map[string]any{
		"summary": req.Summary,
	}
	if req.Description != "" {
		payload["description"] = req.Description
	}
	if req.Due != nil {
		payload["due"] = req.Due
	}
	if req.Start != nil {
		payload["start"] = req.Start
	}
	if len(req.Members) > 0 {
		payload["members"] = req.Members
	}
	if len(req.Tasklists) > 0 {
		payload["tasklists"] = req.Tasklists
	}
	if req.ClientToken != "" {
		payload["client_token"] = req.ClientToken
	}
	if req.CompletedAt != nil {
		payload["completed_at"] = *req.CompletedAt
	}
	if req.Extra != "" {
		payload["extra"] = req.Extra
	}
	if req.RepeatRule != "" {
		payload["repeat_rule"] = req.RepeatRule
	}
	if req.Mode != nil {
		payload["mode"] = *req.Mode
	}
	if req.IsMilestone != nil {
		payload["is_milestone"] = *req.IsMilestone
	}

	apiReq := &larkcore.ApiReq{
		ApiPath:                   "/open-apis/task/v2/tasks",
		HttpMethod:                http.MethodPost,
		PathParams:                larkcore.PathParams{},
		QueryParams:               larkcore.QueryParams{},
		Body:                      payload,
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser},
	}
	if strings.TrimSpace(req.UserIDType) != "" {
		apiReq.QueryParams.Set("user_id_type", req.UserIDType)
	}

	apiResp, err := larkcore.Request(ctx, apiReq, c.coreConfig, option)
	if err != nil {
		return Task{}, err
	}
	if apiResp == nil {
		return Task{}, errors.New("create task failed: empty response")
	}
	resp := &createTaskResponse{ApiResp: apiResp}
	if err := apiResp.JSONUnmarshalBody(resp, c.coreConfig); err != nil {
		return Task{}, err
	}
	if !resp.Success() {
		return Task{}, formatCodeError("create task failed", resp.CodeError, resp.ApiResp)
	}
	if resp.Data == nil {
		return Task{}, nil
	}
	return resp.Data.Task, nil
}

func (c *Client) UpdateTask(ctx context.Context, token string, req UpdateTaskRequest) (Task, error) {
	if !c.available() || c.coreConfig == nil {
		return Task{}, ErrUnavailable
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return Task{}, errors.New("tenant access token is required")
	}
	return c.updateTask(ctx, req, larkcore.WithTenantAccessToken(tenantToken))
}

func (c *Client) UpdateTaskWithUserToken(ctx context.Context, userAccessToken string, req UpdateTaskRequest) (Task, error) {
	if !c.available() || c.coreConfig == nil {
		return Task{}, ErrUnavailable
	}
	userAccessToken = strings.TrimSpace(userAccessToken)
	if userAccessToken == "" {
		return Task{}, errors.New("user access token is required")
	}
	return c.updateTask(ctx, req, larkcore.WithUserAccessToken(userAccessToken))
}

func (c *Client) updateTask(ctx context.Context, req UpdateTaskRequest, option larkcore.RequestOptionFunc) (Task, error) {
	taskGUID := strings.TrimSpace(req.TaskGUID)
	if taskGUID == "" {
		return Task{}, errors.New("task guid is required")
	}
	if len(req.UpdateFields) == 0 {
		return Task{}, errors.New("update_fields is required")
	}
	taskPayload := req.Task
	if taskPayload == nil {
		taskPayload = map[string]any{}
	}

	payload := map[string]any{
		"task":          taskPayload,
		"update_fields": req.UpdateFields,
	}

	apiReq := &larkcore.ApiReq{
		ApiPath:                   "/open-apis/task/v2/tasks/:task_guid",
		HttpMethod:                http.MethodPatch,
		PathParams:                larkcore.PathParams{},
		QueryParams:               larkcore.QueryParams{},
		Body:                      payload,
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser},
	}
	apiReq.PathParams.Set("task_guid", taskGUID)
	if strings.TrimSpace(req.UserIDType) != "" {
		apiReq.QueryParams.Set("user_id_type", req.UserIDType)
	}

	apiResp, err := larkcore.Request(ctx, apiReq, c.coreConfig, option)
	if err != nil {
		return Task{}, err
	}
	if apiResp == nil {
		return Task{}, errors.New("update task failed: empty response")
	}
	resp := &updateTaskResponse{ApiResp: apiResp}
	if err := apiResp.JSONUnmarshalBody(resp, c.coreConfig); err != nil {
		return Task{}, err
	}
	if !resp.Success() {
		return Task{}, formatCodeError("update task failed", resp.CodeError, resp.ApiResp)
	}
	if resp.Data == nil {
		return Task{}, nil
	}
	return resp.Data.Task, nil
}

func (c *Client) GetTask(ctx context.Context, token string, req GetTaskRequest) (Task, error) {
	if !c.available() || c.coreConfig == nil {
		return Task{}, ErrUnavailable
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return Task{}, errors.New("tenant access token is required")
	}
	return c.getTask(ctx, req, larkcore.WithTenantAccessToken(tenantToken))
}

func (c *Client) GetTaskWithUserToken(ctx context.Context, userAccessToken string, req GetTaskRequest) (Task, error) {
	if !c.available() || c.coreConfig == nil {
		return Task{}, ErrUnavailable
	}
	userAccessToken = strings.TrimSpace(userAccessToken)
	if userAccessToken == "" {
		return Task{}, errors.New("user access token is required")
	}
	return c.getTask(ctx, req, larkcore.WithUserAccessToken(userAccessToken))
}

func (c *Client) getTask(ctx context.Context, req GetTaskRequest, option larkcore.RequestOptionFunc) (Task, error) {
	taskGUID := strings.TrimSpace(req.TaskGUID)
	if taskGUID == "" {
		return Task{}, errors.New("task guid is required")
	}

	apiReq := &larkcore.ApiReq{
		ApiPath:                   "/open-apis/task/v2/tasks/:task_guid",
		HttpMethod:                http.MethodGet,
		PathParams:                larkcore.PathParams{},
		QueryParams:               larkcore.QueryParams{},
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser},
	}
	apiReq.PathParams.Set("task_guid", taskGUID)
	if strings.TrimSpace(req.UserIDType) != "" {
		apiReq.QueryParams.Set("user_id_type", req.UserIDType)
	}

	apiResp, err := larkcore.Request(ctx, apiReq, c.coreConfig, option)
	if err != nil {
		return Task{}, err
	}
	if apiResp == nil {
		return Task{}, errors.New("get task failed: empty response")
	}
	resp := &getTaskResponse{ApiResp: apiResp}
	if err := apiResp.JSONUnmarshalBody(resp, c.coreConfig); err != nil {
		return Task{}, err
	}
	if !resp.Success() {
		return Task{}, formatCodeError("get task failed", resp.CodeError, resp.ApiResp)
	}
	if resp.Data == nil {
		return Task{}, nil
	}
	return resp.Data.Task, nil
}

func (c *Client) DeleteTask(ctx context.Context, token string, taskGUID string) error {
	if !c.available() || c.coreConfig == nil {
		return ErrUnavailable
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return errors.New("tenant access token is required")
	}
	return c.deleteTask(ctx, taskGUID, larkcore.WithTenantAccessToken(tenantToken))
}

func (c *Client) DeleteTaskWithUserToken(ctx context.Context, userAccessToken string, taskGUID string) error {
	if !c.available() || c.coreConfig == nil {
		return ErrUnavailable
	}
	userAccessToken = strings.TrimSpace(userAccessToken)
	if userAccessToken == "" {
		return errors.New("user access token is required")
	}
	return c.deleteTask(ctx, taskGUID, larkcore.WithUserAccessToken(userAccessToken))
}

func (c *Client) deleteTask(ctx context.Context, taskGUID string, option larkcore.RequestOptionFunc) error {
	taskGUID = strings.TrimSpace(taskGUID)
	if taskGUID == "" {
		return errors.New("task guid is required")
	}
	apiReq := &larkcore.ApiReq{
		ApiPath:                   "/open-apis/task/v2/tasks/:task_guid",
		HttpMethod:                http.MethodDelete,
		PathParams:                larkcore.PathParams{},
		QueryParams:               larkcore.QueryParams{},
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser},
	}
	apiReq.PathParams.Set("task_guid", taskGUID)

	apiResp, err := larkcore.Request(ctx, apiReq, c.coreConfig, option)
	if err != nil {
		return err
	}
	if apiResp == nil {
		return errors.New("delete task failed: empty response")
	}
	resp := &deleteTaskResponse{ApiResp: apiResp}
	if err := apiResp.JSONUnmarshalBody(resp, c.coreConfig); err != nil {
		return err
	}
	if !resp.Success() {
		return formatCodeError("delete task failed", resp.CodeError, resp.ApiResp)
	}
	return nil
}

func (c *Client) ListTasks(ctx context.Context, userAccessToken string, req ListTasksRequest) (ListTasksResult, error) {
	if !c.available() || c.coreConfig == nil {
		return ListTasksResult{}, ErrUnavailable
	}
	userAccessToken = strings.TrimSpace(userAccessToken)
	if userAccessToken == "" {
		return ListTasksResult{}, errors.New("user access token is required")
	}
	if req.PageSize < 0 {
		return ListTasksResult{}, errors.New("page_size must be greater than or equal to 0")
	}
	if req.PageSize > 100 {
		return ListTasksResult{}, errors.New("page_size must be less than or equal to 100")
	}

	apiReq := &larkcore.ApiReq{
		ApiPath:                   "/open-apis/task/v2/tasks",
		HttpMethod:                http.MethodGet,
		PathParams:                larkcore.PathParams{},
		QueryParams:               larkcore.QueryParams{},
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeUser},
	}
	if req.PageSize > 0 {
		apiReq.QueryParams.Set("page_size", fmt.Sprint(req.PageSize))
	}
	if strings.TrimSpace(req.PageToken) != "" {
		apiReq.QueryParams.Set("page_token", req.PageToken)
	}
	if req.Completed != nil {
		apiReq.QueryParams.Set("completed", fmt.Sprintf("%t", *req.Completed))
	}
	if strings.TrimSpace(req.Type) != "" {
		apiReq.QueryParams.Set("type", req.Type)
	}
	if strings.TrimSpace(req.UserIDType) != "" {
		apiReq.QueryParams.Set("user_id_type", req.UserIDType)
	}

	apiResp, err := larkcore.Request(ctx, apiReq, c.coreConfig, larkcore.WithUserAccessToken(userAccessToken))
	if err != nil {
		return ListTasksResult{}, err
	}
	if apiResp == nil {
		return ListTasksResult{}, errors.New("list tasks failed: empty response")
	}
	resp := &listTasksResponse{ApiResp: apiResp}
	if err := apiResp.JSONUnmarshalBody(resp, c.coreConfig); err != nil {
		return ListTasksResult{}, err
	}
	if !resp.Success() {
		return ListTasksResult{}, formatCodeError("list tasks failed", resp.CodeError, resp.ApiResp)
	}

	result := ListTasksResult{}
	if resp.Data != nil {
		if resp.Data.Items != nil {
			result.Items = resp.Data.Items
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

func (c *Client) ListTasklists(ctx context.Context, userAccessToken string, req ListTasklistsRequest) (ListTasklistsResult, error) {
	if !c.available() || c.coreConfig == nil {
		return ListTasklistsResult{}, ErrUnavailable
	}
	userAccessToken = strings.TrimSpace(userAccessToken)
	if userAccessToken == "" {
		return ListTasklistsResult{}, errors.New("user access token is required")
	}
	if req.PageSize < 0 {
		return ListTasklistsResult{}, errors.New("page_size must be greater than or equal to 0")
	}
	if req.PageSize > 500 {
		return ListTasklistsResult{}, errors.New("page_size must be less than or equal to 500")
	}

	apiReq := &larkcore.ApiReq{
		ApiPath:                   "/open-apis/task/v2/tasklists",
		HttpMethod:                http.MethodGet,
		PathParams:                larkcore.PathParams{},
		QueryParams:               larkcore.QueryParams{},
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeUser},
	}
	if req.PageSize > 0 {
		apiReq.QueryParams.Set("page_size", fmt.Sprint(req.PageSize))
	}
	if strings.TrimSpace(req.PageToken) != "" {
		apiReq.QueryParams.Set("page_token", req.PageToken)
	}
	if strings.TrimSpace(req.UserIDType) != "" {
		apiReq.QueryParams.Set("user_id_type", req.UserIDType)
	}

	apiResp, err := larkcore.Request(ctx, apiReq, c.coreConfig, larkcore.WithUserAccessToken(userAccessToken))
	if err != nil {
		return ListTasklistsResult{}, err
	}
	if apiResp == nil {
		return ListTasklistsResult{}, errors.New("list tasklists failed: empty response")
	}
	resp := &listTasklistsResponse{ApiResp: apiResp}
	if err := apiResp.JSONUnmarshalBody(resp, c.coreConfig); err != nil {
		return ListTasklistsResult{}, err
	}
	if !resp.Success() {
		return ListTasklistsResult{}, formatCodeError("list tasklists failed", resp.CodeError, resp.ApiResp)
	}

	result := ListTasklistsResult{}
	if resp.Data != nil {
		items := resp.Data.Items
		if len(items) == 0 && len(resp.Data.Tasklists) > 0 {
			items = resp.Data.Tasklists
		}
		result.Items = items
		if resp.Data.PageToken != nil {
			result.PageToken = *resp.Data.PageToken
		}
		if resp.Data.HasMore != nil {
			result.HasMore = *resp.Data.HasMore
		}
	}
	return result, nil
}

func (c *Client) CreateTasklist(ctx context.Context, token string, req CreateTasklistRequest) (TaskList, error) {
	if !c.available() || c.coreConfig == nil {
		return TaskList{}, ErrUnavailable
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return TaskList{}, errors.New("tenant access token is required")
	}
	return c.createTasklist(ctx, req, larkcore.WithTenantAccessToken(tenantToken))
}

func (c *Client) CreateTasklistWithUserToken(ctx context.Context, userAccessToken string, req CreateTasklistRequest) (TaskList, error) {
	if !c.available() || c.coreConfig == nil {
		return TaskList{}, ErrUnavailable
	}
	userAccessToken = strings.TrimSpace(userAccessToken)
	if userAccessToken == "" {
		return TaskList{}, errors.New("user access token is required")
	}
	return c.createTasklist(ctx, req, larkcore.WithUserAccessToken(userAccessToken))
}

func (c *Client) createTasklist(ctx context.Context, req CreateTasklistRequest, option larkcore.RequestOptionFunc) (TaskList, error) {
	if strings.TrimSpace(req.Name) == "" {
		return TaskList{}, errors.New("name is required")
	}
	payload := map[string]any{
		"name": req.Name,
	}
	if len(req.Members) > 0 {
		payload["members"] = req.Members
	}

	apiReq := &larkcore.ApiReq{
		ApiPath:                   "/open-apis/task/v2/tasklists",
		HttpMethod:                http.MethodPost,
		PathParams:                larkcore.PathParams{},
		QueryParams:               larkcore.QueryParams{},
		Body:                      payload,
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser},
	}
	if strings.TrimSpace(req.UserIDType) != "" {
		apiReq.QueryParams.Set("user_id_type", req.UserIDType)
	}

	apiResp, err := larkcore.Request(ctx, apiReq, c.coreConfig, option)
	if err != nil {
		return TaskList{}, err
	}
	if apiResp == nil {
		return TaskList{}, errors.New("create tasklist failed: empty response")
	}
	resp := &createTasklistResponse{ApiResp: apiResp}
	if err := apiResp.JSONUnmarshalBody(resp, c.coreConfig); err != nil {
		return TaskList{}, err
	}
	if !resp.Success() {
		return TaskList{}, formatCodeError("create tasklist failed", resp.CodeError, resp.ApiResp)
	}
	if resp.Data == nil {
		return TaskList{}, nil
	}
	return resp.Data.Tasklist, nil
}

func (c *Client) UpdateTasklist(ctx context.Context, token string, req UpdateTasklistRequest) (TaskList, error) {
	if !c.available() || c.coreConfig == nil {
		return TaskList{}, ErrUnavailable
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return TaskList{}, errors.New("tenant access token is required")
	}
	return c.updateTasklist(ctx, req, larkcore.WithTenantAccessToken(tenantToken))
}

func (c *Client) UpdateTasklistWithUserToken(ctx context.Context, userAccessToken string, req UpdateTasklistRequest) (TaskList, error) {
	if !c.available() || c.coreConfig == nil {
		return TaskList{}, ErrUnavailable
	}
	userAccessToken = strings.TrimSpace(userAccessToken)
	if userAccessToken == "" {
		return TaskList{}, errors.New("user access token is required")
	}
	return c.updateTasklist(ctx, req, larkcore.WithUserAccessToken(userAccessToken))
}

func (c *Client) updateTasklist(ctx context.Context, req UpdateTasklistRequest, option larkcore.RequestOptionFunc) (TaskList, error) {
	tasklistGUID := strings.TrimSpace(req.TasklistGUID)
	if tasklistGUID == "" {
		return TaskList{}, errors.New("tasklist guid is required")
	}
	if len(req.UpdateFields) == 0 {
		return TaskList{}, errors.New("update_fields is required")
	}
	tasklistPayload := req.Tasklist
	if tasklistPayload == nil {
		tasklistPayload = map[string]any{}
	}
	payload := map[string]any{
		"tasklist":      tasklistPayload,
		"update_fields": req.UpdateFields,
	}
	if strings.TrimSpace(req.OriginOwnerToRole) != "" {
		payload["origin_owner_to_role"] = req.OriginOwnerToRole
	}

	apiReq := &larkcore.ApiReq{
		ApiPath:                   "/open-apis/task/v2/tasklists/:tasklist_guid",
		HttpMethod:                http.MethodPatch,
		PathParams:                larkcore.PathParams{},
		QueryParams:               larkcore.QueryParams{},
		Body:                      payload,
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser},
	}
	apiReq.PathParams.Set("tasklist_guid", tasklistGUID)
	if strings.TrimSpace(req.UserIDType) != "" {
		apiReq.QueryParams.Set("user_id_type", req.UserIDType)
	}

	apiResp, err := larkcore.Request(ctx, apiReq, c.coreConfig, option)
	if err != nil {
		return TaskList{}, err
	}
	if apiResp == nil {
		return TaskList{}, errors.New("update tasklist failed: empty response")
	}
	resp := &updateTasklistResponse{ApiResp: apiResp}
	if err := apiResp.JSONUnmarshalBody(resp, c.coreConfig); err != nil {
		return TaskList{}, err
	}
	if !resp.Success() {
		return TaskList{}, formatCodeError("update tasklist failed", resp.CodeError, resp.ApiResp)
	}
	if resp.Data == nil {
		return TaskList{}, nil
	}
	return resp.Data.Tasklist, nil
}

func (c *Client) GetTasklist(ctx context.Context, token string, req GetTasklistRequest) (TaskList, error) {
	if !c.available() || c.coreConfig == nil {
		return TaskList{}, ErrUnavailable
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return TaskList{}, errors.New("tenant access token is required")
	}
	return c.getTasklist(ctx, req, larkcore.WithTenantAccessToken(tenantToken))
}

func (c *Client) GetTasklistWithUserToken(ctx context.Context, userAccessToken string, req GetTasklistRequest) (TaskList, error) {
	if !c.available() || c.coreConfig == nil {
		return TaskList{}, ErrUnavailable
	}
	userAccessToken = strings.TrimSpace(userAccessToken)
	if userAccessToken == "" {
		return TaskList{}, errors.New("user access token is required")
	}
	return c.getTasklist(ctx, req, larkcore.WithUserAccessToken(userAccessToken))
}

func (c *Client) getTasklist(ctx context.Context, req GetTasklistRequest, option larkcore.RequestOptionFunc) (TaskList, error) {
	tasklistGUID := strings.TrimSpace(req.TasklistGUID)
	if tasklistGUID == "" {
		return TaskList{}, errors.New("tasklist guid is required")
	}
	apiReq := &larkcore.ApiReq{
		ApiPath:                   "/open-apis/task/v2/tasklists/:tasklist_guid",
		HttpMethod:                http.MethodGet,
		PathParams:                larkcore.PathParams{},
		QueryParams:               larkcore.QueryParams{},
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser},
	}
	apiReq.PathParams.Set("tasklist_guid", tasklistGUID)
	if strings.TrimSpace(req.UserIDType) != "" {
		apiReq.QueryParams.Set("user_id_type", req.UserIDType)
	}

	apiResp, err := larkcore.Request(ctx, apiReq, c.coreConfig, option)
	if err != nil {
		return TaskList{}, err
	}
	if apiResp == nil {
		return TaskList{}, errors.New("get tasklist failed: empty response")
	}
	resp := &getTasklistResponse{ApiResp: apiResp}
	if err := apiResp.JSONUnmarshalBody(resp, c.coreConfig); err != nil {
		return TaskList{}, err
	}
	if !resp.Success() {
		return TaskList{}, formatCodeError("get tasklist failed", resp.CodeError, resp.ApiResp)
	}
	if resp.Data == nil {
		return TaskList{}, nil
	}
	return resp.Data.Tasklist, nil
}

func (c *Client) DeleteTasklist(ctx context.Context, token string, tasklistGUID string) error {
	if !c.available() || c.coreConfig == nil {
		return ErrUnavailable
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return errors.New("tenant access token is required")
	}
	return c.deleteTasklist(ctx, tasklistGUID, larkcore.WithTenantAccessToken(tenantToken))
}

func (c *Client) DeleteTasklistWithUserToken(ctx context.Context, userAccessToken string, tasklistGUID string) error {
	if !c.available() || c.coreConfig == nil {
		return ErrUnavailable
	}
	userAccessToken = strings.TrimSpace(userAccessToken)
	if userAccessToken == "" {
		return errors.New("user access token is required")
	}
	return c.deleteTasklist(ctx, tasklistGUID, larkcore.WithUserAccessToken(userAccessToken))
}

func (c *Client) deleteTasklist(ctx context.Context, tasklistGUID string, option larkcore.RequestOptionFunc) error {
	tasklistGUID = strings.TrimSpace(tasklistGUID)
	if tasklistGUID == "" {
		return errors.New("tasklist guid is required")
	}
	apiReq := &larkcore.ApiReq{
		ApiPath:                   "/open-apis/task/v2/tasklists/:tasklist_guid",
		HttpMethod:                http.MethodDelete,
		PathParams:                larkcore.PathParams{},
		QueryParams:               larkcore.QueryParams{},
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser},
	}
	apiReq.PathParams.Set("tasklist_guid", tasklistGUID)

	apiResp, err := larkcore.Request(ctx, apiReq, c.coreConfig, option)
	if err != nil {
		return err
	}
	if apiResp == nil {
		return errors.New("delete tasklist failed: empty response")
	}
	resp := &deleteTasklistResponse{ApiResp: apiResp}
	if err := apiResp.JSONUnmarshalBody(resp, c.coreConfig); err != nil {
		return err
	}
	if !resp.Success() {
		return formatCodeError("delete tasklist failed", resp.CodeError, resp.ApiResp)
	}
	return nil
}
