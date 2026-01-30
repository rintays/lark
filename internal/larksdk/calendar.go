package larksdk

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

type primaryCalendarResponse struct {
	*larkcore.ApiResp `json:"-"`
	larkcore.CodeError
	Data *primaryCalendarResponseData `json:"data"`
}

type primaryCalendarResponseData struct {
	CalendarID string   `json:"calendar_id"`
	Calendar   Calendar `json:"calendar"`
}

func (r *primaryCalendarResponse) Success() bool {
	return r.Code == 0
}

type listCalendarEventsResponse struct {
	*larkcore.ApiResp `json:"-"`
	larkcore.CodeError
	Data *listCalendarEventsResponseData `json:"data"`
}

type listCalendarEventsResponseData struct {
	Items     []CalendarEvent `json:"items"`
	Events    []CalendarEvent `json:"events"`
	PageToken string          `json:"page_token"`
	HasMore   bool            `json:"has_more"`
	SyncToken string          `json:"sync_token"`
}

func (r *listCalendarEventsResponse) Success() bool {
	return r.Code == 0
}

type createCalendarEventResponse struct {
	*larkcore.ApiResp `json:"-"`
	larkcore.CodeError
	Data *createCalendarEventResponseData `json:"data"`
}

type createCalendarEventResponseData struct {
	Event   CalendarEvent `json:"event"`
	EventID string        `json:"event_id"`
}

func (r *createCalendarEventResponse) Success() bool {
	return r.Code == 0
}

type createCalendarEventAttendeesResponse struct {
	*larkcore.ApiResp `json:"-"`
	larkcore.CodeError
}

func (r *createCalendarEventAttendeesResponse) Success() bool {
	return r.Code == 0
}

func (c *Client) PrimaryCalendar(ctx context.Context, token string) (Calendar, error) {
	if !c.available() || c.coreConfig == nil {
		return Calendar{}, ErrUnavailable
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return Calendar{}, errors.New("tenant access token is required")
	}

	apiReq := &larkcore.ApiReq{
		ApiPath:                   "/open-apis/calendar/v4/calendars/primary",
		HttpMethod:                http.MethodPost,
		PathParams:                larkcore.PathParams{},
		QueryParams:               larkcore.QueryParams{},
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser},
	}

	apiResp, err := larkcore.Request(ctx, apiReq, c.coreConfig, larkcore.WithTenantAccessToken(tenantToken))
	if err != nil {
		return Calendar{}, err
	}
	if apiResp == nil {
		return Calendar{}, errors.New("primary calendar failed: empty response")
	}
	resp := &primaryCalendarResponse{ApiResp: apiResp}
	if err := apiResp.JSONUnmarshalBody(resp, c.coreConfig); err != nil {
		return Calendar{}, err
	}
	if !resp.Success() {
		return Calendar{}, fmt.Errorf("primary calendar failed: %s", resp.Msg)
	}
	if resp.Data == nil {
		return Calendar{}, errors.New("primary calendar response missing data")
	}
	calendar := resp.Data.Calendar
	if calendar.CalendarID == "" {
		calendar.CalendarID = resp.Data.CalendarID
	}
	if calendar.CalendarID == "" {
		return Calendar{}, errors.New("primary calendar response missing calendar_id")
	}
	return calendar, nil
}

func (c *Client) ListCalendarEvents(ctx context.Context, token string, req ListCalendarEventsRequest) (ListCalendarEventsResult, error) {
	if !c.available() || c.coreConfig == nil {
		return ListCalendarEventsResult{}, ErrUnavailable
	}
	if req.CalendarID == "" {
		return ListCalendarEventsResult{}, errors.New("calendar id is required")
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return ListCalendarEventsResult{}, errors.New("tenant access token is required")
	}

	apiReq := &larkcore.ApiReq{
		ApiPath:                   "/open-apis/calendar/v4/calendars/:calendar_id/events",
		HttpMethod:                http.MethodGet,
		PathParams:                larkcore.PathParams{},
		QueryParams:               larkcore.QueryParams{},
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser},
	}
	apiReq.PathParams.Set("calendar_id", req.CalendarID)
	if req.StartTime != "" {
		apiReq.QueryParams.Set("start_time", req.StartTime)
	}
	if req.EndTime != "" {
		apiReq.QueryParams.Set("end_time", req.EndTime)
	}
	if req.PageSize > 0 {
		apiReq.QueryParams.Set("page_size", fmt.Sprintf("%d", req.PageSize))
	}
	if req.PageToken != "" {
		apiReq.QueryParams.Set("page_token", req.PageToken)
	}
	if req.SyncToken != "" {
		apiReq.QueryParams.Set("sync_token", req.SyncToken)
	}

	apiResp, err := larkcore.Request(ctx, apiReq, c.coreConfig, larkcore.WithTenantAccessToken(tenantToken))
	if err != nil {
		return ListCalendarEventsResult{}, err
	}
	if apiResp == nil {
		return ListCalendarEventsResult{}, errors.New("list calendar events failed: empty response")
	}
	resp := &listCalendarEventsResponse{ApiResp: apiResp}
	if err := apiResp.JSONUnmarshalBody(resp, c.coreConfig); err != nil {
		return ListCalendarEventsResult{}, err
	}
	if !resp.Success() {
		return ListCalendarEventsResult{}, fmt.Errorf("list calendar events failed: %s", resp.Msg)
	}

	result := ListCalendarEventsResult{}
	if resp.Data != nil {
		items := resp.Data.Items
		if len(items) == 0 {
			items = resp.Data.Events
		}
		result.Items = items
		result.PageToken = resp.Data.PageToken
		result.HasMore = resp.Data.HasMore
		result.SyncToken = resp.Data.SyncToken
	}
	return result, nil
}

func (c *Client) CreateCalendarEvent(ctx context.Context, token string, req CreateCalendarEventRequest) (CalendarEvent, error) {
	if !c.available() || c.coreConfig == nil {
		return CalendarEvent{}, ErrUnavailable
	}
	if req.CalendarID == "" {
		return CalendarEvent{}, errors.New("calendar id is required")
	}
	if req.Summary == "" {
		return CalendarEvent{}, errors.New("summary is required")
	}
	if req.StartTime == 0 || req.EndTime == 0 {
		return CalendarEvent{}, errors.New("start and end times are required")
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return CalendarEvent{}, errors.New("tenant access token is required")
	}

	payload := map[string]any{
		"summary": req.Summary,
		"start_time": map[string]string{
			"timestamp": fmt.Sprintf("%d", req.StartTime),
		},
		"end_time": map[string]string{
			"timestamp": fmt.Sprintf("%d", req.EndTime),
		},
	}
	if req.Description != "" {
		payload["description"] = req.Description
	}

	apiReq := &larkcore.ApiReq{
		ApiPath:                   "/open-apis/calendar/v4/calendars/:calendar_id/events",
		HttpMethod:                http.MethodPost,
		PathParams:                larkcore.PathParams{},
		QueryParams:               larkcore.QueryParams{},
		Body:                      payload,
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser},
	}
	apiReq.PathParams.Set("calendar_id", req.CalendarID)

	apiResp, err := larkcore.Request(ctx, apiReq, c.coreConfig, larkcore.WithTenantAccessToken(tenantToken))
	if err != nil {
		return CalendarEvent{}, err
	}
	if apiResp == nil {
		return CalendarEvent{}, errors.New("create calendar event failed: empty response")
	}
	resp := &createCalendarEventResponse{ApiResp: apiResp}
	if err := apiResp.JSONUnmarshalBody(resp, c.coreConfig); err != nil {
		return CalendarEvent{}, err
	}
	if !resp.Success() {
		return CalendarEvent{}, fmt.Errorf("create calendar event failed: %s", resp.Msg)
	}
	if resp.Data == nil {
		return CalendarEvent{}, errors.New("create calendar event response missing data")
	}
	result := resp.Data.Event
	if result.EventID == "" {
		result.EventID = resp.Data.EventID
	}
	if result.EventID == "" {
		return CalendarEvent{}, errors.New("create calendar event response missing event_id")
	}
	return result, nil
}

func (c *Client) CreateCalendarEventAttendees(ctx context.Context, token string, req CreateCalendarEventAttendeesRequest) error {
	if !c.available() || c.coreConfig == nil {
		return ErrUnavailable
	}
	if req.CalendarID == "" {
		return errors.New("calendar id is required")
	}
	if req.EventID == "" {
		return errors.New("event id is required")
	}
	if len(req.Attendees) == 0 {
		return errors.New("attendees are required")
	}
	tenantToken := c.tenantToken(token)
	if tenantToken == "" {
		return errors.New("tenant access token is required")
	}

	apiReq := &larkcore.ApiReq{
		ApiPath:                   "/open-apis/calendar/v4/calendars/:calendar_id/events/:event_id/attendees",
		HttpMethod:                http.MethodPost,
		PathParams:                larkcore.PathParams{},
		QueryParams:               larkcore.QueryParams{},
		Body:                      map[string]any{"attendees": req.Attendees},
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser},
	}
	apiReq.PathParams.Set("calendar_id", req.CalendarID)
	apiReq.PathParams.Set("event_id", req.EventID)

	apiResp, err := larkcore.Request(ctx, apiReq, c.coreConfig, larkcore.WithTenantAccessToken(tenantToken))
	if err != nil {
		return err
	}
	if apiResp == nil {
		return errors.New("create calendar event attendees failed: empty response")
	}
	resp := &createCalendarEventAttendeesResponse{ApiResp: apiResp}
	if err := apiResp.JSONUnmarshalBody(resp, c.coreConfig); err != nil {
		return err
	}
	if !resp.Success() {
		return fmt.Errorf("create calendar event attendees failed: %s", resp.Msg)
	}
	return nil
}
