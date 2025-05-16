package landscape

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"resty.dev/v3"
)

type Client struct {
	api      *resty.Client
	endpoint string
	email    *string
	password *string
	account  *string
	key      *string
	secret   *string
	expiry   *time.Duration

	Activities *ActivitiesService
}

// Option is the type for functional options.
type Option func(*Client)

// New creates a new client instance with the specified endpoint and options.
func New(endpoint string, options ...Option) *Client {
	c := &Client{
		api:      resty.New().EnableTrace(),
		endpoint: endpoint,
	}
	c.api.SetBaseURL(endpoint + "/api/v2")
	for _, option := range options {
		option(c)
	}
	c.Activities = &ActivitiesService{client: c}
	return c
}

// WithLoginAuth sets the login authentication for the client; it is used to authenticate
// using the login authentication method.
func WithLoginAuth(email string, password string, account *string) Option {
	return func(c *Client) {
		if email != "" {
			c.email = &email
		}
		if password != "" {
			c.password = &password
		}
		c.account = account
	}
}

// WithSSOAuth sets the SSO authentication for the client; it is used to authenticate
// using the SSO authentication method.
func WithSSOAuth(key string, secret string, expiry *time.Duration) Option {
	return func(c *Client) {
		if key != "" {
			c.key = &key
		}
		if secret != "" {
			c.secret = &secret
		}
		c.expiry = expiry
	}
}

// WithDebug sets the debug mode for the client; it is used to enable or disable
// verbose logging of the HTTP requests and responses.
func WithDebug() Option {
	return func(c *Client) {
		slog.Debug("activating debug mode")
		c.api.SetDebug(true)
	}
}

// WithTrace sets the trace mode for the client; it is used to enable or disable
// detailed logging of the HTTP requests and responses, including the request
// and response headers, body, and status code.
func WithTrace() Option {
	return func(c *Client) {
		slog.Debug("activating trace mode")
		c.api.SetTrace(true)
	}
}

// WithInsecureSkipVerify sets the InsecureSkipVerify option for the TLS client
// configuration. This option is used to skip verification of the server's
// certificate chain and host name. It is useful for testing purposes, but
// should not be used in production environments.
// It is equivalent to setting the `-k` or `--insecure` option in curl.
func WithInsecureSkipVerify() Option {
	return func(c *Client) {
		slog.Debug("skipping TLS certificate verification")
		c.api.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
}

// WithGenerateCurlCommand sets the EnableGenerateCurlCmd option for the client.
// This option is used to enable or disable the generation of curl commands
// for the HTTP requests made by the client. It is useful for debugging and
// testing purposes, as it allows you to see the equivalent curl command for
// the HTTP request being made.
func WithGenerateCurlCommand() Option {
	return func(client *Client) {
		client.api.EnableGenerateCurlCmd()
	}
}

func (c *Client) Login() error {
	if c.email != nil && c.password != nil {
		return c.loginWithPassword()
	}
	if c.key != nil && c.secret != nil {
		return c.loginWithSSO()
	}
	return fmt.Errorf("no authentication method provided")
}

func (c *Client) loginWithPassword() error {
	form := map[string]string{
		"email":    *c.email,
		"password": *c.password,
	}
	if c.account != nil {
		form["account"] = *c.account
	}

	type payload struct {
		Accounts []struct {
			Default bool   `json:"default"`
			Name    string `json:"name"`
			Title   string `json:"title"`
		} `json:"accounts"`
		CurrentAccount string `json:"current_account"`
		Email          string `json:"email"`
		Name           string `json:"name"`
		SelfHosted     bool   `json:"self_hosted"`
		Token          string `json:"token"`
	}

	output := &payload{}

	request := c.api.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetFormData(form).
		SetResult(output).
		SetError(output).AllowMethodDeletePayload

	if c.api.IsDebug() {
		slog.Debug("request", "body", request.Body)
	}

	response, err := request.Post("/api/v2/login")

	if err != nil {
		slog.Error("error logging in", "error", err)
		return err
	}

	if response.StatusCode() != http.StatusOK {
		slog.Error("error logging in", "status", response.StatusCode())
		return fmt.Errorf("error logging in: %s", response.Status())
	}

	if c.api.IsDebug() {
		slog.Debug("response", "body", response, "result", output, "error", err)
	}
	if c.api.IsTrace() {
		slog.Debug("trace", "info", response.Request.TraceInfo())
	}

	return err
}

func (c *Client) loginWithSSO() error {
	return nil
}
