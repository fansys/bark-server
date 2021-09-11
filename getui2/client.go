package getui2

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/mritd/logger"
	"golang.org/x/net/http2"
	"io"
	"net"
	"net/http"
	"time"
)

// Host GeTui urls
const (
	Host = "https://restapi.getui.com"
)

var (
	// TLSDialTimeout is the maximum amount of time a dial will wait for a connect
	// to complete.
	TLSDialTimeout = 20 * time.Second
	// HTTPClientTimeout specifies a time limit for requests made by the
	// HTTPClient. The timeout includes connection time, any redirects,
	// and reading the response body.
	HTTPClientTimeout = 60 * time.Second
	// TCPKeepAlive specifies the keep-alive period for an active network
	// connection. If zero, keep-alives are not enabled.
	TCPKeepAlive = 60 * time.Second
)

// DialTLS is the default dial function for creating TLS connections for
// non-proxied HTTPS requests.
var DialTLS = func(network, addr string, cfg *tls.Config) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   TLSDialTimeout,
		KeepAlive: TCPKeepAlive,
	}
	return tls.DialWithDialer(dialer, network, addr, cfg)
}

// Client represents a connection with the APNs
type Client struct {
	Host        string
	AppId       string
	Certificate tls.Certificate
	Token       *Token
	HTTPClient  *http.Client
}
type connectionCloser interface {
	CloseIdleConnections()
}

// NewClient returns a new Client with an underlying http.Client
//
// If your use case involves multiple long-lived connections, consider using
// the ClientManager, which manages clients for you.
func NewClient(certificate tls.Certificate) *Client {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
	if len(certificate.Certificate) > 0 {
		tlsConfig.BuildNameToCertificate()
	}
	transport := &http2.Transport{
		TLSClientConfig: tlsConfig,
		DialTLS:         DialTLS,
	}
	return &Client{
		HTTPClient: &http.Client{
			Transport: transport,
			Timeout:   HTTPClientTimeout,
		},
		Certificate: certificate,
		Host:        Host,
	}
}

// NewTokenClient returns a new Client with an underlying http.Client configured
// with the correct APNs HTTP/2 transport settings. It does not connect to the APNs
// until the first Notification is sent via the Push method.
//
// As per the Apple APNs Provider API, you should keep a handle on this client
// so that you can keep your connections with APNs open across multiple
// notifications; don’t repeatedly open and close connections. APNs treats rapid
// connection and disconnection as a denial-of-service attack.
func NewTokenClient(token *Token) *Client {
	transport := &http2.Transport{
		DialTLS: DialTLS,
	}
	return &Client{
		Token: token,
		HTTPClient: &http.Client{
			Transport: transport,
			Timeout:   HTTPClientTimeout,
		},
		Host: Host,
	}
}

// Push sends a Notification to the APNs gateway. If the underlying http.Client
// is not currently connected, this method will attempt to reconnect
// transparently before sending the notification. It will return a Response
// indicating whether the notification was accepted or rejected by the APNs
// gateway, or an error if something goes wrong.
//
// Use PushWithContext if you need better cancellation and timeout control.
func (c *Client) Push(n *Notification) (*Response, error) {
	return c.PushWithContext(nil, n)
}

// PushWithContext sends a Notification to the APNs gateway. Context carries a
// deadline and a cancellation signal and allows you to close long running
// requests when the context timeout is exceeded. Context can be nil, for
// backwards compatibility.
//
// If the underlying http.Client is not currently connected, this method will
// attempt to reconnect transparently before sending the notification. It will
// return a Response indicating whether the notification was accepted or
// rejected by the APNs gateway, or an error if something goes wrong.
func (c *Client) PushWithContext(ctx Context, n *Notification) (*Response, error) {
	payload, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%v/v2/%v/push/single/cid", c.Host, c.AppId)
	logger.Infof("GeTui push message: %v", string(payload))
	return c.Execute(ctx, url, payload, true)
}

func (c *Client) GetToken(ctx Context, params map[string]string) (*Response, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%v/v2/%v/auth", c.Host, c.AppId)
	logger.Info("GeTui refresh token start...")
	return c.Execute(ctx, url, payload, false)
}

func (c *Client) Execute(ctx Context, url string, payload []byte, token bool) (*Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	if token {
		c.setTokenHeader(req)
	}
	setHeaders(req)
	httpRes, err := c.requestWithContext(ctx, req)
	if err != nil {
		logger.Errorf("GeTui request error. url: %v, error.: %v", url, err)
		return nil, err
	}
	defer httpRes.Body.Close()

	response := &Response{}
	response.StatusCode = httpRes.StatusCode

	decoder := json.NewDecoder(httpRes.Body)
	if err := decoder.Decode(&response); err != nil && err != io.EOF {
		return &Response{}, err
	}
	return response, nil
}

// CloseIdleConnections closes any underlying connections which were previously
// connected from previous requests but are now sitting idle. It will not
// interrupt any connections currently in use.
func (c *Client) CloseIdleConnections() {
	c.HTTPClient.Transport.(connectionCloser).CloseIdleConnections()
}

func (c *Client) setTokenHeader(r *http.Request) {
	c.Token.GenerateIfExpired(c)
	r.Header.Set("token", c.Token.Token)
}

func setHeaders(r *http.Request) {
	r.Header.Set("Charset", "UTF-8")
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
}
