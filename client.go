package jenius

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	Logger   *log.Logger
}

// NewClient : this function will always be called when the library is in use
func NewClient() Client {
	return Client{
		// LogLevel is the logging level used by the Jenius library
		// 0: No logging
		// 1: Errors only
		// 2: Errors + informational (default)
		// 3: Errors + informational + debug
		LogLevel: 2,
		Logger:   log.New(os.Stderr, "", log.LstdFlags),
	}
}

// ===================== HTTP CLIENT ================================================
var defHTTPTimeout = 15 * time.Second
var httpClient = &http.Client{Timeout: defHTTPTimeout}

// NewRequest : send new request
func (c *Client) NewRequest(method string, fullPath string, headers map[string]string, body io.Reader) (*http.Request, error) {
	logLevel := c.LogLevel
	logger := c.Logger

	req, err := http.NewRequest(method, fullPath, body)
	if err != nil {
		if logLevel > 0 {
			logger.Println("Request creation failed: ", err)
		}
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
	logger := c.Logger

	command, _ := http2curl.GetCurlCommand(req)
	start := time.Now()
	logger.Println("Start requesting: ", req.RequestURI)
	res, err := httpClient.Do(req)
	if err != nil {
		logger.Println("Request failed. Error : %v , Curl Request : %v", err, command)
		return err
	}
	defer res.Body.Close()
	logger.Println("Completed in ", time.Since(start))
	logger.Println("Curl Request: ", command)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Println("Cannot read response body: ", err)
		return err
	}

	logger.Println("Jenius response: ", resBody, " status code : ", res.StatusCode)

	if v != nil && res.StatusCode == 200 {
		if err = json.Unmarshal(resBody, v); err != nil {
			logger.Println("Jenius status code 200 unmarshall error: ", resBody)
			return err
		}
	}

	if x != nil && res.StatusCode != 200 {
		if err = json.Unmarshal(resBody, x); err != nil {
			logger.Println("Jenius status code not 200 unmarshall error: ", resBody)
			return err
		}
	}

	return err
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
