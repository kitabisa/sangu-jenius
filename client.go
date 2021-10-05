package jenius

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"moul.io/http2curl"
)

// Client struct
type Client struct {
	JeniusBaseUrl       string
	JeniusTokenBaseUrl  string
	JeniusClientId      string
	JeniusClientSecret  string
	JeniusApiKey        string
	JeniusApiSecret     string
	JeniusXChannelId    string
	JeniusOauthTokenUrl string
	JeniusPayRequestUrl string
	JeniusPayStatusUrl  string
	JeniusPayRefundUrl  string

	LogLevel int
	Logger   Logger
}


// NewClient : this function will always be called when the library is in use
func NewClient() Client {
	logOption := LogOption{
		Format:          "text",
		Level:           "info",
		TimestampFormat: "2006-01-02T15:04:05-0700",
		CallerToggle:    false,
	}

	logger := *NewLogger(logOption)
	return Client{
		// LogLevel is the logging level used by the Jenius library
		// 0: No logging
		// 1: Errors only
		// 2: Errors + informational (default)
		// 3: Errors + informational + debug
		LogLevel: 2,
		Logger:     logger,
	}
}

// ===================== HTTP CLIENT ================================================
var defHTTPTimeout = 15 * time.Second
var httpClient = &http.Client{Timeout: defHTTPTimeout}

// NewRequest : send new request
func (c *Client) NewRequest(method string, fullPath string, headers map[string]string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, fullPath, body)
	if err != nil {
		c.Logger.Info("Request creation failed: %v ", err)
		return nil, err
	}

	if headers != nil {
		for k, vv := range headers {
			req.Header.Set(k, vv)
		}
	}

	if len(headers) == 1 {
		req.SetBasicAuth(c.JeniusClientId, c.JeniusClientSecret)
	}

	return req, nil
}

// ExecuteRequest : execute request
func (c *Client) ExecuteRequest(req *http.Request, v interface{}, x interface{}) (err error) {
	command, _ := http2curl.GetCurlCommand(req)
	start := time.Now()
	c.Logger.Info("Start requesting: %v ", req.URL)
	res, err := httpClient.Do(req)
	if err != nil {
		c.Logger.Error("Request failed. Error : %v , Curl Request : %v", err, command)
		return
	}
	defer res.Body.Close()
	c.Logger.Info("Completed in %v", time.Since(start))
	c.Logger.Info("Curl Request: %v ", command)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.Logger.Info("Cannot read response body: %v ", err)
		return
	}

	c.Logger.Info("Jenius HTTP status response : %d", res.StatusCode)
	c.Logger.Info("Jenius response body : %s", string(resBody))

	if v != nil && res.StatusCode == 200 {
		if err = json.Unmarshal(resBody, v); err != nil {
			c.Logger.Info("Jenius status code 200 unmarshall error: %v ", resBody)
			return
		}
	}

	if x != nil && res.StatusCode != 200 {
		if err = json.Unmarshal(resBody, x); err != nil {
			c.Logger.Info("Jenius status code not 200 unmarshall error: %v ", resBody)
			return
		}
	}

	return
}

// Call the Jenius API at specific `path` using the specified HTTP `method`. The result will be
// given to `v` if there is no error. If any error occurred, the return of this function is the error
// itself, otherwise nil.
func (c *Client) Call(method, path string, header map[string]string, body io.Reader, v interface{}, x interface{}) error {
	req, err := c.NewRequest(method, path, header, body)
	if err != nil {
		return err
	}

	return c.ExecuteRequest(req, v, x)
}

// ===================== END HTTP CLIENT ================================================
