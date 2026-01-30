package larksdk

import (
	"errors"
	"net/http"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"

	"lark/internal/config"
)

var ErrUnavailable = errors.New("lark sdk unavailable")

// Option configures SDK client initialization.
type Option func(*options)

type options struct {
	httpClient        *http.Client
	tenantAccessToken string
}

// WithHTTPClient overrides the HTTP client used by the SDK.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(o *options) {
		o.httpClient = httpClient
	}
}

// WithTenantAccessToken sets a default tenant access token for requests.
func WithTenantAccessToken(token string) Option {
	return func(o *options) {
		o.tenantAccessToken = token
	}
}

type Client struct {
	sdk               *lark.Client
	coreConfig        *larkcore.Config
	tenantAccessToken string
}

func New(cfg *config.Config, opts ...Option) (*Client, error) {
	if cfg == nil {
		return nil, ErrUnavailable
	}
	if cfg.AppID == "" || cfg.AppSecret == "" {
		return nil, ErrUnavailable
	}
	settings := options{tenantAccessToken: cfg.TenantAccessToken}
	for _, opt := range opts {
		opt(&settings)
	}

	clientOptions := []lark.ClientOptionFunc{
		lark.WithEnableTokenCache(false),
	}
	coreConfig := &larkcore.Config{
		BaseUrl:          lark.FeishuBaseUrl,
		AppId:            cfg.AppID,
		AppSecret:        cfg.AppSecret,
		EnableTokenCache: false,
		AppType:          larkcore.AppTypeSelfBuilt,
	}
	if cfg.BaseURL != "" {
		clientOptions = append(clientOptions, lark.WithOpenBaseUrl(cfg.BaseURL))
		coreConfig.BaseUrl = cfg.BaseURL
	}
	if settings.httpClient != nil {
		clientOptions = append(clientOptions, lark.WithHttpClient(settings.httpClient))
		coreConfig.HttpClient = settings.httpClient
	}

	larkcore.NewLogger(coreConfig)
	larkcore.NewCache(coreConfig)
	larkcore.NewSerialization(coreConfig)
	larkcore.NewHttpClient(coreConfig)

	sdk := lark.NewClient(cfg.AppID, cfg.AppSecret, clientOptions...)
	return &Client{sdk: sdk, coreConfig: coreConfig, tenantAccessToken: settings.tenantAccessToken}, nil
}

func (c *Client) available() bool {
	return c != nil && c.sdk != nil
}

func (c *Client) tenantToken(token string) string {
	if token != "" {
		return token
	}
	if c == nil {
		return ""
	}
	return c.tenantAccessToken
}
