package landscape

import (
	"github.com/dihedron/go-log-facade/logging"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
	key    string
	secret string
	logger logging.Logger
}

// Option is the type for functional options.
type Option func(*Client)

func New(options ...Option) *Client {
	c := &Client{
		client: resty.
			New().
			EnableTrace(),
	}
	return c
}

func WithAPIEndpoint(url string) Option {
	return func(c *Client) {
		c.client.SetBaseURL(url)
	}
}

func WithAuthInfo(key string, secret string) Option {
	return func(c *Client) {
		c.key = key
		c.secret = secret
		c.client.SetBasicAuth(key, secret)
	}
}

// SignRequest creates the signature for a given request. API requests are HTTPS
// requests that use the HTTP verb GET or POST and a query parameter named action.
// To be able to make a request, you’ll need to know which endpoint to use, what
// the action is, and what parameters it takes. All methods take a list of mandatory
// arguments which you need to pass every time:
// action: The name of the method you want to call
// access_key_id: The access key given to you in the Landscape Web UI. You need to go to your settings section in Landscape to be able to generate it along with the secret key.
// signature_method: The method used to signed the request, always HmacSHA256 for now.
// signature_version: The version of the signature mechanism, always 2 for now.
// timestamp: The time in UTC in ISO 8601 format, used to indicate the validity of the signature.
// version: The version of the API, 2011-08-01 being the current one. It’s in the form of a date.
func (c *Client) SignRequest() (string, error) {
	return "", nil
}
